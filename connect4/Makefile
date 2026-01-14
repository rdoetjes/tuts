CC = gcc
objects = connect4.o

all: connect4

connect4: $(objects)
	$(CC) $(objects) -o connect4

connect4.o: connect4.c
	$(CC) -c connect4.c -o connect4.o

.PHONY: clean
clean:
	rm -f connect4.o connect4
