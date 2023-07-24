\ lets make a "knightrider scanner"
\ we will print the value of the LED BAR from kitt
\ as a hex decimal value later on the PiPico we do it with LEDS
\ for now the hex values represent the 8 bit output

\ multiply the value on the stack by two to move to the next LED
\ in affact making the LED move from one to the next
: MUL2  
  DUP 2 * SWAP DROP ;

\ divide the value on the stack by two to move to the next LED
\ in affact moving the LED down from it'a current position
: DIV2  
  DUP 2 / SWAP DROP ;

\ print the hex value from the current value of the stack
: P_HEX
  CR DUP hex. ;

\ multiply the value thats on the top of the stack by 2
\ print out it's hex value and repeat as long as the value
\ is less than or equal to 128
: SCAN_UP
  BEGIN DUP 128 < WHILE MUL2 P_HEX REPEAT ;

\ divide the value thats on the top of the stack by 2
\ print out it's hex value and repeat as long as the value
\ is less than or equal to 1
: SCAN_DOWN
  BEGIN DUP 1 > WHILE DIV2 P_HEX REPEAT ;

\ this will multiply and divide the value on the stack
\ in affact making it look like a single LED (or bit) moving left
\ and right 
: KITT
  1 BEGIN 1 1 = WHILE SCAN_UP SCAN_DOWN REPEAT ;
