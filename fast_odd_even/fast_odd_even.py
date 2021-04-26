#!/usr/bin/env python3

import time

ITS = 10000000

def moduloMethod():
    for i in range (0, ITS):
        if i % 2 == 0:
            pass

def andMethod():
    for i in range (0, ITS):
        if not i & 2:
            pass

def Average(lst):
    return sum(lst) / len(lst)

aC = []
bC = []

for i in range(0, 100):   
    start = time.time_ns()
    andMethod()
    stop  = time.time_ns()
    aC.append(stop-start)

    start = time.time_ns()
    moduloMethod()
    stop  = time.time_ns()
    bC.append(stop-start)

aAvg = Average(aC)
aMax = max(aC)
aMin = min(aC)

bAvg = Average(bC)
bMax = max(bC)
bMin = min(bC)

print("and method average %d max %d min %d" % (aAvg, aMax, aMin))
print("modulo method average %d max %d min %d" % (bAvg, bMax, bMin))

