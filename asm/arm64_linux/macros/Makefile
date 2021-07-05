OBJECTS = binprint.o 

all: binprint

%.o: %.s
	as $< -o $@

binprint: $(OBJECTS)
	ld -o binprint $(OBJECTS)

clean:
	rm *.o
	rm binprint
