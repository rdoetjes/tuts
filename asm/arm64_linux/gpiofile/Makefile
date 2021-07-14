OBJECTS = gpiofile.o 

all: gpiofile

%.o: %.s
	as -g $< -o $@

gpiofile: $(OBJECTS)
	ld -o gpiofile $(OBJECTS)

clean:
	rm *.o
	rm gpiofile
