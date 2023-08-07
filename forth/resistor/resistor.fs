0,0 constant black  
2,0 constant red  
1,0 constant brown  
3,0 constant orange  
4,0 constant yellow  
5,0 constant green  
6,0 constant blue  
7,0 constant violet  
7,0 constant purple  
8,0 constant grey  
9,0 constant white  

: exp ( u1 u2 -- u3 ) \ u3 = u1^u2
   dup
   0 = if rot 2drop 1 exit then
   over swap 1 ?do over * loop nip ; 

: convert_fp
  1,0 * swap ;

: ohm ( u1 u2 u3 -- u4 )
  \ calculate the 3rd ring exponent
  rot rot >r >r       \ store the first 2 args on return stack
  10                  \ exponent of 10
  swap                \ swap so we have 10 3rd_ring_value on stack
  exp                 \ catch the 0 case that should return 1
  r> r> swap          \ pop the first 2 arguments form return stack

  \ calculate the 1st_ring * 10 + 2nd_ring
  10 * +               

  \ if 3rd ring calculated exponent is not 0 than multiply that with the 1st_ring*10+2nd_ring value on the stack
  swap                
  dup                 
  0 > if * else nip then ;                 

  : into_hr_fmt
    dup
    1000000 / 1 >= if 
      convert_fp 
      1000000,0 f/ f. ." MOhm" cr exit then

    dup
    1000 / 1 >= if 
      convert_fp 
      1000,0 f/ f. ." KOhm" cr exit then

    . ." Ohm" cr ;

