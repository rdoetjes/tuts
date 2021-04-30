#!/usr/bin/env python3
import os

class ErrorCheck:
    @staticmethod
    def ThrowTypeErrorIfObjectIsNull(obj):
        if not obj:
            raise TypeError("ErrorCheck object is set to null ")
    
class File:
    @staticmethod
    def __writeDataToFileObjectAndClose(fileDesc, data):
        ErrorCheck.ThrowTypeErrorIfObjectIsNull(fileDesc)
        fileDesc.write(data)
        fileDesc.close()
     
    @staticmethod
    def writeText(fileName, text, mode):
        try:
            f = open(fileName, mode)
            ErrorCheck.ThrowTypeErrorIfObjectIsNull(f)
            File.__writeDataToFileObjectAndClose(f, text)
        except TypeError as e:
            raise TypeError("the file %s could not be open, fileObject is None" % (fileName) )
        except IOError as e:
            raise e

    @staticmethod
    def createFileOrAppendWithText(fileName, text):
        try:
            File.writeText(fileName,text, "a+")
        except IOError as e:
            raise e

    @staticmethod
    def createFileWithText(fileName, text):
        try:
            File.writeText(fileName,text, "w+")
        except IOError as e:
            raise e

#Main program
textToSave = "Hi write me to the file please\n"
okayFile = "okay.txt"

try:
    File.createFileOrAppendWithText(okayFile, textToSave)
except Exception as e:
    print("Uncaught exception %s" % (e))
    exit(1)

