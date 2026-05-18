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
UF2  = forth.uf2
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

$(UF2): $(BIN)
	python3 uf2conv.py $(BIN) -c -b 0x10000000 -o $(UF2)

clean:
	rm -f $(OBJS) $(ELF) $(BIN) $(PREPROCESSED) *.log *.clean *.pre.s

# Interactive test: Run Forth in QEMU and connect your keyboard to the serial port
test: $(ELF)
	@echo "Starting QEMU interactive session..."
	@echo "Press Ctrl-A then X to exit QEMU."
	$(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial mon:stdio

# Run all test files individually and assert against .expected files if present
auto-test: $(ELF) $(TEST_FILES)
	@echo "All tests completed."

$(TEST_FILES): $(ELF)
	@echo "------------------------------------------------------------"
	@echo "Running test: $@"
	@LOG_FILE=$@.log; \
	CLEAN_LOG=$$LOG_FILE.clean; \
	EXPECTED=$@.expected; \
	(echo " "; cat $@; sleep 10; echo " ") | $(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) -nographic -serial stdio -monitor none > $$LOG_FILE 2>&1 & \
	SUB_PID=$$!; \
	sleep 20; \
	kill -HUP $$SUB_PID 2>/dev/null || true; \
	sleep 1; \
	pkill -P $$SUB_PID 2>/dev/null || true; \
	sync; \
	grep -v "qemu:" $$LOG_FILE | grep -v "R00=" | grep -v "XPSR=" | grep -v "FPSCR:" | grep -v "terminating on signal" | grep -v "^[[:space:]]*ok[[:space:]]*$$" | sed -e :a -e '/^\s*$$/{$$d;N;ba' -e '}' > $$CLEAN_LOG; \
	if [ -f $$EXPECTED ]; then \
		echo "Asserting against $$EXPECTED..."; \
		grep -v "^[[:space:]]*ok[[:space:]]*$$" $$EXPECTED | sed -e :a -e '/^\s*$$/{$$d;N;ba' -e '}' > $$EXPECTED.tmp; \
		if diff -buBw $$EXPECTED.tmp $$CLEAN_LOG; then \
			echo "PASS: $@ matches $$EXPECTED"; \
			rm $$EXPECTED.tmp; \
		else \
			echo "FAIL: $@ differs from $$EXPECTED"; \
			rm $$EXPECTED.tmp; \
			exit 1; \
		fi \
	else \
		echo "No expected file found. Results for $@:"; \
		cat $$CLEAN_LOG; \
	fi
	@echo "------------------------------------------------------------"

# Helper to show word definitions in the binary
disasm: $(ELF)
	arm-none-eabi-objdump -d $(ELF) | less
