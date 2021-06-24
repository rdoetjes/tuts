#!/usr/bin/env python3
import os

def writeDataToFile(fileName, data):
    try:
        f = open(fileName, "a+")
        f.write(data)
        f.close()
    except IOError:
        raise IOError

#Main program
data = "Hi write me to the file please"
okayFile = "okay.txt"

try:
    writeDataToFile(okayFile, data)
except IOError as  e:
    print("Error %s occured, program halted" % {e.strerror})
    exit(1)

