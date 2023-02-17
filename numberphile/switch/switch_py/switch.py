#!/usr/bin/env python3

def print_grid(switches):
    n = len(switches) + 1
    for i in range(1, n):
        print(switches[i-1], end='')
        if (i % 10 == 0):
            print()
    print()

def switch_logic(switches):
    n = len(switches) + 1
    for person in range(1,n):
        print("person: "+str(person))
        for button in range(1,n):
            if (button % person == 0):
                switches[button-1] = switches[button-1] ^ 1
        print_grid(switches)

if __name__ == "__main__":
    n = 100
    switches=[0] * n
    switch_logic(switches)
