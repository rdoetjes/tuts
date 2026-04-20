#!/usr/bin/env python3
"""
crypt/tools.py

Helpers to inspect recovered key material and to decrypt ciphertexts
using a recovered key for the repository's XOR+rotate obfuscation.

Provides two subcommands:
  - show-key    : display a recovered key in readable string form
  - decrypt     : decrypt a ciphertext using a recovered key

Usage examples:
  python3 tools.py show-key recovered_key.bin
  python3 tools.py show-key --text recovered_key_text recovered_key.bin
  python3 tools.py decrypt recovered_key.bin cipher.bin > plaintext.txt
  python3 tools.py decrypt --text recovered_key_text recovered_key.bin cipher.bin -o plaintext.bin

Notes:
  - The key format produced by the recovery tool is a binary file of length
    key_len. Ambiguous positions are written as 0x00. If you used the text
    output mode from the recovery tool (human readable candidate list), pass
    that via --text so ambiguous positions (??) are honored.
  - This code implements the same inverse transform as the repo's decrypt:
      start = len(cipher) % key_len
      for each message index i:
          kidx = (start + i) % key_len
          shift = (key[kidx] + i) % 8
          tmp = ror8(cipher[i], shift)
          plain[i] = tmp ^ key[kidx]
"""

from __future__ import annotations

import argparse
import os
import stat
import sys
from typing import List, Optional, Tuple


def ror8(v: int, s: int) -> int:
    s &= 7
    if s == 0:
        return v & 0xFF
    v &= 0xFF
    return ((v >> s) | ((v << (8 - s)) & 0xFF)) & 0xFF


def rol8(v: int, s: int) -> int:
    s &= 7
    if s == 0:
        return v & 0xFF
    v &= 0xFF
    return ((v << s) | (v >> (8 - s))) & 0xFF


def read_file_bytes(path: str) -> bytes:
    with open(path, "rb") as f:
        return f.read()


def parse_recovered_text(path: str) -> Optional[List[Optional[int]]]:
    """
    Parse a text candidate output produced by recover_key (the -t/--text file).
    The file format written by recover_key first line contains hex tokens or '??'
    Example first line:
      54 ?? 4A 6B ??
    We try to find the first non-empty non-comment line and split tokens.

    Returns list of len = key_len with int values for known candidates or None for ambiguous.
    """
    try:
        with open(path, "r", encoding="utf-8", errors="replace") as f:
            for raw in f:
                line = raw.strip()
                if not line or line.startswith("#"):
                    continue
                # first meaningful line may contain tokens separated by space
                tokens = line.split()
                # tokens might be in form "54 ?? 4A"
                parsed: List[Optional[int]] = []
                valid = True
                for t in tokens:
                    if t == "??":
                        parsed.append(None)
                    else:
                        # allow tokens like "0A" or "0a" or "54"
                        try:
                            val = int(t, 16)
                            if not (0 <= val <= 255):
                                valid = False
                                break
                            parsed.append(val)
                        except Exception:
                            valid = False
                            break
                if valid and parsed:
                    return parsed
                # if not valid, keep scanning (maybe the file format differs)
    except FileNotFoundError:
        return None
    return None


def show_key_bin(key_bin: bytes, key_text: Optional[List[Optional[int]]] = None) -> str:
    """
    Return a readable string representation of the recovered key.
    If key_text is provided, it is a list with either int (0..255) or None for ambiguous,
    and it will be used to mark ambiguous positions.
    Otherwise any 0x00 in the binary key will be shown as '?' (ambiguous).
    Printable ASCII bytes are shown as characters, others as \\xHH.
    """
    out_chars = []
    key_len = len(key_bin)
    for i in range(key_len):
        kb = key_bin[i]
        ambiguous = False
        if key_text is not None:
            if i < len(key_text) and key_text[i] is None:
                ambiguous = True
        else:
            # interpret 0x00 as ambiguous marker (recover_key writes 0x00 for ambiguous)
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


def decrypt_with_key(key_bytes: bytes, cipher_bytes: bytes) -> bytes:
    """
    Decrypt cipher_bytes with provided key_bytes, following the repo's transform.

    Note: If key contains ambiguous positions (0x00 that are actually unknown),
    the output will be incorrect at those positions.
    """
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


