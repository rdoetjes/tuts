# Makefile for RP2040 Forth Implementation
# Designed to build for RP2040 hardware and test core logic via QEMU

AS = arm-none-eabi-as
LD = arm-none-eabi-ld
OBJCOPY = arm-none-eabi-objcopy
QEMU = qemu-system-arm

# RP2040 is Cortex-M0+, but for QEMU testing we often use mps2-an385 (Cortex-M3)
# as it has better generic support for basic Thumb-1 execution.
ASFLAGS = -mcpu=cortex-m0plus -mthumb -g
LDFLAGS = -T linker.ld

# Files
SRCS = forth.s
OBJS = forth.o
ELF  = forth.elf
BIN  = forth.bin
TEST_SCRIPT = simple_test.fth

.PHONY: all clean test auto-test

all: $(BIN)

$(ELF): $(OBJS) linker.ld
	$(LD) $(LDFLAGS) $(OBJS) -o $@

$(OBJS): forth.s rp2040.inc
	$(AS) $(ASFLAGS) forth.s -o $@

$(BIN): $(ELF)
	$(OBJCOPY) -O binary $< $@

clean:
	rm -f $(OBJS) $(ELF) $(BIN) test_output.log

# Interactive test: Run Forth in QEMU and connect your keyboard to the serial port
test: $(ELF)
	@echo "Starting QEMU interactive session..."
	@echo "Press Ctrl-A then X to exit QEMU."
	$(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial mon:stdio

# Auto-test: Feed the test script into Forth and capture output
# Note: Since the hardware UART addresses differ between RP2040 and the QEMU mps2 machine,
# this assumes the code is either patched for the test or running in a compatible environment.
auto-test: $(ELF)
	@echo "Running automatic tests from $(TEST_SCRIPT)..."
	@if [ -f $(TEST_SCRIPT) ]; then \
		(echo " "; cat $(TEST_SCRIPT); sleep 1; echo " ") | $(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial stdio -monitor none > test_output.log 2>&1 & \
		sleep 4; \
		pkill -P $$! || true; \
		echo "--- Test Output Start ---"; \
		grep -v "qemu:" test_output.log | grep -v "R00=" | grep -v "XPSR=" | grep -v "FPSCR:" || cat test_output.log; \
		echo "--- Test Output End ---"; \
	else \
		echo "Error: $(TEST_SCRIPT) not found."; \
		exit 1; \
	fi

# Helper to show word definitions in the binary
disasm: $(ELF)
	arm-none-eabi-objdump -d $(ELF) | less
