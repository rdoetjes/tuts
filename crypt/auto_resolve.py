#!/usr/bin/env python3
"""
crypt/auto_resolve.py

Auto-resolve ambiguous recovered key positions by preferring printable ASCII
candidates, then single-candidate positions, and fall back to the best-effort
binary key where no candidates can be chosen.

This script uses the same transform logic as the project's tools (rotate + XOR)
to decrypt a ciphertext once a key has been auto-resolved.

Usage:
  python3 auto_resolve.py \
      --key recovered_key \
      --candidates recovered_key_final_text \
      --cipher cipher.bin \
      --out-key recovered_key_auto.bin \
      --out-plain recovered_plain.bin

Defaults:
  key:        recovered_key
  candidates: recovered_key_text
  cipher:     cipher.bin
  out-key:    recovered_key_auto.bin
  out-plain:  recovered_plain.bin

Notes:
- The candidates file is the human-readable output produced by the recovery
  tool; it contains per-position candidate lists like:
    key[  44] (  3): 0x2B, 0x59, 0xFA
  The script prefers printable ASCII values among the candidates for each index.
- If the external helpers in tools.py are available in the same directory, this
  script will reuse `decrypt_with_key` and `show_key_bin` from that module;
  otherwise it uses an equivalent local implementation.
"""

from __future__ import annotations

import argparse
import os
import re
import sys
from typing import Dict, List, Optional

# Try to import helpers from tools.py if present; otherwise define local equivalents.
try:
    import tools as tools_mod  # type: ignore

    decrypt_with_key = tools_mod.decrypt_with_key  # type: ignore
    show_key_bin = tools_mod.show_key_bin  # type: ignore
except Exception:
    # Local fallback: implement ror8 and decrypt_with_key consistent with tools.py
    def ror8(v: int, s: int) -> int:
        s &= 7
        if s == 0:
            return v & 0xFF
        v &= 0xFF
        return ((v >> s) | ((v << (8 - s)) & 0xFF)) & 0xFF

    def decrypt_with_key(key_bytes: bytes, cipher_bytes: bytes) -> bytes:
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

    def show_key_bin(
        key_bin: bytes, key_text: Optional[List[Optional[int]]] = None
    ) -> str:
        out_chars = []
        key_len = len(key_bin)
        for i in range(key_len):
            kb = key_bin[i]
            ambiguous = False
            if key_text is not None:
                if i < len(key_text) and key_text[i] is None:
                    ambiguous = True
            else:
                if kb == 0x00:
                    ambiguous = True
            if ambiguous:
                out_chars.append("?")
            else:
                if 32 <= kb <= 126:
                    out_chars.append(chr(kb))
                else:
                    out_chars.append("\\x%02X" % kb)
        return "".join(out_chars)


def parse_detailed_candidates(text: str) -> Dict[int, List[int]]:
    """
    Parse detailed per-position candidate lines from the recover_key text output.
    Returns a dict mapping position -> list of candidate byte values (ints 0..255).

    Expected line format (examples):
      key[   0] (  1): 0x54
      key[  44] (  3): 0x2B, 0x59, 0xFA
    """
    out: Dict[int, List[int]] = {}
    # Regex to capture index and the rest of the candidate list
    re_line = re.compile(r"^key\[\s*(\d+)\]\s*\(\s*\d+\s*\)\s*:\s*(.*)$")
    # Candidate hex values like "0x54" optionally separated by commas/spaces
    re_hex = re.compile(r"0x([0-9A-Fa-f]{1,2})")
    for raw in text.splitlines():
        m = re_line.match(raw.strip())
        if not m:
            continue
        idx = int(m.group(1))
        rest = m.group(2)
        vals: List[int] = []
        for hm in re_hex.finditer(rest):
            hv = hm.group(1)
            vals.append(int(hv, 16))
        # If we found hex values, record them
        if vals:
            out[idx] = vals
    return out


def parse_first_token_line(text: str) -> Optional[List[Optional[int]]]:
    """
    Recover the initial token line (the compact line near top) from the candidates file.
    That line contains tokens like: 54 ?? ?? 55 48 ?? ...
    We return a list aligned with key positions: integer value or None for '??'.
    This is a fallback and may be shorter than the detailed list.
    """
    for raw in text.splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        # The first non-empty non-comment line should be the compact token line.
        tokens = line.split()
        out: List[Optional[int]] = []
        ok = True
        for t in tokens:
            if t == "??":
                out.append(None)
            else:
                try:
                    # allow tokens both like "54" or "0x54"
                    if t.startswith("0x") or t.startswith("0X"):
                        v = int(t, 16)
                    else:
                        v = int(t, 16)
                    out.append(v & 0xFF)
                except Exception:
                    ok = False
                    break
        if ok and out:
            return out
    return None


