OBJECTS = knightrider.o system.o fileio.o gpio.o print.o

all: knightrider

%.o : %.S
	gcc -c $< -o $@

%.o: %.s
	as $< -o $@

knightrider: $(OBJECTS)
	ld -o knightrider $(OBJECTS)

clean:
	rm *.o
	rm knightrider
