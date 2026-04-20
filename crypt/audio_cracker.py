#!/usr/bin/env python3
"""
crypt/audio_cracker.py

Automatic detector + beam solver for XOR+rotate-style obfuscation targeted at
audio/plaintext-like signals. Designed to work as a heuristic, time-limited,
multi-core tool that tries small per-byte hypotheses and assembles promising
full-key candidates using a beam search guided by audio scoring.

Usage (example):
  python3 audio_cracker.py \
    --cipher cipher.bin \
    --key recovered_key_final \
    --candidates recovered_key_final_text \
    --min-k 16 --max-k 256 \
    --cores 4 --time 600 --beam 64 --keep 3

Important notes:
- This is a heuristic attacker: it tries many small hypotheses (key bytes 0..255,
  rotation derived from key byte + message index mod 8) and scores decryptions
  for "speech-likeness" rather than English text.
- It prefers printable candidates if `--prefer-printable` is set, which is useful
  if keys were human-readable.
- The script will write best candidate keys and decrypted raw outputs in the
  current directory named like `best_key_<klen>_<rank>.bin` and
  `best_plain_<klen>_<rank>.raw`.
- It tries to be robust if numpy is not installed; however numpy greatly speeds
  up spectral computations.

This is meant to run locally on your machine and operate on files you provide.
Do not run against data you are not authorized to analyze.
"""

from __future__ import annotations

import argparse
import itertools
import math
import os
import struct
import sys
import time
from concurrent.futures import ProcessPoolExecutor, ThreadPoolExecutor, as_completed
from typing import Dict, List, Optional, Tuple

# Try to import numpy for faster spectral features; it's optional.
try:
    import numpy as np  # type: ignore

    HAVE_NUMPY = True
except Exception:
    HAVE_NUMPY = False


# ------------------------------
# Utility functions (byte ops)
# ------------------------------
def rol8(v: int, s: int) -> int:
    s &= 7
    v &= 0xFF
    if s == 0:
        return v
    return ((v << s) | (v >> (8 - s))) & 0xFF


def ror8(v: int, s: int) -> int:
    s &= 7
    v &= 0xFF
    if s == 0:
        return v
    return ((v >> s) | (v << (8 - s))) & 0xFF


# ------------------------------
# Audio sample conversion
# ------------------------------
def interpret_samples(raw: bytes, fmt: str) -> List[int]:
    """
    Convert raw ciphertext bytes into integer sample values for scoring.

    fmt: '8u' (unsigned 8-bit PCM), '8s' (signed 8-bit), '16le' (signed 16-bit little endian)
    For 'mulaw' and 'alaw' we attempt a simple expand to linear (mu-law).
    """
    if fmt == "8u":
        # Convert unsigned 0..255 to signed centered (-128..127)
        return [b - 128 for b in raw]
    elif fmt == "8s":
        return [struct.unpack("b", bytes([b]))[0] for b in raw]
    elif fmt == "16le":
        if len(raw) % 2 != 0:
            raw = raw[:-1]
        res = []
        for i in range(0, len(raw), 2):
            res.append(struct.unpack("<h", raw[i : i + 2])[0])
        return res
    elif fmt == "mulaw":
        # µ-law decoding table per ITU G.711
        return [mulaw_decode_byte(b) for b in raw]
    elif fmt == "alaw":
        return [alaw_decode_byte(b) for b in raw]
    else:
        # Fallback: treat as 8-bit unsigned
        return [b - 128 for b in raw]


# Basic µ-law and A-law decoding implementations
def mulaw_decode_byte(byte_val: int) -> int:
    # Convert 0..255 mu-law to signed 16-bit
    byte_val = ~byte_val & 0xFF
    sign = byte_val & 0x80
    exponent = (byte_val >> 4) & 0x07
    mantissa = byte_val & 0x0F
    sample = ((mantissa << 3) + 0x84) << exponent
    sample = sample - 0x84
    if sign != 0:
        sample = -sample
    return int(sample) >> 2  # scale down to keep values reasonable


