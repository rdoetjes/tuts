# Switch Problem

Numberphile posed an interesting mathematical riddle in this video:
[![Watch the video](https://img.youtube.com/vi/-UBDRX6bk-A/maxresdefault.jpg)](https://youtu.be/-UBDRX6bk-A)
I didn't know the answer so I paused the video and wrote a little application in Python to get to the final answer before watching it and learning the actual maths behind it.

## A good beginner programming task

This is actually a fun little beginner programming task. It will show if a developer has concepts of structuring code in single purpose functions (or classes if you so desire). And it will also reveal if a developer knows about bitwise operations. The most efficient way to switch the light switches is by using an XOR, which does not incur branching operations and speeds up the program when dealling with enormous amounts of switching operations.<br/>
It will also test to see if a developer sees that in certain cases iterating starting from 1 instead of 0 is more beneficial (in this case since we are dealing with a modulo so starting from 0 would cause a divide by zero. So it's easier to start by 1). It will  also show if the developer remains consistent with their choice to start from either 1 or 0, when looping through the persons and switches. This consistency in choice in a file is important for maintenance and readability. 

