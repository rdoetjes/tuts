# Makefile for tone_dialer (ALSA DTMF / CCITT5 / MF tone generator)
# Updated to build multiple .cpp translation units (serial.cpp, dtmf.cpp, c5.cpp, parsing.cpp, playing.cpp, etc.)
# Usage:
#   make           -> build optimized release binary
#   make debug     -> build debug binary (named tone_dialer-debug)
#   make clean
#   make run       -> build release and run with example_sequence.txt
#
TARGET     := tone_dialer
SRCS       := $(wildcard *.cpp)
OBJS       := $(SRCS:.cpp=.o)

CXX        := g++
CXXFLAGS   := -std=c++17 -Wall -Wextra -pedantic
LDFLAGS    := -lasound -lm

# Debug build flags
DEBUG_CXXFLAGS   := -g -O0 -DDEBUG
# Release/optimized build flags
RELEASE_CXXFLAGS := -O2 -DNDEBUG -march=native

# Default target
all: release

# Build debug target (produces tone_dialer-debug)
debug: CXXFLAGS += $(DEBUG_CXXFLAGS)
debug: CXXFLAGS += -fno-inline
debug: $(TARGET)-debug

# Build release target (produces tone_dialer)
release: CXXFLAGS += $(RELEASE_CXXFLAGS)
release: $(TARGET)

# Link rules
$(TARGET): $(OBJS)
	$(CXX) $(CXXFLAGS) $(OBJS) -o $@ $(LDFLAGS)

$(TARGET)-debug: $(OBJS)
	$(CXX) $(CXXFLAGS) $(OBJS) -o $@ $(LDFLAGS)

# Compile .cpp -> .o
%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Utility targets
clean:
	rm -f $(TARGET) $(TARGET)-debug $(OBJS) *.o core

rebuild: clean all

# Run helpers (example file name - replace with your sequence file as needed)
run: release
	./$(TARGET) example_sequence.txt

run-debug: debug
	./$(TARGET)-debug example_sequence.txt

help:
	@echo "Available targets:"
	@echo "  make              → build optimized release version"
	@echo "  make debug        → build debug version (tone_dialer-debug)"
	@echo "  make release      → build optimized release version (same as default)"
	@echo "  make clean        → remove built binaries and objects"
	@echo "  make rebuild      → clean + all"
	@echo "  make run          → build release and run with example file"
	@echo "  make run-debug    → build debug and run with example file"

.PHONY: all debug release clean rebuild run run-debug help