def pick_printable_candidate(cands: List[int]) -> Optional[int]:
    """
    Prefer printable ASCII (32..126). If multiple printable, prefer the first
    printable that looks alphanumeric or typical punctuation. Otherwise return
    the first printable. If no printable candidate, return None.
    """
    printable = [c for c in cands if 32 <= c <= 126]
    if not printable:
        return None
    # Prefer alnum/punctuation that matches typical key chars
    for c in printable:
        ch = chr(c)
        if ch.isalnum() or ch in "-_./:@&^*%$#()[]{}!+,":
            return c
    # else return first printable
    return printable[0]


def auto_resolve(key_bin_path: str, candidates_path: str, out_key_path: str) -> bytes:
    """
    Build an auto-resolved key by reading the binary recovered key and the
    candidate text file, preferring printable candidates when ambiguous.
    Writes the resulting key to out_key_path and returns the key bytes.
    """
    if not os.path.exists(key_bin_path):
        raise FileNotFoundError(f"Binary key file not found: {key_bin_path}")
    if not os.path.exists(candidates_path):
        raise FileNotFoundError(f"Candidates text file not found: {candidates_path}")

    with open(key_bin_path, "rb") as f:
        orig_key = bytearray(f.read())

    with open(candidates_path, "r", encoding="utf-8", errors="replace") as f:
        cand_text = f.read()

    # Parse detailed candidates and compact token line
    detailed = parse_detailed_candidates(cand_text)
    compact = parse_first_token_line(cand_text)

    key_len = len(orig_key)
    resolved = bytearray(orig_key)  # start with binary best-effort (0x00 at ambiguous)

    # Iterate through positions and try to pick a candidate
    for i in range(key_len):
        chosen: Optional[int] = None
        if i in detailed:
            cands = detailed[i]
            if len(cands) == 1:
                chosen = cands[0]
            else:
                chosen = pick_printable_candidate(cands)
                if chosen is None:
                    # no printable candidate: if any candidate equals orig_bin non-zero, prefer it
                    if orig_key[i] != 0x00 and orig_key[i] in cands:
                        chosen = orig_key[i]
                    # otherwise leave ambiguous for now
        elif compact is not None and i < len(compact):
            # Fall back to compact token line: if single token != None pick it
            tok = compact[i]
            if tok is not None:
                chosen = tok
        # If we picked a candidate, set it; otherwise leave orig (may be 0x00)
        if chosen is not None:
            resolved[i] = chosen & 0xFF

    # Write out resolved key
    with open(out_key_path, "wb") as f:
        f.write(bytes(resolved))

    return bytes(resolved)


def main(argv: Optional[List[str]] = None) -> int:
    ap = argparse.ArgumentParser(
        description="Auto-resolve recovered key candidates preferring printable bytes"
    )
    ap.add_argument(
        "--key",
        "-k",
        default="recovered_key",
        help="binary recovered key file (default recovered_key)",
    )
    ap.add_argument(
        "--candidates",
        "-c",
        default="recovered_key_text",
        help="human-readable candidate text file",
    )
    ap.add_argument(
        "--cipher", "-C", default="cipher.bin", help="ciphertext file to decrypt"
    )
    ap.add_argument(
        "--out-key",
        "-o",
        default="recovered_key_auto.bin",
        help="path to write auto-resolved binary key",
    )
    ap.add_argument(
        "--out-plain",
        "-p",
        default="recovered_plain.bin",
        help="path to write decrypted plaintext",
    )
    ap.add_argument(
        "--show",
        action="store_true",
        help="print readable key string to stdout after resolution",
    )
    args = ap.parse_args(argv)

    try:
        resolved_key = auto_resolve(args.key, args.candidates, args.out_key)
    except Exception as e:
        print("Error during auto-resolve:", e, file=sys.stderr)
        return 2

    print(f"Wrote auto-resolved key to: {args.out_key} (length {len(resolved_key)})")

    # Show readable string
    if args.show:
        try:
            key_text_parsed = None
            # try to reuse parse_first_token_line to mark ambiguous positions; fallback to None
            with open(args.candidates, "r", encoding="utf-8", errors="replace") as f:
                key_text_parsed = parse_first_token_line(f.read())
        except Exception:
            key_text_parsed = None
        print("Readable key:", show_key_bin(resolved_key, key_text_parsed))

    # Decrypt ciphertext if available and write
    if os.path.exists(args.cipher):
        try:
            with open(args.cipher, "rb") as f:
                cipher_bytes = f.read()
            plain = decrypt_with_key(resolved_key, cipher_bytes)
            with open(args.out_plain, "wb") as fo:
                fo.write(plain)
            print(
                f"Wrote decrypted plaintext to: {args.out_plain} (length {len(plain)})"
            )
        except Exception as e:
            print("Decryption failed:", e, file=sys.stderr)
            return 3
    else:
        print(f"Cipher file not found, skipping decryption: {args.cipher}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
