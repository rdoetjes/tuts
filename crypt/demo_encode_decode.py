#!/usr/bin/env python3
"""
crypt/demo_encode_decode.py

Demo script that:
 - synthesizes a short spoken quote using the system TTS (tries macOS `say` then `espeak`)
 - converts the generated audio to 8-bit unsigned PCM (simulating a simple ADC)
 - encrypts the 8-bit audio bytes with the repository XOR+rotate obfuscation (same logic as main.c)
 - writes a "sniffer" WAV containing the ciphertext bytes (as 8-bit unsigned PCM)
 - decrypts the ciphertext and writes a decoded WAV so you can compare

Usage:
  python3 demo_encode_decode.py [--quote "text"] [--sr 20000] [--outdir ./]

Notes:
 - This script tries to call the local TTS tools to synthesize a WAV/AIFF.
 - It is written to be robust: it supports AIFF output from macOS `say` and WAV from `espeak`.
 - It requires Python's standard audio modules (wave, aifc, audioop) available in CPython.
 - The encryption used is the same as the project's `encrypt`/`decrypt` functions:
     start = message_len % key_len
     kidx  = (start + i) % key_len
     shift = (keybyte + i) % 8
     tmp   = plaintext[i] ^ keybyte
     cipher[i] = rol8(tmp, shift)
   and decrypt reverses these steps.
"""

from __future__ import annotations

# Import aifc conditionally — not all Python distributions provide it (some minimal builds).
# If aifc is not available we will fall back to using system tools (afconvert or sox)
# to convert AIFF to WAV when necessary.
try:
    import aifc  # type: ignore

    HAVE_AIFC = True
except Exception:
    aifc = None  # type: ignore
    HAVE_AIFC = False

import argparse

# audioop may not be available in some Python builds (minimal installations).
# Provide a fallback implementation that supplies the small subset we need:
# - audioop.lin2lin(data, inwidth, outwidth)
# - audioop.tomono(data, width, lfac, rfac)
try:
    import audioop  # type: ignore
except Exception:

    class _AudioOpFallback:
        @staticmethod
        def lin2lin(data: bytes, inwidth: int, outwidth: int) -> bytes:
            # Fast paths
            if inwidth == outwidth:
                return data
            # 16-bit signed LE to 8-bit signed
            if inwidth == 2 and outwidth == 1:
                out = bytearray()
                for i in range(0, len(data), 2):
                    if i + 2 > len(data):
                        break
                    val = int.from_bytes(data[i : i + 2], "little", signed=True)
                    # reduce to signed 8-bit by shifting
                    out.append((val >> 8) & 0xFF)
                return bytes(out)
            # 8-bit signed to 16-bit signed LE
            if inwidth == 1 and outwidth == 2:
                out = bytearray()
                for b in data:
                    sval = b if b < 128 else b - 256
                    out += int(sval << 8).to_bytes(2, "little", signed=True)
                return bytes(out)
            # Generic fallback: scalar conversion using signed interpretation
            out = bytearray()
            for i in range(0, len(data), inwidth):
                chunk = data[i : i + inwidth]
                if len(chunk) < inwidth:
                    break
                if inwidth == 1:
                    sval = chunk[0] if chunk[0] < 128 else chunk[0] - 256
                else:
                    sval = int.from_bytes(chunk, "little", signed=True)
                if outwidth == 1:
                    out.append((sval >> 8) & 0xFF)
                else:
                    out += int(sval).to_bytes(outwidth, "little", signed=True)
            return bytes(out)

        @staticmethod
        def tomono(data: bytes, width: int, lfac: float, rfac: float) -> bytes:
            # width: sample width in bytes for each channel (1 or 2)
            out = bytearray()
            step = width * 2
            if width not in (1, 2):
                # Not supported widths: return original
                return data
            for i in range(0, len(data), step):
                if i + step > len(data):
                    break
                if width == 1:
                    left = data[i]
                    right = data[i + 1]
                    left = left if left < 128 else left - 256
                    right = right if right < 128 else right - 256
                else:
                    left = int.from_bytes(data[i : i + 2], "little", signed=True)
                    right = int.from_bytes(data[i + 2 : i + 4], "little", signed=True)
                mono = int(left * lfac + right * rfac)
                if width == 1:
                    out.append((mono + 256) & 0xFF if mono < 0 else mono & 0xFF)
                else:
                    out += int(mono).to_bytes(2, "little", signed=True)
            return bytes(out)

    audioop = _AudioOpFallback()
