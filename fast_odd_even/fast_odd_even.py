#!/usr/bin/env python3

import time

ITS = 10000000

def slow():
	for i in range (0, ITS):
		if i % 2 == 0:
			pass

def fast():
	for i in range (0, ITS):
		if i & 2:
			pass


start = time.time_ns()
slow()
stop  = time.time_ns()
print("modulo took: %d " % ( (stop - start) / 1000 ) )

start = time.time_ns()
slow()
stop  = time.time_ns()
print("and took: %d " % ( (stop - start) /1000 ) )
