# RP2040 Forth (Subroutine Threaded Code)

A minimal, high-performance Forth implementation for the Raspberry Pi RP2040 (Cortex-M0+) written in ARM Thumb assembly.

## Features

- **Subroutine Threaded Code (STC)**: Every Forth word is a native Thumb subroutine.
- **Assembly Shorthand**: Use Forth words directly in assembly (e.g., `DUP`, `PLUS`, `EXIT`) via macros.
- **Interactive REPL**: Accessible over UART0 (Pins 0/1) at 115,200 baud.
- **Runtime Compiler**: Define new words in RAM using `:` and `;`.
- **Dual Compatibility**: Includes drivers for both RP2040 hardware and QEMU (MPS2 AN385) for testing.

## Prerequisites

You will need the ARM GNU Toolchain and QEMU installed:

```bash
# MacOS (Homebrew)
brew install --cask gcc-arm-embedded
brew install qemu

# Ubuntu/Debian
sudo apt-get install gcc-arm-none-eabi qemu-system-arm
```

## Building

To assemble and link the Forth image:

```bash
cd rp2040
make
```

This will produce `forth.bin` (for flashing to hardware) and `forth.elf` (for debugging/QEMU).

## Testing with QEMU

You can run the Forth REPL locally using QEMU. 

### Interactive Mode
Run the following command to start an interactive session:

```bash
make test
```
Type `WORDS` to see the dictionary or `: DOUBLE DUP + ;` to define a new word. Press `Ctrl-A` then `X` to exit.

### Automated Test Script
To run the included `test_words.fth` script and verify functionality:

```bash
make auto-test
```

## Hardware Setup

1.  **Wiring**: Connect a USB-to-UART adapter to the RP2040:
    - RP2040 GPIO 0 (TX) -> Adapter RX
    - RP2040 GPIO 1 (RX) -> Adapter TX
    - GND -> GND
2.  **Flashing**: Use `elf2uf2` or `picotool` to load `forth.elf` onto your Pico/RP2040, or convert `forth.bin` to a UF2 file.
3.  **Terminal**: Connect at 115,200 bps, 8N1.

## Dictionary Shorthand

When writing new words in `forth.s`, you can use the following macros to avoid typing `bl` manually:

```asm
defword "2DUP", f_2dup
    OVER
    OVER
    EXIT
```

Available macros include: `DUP`, `DROP`, `SWAP`, `OVER`, `PLUS`, `MINUS`, `STAR`, `ANDW`, `ORW`, `FETCH`, `STORE`, `KEY`, `EMIT`, and `EXIT`.