import shutil
import subprocess
import sys
import tempfile
import wave
from pathlib import Path
from typing import Optional

# KEY copied from main.c (the demo's shared secret)
KEY = (
    "ThKUHIT678t2bkhabxI^&TI*gyhu2lh3jlhbiut987OYHLIUJBLWEUDG3278OYGHLIUQHBSQlHUTI"
    "&^tgJHBHJ6TGKUYGKJYiguykjgiyuti7^TIGyukhb"
)


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


def encrypt_bytes(message: bytes) -> bytes:
    key_len = len(KEY)
    start = len(message) % key_len
    out = bytearray(len(message))
    for i in range(len(message)):
        kidx = (start + i) % key_len
        keybyte = ord(KEY[kidx]) & 0xFF
        plain = message[i] & 0xFF
        shift = (keybyte + (i & 0xFF)) % 8
        tmp = plain ^ keybyte
        out[i] = rol8(tmp, shift)
    return bytes(out)


def decrypt_bytes(message: bytes) -> bytes:
    key_len = len(KEY)
    start = len(message) % key_len
    out = bytearray(len(message))
    for i in range(len(message)):
        kidx = (start + i) % key_len
        keybyte = ord(KEY[kidx]) & 0xFF
        shift = (keybyte + (i & 0xFF)) % 8
        cipher = message[i] & 0xFF
        tmp = ror8(cipher, shift)
        out[i] = tmp ^ keybyte
    return bytes(out)


def synthesize_tts_to_file(text: str, out_path: Path, sr: int) -> None:
    """
    Try macOS 'say' first (produces AIFF), then 'espeak' (produces WAV).
    Writes to out_path. Raises RuntimeError if no TTS available.
    """
    # Try macOS 'say'
    say_exe = shutil.which("say")
    if say_exe:
        # 'say' can write AIFF. We'll request 16-bit LE at requested sample rate if possible.
        # Using data-format requires macOS 'say' support; form the command conservatively.
        # Write to a temporary .aiff then convert via aifc.
        tmp_aiff = out_path.with_suffix(".aiff.tmp")
        cmd = [say_exe, "-o", str(tmp_aiff), "--data-format=LEI16@" + str(sr), text]
        try:
            subprocess.run(
                cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )
            # Convert AIFF to WAV content by reading with aifc and writing via wave
            with aifc.open(str(tmp_aiff), "r") as af:
                nchan = af.getnchannels()
                sampwidth = af.getsampwidth()
                fr = af.getframerate()
                frames = af.readframes(af.getnframes())
            # Write WAV at out_path
            with wave.open(str(out_path), "wb") as wf:
                wf.setnchannels(nchan)
                wf.setsampwidth(sampwidth)
                wf.setframerate(fr)
                wf.writeframes(frames)
            tmp_aiff.unlink(missing_ok=True)
            return
        except Exception:
            # fall through to try espeak
            if tmp_aiff.exists():
                tmp_aiff.unlink(missing_ok=True)

    # Try espeak
    espeak_exe = shutil.which("espeak")
    if espeak_exe:
        # espeak -w out.wav "text"
        cmd = [espeak_exe, "-w", str(out_path), text]
        try:
            subprocess.run(
                cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )
            return
        except Exception:
            pass

    raise RuntimeError(
        "No suitable TTS engine found ('say' or 'espeak'). Please install one."
    )


