OBJECTS = decimal.o 

all: decimal

%.o: %.s
	as $< -o $@

decimal: $(OBJECTS)
	ld -o decimal $(OBJECTS)

clean:
	rm *.o
	rm decimal
