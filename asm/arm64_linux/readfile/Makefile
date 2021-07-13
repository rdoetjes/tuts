OBJECTS = readfile.o 

all: readfile

%.o: %.s
	as -g $< -o $@

readfile: $(OBJECTS)
	ld -o readfile $(OBJECTS)

clean:
	rm *.o
	rm readfile
