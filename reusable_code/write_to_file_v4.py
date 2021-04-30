#!/usr/bin/env python3
import os

def ThrowErrorIfObjectIsNull(obj):
    if not obj:
        raise TypeError("Object is set to null")

def writeToFileAndClose(fileDesc, data):
    ThrowErrorIfObjectIsNull(fileDesc)
    fileDesc.write(data)
    fileDesc.close()

def writeDataToFile(fileName, data):
    try:
        f = open(fileName, "w+")
        writeToFileAndClose(f, data)
    except IOError as e:
        raise e

#Main program
data = "Hi write me to the file please"
okayFile = "okay.txt"

try:
    writeDataToFile(okayFile, data)
except Exception as e:
    print("Uncaught exception %s" % (e))
    exit(1)

