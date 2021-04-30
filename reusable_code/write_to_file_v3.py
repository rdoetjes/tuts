#!/usr/bin/env python3
import os

def ThrowErrorIfObjectIsNull(obj):
    if not obj:
        raise TypeError("Object is set to Null")

def writeDataToFile(fileName, data):
    try:
        f = open(fileName, "w+")
        f = None
        ThrowErrorIfObjectIsNull(f)
        f.write(data)
        f.close()
    except Exception as e:
        raise e

#Main program
data = "Hi write me to the file please"
okayFile = "okay.txt"

try:
    writeDataToFile(okayFile, data)
except Exception as e:
    print("Uncaught exception %s" % (e))
    exit(1)

