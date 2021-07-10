OBJECTS = binprint.o 

all: binprint

%.o: %.s
	as -g $< -o $@

binprint: $(OBJECTS)
	ld -o binprint $(OBJECTS)

clean:
	rm *.o
	rm binprint