def alaw_decode_byte(byte_val: int) -> int:
    byte_val ^= 0x55
    sign = byte_val & 0x80
    exponent = (byte_val >> 4) & 0x07
    mantissa = byte_val & 0x0F
    sample = (mantissa << 4) + 8
    if exponent != 0:
        sample = (sample + 0x100) << (exponent - 1)
    if sign != 0:
        sample = -sample
    return int(sample) >> 2


# ------------------------------
# Scoring functions (audio heuristics)
# ------------------------------
def zero_crossing_rate(samples: List[int]) -> float:
    if not samples:
        return 0.0
    zc = 0
    for i in range(1, len(samples)):
        if (samples[i - 1] >= 0 and samples[i] < 0) or (
            samples[i - 1] < 0 and samples[i] >= 0
        ):
            zc += 1
    return zc / len(samples)


def mean_abs_amplitude(samples: List[int]) -> float:
    if not samples:
        return 0.0
    return sum(abs(x) for x in samples) / len(samples)


def variance(samples: List[int]) -> float:
    if not samples:
        return 0.0
    m = sum(samples) / len(samples)
    return sum((x - m) ** 2 for x in samples) / len(samples)


def spectral_flatness(samples: List[int]) -> float:
    """
    Compute spectral flatness on samples using numpy if available; fallback to
    a cheap approximation (ratio of geometric mean to arithmetic mean of
    magnitudes using an FFT-like window).
    """
    if len(samples) < 8:
        return 1.0
    try:
        if HAVE_NUMPY:
            arr = np.array(samples, dtype=np.float32)
            # apply small window to get stable estimate
            N = min(512, len(arr))
            arr = arr[:N]
            # Hanning window
            w = np.hanning(N)
            spec = np.abs(np.fft.rfft(arr * w))
            spec = spec + 1e-12
            geo = np.exp(np.mean(np.log(spec)))
            arith = np.mean(spec)
            return float(geo / arith)
        else:
            # cheap approximation: compute DFT-like using small subset
            N = min(128, len(samples))
            s = samples[:N]
            # compute magnitudes at few freq bins
            mags = []
            for k in range(1, 9):
                re = 0.0
                im = 0.0
                for n in range(N):
                    angle = 2.0 * math.pi * k * n / N
                    re += s[n] * math.cos(angle)
                    im -= s[n] * math.sin(angle)
                mags.append(math.hypot(re, im) + 1e-12)
            geo = math.exp(sum(math.log(m) for m in mags) / len(mags))
            arith = sum(mags) / len(mags)
            return geo / arith
    except Exception:
        return 1.0


def audio_score(samples: List[int]) -> float:
    """
    Compute a heuristic score where higher is better (more speech-like).

    Heuristic components (weights chosen empirically for speech):
     - spectral_flatness: lower for tonal/speech -> use negative weight
     - zero_crossing_rate: moderate rates preferred (not all zero or all random)
     - mean_abs_amplitude: moderate amplitude preferred (not all zero)
     - variance: some variance is expected (not constant)
    """
    if not samples:
        return -1e6
    mf = mean_abs_amplitude(samples)
    var = variance(samples)
    zcr = zero_crossing_rate(samples)
    sf = spectral_flatness(samples)

    # Normalize heuristics to rough expected ranges
    # these constants are heuristic; they worked reasonably in tests
    score = 0.0
    # prefer some amplitude but not huge
    score += -abs(mf - 20.0) * 0.5
    # prefer variance not tiny
    score += math.log(var + 1.0) * 0.6
    # prefer moderate zcr (~0.05..0.3)
    score += -abs(zcr - 0.12) * 5.0
    # penalize high spectral flatness (white noise)
    score += -sf * 8.0
    # small bias to prefer non-zero signals
    if mf < 1.0:
        score -= 20.0
    return score


