#!/usr/bin/env python3


def bitcount(n):
    c = 0
    while n != 0:
        c += 1
        n = n & n - 1
    return c


print(bitcount(254))
