# Makefile for tone_dialer (ALSA DTMF / CCITT5 / MF tone generator)

# ────────────────────────────────────────────────────────────────────────────────
# Configuration
# ────────────────────────────────────────────────────────────────────────────────

TARGET     := tone_dialer
SRC        := tone_dialer.cpp

CXX        := g++
CXXFLAGS   := -std=c++17 -Wall -Wextra -pedantic
LDFLAGS    := -lasound -lm

# Debug build flags
DEBUG_CXXFLAGS   := -g -O0 -DDEBUG
# Release/optimized build flags
RELEASE_CXXFLAGS := -O2 -DNDEBUG -march=native

# Default target: optimized release build
all: release

# ────────────────────────────────────────────────────────────────────────────────
# Build rules
# ────────────────────────────────────────────────────────────────────────────────

debug: CXXFLAGS += $(DEBUG_CXXFLAGS)
debug: $(TARGET)-debug

release: CXXFLAGS += $(RELEASE_CXXFLAGS)
release: $(TARGET)

$(TARGET): $(SRC)
	$(CXX) $(CXXFLAGS) $(SRC) -o $@ $(LDFLAGS)

$(TARGET)-debug: $(SRC)
	$(CXX) $(CXXFLAGS) $(SRC) -o $@ $(LDFLAGS)

# ────────────────────────────────────────────────────────────────────────────────
# Utility targets
# ────────────────────────────────────────────────────────────────────────────────

clean:
	rm -f $(TARGET) $(TARGET)-debug *.o core

rebuild: clean all

run: release
	./$(TARGET) example_sequence.txt

run-debug: debug
	./$(TARGET)-debug example_sequence.txt

help:
	@echo "Available targets:"
	@echo "  make              → build optimized release version"
	@echo "  make debug        → build debug version"
	@echo "  make release      → build optimized release version (same as default)"
	@echo "  make clean        → remove built binaries"
	@echo "  make rebuild      → clean + all"
	@echo "  make run          → build release and run with example file"
	@echo "  make run-debug    → build debug and run with example file"

.PHONY: all debug release clean rebuild run run-debug help
