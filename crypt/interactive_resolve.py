#!/usr/bin/env python3
"""
crypt/interactive_resolve.py

Interactive command-line tool to disambiguate recovered key bytes and preview
decryption results.

Purpose
-------
- Load a binary "recovered key" file (best-effort key where ambiguous positions
  are 0x00) and a human-readable candidate file produced by `recover_key` (text).
- Load a ciphertext file (the message you want to decrypt).
- Allow you to interactively inspect ambiguous key byte positions, select a
  candidate for a given position, undo selections, auto-pick printable
  candidates, preview the decrypted plaintext with the current key choices,
  and save the final key and plaintext.

Usage
-----
$ python3 interactive_resolve.py \
    --key recovered_key_final \
    --candidates recovered_key_final_text \
    --cipher cipher.bin

If you omit options, defaults are:
  key        -> recovered_key
  candidates -> recovered_key_text
  cipher     -> cipher.bin

Commands (type "help" in the REPL for details)
- help
- status                      Show summary of resolved / ambiguous positions
- list [start] [end]          List candidate summaries for indexes in range
- show IDX                    Show full candidate list for key position IDX
- pick IDX VAL                Pick candidate value VAL (hex like 0x4A or decimal 74)
                              or use an index into the candidate list: "pick IDX #N"
- auto_printable              Auto-pick a printable ASCII candidate where available
- undo IDX                    Undo manual pick for IDX (revert to original recovered byte)
- preview [N]                 Preview decrypted plaintext (first N bytes, default 200)
- decrypt OUTFILE             Write decrypted plaintext to OUTFILE using current key
- savekey OUTFILE             Write current key (binary) to OUTFILE
- showkey                     Print readable key string for current key
- quit / exit                 Exit the interactive session

Notes
-----
- This tool is heuristic and interactive. Use multiple cribs / ciphertexts
  passed to the recovery tool to improve candidate lists before running this.
- This script reimplements the same inverse transform used by the project's
  decrypt logic:
      start = msg_len % key_len
      kidx = (start + i) % key_len
      shift = (keybyte + i) % 8
      tmp = ror8(cipher[i], shift)
      plain = tmp ^ keybyte

Author: interactive helper for the crypt project
"""

from __future__ import annotations

import argparse
import os
import re
import sys
from typing import Dict, List, Optional, Tuple

# Constants
DEFAULT_KEY_BIN = "recovered_key"
DEFAULT_CAND_TEXT = "recovered_key_text"
DEFAULT_CIPHER = "cipher.bin"
MAX_PREVIEW = 200
CAND_LINE_RE = re.compile(r"^key\[\s*(\d+)\]\s*\(\s*\d+\s*\)\s*:\s*(.*)$")
HEX_RE = re.compile(r"0x([0-9A-Fa-f]{1,2})")
TOKEN_LINE_RE = re.compile(
    r"^(?:[0-9A-Fa-f]{1,2}|0x[0-9A-Fa-f]{1,2}|\?\?)(?:\s+(?:[0-9A-Fa-f]{1,2}|0x[0-9A-Fa-f]{1,2}|\?\?))*$"
)


def ror8(v: int, s: int) -> int:
    s &= 7
    v &= 0xFF
    if s == 0:
        return v
    return ((v >> s) | ((v << (8 - s)) & 0xFF)) & 0xFF


def rol8(v: int, s: int) -> int:
    s &= 7
    v &= 0xFF
    if s == 0:
        return v
    return ((v << s) | (v >> (8 - s))) & 0xFF


def parse_detailed_candidates(text: str) -> Dict[int, List[int]]:
    """
    Parse the detailed per-position candidate lines from the recover_key text output.
    Return map: index -> list of candidate byte ints.
    """
    out: Dict[int, List[int]] = {}
    for raw in text.splitlines():
        m = CAND_LINE_RE.match(raw.strip())
        if not m:
            continue
        idx = int(m.group(1))
        rest = m.group(2)
        vals: List[int] = []
        for hm in HEX_RE.finditer(rest):
            vals.append(int(hm.group(1), 16) & 0xFF)
        if vals:
            out[idx] = vals
    return out