def read_wav_as_mono_16(path: Path) -> tuple[bytes, int]:
    """
    Read a WAV or AIFF (I/O handled by wave/aifc) and return mono signed 16-bit little-endian samples bytes and sample rate.
    """
    # Try wave first
    try:
        with wave.open(str(path), "rb") as wf:
            nchan = wf.getnchannels()
            sampwidth = wf.getsampwidth()
            fr = wf.getframerate()
            frames = wf.readframes(wf.getnframes())
    except wave.Error:
        # Try to handle AIFF. Prefer Python's aifc if available, otherwise
        # attempt to convert the input file to WAV using a system tool
        # (macOS `afconvert` or `sox`) and then read the resulting WAV.
        if HAVE_AIFC and aifc is not None:
            with aifc.open(str(path), "rb") as af:
                nchan = af.getnchannels()
                sampwidth = af.getsampwidth()
                fr = af.getframerate()
                frames = af.readframes(af.getnframes())
        else:
            # Create a temporary WAV path and try to convert using available system tools.
            tmp_wav = Path(tempfile.mktemp(suffix=".wav"))
            try:
                afconvert = shutil.which("afconvert")
                if afconvert:
                    subprocess.run([afconvert, str(path), str(tmp_wav)], check=True)
                else:
                    sox = shutil.which("sox")
                    if sox:
                        subprocess.run([sox, str(path), str(tmp_wav)], check=True)
                    else:
                        raise RuntimeError(
                            "No aifc available and neither 'afconvert' nor 'sox' were found to convert AIFF"
                        )
                with wave.open(str(tmp_wav), "rb") as wf:
                    nchan = wf.getnchannels()
                    sampwidth = wf.getsampwidth()
                    fr = wf.getframerate()
                    frames = wf.readframes(wf.getnframes())
            finally:
                # Best-effort cleanup
                try:
                    if tmp_wav.exists():
                        tmp_wav.unlink()
                except Exception:
                    pass

    # Convert to 2-byte signed little-endian if necessary
    if sampwidth == 1:
        # 8-bit WAV likely unsigned; convert to signed 16-bit
        # audioop.lin2lin converts widths but treats 1-byte as signed; for unsigned 8-bit we offset.
        # First convert unsigned->signed by subtracting 128
        signed8 = bytes(((b - 128) & 0xFF) for b in frames)
        frames16 = audioop.lin2lin(signed8, 1, 2)
    elif sampwidth == 2:
        # if system endianness is not little-endian, wave gives little-endian by spec
        frames16 = frames
    else:
        # convert other widths if possible
        frames16 = audioop.lin2lin(frames, sampwidth, 2)

    # If stereo, convert to mono by averaging channels
    if nchan > 1:
        frames_mono = audioop.tomono(frames16, 2, 0.5, 0.5)
    else:
        frames_mono = frames16

    # Ensure output is little-endian signed 16-bit
    return frames_mono, fr


def convert_16bit_to_8bit_unsigned(frames16: bytes) -> bytes:
    """
    Convert signed 16-bit little-endian PCM to unsigned 8-bit PCM (0..255).
    Simple approach: convert 16-bit -> signed 8-bit (lin2lin), then add 128 to each byte.
    """
    signed8 = audioop.lin2lin(frames16, 2, 1)  # gives signed 8-bit
    unsigned8 = bytes(((b + 128) & 0xFF) for b in signed8)
    return unsigned8


def convert_8bit_unsigned_to_16bit_signed(unsigned8: bytes) -> bytes:
    """
    Convert unsigned 8-bit PCM (0..255) back to signed 16-bit little-endian PCM.
    """
    signed8 = bytes(((b - 128) & 0xFF) for b in unsigned8)
    frames16 = audioop.lin2lin(signed8, 1, 2)
    return frames16


