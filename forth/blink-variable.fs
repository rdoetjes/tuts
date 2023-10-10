led import

: blink ( n-- ) \ n = the amount of ms led is on
  begin 
    dup 
    green toggle-led 
    ms 
    key? until ;


