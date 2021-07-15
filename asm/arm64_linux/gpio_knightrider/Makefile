OBJECTS = knightrider.o 

all: knightrider

%.o: %.s
	as -g $< -o $@

knightrider: $(OBJECTS)
	ld -o knightrider $(OBJECTS)

clean:
	rm *.o
	rm knightrider
