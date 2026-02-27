CC = gcc
CFLAGS = -std=gnu99 -Wall -Wextra
RM = rm -f

SRCS = main.c serial.c
OBJS = $(SRCS:.c=.o)
EXE = brute-force

$(EXE): $(OBJS)
	$(CC) $(CFLAGS) $^ -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) $(OBJS)
	$(RM) $(EXE)
