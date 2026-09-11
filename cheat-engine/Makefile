# Project variables
BINARY_NAME=cheat-engine

# Determine OS using a method compatible with both GNU and BSD make
# BSD make ignores lines starting with 'ifdef', GNU make follows them.
# However, for simplicity and to avoid the "shell uname" warning in BSD make,
# we'll use a shell-based detection inside a single variable assignment.
OS != uname

# Compiler settings
CC_Darwin = clang
CC_FreeBSD = clang
CC_Linux = gcc

SIGN_Darwin = codesign -s - --entitlements entitlements.plist --force

# Selected compiler based on detected OS
CC = $(CC_$(OS))
SIGN = $(SIGN_$(OS))

.PHONY: all build clean run help

all: build

build:
	@echo "Building for $(OS)..."
	CC=$(CC) CGO_ENABLED=1 go build -o $(BINARY_NAME) .
	@if [ "$(OS)" = "Darwin" ]; then \
		if [ -f entitlements.plist ]; then \
			echo "Signing binary with entitlements..."; \
			$(SIGN) $(BINARY_NAME); \
		else \
			echo "Warning: entitlements.plist not found."; \
		fi; \
	fi

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	go clean

run: build
	@if [ -z "$(PROC)" ]; then \
		echo "Usage: make run PROC=<process_name>"; \
	else \
		sudo ./$(BINARY_NAME) $(PROC); \
	fi

help:
	@echo "Available commands:"
	@echo "  make build      - Build the binary"
	@echo "  make run PROC=x - Build and run as root against process 'x'"
	@echo "  make clean      - Remove binary and build cache"
