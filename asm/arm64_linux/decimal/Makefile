OBJECTS = decimal.o common.o input.o conversion.o

all: decimal

%.o: %.s
	as -g $< -o $@

decimal: $(OBJECTS)
	ld -o decimal $(OBJECTS)

clean:
	rm *.o
	rm decimal
