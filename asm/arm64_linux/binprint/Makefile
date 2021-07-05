binprint: binprint.o
	ld -s -o binprint binprint.o

binprint.o: binprint.s
	as -o binprint.o binprint.s

clean:
	rm *.o
	rm binprint
