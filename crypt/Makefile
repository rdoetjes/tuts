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

.PHONY: all build tests test clean show

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

clean:
	-rm -f $(BIN) $(TEST_BIN)

# convenience target to show compile commands used
show:
	@printf "CC = %s\nCFLAGS = %s\nSRC_MAIN = %s\nSRC_TEST = %s\nBIN = %s\nTEST_BIN = %s\n" \
		"$(CC)" "$(CFLAGS)" "$(SRC_MAIN)" "$(SRC_TEST)" "$(BIN)" "$(TEST_BIN)"
