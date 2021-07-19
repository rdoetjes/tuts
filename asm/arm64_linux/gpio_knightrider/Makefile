OBJECTS = knightrider.o system.o fileio.o gpio.o print.o

all: knightrider

%.o : %.S
	gcc $(DEBUGFLGS) -c $< -o $@

%.o: %.s
	as -g $< -o $@

knightrider: $(OBJECTS)
	ld -o knightrider $(OBJECTS)

clean:
	rm *.o
	rm knightrider