def write_wav_8u(path: Path, data8u: bytes, sr: int) -> None:
    with wave.open(str(path), "wb") as wf:
        wf.setnchannels(1)
        wf.setsampwidth(1)
        wf.setframerate(sr)
        wf.writeframes(data8u)


def write_wav_16le(path: Path, frames16: bytes, sr: int) -> None:
    with wave.open(str(path), "wb") as wf:
        wf.setnchannels(1)
        wf.setsampwidth(2)
        wf.setframerate(sr)
        wf.writeframes(frames16)


def main():
    ap = argparse.ArgumentParser(
        description="Demo encode/decode audio using XOR+rotate obfuscation"
    )
    ap.add_argument(
        "--quote",
        "-q",
        default="To be, or not to be, that is the question.",
        help="Quote to synthesize",
    )
    ap.add_argument(
        "--sr",
        type=int,
        default=20000,
        help="sample rate for synthesized audio and output WAVs (default 20000)",
    )
    ap.add_argument("--outdir", default=".", help="output directory for demo files")
    args = ap.parse_args()

    outdir = Path(args.outdir)
    outdir.mkdir(parents=True, exist_ok=True)

    print(
        "Demo: synthesize quote -> convert -> encrypt -> sniff -> decrypt -> write WAVs"
    )
    print("Quote:", args.quote)

    with tempfile.TemporaryDirectory() as td:
        td_path = Path(td)
        synth_path = td_path / "synth.wav"

        print("Synthesizing TTS...")
        try:
            synthesize_tts_to_file(args.quote, synth_path, args.sr)
            print("Synthesized to:", synth_path)
        except Exception as e:
            print("TTS synthesis failed:", e, file=sys.stderr)
            sys.exit(1)

        # Read and normalize to mono 16-bit
        frames16, sr = read_wav_as_mono_16(synth_path)
        print(
            "Read synthesized audio: {} bytes, sample rate {}".format(len(frames16), sr)
        )

        # Convert to 8-bit unsigned PCM (simulate ADC)
        data8u = convert_16bit_to_8bit_unsigned(frames16)
        orig_8u = outdir / "original_8bit.wav"
        write_wav_8u(orig_8u, data8u, sr)
        print("Wrote 8-bit original WAV:", orig_8u)

        # Encrypt 8-bit audio bytes
        print("Encrypting 8-bit audio bytes...")
        cipher_bytes = encrypt_bytes(data8u)
        sniffer_wav = outdir / "sniffer_cipher_8bit.wav"
        write_wav_8u(sniffer_wav, cipher_bytes, sr)
        print("Wrote sniffer WAV (ciphertext as 8-bit PCM):", sniffer_wav)

        # Decrypt back
        print("Decrypting ciphertext...")
        decoded8u = decrypt_bytes(cipher_bytes)
        decoded_wav = outdir / "decoded_8bit.wav"
        write_wav_8u(decoded_wav, decoded8u, sr)
        print("Wrote decoded 8-bit WAV:", decoded_wav)

        # Also write decoded as 16-bit for nicer playback
        decoded16 = convert_8bit_unsigned_to_16bit_signed(decoded8u)
        decoded16_wav = outdir / "decoded_16bit.wav"
        write_wav_16le(decoded16_wav, decoded16, sr)
        print("Wrote decoded 16-bit WAV:", decoded16_wav)

        # Save original 16-bit version for comparison
        orig16_out = outdir / "original_16bit.wav"
        write_wav_16le(orig16_out, frames16, sr)
        print("Wrote original 16-bit WAV (synth source):", orig16_out)

    print("\nDemo complete. Files written to:", outdir.resolve())
    print(
        "Listen to - original_16bit.wav, original_8bit.wav, sniffer_cipher_8bit.wav, decoded_16bit.wav"
    )
    print(
        "Note: decoded_16bit.wav should closely match original_16bit.wav if encryption/decryption are correct."
    )


if __name__ == "__main__":
    main()
