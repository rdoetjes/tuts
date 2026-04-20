#!/usr/bin/env python3
"""
crypt/hackme.py

End-to-end demo wrapper that automates:
  - generating demo ciphertext + crib fragments
  - recovering per-position candidates
  - auto-resolving printable candidates
  - running an audio-focused beam-search cracker
  - converting the best raw plaintexts to WAV for listening

This script expects to live in the repo at `crypt/hackme.py`. It calls the
binaries and scripts present in the same directory (or adjacent): `make_cipher`,
`recover_key`, `auto_resolve.py`, `audio_cracker.py`, `tools.py`.

Usage:
  python3 crypt/hackme.py [--sr SAMPLE_RATE] [--fmt FORMAT] [--time SECONDS]

Default demo behavior:
  - Generate cipher: `cipher.bin` and crib fragments at offsets [0,16,32,48,64,80]
  - Run `recover_key` using the original key length if available (automatically)
  - Run `auto_resolve.py` to produce an auto-resolved key
  - Run `audio_cracker.py` to refine candidate keys (multi-core, time-limited)
  - Produce WAV files for the top candidate plaintexts found for inspection

Notes:
- This is a demo automation wrapper. It runs external programs that must exist
  in the same directory (or be reachable). It prints progress and basic error
  diagnostics to help you run interactively.
- The script is conservative with defaults but exposes a couple of options.
"""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
import sys
import time
import wave
from glob import glob
from pathlib import Path
from typing import Optional

HERE = Path(__file__).resolve().parent


def find_executable(name: str) -> Optional[str]:
    """
    Prefer local copies in the repo directory; otherwise fallback to PATH.
    """
    local = HERE / name
    if local.exists() and os.access(str(local), os.X_OK):
        return str(local)
    # fallback to system path
    exe = shutil.which(name)
    return exe


