bitcount: bitcount.o
	ld -s -o bitcount bitcount.o

bitcount.o: bitcount.s
	as -o bitcount.o bitcount.s

clean:
	rm *.o
	rm bitcount