# ------------------------------
# Per-position candidate evaluation
# ------------------------------
def evaluate_keybyte_for_positions(args: Tuple) -> Tuple[int, List[Tuple[int, float]]]:
    """
    Worker function used in parallel to evaluate a single key index for all 256
    possible key byte values.

    Args tuple contains:
      (key_index, positions_list, cipher_bytes, key_len, prefer_printable_bool, fmt)
    Returns:
      (key_index, list of (candidate_byte, score) sorted descending)
    """
    key_index, positions, cipher, key_len, prefer_printable, sample_fmt = args
    # positions: list of message indices mapping to this key index
    # cipher: raw bytes
    # For each candidate k in 0..255, decrypt the bytes at positions and collect sample stream
    candidates: List[Tuple[int, float]] = []
    # Pre-extract bytes for these positions
    cbytes = [cipher[pos] for pos in positions]
    for k in range(256):
        # build sample stream as decrypted bytes for these positions
        samples_bytes = []
        for idx, c in zip(positions, cbytes):
            shift = (k + (idx & 0xFF)) % 8
            tmp = ror8(c, shift)
            plain_b = tmp ^ k
            samples_bytes.append(plain_b & 0xFF)
        # convert to signed samples for scoring
        samples = interpret_samples(bytes(samples_bytes), sample_fmt)
        sc = audio_score(samples)
        # apply bias if prefer_printable and k is printable ASCII
        if prefer_printable:
            if 32 <= k <= 126:
                sc += 2.0
        candidates.append((k, sc))
    # sort and return top candidates
    candidates.sort(key=lambda x: x[1], reverse=True)
    # return key_index and top candidate list (we'll limit outside)
    return (key_index, candidates)


# ------------------------------
# Beam search combining candidates
# ------------------------------
def beam_search_combine(
    cipher: bytes,
    key_len: int,
    per_index_candidates: Dict[int, List[int]],
    original_key: bytes,
    beam_width: int,
    time_budget_seconds: float,
    sample_fmt: str,
    top_keep: int = 3,
) -> List[Tuple[float, bytes]]:
    """
    Combine per-index candidate lists into full keys using a beam search.
    Returns list of top (score, key_bytes) tuples.
    """
    start_time = time.monotonic()

    # Order indices by descending ambiguity (more candidates first)
    indices = list(range(key_len))
    # sort by len(candidates) descending (unknown positions count as large)
    indices.sort(key=lambda i: -len(per_index_candidates.get(i, [original_key[i]])))

    # initial beam: use original_key as base
    base_key = bytearray(original_key)
    beams: List[Tuple[float, bytearray]] = []
    # initial score is score of decrypt with base_key
    try:
        base_plain = decrypt_with_key_bytes(base_key, cipher, sample_fmt)
        base_score = audio_score(
            interpret_samples(base_plain[: min(1024, len(base_plain))], sample_fmt)
        )
    except Exception:
        base_score = -1e9
    beams.append((base_score, base_key))

    for idx in indices:
        # time check
        if time.monotonic() - start_time > time_budget_seconds:
            break
        new_beams: List[Tuple[float, bytearray]] = []
        candidates = per_index_candidates.get(idx, [base_key[idx]])
        # If single candidate and same as base, maybe skip heavy scoring
        if len(candidates) == 1 and candidates[0] == base_key[idx]:
            # no change, keep beams unchanged
            continue
        # Expand each beam with each candidate (prune candidates to top K to control branching)
        top_k = min(len(candidates), 16)
        cand_list = candidates[:top_k]
        for score, kb in beams:
            for cand in cand_list:
                new_kb = bytearray(kb)
                new_kb[idx] = cand
                # compute global score by decrypting (or partial) and scoring
                try:
                    plain = decrypt_with_key_bytes(new_kb, cipher, sample_fmt)
                    sc = audio_score(
                        interpret_samples(plain[: min(2048, len(plain))], sample_fmt)
                    )
                except Exception:
                    sc = -1e9
                new_beams.append((sc, new_kb))
        # Prune beams to top beam_width
        new_beams.sort(key=lambda x: x[0], reverse=True)
        beams = new_beams[:beam_width]
    # Return top results (coerce to bytes)
    results: List[Tuple[float, bytes]] = [
        (sc, bytes(kb)) for sc, kb in beams[:top_keep]
    ]
    return results


