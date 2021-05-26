helloworld: helloworld.o
	ld -s -o helloworld helloworld.o

helloworld.o: helloworld.s
	as -o helloworld.o helloworld.s

clean:
	rm *.o
	rm helloworld