def run_cmd(cmd, cwd: Optional[Path] = None, timeout: Optional[float] = None):
    print("+", " ".join(cmd))
    try:
        proc = subprocess.run(
            cmd,
            cwd=str(cwd) if cwd is not None else None,
            timeout=timeout,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
        if proc.stdout:
            print(proc.stdout.strip())
        if proc.stderr:
            print(proc.stderr.strip(), file=sys.stderr)
        return proc
    except subprocess.CalledProcessError as e:
        print("Command failed:", e, file=sys.stderr)
        if e.stdout:
            print("stdout:", e.stdout, file=sys.stderr)
        if e.stderr:
            print("stderr:", e.stderr, file=sys.stderr)
        raise
    except subprocess.TimeoutExpired as e:
        print("Command timed out:", e, file=sys.stderr)
        raise


def to_wav_from_raw(
    in_raw: Path, out_wav: Path, sample_rate: int = 20000, fmt: str = "8u"
):
    """
    Convert a raw plaintext file (as produced by our tools) into a WAV file.
    Supports formats: '8u' (unsigned 8-bit PCM), '8s' (signed 8-bit), '16le'.
    """
    print(f"Converting {in_raw} -> {out_wav} (fmt={fmt}, sr={sample_rate})")
    data = in_raw.read_bytes()
    if fmt == "8u":
        # wave module expects bytes for sample width 1 as unsigned bytes
        with wave.open(str(out_wav), "wb") as w:
            w.setnchannels(1)
            w.setsampwidth(1)
            w.setframerate(sample_rate)
            w.writeframes(data)
    elif fmt == "8s":
        # convert signed ints to unsigned bytes by adding 128
        converted = bytes(((b + 128) & 0xFF) for b in data)
        with wave.open(str(out_wav), "wb") as w:
            w.setnchannels(1)
            w.setsampwidth(1)
            w.setframerate(sample_rate)
            w.writeframes(converted)
    elif fmt == "16le":
        # Expect raw is little-endian signed 16-bit (length should be even)
        if len(data) % 2 != 0:
            print("16-bit raw length is odd, truncating last byte.")
            data = data[:-1]
        with wave.open(str(out_wav), "wb") as w:
            w.setnchannels(1)
            w.setsampwidth(2)
            w.setframerate(sample_rate)
            w.writeframes(data)
    else:
        raise ValueError("Unsupported fmt: " + fmt)


def locate_or_fail(name: str):
    exe = find_executable(name)
    if not exe:
        print(
            f"Required tool '{name}' not found (expected in {HERE} or PATH).",
            file=sys.stderr,
        )
        sys.exit(2)
    return exe


def main():
    ap = argparse.ArgumentParser(
        description="End-to-end demo: generate cipher, recover, crack, produce WAVs"
    )
    ap.add_argument(
        "--sr",
        type=int,
        default=20000,
        help="sample rate to write WAV files (default 20000)",
    )
    ap.add_argument(
        "--fmt",
        choices=["8u", "8s", "16le", "mulaw", "alaw"],
        default="8u",
        help="assumed sample format for producing WAVs and scoring (default 8u)",
    )
    ap.add_argument(
        "--time",
        type=int,
        default=600,
        help="total time budget (seconds) for audio cracker (default 600)",
    )
    ap.add_argument(
        "--cores", type=int, default=0, help="cores to use (0 => auto / all cores)"
    )
    ap.add_argument(
        "--keylen-min", type=int, default=32, help="min key length to try (default 32)"
    )
    ap.add_argument(
        "--keylen-max",
        type=int,
        default=160,
        help="max key length to try (default 160)",
    )
    ap.add_argument(
        "--prefer-printable",
        action="store_true",
        help="prefer printable keys in auto-resolve",
    )
    ap.add_argument(
        "--no-clean", action="store_true", help="don't remove intermediate files"
    )
    opts = ap.parse_args()

    # Locate tools
    make_cipher = locate_or_fail("make_cipher")
    recover_key = locate_or_fail("recover_key")
    # python scripts: we expect them in THIS directory
    auto_resolve = HERE / "auto_resolve.py"
    audio_cracker = HERE / "audio_cracker.py"
    tools_py = HERE / "tools.py"
    for p in (auto_resolve, audio_cracker, tools_py):
        if not p.exists():
            print(f"Required script missing: {p}", file=sys.stderr)
            sys.exit(2)

    # Step 1: generate demo cipher + cribs (multiple offsets)
    cipher_out = HERE / "cipher.bin"
    crib_offsets = [0, 16, 32, 48, 64, 80]
    crib_len = 32
    print("Step 1: Generating demo cipher and cribs...")
    # Use the same cipher_out across calls (make_cipher writes same plaintext -> cipher each time)
    for off in crib_offsets:
        crib_out = HERE / f"crib{off}.bin"
        cmd = [
            str(make_cipher),
            str(cipher_out),
            str(crib_out),
            str(off),
            str(crib_len),
        ]
        run_cmd(cmd, cwd=HERE)

    # Step 2: Determine key length to try for recover_key: use original_key.bin if present
    original_key_bin = HERE / "original_key.bin"
    keylen = None
    if original_key_bin.exists():
        try:
            keylen = len(original_key_bin.read_bytes())
            print("Found original_key.bin; using its length for recover_key:", keylen)
        except Exception:
            keylen = None

    # Step 3: run recover_key with the generated crib pairs
    print("Step 2: Running recover_key to produce candidate lists...")
    # build the pair tokens: cipher:crib:offset
    pair_tokens = [
        f"{cipher_out.name}:{'crib' + str(off) + '.bin'}:{off}" for off in crib_offsets
    ]
    # The recover_key binary expects paths relative to its cwd; ensure we call with cwd=HERE
    recover_args = [str(recover_key)]
    if keylen is not None:
        recover_args += ["-k", str(keylen)]
    recover_args += ["-o", "recovered_key_final", "-t", "recovered_key_final_text"]
    recover_args += pair_tokens
    run_cmd(recover_args, cwd=HERE)

    # Step 4: auto-resolve printable candidates (auto_resolve.py)
    print(
        "Step 3: Running auto_resolve to produce an auto-resolved key and plaintext..."
    )
    ar_cmd = [
        sys.executable,
        str(auto_resolve),
        "--key",
        "recovered_key_final",
        "--candidates",
        "recovered_key_final_text",
        "--out-key",
        "recovered_key_auto",
        "--out-plain",
        "recovered_plain",
        "--show",
    ]
    run_cmd(ar_cmd, cwd=HERE)

    # Step 5: run audio_cracker.py to attempt to refine (time-limited, multithreaded)
    print("Step 4: Running audio_cracker (heuristic beam search for audio)...")
    cores_arg = str(opts.cores) if opts.cores > 0 else "0"
    cracker_cmd = [
        sys.executable,
        str(HERE / "audio_cracker.py"),
        "--cipher",
        str(cipher_out),
        "--key",
        "recovered_key_final",
        "--candidates",
        "recovered_key_final_text",
        "--fmt",
        opts.fmt,
        "--min-k",
        str(opts.keylen_min),
        "--max-k",
        str(opts.keylen_max),
        "--cores",
        cores_arg,
        "--time",
        str(opts.time),
        "--beam",
        "64",
        "--keep",
        "3",
    ]
    if opts.prefer_printable:
        cracker_cmd.append("--prefer-printable")
    # always pass per-index-limit=8 (do not use walrus assignment to an attribute)
    cracker_cmd += ["--per-index-limit", "8"]
    run_cmd(cracker_cmd, cwd=HERE, timeout=opts.time + 30)

    # Step 6: Collect best outputs for the key length(s) we most care about
    # Prefer using the auto_resolve key length if available; else use original key length; else try candidate files produced by cracker.
    candidates = list(HERE.glob("best_plain_*_0.raw"))
    if not candidates:
        print(
            "No best_plain_*.raw files found from audio_cracker. Trying recovered_plain ..."
        )
        if (HERE / "recovered_plain").exists():
            candidates = [HERE / "recovered_plain"]
    else:
        print(f"Found {len(candidates)} candidate raw plaintexts from audio_cracker.")

    # Convert the top few candidate raw plaintexts to WAV for listening.
    tops = sorted(candidates)[:3]
    out_wavs = []
    for raw in tops:
        # derive output name
        stem = raw.stem
        wav_out = HERE / f"{stem}.wav"
        try:
            to_wav_from_raw(raw, wav_out, sample_rate=opts.sr, fmt=opts.fmt)
            out_wavs.append(wav_out)
        except Exception as e:
            print("Failed to convert to WAV:", e, file=sys.stderr)

    print("\n=== DEMO COMPLETE ===")
    print("Generated files (examples):")
    print("  cipher:", cipher_out.name)
    print("  cribs:", ", ".join(f"crib{off}.bin" for off in crib_offsets))
    print("  recovered candidates:", "recovered_key_final, recovered_key_final_text")
    print("  auto-resolved key:", "recovered_key_auto")
    print("  auto-resolved plaintext:", "recovered_plain")
    print("  audio cracker best keys/plain files: best_key_*.bin, best_plain_*.raw")
    if out_wavs:
        print("  converted WAVs for quick listening:")
        for w in out_wavs:
            print("   ", w.name)
    else:
        print("  (no WAVs produced; check the raw outputs in the repo root)")

    if not opts.no_clean:
        print(
            "\nYou can keep the outputs or run `git clean -f` if you want to remove them."
        )
    print("To inspect a recovered key textually:")
    print(
        "  python3 tools.py show-key best_key_<klen>_0.bin --text recovered_key_final_text"
    )
    print(
        "To play a produced WAV on macOS: `afplay <wav>`; on Linux: `aplay <wav>` or `play <wav>` (sox)."
    )

    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("\nInterrupted by user", file=sys.stderr)
        sys.exit(1)
