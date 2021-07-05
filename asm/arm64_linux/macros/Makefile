OBJECTS = binprint.o print.o

ifdef DEBUG
	DEBUGFLGAS = -g
else
	DEBUGFLAGS= 
endif

all: binprint

%o: %.S
	gcc $(DEBUGFLAGS) -c $< -o $@

%.o: %.s
	as $(DEBUGFLAGS) $< -o $@

binprint: $(OBJECTS)
	ld -o binprint $(OBJECTS)

clean:
	rm *.o
	rm binprint
