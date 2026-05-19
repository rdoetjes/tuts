# Makefile for RP2040 Forth Implementation

AS       = arm-none-eabi-as
LD       = arm-none-eabi-ld
OBJCOPY  = arm-none-eabi-objcopy
QEMU     = qemu-system-arm
PYTHON   = python3

ASFLAGS  = -mcpu=cortex-m0plus -mthumb -g
LDFLAGS  = -T linker.ld
PICO_ASFLAGS = $(ASFLAGS) --defsym PICO=1

# --------------------------------------------------
# Common files
# --------------------------------------------------

ELF  = forth.elf
BIN  = forth.bin

PICO_ELF = forth_pico.elf
PICO_BIN = forth_pico.bin
PICO_UF2 = forth_pico.uf2

TEST_FILES = $(wildcard *.fth)

# --------------------------------------------------
# Normal build
# --------------------------------------------------

OBJS = boot.o forth.o

.PHONY: all clean test auto-test pico disasm

all: $(BIN)

boot.o: boot.s
	$(AS) $(ASFLAGS) -c $< -o $@

forth.o: forth.s
	$(AS) $(ASFLAGS) -c $< -o $@

$(ELF): $(OBJS) linker.ld
	$(LD) $(LDFLAGS) $(OBJS) -o $@

$(BIN): $(ELF)
	$(OBJCOPY) -O binary $< $@

$(PICO_UF2): $(PICO_BIN)
	cp $(PICO_BIN) flash_image.bin
	$(PYTHON) uf2conv.py \
		-f 0xe48bff56 \
		-b 0x10000000 \
		-c flash_image.bin \
		-o $@

pico: $(PICO_UF2)

# --------------------------------------------------
# Pico build
# --------------------------------------------------

PICO_OBJS = boot_pico.o forth_pico.o

boot_pico.o: boot.s rp2040.inc
	$(AS) $(PICO_ASFLAGS) -c $< -o $@

forth_pico.o: forth.s rp2040.inc
	$(AS) $(PICO_ASFLAGS) -c $< -o $@

$(PICO_ELF): $(PICO_OBJS) linker_pico.ld
	$(LD) -T linker_pico.ld $(PICO_OBJS) -o $@

$(PICO_BIN): $(PICO_ELF)
	$(OBJCOPY) -O binary $< $@

$(PICO_UF2): $(PICO_BIN)
	cp $(PICO_BIN) flash_image.bin
	$(PYTHON) uf2conv.py \
		-f 0xe48bff56 \
		-b 0x10000000 \
		-c flash_image.bin \
		-o $@


# -------------------------
# Clean
# -------------------------

clean:
	rm -f $(OBJS) $(ELF) $(BIN) \
	      boot_pico.o forth_pico.o forth_pico.elf forth_pico.bin forth_pico.uf2 \
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