def decrypt_with_key_bytes(
    key_bytes: bytes, cipher_bytes: bytes, sample_fmt: str
) -> bytes:
    """
    Decrypt entire cipher_bytes using the per-message mapping in the repository:
      start = len(cipher_bytes) % key_len
      kidx = (start + i) % key_len
      shift = (key[kidx] + i) % 8
      tmp = ror8(cipher[i], shift)
      plain[i] = tmp ^ key[kidx]
    Returns a raw bytes object of plaintext bytes.
    """
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


# ------------------------------
# Orchestration
# ------------------------------
def run_for_keylen(
    cipher: bytes,
    original_key: bytes,
    cand_text: str,
    key_len: int,
    cores: int,
    time_budget_seconds: float,
    prefer_printable: bool,
    sample_fmt: str,
    per_index_candidate_limit: int = 8,
    beam_width: int = 64,
    keep_top: int = 3,
) -> Tuple[int, List[Tuple[float, bytes]]]:
    """
    Run detection + per-index brute-force + beam combine for a single key length.
    Returns (key_len, list of top (score, key_bytes)).
    """
    start_time = time.monotonic()
    msg_len = len(cipher)
    start = msg_len % key_len

    # map key index -> list of message indices that map to it
    index_positions: Dict[int, List[int]] = {i: [] for i in range(key_len)}
    for i in range(msg_len):
        kidx = (start + i) % key_len
        index_positions[kidx].append(i)

    # Build candidate map from candidate text (detailed lists)
    # We will use the detailed candidate lists if present; otherwise we'll compute via brute force
    detailed_candidates: Dict[int, List[int]] = {}
    # Parse detailed lines
    for raw in cand_text.splitlines():
        raw = raw.strip()
        if raw.startswith("key["):
            # attempt parse
            try:
                left, right = raw.split(":", 1)
                idx = int(left[left.find("[") + 1 : left.find("]")])
                # find hex tokens in right
                toks = []
                for part in right.strip().split():
                    part = part.strip().strip(",")
                    if part.startswith("0x") or part.startswith("0X"):
                        toks.append(int(part, 16) & 0xFF)
                if toks:
                    detailed_candidates[idx] = toks
            except Exception:
                continue

    # Prepare per-index evaluation tasks for indices that need brute force
    indices = list(range(key_len))
    tasks = []
    for idx in indices:
        if idx in detailed_candidates and len(detailed_candidates[idx]) >= 1:
            # we will use these as possible candidates instead of brute-forcing all 256
            # but still we might want to score them more accurately later
            continue
        # else schedule brute-force for this index
        tasks.append(
            (idx, index_positions[idx], cipher, key_len, prefer_printable, sample_fmt)
        )

    # Evaluate tasks in parallel
    per_index_results: Dict[int, List[Tuple[int, float]]] = {}
    if tasks:
        if cores is None or cores <= 1:
            # sequential
            for t in tasks:
                idx, cand_scores = evaluate_keybyte_for_positions(t)
                per_index_results[idx] = cand_scores
                # time check
                if time.monotonic() - start_time > time_budget_seconds:
                    break
        else:
            with ProcessPoolExecutor(max_workers=cores) as ex:
                futures = {
                    ex.submit(evaluate_keybyte_for_positions, t): t[0] for t in tasks
                }
                for fut in as_completed(futures):
                    try:
                        idx, cand_scores = fut.result()
                        per_index_results[idx] = cand_scores
                    except Exception as e:
                        # fallback: empty
                        per_index_results[futures[fut]] = []
                    if time.monotonic() - start_time > time_budget_seconds:
                        break

    # Now build final per-index candidate lists (ints) limited to per_index_candidate_limit
    per_index_candidates: Dict[int, List[int]] = {}
    for idx in indices:
        if idx in detailed_candidates:
            per_index_candidates[idx] = detailed_candidates[idx][
                :per_index_candidate_limit
            ]
        elif idx in per_index_results:
            # per_index_results[idx] is list of (k, score)
            top = [k for k, s in per_index_results[idx][:per_index_candidate_limit]]
            per_index_candidates[idx] = top
        else:
            # fallback to original key byte if present, else all possible (practical no)
            if idx < len(original_key) and original_key[idx] != 0x00:
                per_index_candidates[idx] = [original_key[idx]]
            else:
                # give a small default set (e.g., common ASCII) to reduce search
                per_index_candidates[idx] = [ord(" "), ord("T"), ord("e"), ord("a")]

    # Optionally prefer printable candidates by reordering per-index lists
    if prefer_printable:
        for idx, vals in list(per_index_candidates.items()):
            printable = [v for v in vals if 32 <= v <= 126]
            nonprint = [v for v in vals if not (32 <= v <= 126)]
            per_index_candidates[idx] = printable + nonprint

    # Run beam search combine
    remaining_time = max(0.1, time_budget_seconds - (time.monotonic() - start_time))
    results = beam_search_combine(
        cipher=cipher,
        key_len=key_len,
        per_index_candidates=per_index_candidates,
        original_key=original_key,
        beam_width=beam_width,
        time_budget_seconds=remaining_time,
        sample_fmt=sample_fmt,
        top_keep=keep_top,
    )
    return (key_len, results)


