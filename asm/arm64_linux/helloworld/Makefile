hello64: hello64.o
	ld -s -o hello64 hello64.o

hello64.o: hello64.s
	as -o hello64.o hello64.s -gdwarf64 -g3

clean:
	rm *.o
	rm hello64
