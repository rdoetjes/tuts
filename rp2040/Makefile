# Makefile for RP2040 Forth Implementation
# Designed to build for RP2040 hardware and test core logic via QEMU

AS = arm-none-eabi-as
LD = arm-none-eabi-ld
OBJCOPY = arm-none-eabi-objcopy
QEMU = qemu-system-arm
PYTHON = python3

# RP2040 is Cortex-M0+, but for QEMU testing we often use mps2-an385 (Cortex-M3)
ASFLAGS = -mcpu=cortex-m0plus -mthumb -g
LDFLAGS = -T linker.ld

# Files
SRCS = forth.s
OBJS = forth.o
ELF  = forth.elf
BIN  = forth.bin
UF2  = forth.uf2
TEST_FILES = $(wildcard *.fth)

.PHONY: all clean test auto-test pico

all: $(BIN)

pico: forth_pico.uf2

# -------------------------
# Core build
# -------------------------

$(OBJS): $(SRCS)
	$(AS) $(ASFLAGS) -c $< -o $@

$(ELF): $(OBJS) linker.ld
	$(LD) $(LDFLAGS) $(OBJS) -o $@

$(BIN): $(ELF)
	$(OBJCOPY) -O binary $< $@

# -------------------------
# Pico build
# -------------------------

forth_pico.o: $(SRCS) rp2040.inc
	$(AS) $(ASFLAGS) --defsym PICO=1 -c $(SRCS) -o $@

$(BOOT2_OBJS): $(BOOT2_SRCS)
	$(AS) $(BOOT2_ASFLAGS) -c $< -o $@

forth_pico.elf: $(BOOT2_OBJS) forth_pico.o linker_pico.ld
	$(LD) $(LDFLAGS) $(BOOT2_OBJS) forth_pico.o -o $@

forth_pico.bin: forth_pico.elf
	$(OBJCOPY) -O binary $< $@

forth_pico.uf2: forth_pico.bin pad_checksum.py
	#@echo "Calculating and injecting Boot Stage 2 CRC..."
	#$(PYTHON) pad_checksum.py $<
	cat boot2.bin forth_pico.bin > flash_image.bin
	$(PYTHON)  uf2conv.py -f 0xe48bff56 -b 0x10000000 -c flash_image.bin -o forth_pico.uf2

# -------------------------
# Clean
# -------------------------

clean:
	rm -f $(OBJS) $(ELF) $(BIN) \
	      forth_pico.o forth_pico.elf forth_pico.bin forth_pico.uf2 \
	      $(BOOT2_OBJS) *.log *.clean *.pre.s *.uf2

# -------------------------
# QEMU test
# -------------------------

test: $(ELF)
	@echo "Starting QEMU interactive session..."
	@echo "Press Ctrl-A then X to exit QEMU."
	$(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) \
		-nographic -serial mon:stdio

auto-test: $(ELF) $(TEST_FILES)
	@echo "All tests completed."

$(TEST_FILES): $(ELF)
	@echo "------------------------------------------------------------"
	@echo "Running test: $@"
	@LOG_FILE=$@.log; \
	CLEAN_LOG=$$LOG_FILE.clean; \
	EXPECTED=$@.expected; \
	(echo " "; cat $@; sleep 10; echo " ") | \
	$(QEMU) -machine mps2-an385 -cpu cortex-m3 -kernel $(ELF) \
	-nographic -serial stdio -monitor none > $$LOG_FILE 2>&1 & \
	SUB_PID=$$!; \
	sleep 20; \
	kill -HUP $$SUB_PID 2>/dev/null || true; \
	sleep 1; \
	pkill -P $$SUB_PID 2>/dev/null || true; \
	sync; \
	grep -v "qemu:" $$LOG_FILE | \
	grep -v "R00=" | \
	grep -v "XPSR=" | \
	grep -v "FPSCR:" | \
	grep -v "terminating on signal" | \
	grep -v "^[[:space:]]*ok[[:space:]]*$$" | \
	sed -e :a -e '/^\s*$$/{$$d;N;ba' -e '}' > $$CLEAN_LOG; \
	if [ -f $$EXPECTED ]; then \
		echo "Asserting against $$EXPECTED..."; \
		grep -v "^[[:space:]]*ok[[:space:]]*$$" $$EXPECTED | \
		sed -e :a -e '/^\s*$$/{$$d;N;ba' -e '}' > $$EXPECTED.tmp; \
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

# -------------------------
# Disassembly helper
# -------------------------

disasm: $(ELF)
	arm-none-eabi-objdump -d $(ELF) | less