def parse_compact_token_line(text: str) -> Optional[List[Optional[int]]]:
    """
    Parse the compact first token line (e.g. "54 ?? ?? 55 48 ?? ...")
    Return list of ints or None for '??'.
    """
    for raw in text.splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        # simple check to see if line looks like the compact token line
        if TOKEN_LINE_RE.match(line):
            toks = line.split()
            out: List[Optional[int]] = []
            ok = True
            for t in toks:
                if t == "??":
                    out.append(None)
                else:
                    try:
                        # allow hex or bare hex
                        v = int(t, 16)
                        out.append(v & 0xFF)
                    except Exception:
                        ok = False
                        break
            if ok:
                return out
    return None


def read_file_text(path: str) -> str:
    with open(path, "r", encoding="utf-8", errors="replace") as f:
        return f.read()


def read_file_bytes(path: str) -> bytes:
    with open(path, "rb") as f:
        return f.read()


def build_candidate_map(cand_text: str, key_len: int) -> Dict[int, List[int]]:
    """
    Build candidate lists for positions 0..key_len-1.
    Use detailed parses when available, else fallback to the compact token line.
    """
    detailed = parse_detailed_candidates(cand_text)
    compact = parse_compact_token_line(cand_text)
    res: Dict[int, List[int]] = {}

    # initialize with detailed where available
    for idx, vals in detailed.items():
        if 0 <= idx < key_len:
            res[idx] = vals.copy()

    # If compact exists and some positions missing, use it as fallback
    if compact is not None:
        for i in range(min(len(compact), key_len)):
            if i not in res:
                tok = compact[i]
                if tok is None:
                    # ambiguous; leave absent -> will be treated as unknown
                    continue
                else:
                    res[i] = [tok]

    return res


def count_single(cmap: Dict[int, List[int]]) -> int:
    return sum(1 for vals in cmap.values() if len(vals) == 1)


def count_ambiguous(cmap: Dict[int, List[int]], key_len: int, orig_key: bytes) -> int:
    total = 0
    for i in range(key_len):
        if i in cmap:
            if len(cmap[i]) > 1:
                total += 1
        else:
            # if not present in cmap, check if orig_key non-zero (treated as resolved) or zero (ambiguous)
            if orig_key[i] == 0x00:
                total += 1
    return total


def show_key_readable(key_bytes: bytes) -> str:
    parts = []
    for b in key_bytes:
        if b == 0x00:
            parts.append("?")
        elif 32 <= b <= 126:
            parts.append(chr(b))
        else:
            parts.append("\\x%02X" % b)
    return "".join(parts)


def decrypt_with_key_bytes(key_bytes: bytes, cipher_bytes: bytes) -> bytes:
    if len(key_bytes) == 0:
        raise ValueError("key length must be > 0")
    key_len = len(key_bytes)
    msg_len = len(cipher_bytes)
    start = msg_len % key_len
    out = bytearray(msg_len)
    for i in range(msg_len):
        kidx = (start + i) % key_len
        k = key_bytes[kidx]
        shift = (k + (i & 0xFF)) % 8
        tmp = ror8(cipher_bytes[i], shift)
        p = tmp ^ k
        out[i] = p & 0xFF
    return bytes(out)


def safe_print_preview(plain: bytes, n: int = MAX_PREVIEW) -> None:
    show = plain[:n]
    out_chars = []
    for b in show:
        if b in (9, 10, 13) or 32 <= b <= 126:
            out_chars.append(chr(b))
        else:
            out_chars.append(".")
    print("".join(out_chars))
    if len(plain) > n:
        print("... (truncated, showing first %d bytes)" % n)


