# Makefile for RP2040 Forth Implementation
# Designed to build for RP2040 hardware and test core logic via QEMU

AS = arm-none-eabi-as
LD = arm-none-eabi-ld
OBJCOPY = arm-none-eabi-objcopy
QEMU = qemu-system-arm
PYTHON = python3

# RP2040 is Cortex-M0+, but for QEMU testing we often use mps2-an385 (Cortex-M3)
# as it has better generic support for basic Thumb-1 execution.
ASFLAGS = -mcpu=cortex-m0plus -mthumb -g
LDFLAGS = -T linker.ld

# Files
SRCS = forth.s
PREPROCESSED = forth.pre.s
OBJS = forth.o
ELF  = forth.elf
BIN  = forth.bin
TEST_FILES = $(wildcard *.fth)

.PHONY: all clean test auto-test $(TEST_FILES)

all: $(BIN)

$(PREPROCESSED): $(SRCS) rp2040.inc preprocess_forth.py
	@echo "Preprocessing Forth source..."
	$(PYTHON) preprocess_forth.py $(SRCS) > $@

$(OBJS): $(PREPROCESSED)
	$(AS) $(ASFLAGS) $(PREPROCESSED) -o $@

$(ELF): $(OBJS) linker.ld
	$(LD) $(LDFLAGS) $(OBJS) -o $@

$(BIN): $(ELF)
	$(OBJCOPY) -O binary $< $@

clean:
	rm -f $(OBJS) $(ELF) $(BIN) $(PREPROCESSED) *.log *.pre.s

# Interactive test: Run Forth in QEMU and connect your keyboard to the serial port
test: $(ELF)
	@echo "Starting QEMU interactive session..."
	@echo "Press Ctrl-A then X to exit QEMU."
	$(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial mon:stdio

# Run all test files individually
auto-test: $(ELF) $(TEST_FILES)

$(TEST_FILES): $(ELF)
	@echo "------------------------------------------------------------"
	@echo "Running test: $@"
	@LOG_FILE=$@.log; \
	(echo " "; cat $@; sleep 1; echo " ") | $(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial stdio -monitor none > $$LOG_FILE 2>&1 & \
	SLEEP_PID=$$!; \
	sleep 3; \
	pkill -P $$SLEEP_PID || true; \
	grep -v "qemu:" $$LOG_FILE | grep -v "R00=" | grep -v "XPSR=" | grep -v "FPSCR:" > $$LOG_FILE.clean; \
	echo "--- Results for $@ ---"; \
	cat $$LOG_FILE.clean; \
	echo "------------------------------------------------------------"

# Helper to show word definitions in the binary
disasm: $(ELF)
	arm-none-eabi-objdump -d $(ELF) | less
