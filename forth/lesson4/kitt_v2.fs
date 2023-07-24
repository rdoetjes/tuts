variable high_value
128 high_value !

\ use logical shift left(d2*) to multiple by two, which is more efficient
: MUL2  
  DUP d2* DROP ;

\ use logical shift right(d2/) to multiple by two, which is more efficient
: DIV2  
  DUP d2/ DROP ;

\ print the value on the stack to the sreen in hex format
: P_HEX
  CR DUP hex. ;

\ take the number on top of the stack multiply it by 2, display the value in hex, wait 20 ms
\ repeat this as long as the value on the is less than high_value
: SCAN_UP
  BEGIN 
    DUP high_value @ < WHILE 
      MUL2 
      P_HEX 
      20 ms 
  REPEAT ;

\ take the number on top of the stack divide it by 2 it, display the value in hex, wait 20 ms
\ repeat this as long as the value on the is greater than 1
: SCAN_DOWN
  BEGIN 
    DUP 1 > WHILE 
      DIV2 
      P_HEX 
      20 ms 
  REPEAT ;

\ indefinitely, scan the bit up and down just like the knightrider scanner
: KITT
  1 
  BEGIN 
    1 1 = WHILE 
      SCAN_UP 
      SCAN_DOWN 
  REPEAT ;

\ turn key to start the KITT word automatically after the cold start
KITT `cold
