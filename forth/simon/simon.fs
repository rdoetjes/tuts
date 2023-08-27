rng import
pin import
pwm import

31 constant sequence-size \ the maximum amount of steps (+1) in sequence, +1 so we can process logic from 1 instead of 0

5 constant slice \ this slice for pin 10 (5A slice)

variable sequence sequence-size allot \ array that holds the sequence

variable max-steps  \ number of steps for a level.

variable step   \ current step of the sequence

variable speed  \ speed of simon showing the sequence (gets faster every 10 steps)

: add-move ( move --) \ adds the value to sequence array[n]
  sequence + c! ;

: get-move ( index -- move )) \ gets the value from the sequence list at position n
  sequence + c@ ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl) by using logic and 3 on random int
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  sequence-size 0 ?do gen-move i add-move loop ;

: cs ( -- ) \ clear the stack when something is on there (should never need to clean)
  depth 0 > if depth 0 ?do ." aye!" drop loop then ; 

: set-leds-off ( -- ) \ switch the 4 leds off, this is the starting state of the game
  6 2 ?do 0 i pin! loop ;

: setup ( -- ) \ sets up the pins in such away we can uses maths to calculate step value, led pin and switch from one or the other
  cs
  2 output-pin 
  3 output-pin
  4 output-pin
  5 output-pin

  6 input-pin 6 pull-up-pin
  7 input-pin 7 pull-up-pin
  8 input-pin 8 pull-up-pin
  9 input-pin 9 pull-up-pin 

  \ pwn setup for sound
  slice bit disable-pwm 
  10 pwm-pin 
  false slice pwm-phase-correct! 
  6000 slice pwm-counter-compare-a! 
  0 4 slice pwm-clock-div! ;

: reset-game ( -- ) \ reset the game, set step to 0 and gen. new sequence
  cs              \ clear stack just in case
  set-leds-off    \ set everything to a known state
  0 step !        \ set step to 0, we increment by one in beginning of game loop 
  gen-move-seq ;  \ generate the 30 random steps whichs make up the sequence

: play-beep ( move -- ) \ play a sound corresponding to the move
  0 4 slice pwm-clock-div! \ set octave higher
  dup 0 = if 45000 slice pwm-top! then
  dup 1 = if 55000 slice pwm-top! then
  dup 2 = if 60000 slice pwm-top! then
  dup 3 = if 65500 slice pwm-top! then

  \ 0 slice pwm-counter! 
  slice bit enable-pwm
  drop ;

: stop-beep
  slice bit disable-pwm ;

: toggle-pin-ms ( move ms -- ) \ sets the pin corresponding to move (+2 to get gpio) on and of for ms millisec
  swap dup 2 + toggle-pin
  dup play-beep
  swap dup ms
  swap 2 + toggle-pin
  stop-beep
  ms ;

: simons-move ( n -- ) \ plays the steps from 0 to n
  0 ?do 
    i get-move speed @ toggle-pin-ms \ get the move, get the speed in ms toggle-pin which is move+2 to get gpio
  loop ;

: key-is-down ( switch_pin -- n ) \ lights up the led as long as pressed and returns the step value calculated from switch number (- 6)
  dup 4 - toggle-pin     \ subtract 4 from the pressed switch to turn on correponding led
  dup 6 - play-beep      \ subtract 4 from the pressed switch that corresponds with the sound of the led
  begin dup pin@ -1 = until \ loop as long as we hold
  stop-beep              \ stop the pwm to stop the sound
  dup 4 - toggle-pin     \ subtract 4 from pressed switch to turn off corresponding led
  6 -                    \ subtract 6 from switch pin to get the step value of the sequence
  50 ms ;                \ debounce the release

: poll-keys ( -- ) \ checks for a keypress and timeout after ~1500ms
  0         \ counter for timeout
  begin 
    10 ms   \ delay to add to time out (150 * 10 ms is 1500ms and then exit with -1) doubles as debounce
    1 + dup 150 = if drop -1 exit then 
    0 6 pin@ = if drop 6 key-is-down exit then 
    0 7 pin@ = if drop 7 key-is-down exit then 
    0 8 pin@ = if drop 8 key-is-down exit then 
    0 9 pin@ = if drop 9 key-is-down exit then 
  again ; 

: game-over-sound ( -- ) \ plays deep buzz
  0 10 slice pwm-clock-div! 
  65500 slice pwm-top!
  6000 slice pwm-counter-compare-a! 
  0 slice pwm-counter! 
  slice bit enable-pwm ;

: game-over ( -- ) \ turn all leds on to indicate game over
  game-over-sound
  10 0 ?do 
    2 toggle-pin 
    3 toggle-pin 
    4 toggle-pin 
    5 toggle-pin 
    100 ms
 loop 
 stop-beep ;

: you-beat-the-game ( -- ) \ a little light show
  14 0 ?do 
    2 toggle-pin 0 play-beep
    4 toggle-pin 1 play-beep
    100 ms
    3 toggle-pin 2 play-beep
    5 toggle-pin 3 play-beep
    100 ms
  loop stop-beep ;

: players-move ( step -- n) \ reads the buttons and compare the value against simon's sequence at that step; 
  0 ?do                     \ returns value at step or -1 when player was wrong
    poll-keys
    dup -1 = if exit then   \ if -1 then timeout accord and return -1 for gameover
    dup i get-move <> if drop -1 exit then \ value from poll-keys didn't match the value from get-move, return -1 for gameover
    drop
  loop ;

: wait-for-level-select ( -- ) \ wait for one of the 4 buttons to be pressed to set level and start game
  begin
    10 ms
    0 6 pin@ = if drop 10 max-steps ! exit then 
    0 7 pin@ = if drop 15 max-steps ! exit then 
    0 8 pin@ = if drop 20 max-steps ! exit then
    0 9 pin@ = if drop 30 max-steps ! exit then
  again ;

: set-speed ( n -- ) \ n = step, the procedure will update global var speed
  \ speed up every 10 steps 
  dup 5 < if 300 speed ! then    \ 1..5 300ms
  dup 4 > if 250 speed ! then    \ 5..10 250ms
  dup 10 > if 200 speed ! then   \ 11 .. 15 200ms
  dup 15 > if 175 speed ! then   \ 16..20 175ms
  dup 20 > if 150 speed ! then   \ 21..30 150ms
  24 > if 130 speed ! then ;     \ 25..30 130ms (don't dup we don't return value)

: game-loop ( -- n ) \ loop from 1 to sequence-size if you don't get to sequence-size then returns -1 else returns 100
  begin
    step @ 1 + step !     \ go to next step in the sequence (steps==0 when game is reset)
    step @ set-speed      \ set simon's speed depedning on step
    step @ simons-move    \ simon plays the sequence until step
    step @ players-move   \ user repeats simons sequence
    -1 = if -1 exit then  \ we break out when players-move is -1 (indicating gameover)
    800 ms                \ wait 800 ms and then let simon start next sequence
    step @ max-steps @ =  \ did we reach the whole sequence no? continue TODO: victory light show after until
  until 100 ;             \ 100 is to indicate you beat the whole sequemce

: simon ( -- )            \ SIMON game entry point, loops indefinitely
  setup
  begin
    reset-game
    wait-for-level-select
    1000 ms               \ wait 1 second for player to get ready
    game-loop
    dup -1 = if drop game-over then
    dup 100 = if drop you-beat-the-game then
  again ;
