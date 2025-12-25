CC = gcc
CFLAGS = -Wall -Wextra -pedantic -std=c99
TARGET = knightstour
SOURCE = knightstour.c

all: $(TARGET)

$(TARGET): $(SOURCE)
	$(CC) $(CFLAGS) -o $(TARGET) $(SOURCE)

clean:
	rm -f $(TARGET)
