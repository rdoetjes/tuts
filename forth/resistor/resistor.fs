0 constant black  
2 constant red  
1 constant brown  
3 constant orange  
4 constant yellow  
5 constant green  
6 constant blue  
7 constant violet  
7 constant purple  
8 constant grey  
9 constant white  

: exp ( u1 u2 -- u3 ) \ u3 = u1^u2
   dup
   0 = if rot 2drop 1 exit then \ drop the 10 and 0 that are passed
   over swap 1 ?do over * loop nip 
; 

: convert_fp ( u1 -- u1 u2 ) \ turn into a 'fp'
  1,0 * swap 
;

: ohm ( u1 u2 u3 -- u4 )
  \ calculate the 3rd ring exponent
  rot rot >r >r       \ store the first 2 args on return stack
  10                  \ exponent of 10
  swap                \ swap so we have 10 3rd_ring_value on stack
  exp                 \ catch the 0 case that should return 1
  r> r> swap          \ get the first 2 arguments back from return stack

  \ calculate the 1st_ring * 10 + 2nd_ring
  10 * +               

  \ if 3rd ring calculated exponent is not 0 than multiply that with the 1st_ring*10+2nd_ring value on the stack
  swap                
  dup                 
  0 > if 
    * 
  else 
    nip 
  then 
;                 

: print_resistor_value ( u1 -- ) \ takes the ohm value and turns it into ohm, KOhm or MOhm depending on the size
  dup
  1000000 / 1 >= if     \ is the value 1 or higher after dividing by 1e6 than it's mega ohms
    convert_fp          \ since we want a x.xx value we need to conver number into fixed point
    1000000,0 f/        \ divide by a million
    0,05 d+             \ round up for displaying a nice 1 digit begind the decimal place f.i 2,2 instead of 2,199999
    1 f.n ." MΩ" cr exit \ print number in fixed point format
  then 

  dup
  1000 / 1 >= if        \ is the value 1 or higher after dividing by 1e3 than it's Kohm
    convert_fp          \ convert from number into fixed point
    1000,0 f/           \ divide by 1000
    0,05 d+             \ round up for the so that 1.79999999 becomes 1.8
    1 f.n ." KΩ" cr exit \ print number in fixed point format
  then

  . ." Ω" cr          \ number as neither divisable by 1e6 and 1e3 so we treat it as ohms
;

: resistor_is ( u1, u2, u3 -- ) \ calculate the value of a resistor by passing in the colours as u1 u2 and u3
  ohm print_resistor_value   \ convert to ohms and then print into human readdable format
;
