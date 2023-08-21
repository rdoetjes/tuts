rng import
pin import

30 constant max-moves

variable moves max-moves allot \ maximum moves in simon

variable step 
1 step !

variable speed
250 speed !

: add-move ( value n --) \ adds the value to moves array[n]
  moves + c! ;

: get-move ( n -- n) \ gets the value from the moves list at position n
  moves + c@ ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl)
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  max-moves 0 do gen-move i add-move loop ;

: setup ( -- ) \ sets up the pins
  2 output-pin 0 2 pin!
  3 output-pin 0 3 pin!
  4 output-pin 0 4 pin!
  5 output-pin 0 5 pin!
  6 input-pin 6 pull-up-pin
  7 input-pin 7 pull-up-pin
  8 input-pin 8 pull-up-pin
  9 input-pin 9 pull-up-pin ;

: reset-game ( -- ) \ reset the game, set step to 0 and gen. new sequence
  1 step !
  gen-move-seq ;

: toggle-pin-ms ( ms move -- ) \ sets the pin corresponded by move (+2 to get gpio) on and of ms millisec
  swap
  dup
  2 + toggle-pin
  swap
  dup
  ms
  swap
  2 + toggle-pin
  ms
  ;

: simons-move ( n -- ) \ plays the steps from 0 to n
  0 do 
    i get-move speed @ toggle-pin-ms \ get the move, get the speed in ms toggle-pin which is move+2 to get gpio
  loop 
;

: key-down ( n -- n ) \ lights up the led a maximum of ms to corresponding led and returns the move 0,1,2 or 3
  depth 2 = if swap drop then \ sometimes we seem to have a superious entry keybouce?
  dup 4 - toggle-pin
  begin 100 ms dup -1 swap pin@ = until
  dup 4 - toggle-pin
  6 -
  ;

: poll-keys ( -- ) \ checks for a keypress and timesout after 2000ms
  0
  begin 
  10 ms
  1 + dup 150 = if -1 exit then 
    0 6 pin@ = if drop 6 key-down exit then 
    0 7 pin@ = if drop 7 key-down exit then 
    0 8 pin@ = if drop 8 key-down exit then 
    0 9 pin@ = if drop 9 key-down exit then 
  again ; 

: cs depth 0 do . loop ; 

: game-over
  1 2 pin!
  1 3 pin!
  1 4 pin!
  1 5 pin!
  depth 0 > if cs then \ clear stack 
;

: players-move ( step -- n) \ reads the buttons and checks the value
  0 do 
    poll-keys
    dup -1 = if drop -1 exit then 
    dup i get-move <> if -1 exit then
  loop
;

: simon
  setup
  reset-game
  1000 ms
  begin
    step @ 9 <= if 300 speed ! then
    step @ 10 >= if 200 speed ! then
    step @ 20 >= if 150 speed ! then

    \ simon's move
    step @ simons-move

   \ do user
    step @ players-move
    -1 = if game-over exit then

    step @ 1 + step !
    
    800 ms
    step @ max-moves =
  until ;
