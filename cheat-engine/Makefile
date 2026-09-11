# Project variables
BINARY_NAME=cheat-engine
GO_FILES=$(wildcard *.go)
OS := $(shell uname)

# Compiler settings
ifeq ($(OS), Darwin)
	CC=clang
	SIGN=codesign -s - --entitlements entitlements.plist --force
else ifeq ($(OS), FreeBSD)
	CC=clang
else
	CC=gcc
endif

.PHONY: all build clean run help

all: build

build:
	@echo "Building for $(OS)..."
	CC=$(CC) CGO_ENABLED=1 go build -o $(BINARY_NAME) .
ifeq ($(OS), Darwin)
	@if [ -f entitlements.plist ]; then \
		echo "Signing binary with entitlements..."; \
		$(SIGN) $(BINARY_NAME); \
	else \
		echo "Warning: entitlements.plist not found. Binary might be killed on launch."; \
	fi
endif

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	go clean

# Memory scanning usually requires root privileges
run: build
	@if [ -z "$(PROC)" ]; then \
		echo "Usage: make run PROC=<process_name>"; \
	else \
		sudo ./$(BINARY_NAME) $(PROC); \
	fi

help:
	@echo "Available commands:"
	@echo "  make build      - Build the binary (and sign on macOS)"
	@echo "  make run PROC=x - Build and run as root against process 'x'"
	@echo "  make clean      - Remove binary and build cache"
