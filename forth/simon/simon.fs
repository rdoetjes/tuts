rng import

31 constant max-moves

variable moves max-moves allot \ maximum moves in simon

variable step 
0 step !

: add-move ( value n --) \ adds the value to moves array[n]
  moves + c! ;

: get-move ( n -- n) \ gets the value from the moves list at position n
  moves + c@  ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl)
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  31 0 do gen-move i add-move loop ;

: reset-game ( -- ) \ reset the game, set step to 0 and gen. new sequence
  0 step !
  gen-move-seq ;

: simon
  reset-game
  ;
