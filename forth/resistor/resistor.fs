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
  dup 
  1 = if 10 swap drop exit then   \ cases that do not require multiplication
  dup
  0 = if 1 swap drop exit then    \ case that does not require multiplication
  10                      \ multiplier 
  swap                    \ the number of times we multiply by needs to come 1st
  begin  
  dup                     \ dup the number of times the comparison removes it
  1 > while               \ as long as larger than 1 run the multiply
    swap                  \ get the result (or 10 in first iter) on top of the stack
    10 *                  \ multiply result by 10
    swap                  \ get the the number of multiplies on the stack
    1 -                   \ subtract one from the number of multiplies
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


