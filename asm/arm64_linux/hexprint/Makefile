hexprint: hexprint.o
	ld -s -o hexprint hexprint.o

hexprint.o: hexprint.s
	as -o hexprint.o hexprint.s

clean:
	rm *.o
	rm hexprint
