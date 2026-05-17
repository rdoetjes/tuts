\ Comprehensive Test Suite for Division (/)

\ --- Test Positive Integers ---
12 2 / .     \ Expected: 6
100 10 / .   \ Expected: 10
1000 5 / .   \ Expected: 200

\ --- Test Division by One ---
42 1 / .     \ Expected: 42

\ --- Test Zero Dividend ---
0 5 / .      \ Expected: 0

\ --- Test Signed Division (Positive / Negative) ---
10 -2 / .    \ Expected: -5

\ --- Test Signed Division (Negative / Positive) ---
-10 2 / .    \ Expected: -5

\ --- Test Signed Division (Negative / Negative) ---
-20 -4 / .   \ Expected: 5

\ --- Test Remainder Discarding (Floor Division) ---
7 2 / .      \ Expected: 3
11 3 / .     \ Expected: 3

\ --- Test Large Numbers ---
2000000 2 / . \ Expected: 1000000

\ --- Test Division by Zero (Safe check) ---
\ The current implementation returns 0 for division by zero to avoid crashes.
10 0 / .     \ Expected: 0
