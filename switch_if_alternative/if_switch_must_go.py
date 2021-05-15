#!/usr/bin/env pyhton3

def option_1():
    return "This what you do when 1 is chosen"

def option_2():
    return "this was option 2"

def option_3():
    return "and now you chose option 3"

if __name__ == "__main__":
    s = input("Type one, two or three: ")

    #This sucks
    if s == "one" or s == "1":
       print(option_1())
    if s == "two" or s == "2":
       print(option_2())
    if s == "three" or s == "3":
       print(option_3())

    #This is so much shorter and cleaner
    functions = {"one": option_1, "1": option_1, "two": option_2, "2": option_2, "three": option_3, "3": option_3}
    print(functions[s]())