# ------------------------------
# Main CLI
# ------------------------------
def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        prog="audio_cracker.py",
        description="Automatic audio key recovery (detector + beam solver)",
    )
    p.add_argument("--cipher", "-C", required=True, help="ciphertext file (binary)")
    p.add_argument(
        "--key",
        "-k",
        required=True,
        help="binary recovered key file (best-effort) to use as base",
    )
    p.add_argument(
        "--candidates",
        "-c",
        required=True,
        help="human-readable candidate text from recover_key",
    )
    p.add_argument(
        "--fmt",
        default="8u",
        choices=["8u", "8s", "16le", "mulaw", "alaw"],
        help="assumed sample byte format for scoring (default 8u)",
    )
    p.add_argument(
        "--min-k", type=int, default=16, help="minimum key length to try (default 16)"
    )
    p.add_argument(
        "--max-k", type=int, default=256, help="maximum key length to try (default 256)"
    )
    p.add_argument(
        "--cores",
        type=int,
        default=0,
        help="number of worker processes to use (0 => auto / all cores).",
    )
    p.add_argument(
        "--time",
        type=int,
        default=600,
        help="maximum runtime in seconds for the whole job (default 600 seconds).",
    )
    p.add_argument(
        "--beam",
        type=int,
        default=64,
        help="beam width for combining candidates (default 64).",
    )
    p.add_argument(
        "--keep",
        type=int,
        default=3,
        help="how many top keys to keep per key length (default 3).",
    )
    p.add_argument(
        "--prefer-printable",
        action="store_true",
        help="prefer printable ASCII key bytes when disambiguating.",
    )
    p.add_argument(
        "--per-index-limit",
        type=int,
        default=8,
        help="top candidates per index to keep from per-index search.",
    )
    p.add_argument(
        "--out-dir", default=".", help="output directory for best keys and plaintexts."
    )
    return p.parse_args()


