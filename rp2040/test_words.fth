\ RP2040 Forth Test Script
\ Use this to verify the implementation via the serial REPL

\ 1. Test Basic Arithmetic (Outputs in Hex)
\ 10 + 20 = 30 (0x1E)
10 20 + .

\ 5 * 5 = 25 (0x19)
5 5 * .

\ 2. Test Stack Manipulation
\ Should print 00000002 00000001
1 2 SWAP . .

\ 3. Test Dictionary Listing
\ This should show all built-in words: DUP, DROP, SWAP, +, -, *, WORDS, etc.
WORDS

\ 4. Test Runtime Compilation
\ Define a word that doubles a number
: DOUBLE DUP + ;

\ Verify it works (Should print 0000002A for 42)
21 DOUBLE .

\ 5. Verify the new word is in the dictionary
WORDS

\ 6. Test logic
\ 1 AND 2 = 0
1 2 AND .
\ 1 OR 2 = 3
1 2 OR .

\ 7. Test nested definitions
: SQUARE DUP * ;
: SUM-SQUARES SQUARE SWAP SQUARE + ;

\ 3^2 + 4^2 = 9 + 16 = 25 (0x19)
3 4 SUM-SQUARES .
