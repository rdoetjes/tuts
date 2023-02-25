# Frustr8r

Is a puzzle game created in 2006 by Albert Eckhardt in 2006
It takes the NQueens problem and adds to it by already populating the board with
queens on certain (solvable) positions.

## nqueens.c

I decided to create a solver for this nice little affordable puzzle.<br/>
Initially I wanted to use a whole new concept of just dealing with an array with a structure holding the x,y and type of the 8 queens. I soon ran into the complexity of cleaning up backtracking -- although I did like my very small mathematical checker
to check collisions between queens that I could achieve by just simple arithmatic on each queen pair.<br/>
But sometimes you can definitely over engineer something especially something is fast to compute like this.<br/>
Hence I reverted back to using the NxN array and implementing the tried and tested backtracking algorithm and added the fact that certain columns were already populated and should not be touched.

## puzzle config file

The puzzle file parsing is far from robust, but it works!
You can create a simple file that holds one vector per line for each prepopulated queen.<br/>
*format:*
X,Y</br>

Example:
<pre>
3,5
4,8
6,3
7,7
8,2

</pre>

You can pass in the puzzle config file as argument.<br/>

<pre>
./nqueens puzzle.frustr8r
</pre>

![Alt text](./ss.png?raw=true "Screenshot")

Omitting the argument will just solve the puzzle without prepopulated queens

<pre>
./nqueens
</pre>

![Alt text](./ss1.png?raw=true "Screenshot")

## Building the application

There's a Makefile already present. I did not rely on any external libraries (hence the config file reading and parsing is a bit janky :) ) so simply type make in the directory with the Makefile.

<pre>
make
</pre>