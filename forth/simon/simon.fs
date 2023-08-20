rng import

31 constant max-moves

variable moves max-moves cells allot \ maximum moves in simon

variable step 
0 step !

: add-move ( value n --) \ adds the value to moves array[n]
  cells moves + ! ;

: get-move ( n -- n) \ gets the value from the moves list at position n
  cells moves + @  ;

 : gen-move ( -- n) \ generates random number between 0-4 (4 not incl)
   random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  31 0 do gen-move i add-move loop ;


