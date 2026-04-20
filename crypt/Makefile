# Simple Makefile for the crypt project
# Default: build the main binary and run the unit tests.
#
# Targets:
#   make / make all   -> build the main binary and run tests
#   make build        -> build the main binary only
#   make tests        -> build the tests binary and run it
#   make test         -> alias for make tests
#   make clean        -> remove built artifacts
#   make show         -> print build variables
#   make recover      -> build recover tools, generate demo data and run the recovery demo
#
# The tests are compiled with -DTESTING so main.c's main() is excluded
# in the test build and the standalone tests.c main() is used.

# Compiler selection:
# - If you set the CC environment variable it will be used.
# - Otherwise, prefer clang, then GCC variants if available.
CC ?= $(shell for c in clang gcc-13 gcc-12 gcc-11 gcc; do command -v $$c >/dev/null 2>&1 && echo $$c && break; done)

CFLAGS ?= -std=c11 -Wall -Wextra -O2 -g
LDLIBS ?=

SRC_MAIN := main.c
SRC_TEST := tests.c

BIN := crypt_main
# Use a different test binary name to avoid colliding with the `tests` target name.
TEST_BIN := tests_bin

.PHONY: all build tests test clean show recover recover_tools interactive

# Default: build program and run tests
all: build tests

# Build main binary only
build: $(BIN)

$(BIN): $(SRC_MAIN)
	$(CC) $(CFLAGS) -o $@ $(SRC_MAIN) $(LDLIBS)

# Build the test binary (define TESTING so main.c excludes its runtime main)
# Note: TEST_BIN is intentionally not 'tests' to avoid a target-name collision.
$(TEST_BIN): $(SRC_MAIN) $(SRC_TEST)
	$(CC) $(CFLAGS) -DTESTING -o $@ $(SRC_MAIN) $(SRC_TEST) $(LDLIBS)

# Build and run tests
tests: $(TEST_BIN)
	@echo "Running unit tests..."
	./$(TEST_BIN)

# alias to run tests
test: tests

# Build recover tools (recover_key and make_cipher)
recover_tools: recover_key make_cipher

recover_key: recover_key.c
	$(CC) $(CFLAGS) -o recover_key recover_key.c $(LDLIBS)

make_cipher: make_cipher.c
	$(CC) $(CFLAGS) -o make_cipher make_cipher.c $(LDLIBS)

# Run a demo: build tools, generate demo cipher/crib and run recover_key on them.
# Run a demo: build tools, generate demo data and run the recovery demo.
# This target expects the repository binaries to be built in the current directory.
recover: recover_tools
	@echo "Generating demo ciphertext and crib..."
	# produce cipher.bin and crib.bin (offset 10, crib length default)
	./make_cipher cipher.bin crib.bin 10 32
	@echo "Running recover_key on demo data (cipher.bin:crib.bin:10)..."
	# run recovery (no explicit -k => tool will try default key lengths)
	./recover_key cipher.bin:crib.bin:10

# Run the interactive resolver (requires python3). Uses candidate files produced
# by recover_key. You can override defaults by invoking `make` targets or by
# calling the script directly.
interactive: recover_tools
	@echo "Launching interactive key resolver (interactive_resolve.py)..."
	python3 interactive_resolve.py --key recovered_key_final --candidates recovered_key_final_text --cipher cipher.bin

clean:
	-rm -f $(BIN) $(TEST_BIN) recover_key make_cipher cipher.bin crib.bin recovered_key* *.txt

# convenience target to show compile commands used
show:
	@printf "CC = %s\nCFLAGS = %s\nSRC_MAIN = %s\nSRC_TEST = %s\nBIN = %s\nTEST_BIN = %s\n" \
		"$(CC)" "$(CFLAGS)" "$(SRC_MAIN)" "$(SRC_TEST)" "$(BIN)" "$(TEST_BIN)"
