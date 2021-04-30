#!/usr/bin/env python3
import os

def writeDataToFile(fileName, data):
    try:
        f = open(fileName, "w+")
        f.write(data)
    except IOError:
        raise IOError
    finally:
        f.close()

#Main program
data = "Hi write me to the file please"
okayFile = "okay.txt"

try:
    writeDataToFile(okayFile, data)
except IOError as  e:
    print("Error %s occured, program halted" % {e.strerror})
    exit(1)

