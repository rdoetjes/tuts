#!/usr/bin/env python3

import time


def moduloMethod(ITS):
    temp = 0
    for i in range(0, ITS):
        if i % 2 == 0:
            temp += 1
    print(temp)


def andMethod(ITS):
    temp = 0
    for i in range(0, ITS):
        if not i & 2:
            temp += 1
    print(temp)


def Average(lst):
    return sum(lst) / len(lst)


aC = []
bC = []

ITS = 10000000

for i in range(0, 100):
    start = time.time_ns()
    andMethod(ITS)
    stop = time.time_ns()
    aC.append(stop - start)

    start = time.time_ns()
    moduloMethod(ITS)
    stop = time.time_ns()
    bC.append(stop - start)

aAvg = Average(aC) / 1000

bAvg = Average(bC) / 1000

print("and    method average %d uSec" % (aAvg))
print("modulo method average %d uSec" % (bAvg))
