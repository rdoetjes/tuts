rng import
pin import

31 constant max-moves \ the maximum amount of steps (+1) in sequence

variable moves max-moves allot \ array that holds the sequence

variable step   \ current step of the sequence
1 step !

variable speed  \ speed of simon showing the sequence (gets faster every 10 steps)
250 speed !     

: add-move ( value n --) \ adds the value to moves array[n]
  moves + c! ;

: get-move ( n -- n) \ gets the value from the moves list at position n
  moves + c@ ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl) by using logic and 3 on random int
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  max-moves 0 do gen-move i add-move loop ;

: cs ( -- ) \ clear the stack when something is on there (should never need to clean)
  depth 0 > if depth 0 do drop loop then ; 

: set-leds-off ( -- ) \ switch the 4 leds off, this is the starting state of the game
  6 2 do 0 i pin! loop ;

: setup ( -- ) \ sets up the pins in such away we can uses maths to calculate step value, led pin and switch from one or the other
  cs
  2 output-pin 
  3 output-pin
  4 output-pin
  5 output-pin
  6 input-pin 6 pull-up-pin
  7 input-pin 7 pull-up-pin
  8 input-pin 8 pull-up-pin
  9 input-pin 9 pull-up-pin ;

: reset-game ( -- ) \ reset the game, set step to 0 and gen. new sequence
  cs              \ clear stack just in case
  set-leds-off    \ set everything to a known state
  1 step !        \ set step to first step of the sequence
  gen-move-seq ;  \ generate the 30 random steps whichs make up the sequence

: toggle-pin-ms ( ms move -- ) \ sets the pin corresponding to move (+2 to get gpio) on and of for ms millisec
  swap dup 2 + toggle-pin
  swap dup ms
  swap 2 + toggle-pin
  ms ;

: simons-move ( n -- ) \ plays the steps from 0 to n
  0 do 
    i get-move speed @ toggle-pin-ms \ get the move, get the speed in ms toggle-pin which is move+2 to get gpio
  loop ;

: key-down ( n -- n ) \ lights up the led as long as pressed and returns the step value calculated from switch number (- 6)
  dup 4 - toggle-pin  \ subtract 4 from the pressed switch pin to get corresponding led pin
  begin 30 ms dup -1 swap pin@ = until
  dup .
  dup 4 - toggle-pin
  6 -                \ subtract 6 from switch pin to get the step value of the sequence
  dup . cr ;

: poll-keys ( -- ) \ checks for a keypress and timeout after ~1500ms
  0
  begin 
    10 ms           \ delay to add to time out (150 * 10 ms is 1500ms and then exit with -1)
    1 + dup 150 = if drop -1 exit then 
    0 6 pin@ = if drop 6 key-down exit then 
    0 7 pin@ = if drop 7 key-down exit then 
    0 8 pin@ = if drop 8 key-down exit then 
    0 9 pin@ = if drop 9 key-down exit then 
  again ; 

: game-over ( -- ) \ turn all leds on to indicate game over
  10 0 do 
    2 toggle-pin 
    3 toggle-pin 
    4 toggle-pin 
    5 toggle-pin 
    100 ms
 loop ;

: you-beat-the-game ( -- ) \ a little light show
  14 0 do 
    2 toggle-pin 
    4 toggle-pin 
    100 ms
    3 toggle-pin 
    5 toggle-pin 
    100 ms
  loop ;

: players-move ( step -- n) \ reads the buttons and compare the value against simon's sequence at that step; 
                            \ returns value at step or -1 when player was wrong
  0 do 
    poll-keys
    dup -1 = if exit then          \ if -1 then timeout accord and return -1 for gameover
    dup i get-move <> if drop -1 exit then \ value from poll-keys didn't match the value from get-move, return -1 for gameover
    drop
  loop ;

: wait-for-red-button ( -- ) \ wait for red button to be pressed, which doubles as start button
  begin
    10 ms
    0 6 pin@ = 
  until ;

: game-loop ( -- n ) \ loop from 1 to max-moves if you don't get to max-moves then returns -1 else returns 100
  begin
    \ speed up every 10 steps 
    step @  9 <= if 300 speed ! then 
    step @ 10 >= if 200 speed ! then
    step @ 20 >= if 150 speed ! then

    step @ simons-move  \ simon plays the sequence until step
    
    step @ players-move \ user repeats simons sequence
    -1 = if -1 exit then \ we break out when players-move is -1 (indicating gameover)

    step @ 1 + step !   \ go to next step in the sequence
    
    800 ms              \ wait 800 ms and then let simon start next sequence
    step @ max-moves =  \ did we reach the whole sequence no? continue TODO: victory light show after until
  until 100 ;           \ 100 is to indicate you beat the whole sequemce

: simon ( -- )          \ SIMON game entry point, loops indefinitely
  setup
  begin
    reset-game
    wait-for-red-button
    1000 ms
    game-loop
    dup -1 = if drop game-over then
    dup 100 = if drop you-beat-the-game then
  again ;
