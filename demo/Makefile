# Binary name
BINARY_NAME=demo

# Local cache directory to avoid "operation not permitted" errors in the sandbox/macOS
GOCACHE_PATH=$(shell pwd)/.cache

# Environment variables for Go/CGO
GO_ENV=GOCACHE=$(GOCACHE_PATH) CGO_ENABLED=1

# Force use of system Clang (cc) to avoid GCC/macOS SDK header conflicts
export CC=cc
export CXX=c++

# Homebrew paths for Raylib
HOMEBREW_PREFIX=$(shell brew --prefix)
export CGO_CFLAGS=-I$(HOMEBREW_PREFIX)/include -D_DARWIN_C_SOURCE -Wno-implicit-function-declaration
export CGO_LDFLAGS=-L$(HOMEBREW_PREFIX)/lib -lraylib -framework CoreVideo -framework IOKit -framework Cocoa -framework OpenGL

# Local cache directory
GOCACHE_PATH=$(shell pwd)/.cache
GO_ENV=GOCACHE=$(GOCACHE_PATH) CGO_ENABLED=1

.PHONY: all build run clean tidy

all: build

build:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go build -v -tags external -o $(BINARY_NAME) .

run:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go run -tags external .

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(GOCACHE_PATH)

tidy:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go mod tidy

all: build

# Create cache dir and build
build:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go build -v -o $(BINARY_NAME) .

# Run the simulation
run:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go run .

# Clean up
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(GOCACHE_PATH)

# Fix dependencies
tidy:
	mkdir -p $(GOCACHE_PATH)
	$(GO_ENV) go mod tidy
