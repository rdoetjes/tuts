variable high_value
4096 high_value !

: MUL2  
  DUP d2* DROP ;

: DIV2  
  DUP d2/ DROP ;

: P_HEX
  CR DUP hex. ;

: SCAN_UP
  BEGIN 
    DUP high_value @ < WHILE 
      MUL2 
      P_HEX 
      20 ms 
  REPEAT ;

: SCAN_DOWN
  BEGIN 
    DUP 1 > WHILE 
      DIV2 
      P_HEX 
      20 ms 
  REPEAT ;

: KITT
  1 
  BEGIN 
    1 1 = WHILE 
      SCAN_UP 
      SCAN_DOWN 
  REPEAT ;

KITT `cold
