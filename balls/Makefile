# Set the C compiler to clang to avoid issues with Homebrew GCC and macOS SDK headers
CC=clang
CGO_ENABLED=1
BINARY_NAME=balls

.PHONY: all build run clean tidy

all: build

build:
	CC=$(CC) CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY_NAME) *.go

run: build
	./$(BINARY_NAME)

tidy:
	go mod tidy

clean:
	go clean
	rm -f $(BINARY_NAME)
