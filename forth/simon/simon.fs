rng import
pin import

30 constant max-moves

variable moves max-moves allot \ maximum moves in simon

variable step 

: add-move ( value n --) \ adds the value to moves array[n]
  moves + c! ;

: get-move ( n -- n) \ gets the value from the moves list at position n
  moves + c@ ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl)
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  max-moves 0 do gen-move i add-move loop ;

: setup ( -- ) \ sets up the pins
  2 output-pin
  0 2 pin!
  3 output-pin
  0 3 pin!
  4 output-pin
  0 4 pin!
  5 output-pin
  0 5 pin!
  6 input-pin
  6 pull-up-pin
  7 input-pin
  7 pull-up-pin
  8 input-pin
  8 pull-up-pin
  9 input-pin
  9 pull-up-pin ;

: reset-game ( -- ) \ reset the game, set step to 0 and gen. new sequence
  0 step !
  gen-move-seq ;

: show-moves ( -- ) \ list all the moves to sceen, for  debugging
  max-moves 0 do i get-move . loop ;

: toggle-pin-ms ( ms move -- ) \ sets the pin corresponded by move (+2 to get gpio) on and of ms millisec
  swap
  dup
  2 + 1 swap pin!
  swap
  dup
  ms
  swap
  2 + 0 swap pin!
  ms
  ;

: simon
  setup
  reset-game
  4 1 do 
    cr ." level: " i . cr
    i 10 * 0 do 
      i get-move 250 toggle-pin-ms
      \ 300 2 + 1 swap pin!
      \ 300 ms
      \ 2 + 0 swap pin!
      \ 300 ms
    loop 
  loop ;
