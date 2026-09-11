# Project variables
BINARY_NAME=cheat-engine

# Determine OS using a method compatible with both GNU and BSD make
OS != uname

# Compiler settings
CC_Darwin = clang
CC_FreeBSD = clang
CC_Linux = gcc

SIGN_Darwin = codesign -s - --entitlements entitlements.plist --force

# We use shell expansion inside the commands to get the right variables
# because BSD make and GNU make handle nested expansion differently.

.PHONY: all build clean run help

all: build

build:
	@echo "Building for $(OS)..."
	@case "$(OS)" in \
		"Darwin")  CC=clang; SIGN="codesign -s - --entitlements entitlements.plist --force" ;; \
		"FreeBSD") CC=clang; SIGN="true" ;; \
		*)         CC=gcc; SIGN="true" ;; \
	esac; \
	CC=$$CC CGO_ENABLED=1 go build -o $(BINARY_NAME) .; \
	if [ "$(OS)" = "Darwin" ]; then \
		if [ -f entitlements.plist ]; then \
			echo "Signing binary with entitlements..."; \
			$$SIGN $(BINARY_NAME); \
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
