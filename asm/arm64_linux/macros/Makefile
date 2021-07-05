OBJECTS = binprint.o print.o

ifdef DEBUG
	DEBUGFLAGS = -g
else
	DEBUGFLAGS= 
endif

all: binprint

%.o: %.s
	as $(DEBUGFLAGS) $< -o $@

binprint: $(OBJECTS)
	ld -o binprint $(OBJECTS)

clean:
	rm *.o
	rm binprint
