0 constant black  
2 constant red  
1 constant brown  
3 constant orange  
4 constant yellow  
5 constant green  
6 constant blue  
7 constant violet  
8 constant grey  
9 constant white  

: 10power ( n -- n ) 
  dup                     \ deuplicate since the = if removes the it from the stack
  0 = if 1 nip exit then  \ zero case
  10                      \ multiplier 
  swap                    \ the number of times we multiply by needs to come 1st
  begin  
  dup                     \ dup the number of times the comparison removes it
  1 > while               \ as long as larger than 1 run the multiply
    swap 10 *             \ get the current result on top and multiply by 10     
    swap 1 -              \ get the counter on top and subtract 1     
  repeat   
  drop ;                  \ drop the number of multiplies (1) from stack
    
: r ( n1 n2 n3 -- n4 )
  10power                 \ only run the 10 power when arg is not 0
  rot 10 *                \ multiply that second argument (2nd ring) by 10
  rot +                   \ add the the rings together
  swap                    \ get the 3rd ring (multiples of 10) on top
  dup                     \ copy it because 0 > removes it
  0 > if                  \ if larger than 0 then we multiply the rd ring with the 2rings
    *
  else
    drop                  \ since we do not multiply we need to drop the 0 here
  then 
  ;


