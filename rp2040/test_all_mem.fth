\ Test Script for Forth Memory and Variable Words

\ --- Test VARIABLE, @, and ! ---
VARIABLE X
VARIABLE Y
5 X !
10 Y !
X @ .   \ Should print 5
Y @ .   \ Should print 10
X @ Y @ + X !
X @ .   \ Should print 15

\ --- Test C@ and C! ---
VARIABLE CH
255 CH !
65 CH C! \ 'A'
CH @ .   \ Should print 65 (since it was stored at the start of the word)
CH C@ .  \ Should print 65

\ --- Test HERE, ALLOT, and ALIGN ---
HERE .   \ Print current HERE
10 ALLOT
ALIGN
HERE .   \ Print new HERE (should be at least 12 bytes further and 4-aligned)

\ --- Test , (comma) ---
HERE     \ Save address
123 ,
456 ,
789 ,
DUP @ .      \ Should print 123
DUP 4 + @ .  \ Should print 456
8 + @ .      \ Should print 789

\ --- Test multiple variables in a loop ---
VARIABLE COUNTER
0 COUNTER !
10 0 DO
  COUNTER @ 1 + COUNTER !
LOOP
COUNTER @ . \ Should print 10
