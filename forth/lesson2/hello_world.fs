\ a loop that counts from 1 to the value on the stack 
\ and prints the output to the screen
\ and prints a new line (CR) after that
: SLEEP 1 DO i . i CR LOOP ;

\ set a 1001 on the stack and call sleep ( you could read this as sleep(1001) ) 
1001 SLEEP

\ bugger of the the OS
BYE