def cli_show_key(args: argparse.Namespace):
    key_path = args.keyfile
    key_text_path = args.text
    if not os.path.exists(key_path):
        print(f"key file not found: {key_path}", file=sys.stderr)
        sys.exit(2)
    key_bin = read_file_bytes(key_path)
    key_text = None
    if key_text_path:
        parsed = parse_recovered_text(key_text_path)
        if parsed is None:
            print(
                f"warning: could not parse text file {key_text_path}; falling back to binary-only view",
                file=sys.stderr,
            )
        else:
            key_text = parsed
    human = show_key_bin(key_bin, key_text)
    print("Key length:", len(key_bin))
    print("Key (printable chars shown directly, others as \\xHH, ambiguous as '?'):")
    print(human)


def cli_decrypt(args: argparse.Namespace):
    key_path = args.keyfile
    cipher_path = args.cipherfile
    out_path = args.output

    if not os.path.exists(key_path):
        print(f"key file not found: {key_path}", file=sys.stderr)
        sys.exit(2)
    if not os.path.exists(cipher_path):
        print(f"cipher file not found: {cipher_path}", file=sys.stderr)
        sys.exit(2)

    key_bin = read_file_bytes(key_path)
    cipher_bin = read_file_bytes(cipher_path)

    # If a text candidate file is provided, prefer its info for ambiguous positions:
    key_text = None
    if args.text:
        key_text = parse_recovered_text(args.text)
        # If parse succeeded, build key_bytes where ambiguous positions remain 0x00
        if key_text is not None:
            # If text file contains explicit values for some positions, use them;
            # otherwise fall back to key_bin bytes.
            new_key = bytearray(len(key_bin))
            for i in range(len(key_bin)):
                if i < len(key_text) and key_text[i] is not None:
                    new_key[i] = key_text[i]
                else:
                    new_key[i] = key_bin[i]
            key_bin = bytes(new_key)

    try:
        plain = decrypt_with_key(key_bin, cipher_bin)
    except ValueError as e:
        print("Error:", e, file=sys.stderr)
        sys.exit(1)

    if out_path:
        with open(out_path, "wb") as f:
            f.write(plain)
        print(f"Wrote plaintext to {out_path}")
    else:
        # print to stdout; be careful not to write binary to terminal if interactive
        if sys.stdout.isatty():
            # Show a safe textual representation: printable chars, '.' for others
            safe = []
            for b in plain[:2048]:
                if 32 <= b <= 126:
                    safe.append(chr(b))
                elif b in (9, 10, 13):
                    safe.append(chr(b))
                else:
                    safe.append(".")
            print("".join(safe), end="")
            if len(plain) > 2048:
                print("\n... (truncated binary output; use -o to write full plaintext)")
        else:
            # pipe binary to stdout
            sys.stdout.buffer.write(plain)


def build_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(
        prog="tools.py", description="Key inspection and decryption helpers"
    )
    sub = p.add_subparsers(dest="cmd", required=True)

    ps = sub.add_parser("show-key", help="display a recovered key in readable form")
    ps.add_argument(
        "keyfile",
        help="binary recovered key file (e.g. recovered_key or recovered_key.<len>.bin)",
    )
    ps.add_argument(
        "--text",
        "-t",
        dest="text",
        help="optional human-readable recovered_key text file (contains tokens or '??')",
    )
    ps.set_defaults(func=cli_show_key)

    pd = sub.add_parser("decrypt", help="decrypt a ciphertext using a recovered key")
    pd.add_argument("keyfile", help="binary recovered key file")
    pd.add_argument("cipherfile", help="ciphertext file to decrypt (binary)")
    pd.add_argument(
        "--text",
        "-t",
        dest="text",
        help="optional human-readable recovered_key text file (use to override ambiguous positions)",
    )
    pd.add_argument(
        "-o",
        "--output",
        dest="output",
        help="write plaintext to file (binary). If omitted, prints to stdout (text-safe)",
    )
    pd.set_defaults(func=cli_decrypt)

    return p


def main():
    parser = build_parser()
    args = parser.parse_args()
    args.func(args)


if __name__ == "__main__":
    main()
