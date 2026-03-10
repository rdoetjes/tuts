# bluebox-alsa

A small ALSA-based BlueBox for producing DTMF and CCITT (C5) MF tones. 

The program reads a simple tab-delimited "sequence" file and plays tones via the default ALSA playback device. It's useful for scripting tone sequences such as dialing procedures, signaling (KP1/KP2/ST), and other multi-frequency tones.

## Highlights / Features

- Generates DTMF and  CCITT-5 (C5-style MF) tones.
- Uses ALSA for low-latency audio output.
- Simple, tab-delimited sequence format that supports:
  - `D` (DTMF single-character tones)
  - `C` (C5 / CCITT5 codes and named tones)
  - `~` (wait / sleep / interactive wait-for-enter)
- Ability to play tone for a set duration and wait for set duration to play the next tone, to have ultimate timing flexibility.
- Small, single-file C++ implementation (`tone_dialer.cpp`) and `Makefile` for building.

## Requirements

- Linux with ALSA development libraries (libasound2-dev or equivalent).
- A C++17-capable compiler (g++ recommended).
- You may need permission to access the audio device; add your user to the `audio` group or run the binary with sufficient privileges if ALSA open fails.

The program uses:
- Sample rate: 8000 Hz
- Format: signed 16-bit little-endian mono (S16_LE)

## Build

From the project directory run:

```/dev/null/USAGE.md#L1-4
# Build (optimized release)
make

# Build debug
make debug

# Clean
make clean
```

The default build target here is `release`. The `Makefile` links against ALSA (`-lasound`) — see `Makefile` for details.

## Run

Usage:

```/dev/null/USAGE.md#L1-4
./tone_dialer sequence_file.txt
```

You can also use the Makefile helper targets:

```/dev/null/USAGE.md#L1-4
make run        # builds release and runs a (local) example filename if present
make run-debug  # builds debug and runs a (local) example filename if present
```

Note: The `make run` target in the included `Makefile` attempts to run `./tone_dialer nation_call_kp1`. Replace that example filename with any sequence file you want to play, e.g.:

```/dev/null/USAGE.md#L1-2
./tone_dialer national_call_kp1
```

## Sequence file format

Sequence files are tab-delimited text files. Each non-comment line describes one step:

Type<TAB>Tone<TAB>Duration_ms<TAB>Pause_ms

- Fields are separated by a tab character.
- Lines starting with `#` or `;` are comments and ignored.
- Example line:
  - `D<TAB>5<TAB>120<TAB>80` — play DTMF digit `5` for 120 ms, then pause 80 ms.
  - `C<TAB>KP1<TAB>100<TAB>50` — play C5 code `KP1` for 100 ms, pause 50 ms.
  - `~<TAB><TAB><TAB>` — a wait entry. If the duration field is is empty it waits for you to press Enter. 
  - `~<TAB>300<TAB><TAB>` — wait for 300 milliseconds before continuing to the next stop of the sequence.

Detailed meaning of fields:
- `Type`:
  - `D` — DTMF digit. `Tone` must be a single character: `0-9`, `*`, `#`, or `A-D`.
  - `C` — C5 / CCITT5 tone. `Tone` may be a single digit (`0-9`) mapped to preset freq pairs, or a named tone like `KP1`, `KP2`, `ST`, `SEIZE`, etc.
  - `~` — Wait/interactive. If a duration is provided and >0, the program sleeps that many ms. If duration is 0 or missing, the program waits for you to press Enter.
- `Tone`:
  - For `D`: single DTMF character.
  - For `C`: numeric 0–9 or strings like `KP1`, `KP2`, `ST`, `CODE11`, `CODE12`, `SEIZE`, `PROCEED` (or `PK`), `ANSWER`, `BUSY` (or `BUSYFLASH`), `CLEARBACK`, `CLEARFWD` (or `RELGUARD`).
- `Duration_ms`: how long the tone is played (milliseconds).
- `Pause_ms`: how long of silence to play after the tone (milliseconds). The player enforces a small minimum silence (5 ms) after tones to avoid concatenation glitches.

If the parser encounters an invalid tone it will skip playing that tone and continue.

## Example sequence files

Several example sequences are included in the repository. Here is one ("Mary had a little lamb" style dialing sequence):

```bluebox-alsa/mary_had_alittle_lamb#L1-200
#Dial the toll free number that ends in the USA
D	6	150	50
D	0	150	50
D	4	150	50
D	0	150	50
D	6	150	50
D	6	150	50
D	6	350	50
D	2	150	50
D	2	150	50
D	2	250	50
D	6	150	50
D	6	150	50
D	6	350	50
D	6	150	50
D	0	150	50
D	4	150	50
D	0	150	50
D	6	150	50
D	6	150	50
D	6	150	50
D	8	150	50
D	8	150	50
D	6	150	50
D	8	150	50
D	4	350	50
D	4	350	50
```

An example that mixes interactive waits, clears and seize signals for a national call:

```bluebox-alsa/national_call_kp1#L1-200
#Dial toll free number that ends in the USA
D	0	150	50
D	6	150	50
D	0	150	50
D	2	150	50
D	2	150	50
D	0	150	50
D	6	150	50
D	3	150	50
D	0	150	50
#Wait for answer then press enter
~
#Clear the trunk
C	CLEARFWD	250	50
#Seize the trunk wait for 2600Hz for 300 milliseconds
#This is very timing specific you could use only ~ to wait and listen for 2600Hz confirmation tone
C	SEIZE	250	300
#Signal domestic call
C	KP1	100	50
C	0	55	50
#Dial area code
C	7	55	50
C	1	55	50
C	7	55	50
#Dial the target number in that area
C	8	55	50
C	7	55	50
C	9	55	50
C	3	55	50
C	6	55	50
C	6	55	50
C	0	55	50
C	ST	55	50
```

There are also international examples: `international_transit_kp2`, `international_transit_multi_hop_kp2`, `international_rusty-n-edy_kp1` — inspect them to learn multi-step sequences and timing for transit/long-distance scenarios.

## Implementation notes

- The implementation is in `tone_dialer.cpp`. It generates plain sinusoids (no envelope) for the two frequencies per tone, sums them, and writes signed 16-bit frames to ALSA.
- The code sets up ALSA with:
  - Sample rate 8000 Hz
  - Mono
  - Period size ~128 frames
  - Buffer size = period * 4
- The generator clamps to the int16 range and recovers from ALSA underruns.

## Troubleshooting

- "Cannot open PCM device": ensure ALSA is installed and the default device exists. Check permissions (audio group) or try running with elevated permissions.
- Distorted sound / clipping: the code uses a conservative amplitude (about 40% of max) to mitigate clipping. If you modify amplitudes or mix with other audio, be careful.
- Timing-sensitive sequences (e.g., waiting for a 2600 Hz line): these may require human observation; you can use `~` to wait for a keypress or specify a small sleep duration.

## Contributing / Modifying

- The code is intentionally simple and contained in `tone_dialer.cpp`. If you want:
  - Add amplitude envelopes (fade-in/out) to reduce clicks,
  - Support additional tone sets,
  - Add file output (WAV) instead of ALSA playback,
  - Add command-line options for device selection, sample rate, or volume.

## Files of interest

- `tone_dialer.cpp` — main implementation.
- `Makefile` — build and run helpers.
- Example sequence files:
  - `mary_had_alittle_lamb`
  - `national_call_kp1`
  - `international_transit_kp2`
  - `international_transit_multi_hop_kp2`
  - `international_rusty-n-edy_kp1`

---
