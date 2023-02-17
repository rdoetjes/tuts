#!/usr/bin/env python3

def print_grid(buttons):
    n = len(buttons) + 1
    for i in range(1, n):
        print(buttons[i-1], end='')
        if (i % 10 == 0):
            print()
    print()

def switch_logic(buttons):
    n = len(buttons) + 1
    for person in range(1,n):
        print("person: "+str(person))
        for button in range(1,n):
            if (button % person == 0):
                buttons[button-1] = buttons[button-1] ^ 1
        print_grid(buttons)

if __name__ == "__main__":
    n = 100
    buttons=[0] * n
    
    switch_logic(buttons)
