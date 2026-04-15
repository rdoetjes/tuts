# bluebox-alsa

This project is an ALSA-based tone generator originally written as a single file. It has been refactored into small, focused translation units to improve maintainability and make it easier to extend.

What changed in the refactor
- Code split into logical modules:
  - `main.cpp` — program entry and orchestration.
  - `parsing.cpp` / `parsing.h` — sequence file parsing and the `SequenceStep` type.
  - `dtmf.cpp` / `dtmf.h` — DTMF frequency lookup helpers.
  - `c5.cpp` / `c5.h` — CCITT/C5 code lookup helpers.
  - `playing.cpp` / `playing.h` — ALSA setup, audio generation, and sequence execution.
  - `serial.cpp` / `serial.h` — simple serial command/response helpers used by the `H` command.
- `Makefile` updated to build all translation units and produce the `tone_dialer` executable.

High-level behavior
- The program reads a tab-delimited sequence file and executes steps in order:
  - `D` — play a DTMF digit (single character tone).
  - `C` — play a CCITT / C5 code (digit or named code such as `KP1`, `KP2`, `ST`, etc.).
  - `H` — send the single ASCII byte `H` over a serial device and wait up to 1000 ms for a response:
    - Responds with ASCII `'1'` → interpreted as OK (returns 1).
    - Responds with ASCII `'0'` → interpreted as Failed (returns 0).
    - Timeout, error, or unexpected byte → treated as timeout/error (returns -1).
  - `~` — wait step: if a duration is provided (>0) it sleeps that many milliseconds; if zero or empty it waits for the user to press Enter.
- Serial device:
  - The CLI accepts an optional serial device path. If not provided it defaults to `/dev/USB00` (see `serial.h`).
  - If the serial device cannot be opened the program continues but skips `H` steps with a warning.

Requirements
- Linux with ALSA development libraries installed (e.g. `libasound2-dev` / `alsa-lib`).
- A C++17-capable compiler (g++ recommended).
- Permission to access audio device (join `audio` group or run with appropriate privileges).

Build
- From the `bluebox-alsa` project directory run:

```bluebox-alsa/Makefile#L1-200
make
```

- Build debug variant:

```bluebox-alsa/Makefile#L1-200
make debug
```

- Clean:

```bluebox-alsa/Makefile#L1-200
make clean
```

Run
- Basic usage:

```/dev/null/USAGE.md#L1-2
./tone_dialer sequence.txt [serial_device]
```

- Example (use one of the provided sequence files):

```/dev/null/USAGE.md#L1-2
./tone_dialer national_call_kp1 /dev/USB00
```

Sequence file format (tab-delimited)
- Each non-comment line must be one of:
  - Type<TAB>Tone<TAB>Duration_ms<TAB>Pause_ms
- Supported `Type` values:
  - `D` — DTMF digit. `Tone` must be a single character `0-9`, `*`, `#`, or `A-D`.
  - `C` — CCITT/C5 code. `Tone` may be `0-9` or named codes like `KP1`, `KP2`, `ST`, `SEIZE`, `PROCEED`, `ANSWER`, `BUSY`, `CLEARFWD`, etc.
  - `H` — Serial host command (sends ASCII `H`); `Tone` is ignored. `Duration_ms` and `Pause_ms` may be used to tune behavior or provide pauses after the command.
  - `~` — Wait: if `Duration_ms` is `0` or absent the player waits for Enter; if >0 it sleeps that many milliseconds.
- Comments: lines beginning with `#` or `;` are ignored.
- Example lines:
  - `D<TAB>5<TAB>120<TAB>80` — play DTMF `5` for 120 ms and pause 80 ms.
  - `C<TAB>KP1<TAB>100<TAB>50` — play C5 `KP1` for 100 ms and pause 50 ms.
  - `H<TAB><TAB>0<TAB>100` — send `H`, wait for response (1000 ms internally), then pause 100 ms.
  - `~<TAB><TAB><TAB>` — wait for Enter.

Included example sequence files
- `mary_had_alittle_lamb` — example DTMF sequence.
- `national_call_kp1` — example mixing D and C steps, with an interactive wait and `SEIZE`/`KP1` usage.
- `international_transit_kp2`, `international_transit_multi_hop_kp2`, `international_rusty-n-edy_kp1` — international dialing examples.

Files of interest
- `main.cpp` — CLI and orchestration.
- `parsing.h` / `parsing.cpp` — sequence parsing and `SequenceStep`.
- `dtmf.h` / `dtmf.cpp` — DTMF frequency table.
- `c5.h` / `c5.cpp` — C5/CCITT codes mapping.
- `playing.h` / `playing.cpp` — ALSA setup, tone generation, and sequence execution.
- `serial.h` / `serial.cpp` — serial open/send-and-wait helpers.
- `Makefile` — build rules for multiple translation units.

Notes and suggestions
- The `serial` module currently configures the device to 9600 8N1 with no hardware flow control. If your device uses different settings, adjust `serial.cpp`.
- The tone generator produces plain sinusoids without an amplitude envelope; you may hear small clicks at tone boundaries. Consider adding short fade-in/out envelopes for smoother transitions.
- If you want the program to emit WAV files instead of playing to ALSA, extend `playing.cpp` with a WAV writer and call it instead of ALSA output.
- For diagnostic logging of serial traffic or to support more complex multi-byte responses, extend `serial.cpp::send_and_wait_serial` accordingly.

If you'd like, I can:
- Add a small `example_sequence.txt` and update `Makefile` `run` target to use it.
- Add command-line flags for serial baud rate, ALSA device selection, or volume control.
- Add a table in this README enumerating all supported `C` named codes and their frequencies.

Which of those would you like next?