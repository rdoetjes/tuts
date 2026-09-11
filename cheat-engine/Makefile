# Project variables
BINARY_NAME=cheat-engine
# Use a simple way to detect OS that works in both GNU and BSD make
# GNU make uses $(shell ...), BSD make uses !=
OS_GNU := $(shell uname)
OS_BSD != uname
OS = $(OS_GNU)$(OS_BSD)

# Determine CC and SIGN based on OS
# We use shell scripts inside the targets to remain compatible
CC_DARWIN = clang
CC_FREEBSD = clang
CC_LINUX = gcc

SIGN_DARWIN = codesign -s - --entitlements entitlements.plist --force

.PHONY: all build clean run help

all: build

build:
	@echo "Building for $(OS)..."
	@if [ "$(OS)" = "Darwin" ]; then \
		CC=$(CC_DARWIN) CGO_ENABLED=1 go build -o $(BINARY_NAME) .; \
		if [ -f entitlements.plist ]; then \
			echo "Signing binary with entitlements..."; \
			$(SIGN_DARWIN) $(BINARY_NAME); \
		else \
			echo "Warning: entitlements.plist not found."; \
		fi; \
	elif [ "$(OS)" = "FreeBSD" ]; then \
		CC=$(CC_FREEBSD) CGO_ENABLED=1 go build -o $(BINARY_NAME) .; \
	else \
		CC=$(CC_LINUX) CGO_ENABLED=1 go build -o $(BINARY_NAME) .; \
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