class InteractiveResolver:
    def __init__(self, key_bin: str, cand_txt: str, cipher_file: str) -> None:
        if not os.path.exists(key_bin):
            raise FileNotFoundError("key binary not found: %s" % key_bin)
        if not os.path.exists(cand_txt):
            raise FileNotFoundError("candidates text not found: %s" % cand_txt)
        if not os.path.exists(cipher_file):
            raise FileNotFoundError("cipher file not found: %s" % cipher_file)

        self.key_bin = key_bin
        self.cand_txt = cand_txt
        self.cipher_file = cipher_file

        self.orig_key = bytearray(read_file_bytes(key_bin))
        self.key_len = len(self.orig_key)
        self.cand_text = read_file_text(cand_txt)
        self.candidates = build_candidate_map(self.cand_text, self.key_len)
        self.cipher = read_file_bytes(cipher_file)

        # current_key holds user choices; start from orig_key (0x00 for ambiguous)
        self.current_key = bytearray(self.orig_key)
        # keep a manual history: map idx -> list of previous bytes (for undo)
        self.history: Dict[int, List[int]] = {}

        # track which positions were manually set by user (vs original)
        self.manual_set: Dict[int, bool] = {}

    def status(self) -> None:
        resolved = sum(
            1
            for i in range(self.key_len)
            if (i in self.candidates and len(self.candidates[i]) == 1)
            or (i not in self.candidates and self.orig_key[i] != 0x00)
            or (i in self.manual_set and self.manual_set[i])
        )
        ambiguous = count_ambiguous(self.candidates, self.key_len, self.orig_key)
        print("Key length:", self.key_len)
        print("Resolved (single-candidate or chosen): %d" % resolved)
        print("Ambiguous positions remaining: %d" % ambiguous)
        print("Manual picks made: %d" % len(self.manual_set))
        print("Ciphertext length: %d" % len(self.cipher))

    def list_range(self, start: int = 0, end: Optional[int] = None) -> None:
        if end is None:
            end = min(start + 30, self.key_len - 1)
        end = min(end, self.key_len - 1)
        for i in range(start, end + 1):
            if i in self.candidates:
                vals = self.candidates[i]
                if len(vals) == 1:
                    v = vals[0]
                    display = "0x%02X '%s'" % (v, chr(v) if 32 <= v <= 126 else ".")
                    print("key[%3d]: single %s" % (i, display))
                else:
                    sample = ", ".join("0x%02X" % v for v in vals[:8])
                    if len(vals) > 8:
                        sample += ", ..."
                    print("key[%3d]: %d candidates: %s" % (i, len(vals), sample))
            else:
                if self.orig_key[i] != 0x00:
                    b = self.orig_key[i]
                    print(
                        "key[%3d]: original single 0x%02X '%s'"
                        % (i, b, chr(b) if 32 <= b <= 126 else ".")
                    )
                else:
                    print("key[%3d]: (no candidates, ambiguous; original=0x00)" % (i,))

    def show_candidates(self, idx: int) -> None:
        if idx < 0 or idx >= self.key_len:
            print("Index out of range")
            return
        if idx in self.candidates:
            vals = self.candidates[idx]
            print("Candidates for key[%d] (%d):" % (idx, len(vals)))
            for n, v in enumerate(vals):
                printable = chr(v) if 32 <= v <= 126 else "."
                print("  #%02d : 0x%02X  '%s'" % (n, v, printable))
        else:
            if self.orig_key[idx] != 0x00:
                v = self.orig_key[idx]
                print(
                    "No candidate list; original has single byte: 0x%02X '%s'"
                    % (v, chr(v) if 32 <= v <= 126 else ".")
                )
            else:
                print("No candidate list and original was ambiguous (0x00).")

        cur = self.current_key[idx]
        print(
            "Current chosen value for key[%d]: 0x%02X %s"
            % (idx, cur, ("(manual)" if idx in self.manual_set else ""))
        )

    def pick_candidate_by_value(self, idx: int, val: int) -> None:
        if idx < 0 or idx >= self.key_len:
            print("Index out of range")
            return
        if idx not in self.history:
            self.history[idx] = []
        self.history[idx].append(self.current_key[idx])
        self.current_key[idx] = val & 0xFF
        self.manual_set[idx] = True
        print("Set key[%d] = 0x%02X" % (idx, val & 0xFF))

    def pick_candidate_by_list_index(self, idx: int, list_index: int) -> None:
        if idx not in self.candidates:
            print("No candidate list for this index.")
            return
        vals = self.candidates[idx]
        if list_index < 0 or list_index >= len(vals):
            print("Candidate list index out of range (0..%d)" % (len(vals) - 1))
            return
        self.pick_candidate_by_value(idx, vals[list_index])

    def undo(self, idx: int) -> None:
        if idx not in self.history or not self.history[idx]:
            print("No history to undo for index %d" % idx)
            return
        prev = self.history[idx].pop()
        self.current_key[idx] = prev
        if not self.history[idx]:
            # if no more history entries, remove manual flag if it wasn't really manual before
            if idx in self.manual_set:
                del self.manual_set[idx]
        print("Reverted key[%d] to 0x%02X" % (idx, prev))

    def auto_pick_printable(self) -> None:
        picked = 0
        for idx, vals in self.candidates.items():
            if len(vals) <= 1:
                continue
            # find first printable candidate
            chosen = None
            for v in vals:
                if 32 <= v <= 126:
                    # prefer letters, digits or common punctuation
                    ch = chr(v)
                    if ch.isalnum() or ch in "-_./:@&^*%$#()[]{}!+,":
                        chosen = v
                        break
            if chosen is None:
                # fallback to any printable
                for v in vals:
                    if 32 <= v <= 126:
                        chosen = v
                        break
            if chosen is not None:
                # apply pick
                if idx not in self.history:
                    self.history[idx] = []
                self.history[idx].append(self.current_key[idx])
                self.current_key[idx] = chosen
                self.manual_set[idx] = True
                picked += 1
        print("Auto-picked printable candidates for %d positions" % picked)

    def preview(self, n: int = MAX_PREVIEW) -> None:
        try:
            plain = decrypt_with_key_bytes(bytes(self.current_key), self.cipher)
        except Exception as e:
            print("Decryption error:", e)
            return
        safe_print_preview(plain, n)

    def decrypt_to_file(self, outfile: str) -> None:
        try:
            plain = decrypt_with_key_bytes(bytes(self.current_key), self.cipher)
            with open(outfile, "wb") as fo:
                fo.write(plain)
            print("Wrote plaintext to", outfile)
        except Exception as e:
            print("Failed to decrypt/write:", e)

    def save_key(self, outfile: str) -> None:
        with open(outfile, "wb") as fo:
            fo.write(bytes(self.current_key))
        print("Wrote current key to", outfile)

    def show_readable_key(self) -> None:
        print("Readable key:")
        print(show_key_readable(bytes(self.current_key)))

    # REPL
    def repl(self) -> None:
        print("Interactive key resolver")
        print("Key file:", self.key_bin)
        print("Candidates file:", self.cand_txt)
        print("Cipher file:", self.cipher_file)
        print("Key length:", self.key_len)
        print("Type 'help' for commands.")
        while True:
            try:
                raw = input("ir> ").strip()
            except (EOFError, KeyboardInterrupt):
                print()
                print("Exiting.")
                break
            if not raw:
                continue
            parts = raw.split()
            cmd = parts[0].lower()
            args = parts[1:]

            if cmd in ("quit", "exit"):
                print("Bye.")
                break
            elif cmd == "help":
                self.print_help()
            elif cmd == "status":
                self.status()
            elif cmd == "list":
                if len(args) == 0:
                    self.list_range(0, min(30, self.key_len - 1))
                elif len(args) == 1:
                    try:
                        s = int(args[0], 0)
                        self.list_range(s, min(s + 30, self.key_len - 1))
                    except Exception:
                        print("Invalid start index")
                elif len(args) >= 2:
                    try:
                        s = int(args[0], 0)
                        e = int(args[1], 0)
                        self.list_range(s, e)
                    except Exception:
                        print("Invalid start/end")
            elif cmd == "show":
                if len(args) != 1:
                    print("Usage: show IDX")
                    continue
                try:
                    idx = int(args[0], 0)
                except Exception:
                    print("Invalid index")
                    continue
                self.show_candidates(idx)
            elif cmd == "pick":
                if len(args) < 2:
                    print(
                        "Usage: pick IDX VAL   or  pick IDX #N   (N = candidate list index)"
                    )
                    continue
                try:
                    idx = int(args[0], 0)
                except Exception:
                    print("Invalid index")
                    continue
                valtoken = args[1]
                if valtoken.startswith("#"):
                    try:
                        li = int(valtoken[1:], 0)
                        self.pick_candidate_by_list_index(idx, li)
                    except Exception:
                        print("Invalid list index")
                else:
                    # allow hex like 0x4A or decimal
                    try:
                        v = int(valtoken, 0)
                        self.pick_candidate_by_value(idx, v)
                    except Exception:
                        print("Invalid value token")
            elif cmd == "undo":
                if len(args) != 1:
                    print("Usage: undo IDX")
                    continue
                try:
                    idx = int(args[0], 0)
                except Exception:
                    print("Invalid index")
                    continue
                self.undo(idx)
            elif cmd == "auto_printable":
                self.auto_pick_printable()
            elif cmd == "preview":
                if len(args) == 0:
                    self.preview()
                else:
                    try:
                        n = int(args[0], 0)
                        self.preview(n)
                    except Exception:
                        print("Invalid number")
            elif cmd == "decrypt":
                if len(args) != 1:
                    print("Usage: decrypt OUTFILE")
                    continue
                self.decrypt_to_file(args[0])
            elif cmd == "savekey":
                if len(args) != 1:
                    print("Usage: savekey OUTFILE")
                    continue
                self.save_key(args[0])
            elif cmd == "showkey":
                self.show_readable_key()
            elif cmd == "helpcmds":
                # alias for compatibility
                self.print_help()
            else:
                print("Unknown command:", cmd)
                print("Type 'help' for a list of commands.")

    def print_help(self) -> None:
        print(
            """
Available commands:
  help                 Show this help
  status               Show summary: resolved / ambiguous positions
  list [start] [end]   List candidate summaries for key indexes in range
  show IDX             Show full candidate list for key index IDX
  pick IDX VAL         Pick candidate value VAL (hex 0xHH or decimal)
                       or pick IDX #N to pick N-th candidate from list
  auto_printable       Auto-pick printable ASCII candidate where available
  undo IDX             Undo manual pick for IDX (revert to previous value)
  preview [N]          Preview decrypted plaintext (first N bytes, default 200)
  decrypt OUTFILE      Write decrypted plaintext to OUTFILE using current key
  savekey OUTFILE      Write current key (binary) to OUTFILE
  showkey              Print a readable representation of the current key
  quit / exit          Exit the interactive session
"""
        )


def main(argv: Optional[List[str]] = None) -> int:
    ap = argparse.ArgumentParser(description="Interactive key disambiguation & preview")
    ap.add_argument(
        "--key", "-k", default=DEFAULT_KEY_BIN, help="binary recovered key file"
    )
    ap.add_argument(
        "--candidates",
        "-c",
        default=DEFAULT_CAND_TEXT,
        help="human-readable candidate text file",
    )
    ap.add_argument(
        "--cipher", "-C", default=DEFAULT_CIPHER, help="ciphertext file to decrypt"
    )
    args = ap.parse_args(argv)

    try:
        ir = InteractiveResolver(args.key, args.candidates, args.cipher)
    except Exception as e:
        print("Initialization error:", e, file=sys.stderr)
        return 2

    try:
        ir.repl()
    except KeyboardInterrupt:
        print("\nInterrupted, exiting.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
