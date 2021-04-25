#!/usr/bin/env python3

def one(v):
    print("one %d" % (v) )

def two(v):
    print("one %d" % (v) )

def three(v):
    print("one %d" % (v) )

if __name__ == "__main__":
    cmnds = {"one" : one, "two" : two, "three" : three}

    while True:
        cmd = input("enter command: one, two or three:\n")
        if not cmd in cmnds:
            print("not a valid command")
            continue

        cmnds[cmd](12)
