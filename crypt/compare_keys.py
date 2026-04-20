#!/usr/bin/env python3
r"""
compare_keys.py

Sanitized comparison tool for the crypt project.

This script compares an "original" key file and a "recovered" key file
and prints a concise summary of differences plus a readable representation
of the recovered key.

Notes:
- This module uses a raw module docstring so sequences like \xHH are not
  interpreted by Python during parsing.
- By default it expects the following files in the current directory:
    original_key.bin   (original key bytes extracted from main.c)
    recovered_key      (binary best-effort key produced by recover_key)
- You can override the defaults via CLI options.

The recovered key binary format convention used by the repository:
- The recovery utility writes a binary file whose length equals the
  assumed key length.
- Ambiguous positions are represented by byte value 0x00 in the
  "best-effort" binary output (this is a convention used by the tooling).
- The human-readable recover_key text output lists '??' for ambiguous
  positions and hex tokens for single-candidate bytes; this script
  focuses on binary comparison and prints a one-line human-readable
  recovered key where:
    * printable ASCII bytes are shown as characters
    * 0x00 bytes (ambiguous) are shown as '?'
    * other non-printable bytes are shown as '\xHH' (literal backslash)
"""

from __future__ import annotations

import argparse
import os
import sys
from typing import Tuple


def read_bytes(path: str) -> bytes:
    """Read and return the full contents of the file as bytes."""
    try:
        with open(path, "rb") as f:
            return f.read()
    except FileNotFoundError as e:
        raise FileNotFoundError(f"File not found: {path}") from e
    except OSError as e:
        raise OSError(f"Error reading {path}: {e}") from e


def detailed_byte_repr(b: int) -> str:
    """Return detailed representation of a byte: 0xHH and printable char if any."""
    if 32 <= b <= 126:
        return "0x%02X '%c'" % (b, b)
    else:
        return "0x%02X" % b


def recovered_key_string(rec: bytes) -> str:
    """
    Return a compact, human-readable string for a recovered key.

    - Printable ASCII bytes are shown as characters.
    - Bytes equal to 0x00 are shown as '?' (ambiguous marker).
    - Other non-printable bytes are shown as '\\xHH' (literal backslash).
    """
    parts = []
    for b in rec:
        if b == 0x00:
            parts.append("?")
        elif 32 <= b <= 126:
            parts.append(chr(b))
        else:
            parts.append("\\x%02X" % b)
    return "".join(parts)


def compare_keys(orig: bytes, rec: bytes, verbose: bool = False) -> int:
    """
    Compare original and recovered keys and print results.

    Returns the number of differing byte positions (within overlapping length).
    """
    n_orig = len(orig)
    n_rec = len(rec)
    print(f"Original key length: {n_orig}")
    print(f"Recovered key length: {n_rec}")
    print()

    min_len = min(n_orig, n_rec)
    diffs = 0
    for i in range(min_len):
        if orig[i] != rec[i]:
            diffs += 1
            print(
                f"pos {i:4d}: orig={detailed_byte_repr(orig[i]):12s}  rec={detailed_byte_repr(rec[i])}"
            )

    if n_orig != n_rec:
        print()
        if n_orig > n_rec:
            print(
                f"Original key is longer by {n_orig - n_rec} bytes (extra original bytes below):"
            )
            limit = n_orig - n_rec
            if limit <= 64 or verbose:
                for i in range(n_rec, n_orig):
                    print(f"  orig[{i:4d}] = {detailed_byte_repr(orig[i])}")
        else:
            print(
                f"Recovered key is longer by {n_rec - n_orig} bytes (extra recovered bytes below):"
            )
            limit = n_rec - n_orig
            if limit <= 64 or verbose:
                for i in range(n_orig, n_rec):
                    print(f"  rec[{i:4d}]  = {detailed_byte_repr(rec[i])}")

    print()
    print(f"Total differing positions (within first {min_len} bytes): {diffs}")
    print()
    print(
        "Recovered key (printable chars shown, ambiguous 0x00 as '?', others as \\xHH):"
    )
    print(recovered_key_string(rec))
    print()

    if verbose:
        print("Original key (printable chars shown, others as \\xHH):")

        def byte_repr(b: int) -> str:
            if 32 <= b <= 126:
                return chr(b)
            else:
                return "\\x%02X" % b

        print("".join(byte_repr(b) for b in orig))
        print()

    return diffs


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        description="Compare original key and recovered key files."
    )
    p.add_argument(
        "-o",
        "--original",
        dest="original",
        default="original_key.bin",
        help="Path to original key file (default: original_key.bin)",
    )
    p.add_argument(
        "-r",
        "--recovered",
        dest="recovered",
        default="recovered_key",
        help="Path to recovered key file (default: recovered_key)",
    )
    p.add_argument(
        "-v",
        "--verbose",
        dest="verbose",
        action="store_true",
        help="Verbose output (show extra context)",
    )
    return p.parse_args()


def main() -> int:
    args = parse_args()

    try:
        orig = read_bytes(args.original)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 2

    try:
        rec = read_bytes(args.recovered)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 2

    compare_keys(orig, rec, verbose=args.verbose)
    return 0


if __name__ == "__main__":
    sys.exit(main())
