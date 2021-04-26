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

for i in range(0, 10):   
    start = time.time_ns()
    andMethod()
    stop  = time.time_ns()
    aC.append(stop-start)

    start = time.time_ns()
    moduloMethod()
    stop  = time.time_ns()
    bC.append(stop-start)

aAvg = Average(aC) / 1000
aMax = max(aC) / 1000 
aMin = min(aC) / 1000

bAvg = Average(bC) / 1000
bMax = max(bC) / 1000
bMin = min(bC) / 1000

print("and    method average %d uSec max %d uSec min %d uSec"  % (aAvg, aMax, aMin))
print("modulo method average %d uSec max %d uSec min %d uSec " % (bAvg, bMax, bMin))