def main() -> int:
    args = parse_args()
    cipher_path = args.cipher
    key_path = args.key
    cand_path = args.candidates
    sample_fmt = args.fmt

    if not os.path.exists(cipher_path):
        print("Cipher file not found:", cipher_path, file=sys.stderr)
        return 2
    if not os.path.exists(key_path):
        print("Key file not found:", key_path, file=sys.stderr)
        return 2
    if not os.path.exists(cand_path):
        print("Candidates file not found:", cand_path, file=sys.stderr)
        return 2

    # Read inputs
    with open(cipher_path, "rb") as f:
        cipher = f.read()
    with open(key_path, "rb") as f:
        original_key = f.read()
    with open(cand_path, "r", encoding="utf-8", errors="replace") as f:
        cand_text = f.read()

    cores = args.cores if args.cores > 0 else (os.cpu_count() or 1)
    time_budget = float(max(1, args.time))
    min_k = max(1, args.min_k)
    max_k = max(min_k, args.max_k)
    beam = max(1, args.beam)
    keep_top = max(1, args.keep)
    per_index_limit = max(1, args.per_index_limit)
    prefer_printable = bool(args.prefer_printable)

    print("audio_cracker: starting")
    print(" cipher:", cipher_path)
    print(" key:", key_path)
    print(" candidates:", cand_path)
    print(" format:", sample_fmt)
    print(" key-lengths:", min_k, "to", max_k)
    print(" cores:", cores, "time budget:", time_budget, "s beam:", beam)
    print(" prefer_printable:", prefer_printable)
    print(" per-index candidate limit:", per_index_limit)
    print(" out-dir:", args.out_dir)
    print(" numpy available:", HAVE_NUMPY)
    sys.stdout.flush()

    start_all = time.monotonic()
    remaining = time_budget
    results_all: List[Tuple[int, List[Tuple[float, bytes]]]] = []

    # iterate key lengths; simple strategy: try a subset if range large to keep within budget
    klens = list(range(min_k, max_k + 1))
    # heuristics: if many key lengths, sample them coarsely then refine; but here we'll iterate
    for klen in klens:
        elapsed = time.monotonic() - start_all
        if elapsed >= time_budget:
            print("Time budget exhausted; stopping key-length loop.")
            break
        # compute allocated time for this keylen - distribute remaining across remaining keylens
        remaining_keylens = max(1, len(klens) - (klens.index(klen)))
        alloc = max(10.0, (time_budget - elapsed) / remaining_keylens)
        alloc = min(alloc, time_budget - elapsed)
        print(f"\n=== Trying key length {klen} (alloc {alloc:.1f}s) ===")
        try:
            klen, bests = run_for_keylen(
                cipher=cipher,
                original_key=original_key,
                cand_text=cand_text,
                key_len=klen,
                cores=cores,
                time_budget_seconds=alloc,
                prefer_printable=prefer_printable,
                sample_fmt=sample_fmt,
                per_index_candidate_limit=per_index_limit,
                beam_width=beam,
                keep_top=keep_top,
            )
            results_all.append((klen, bests))
            # write best results immediately
            for rank, (sc, key_bytes) in enumerate(bests):
                safe_name = f"best_key_{klen}_{rank}.bin"
                out_key = os.path.join(args.out_dir, safe_name)
                with open(out_key, "wb") as fo:
                    fo.write(key_bytes)
                # decrypt full plaintext
                plain = decrypt_with_key_bytes(key_bytes, cipher, sample_fmt)
                out_plain = os.path.join(args.out_dir, f"best_plain_{klen}_{rank}.raw")
                with open(out_plain, "wb") as fo:
                    fo.write(plain)
                print(f"Wrote key {out_key}, plain {out_plain}, score {sc:.3f}")
        except Exception as e:
            print("Error while processing keylen", klen, ":", e, file=sys.stderr)
            continue

    total_elapsed = time.monotonic() - start_all
    print("\n=== Finished ===")
    print(f"Total time: {total_elapsed:.1f} s")
    # Summarize results
    for klen, bests in results_all:
        print(f"Keylen {klen}:")
        for rank, (sc, kb) in enumerate(bests):
            # try to display readable key representation (printable where possible)
            readable = "".join(
                chr(b) if 32 <= b <= 126 else "\\" + "x%02X" % b for b in kb
            )
            print(f"  rank {rank}: score={sc:.3f} key={readable}")

    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print("Interrupted by user", file=sys.stderr)
        sys.exit(1)
