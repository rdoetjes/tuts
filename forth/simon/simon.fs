rng import
pin import
pwm import

31 constant sequence-size \ the maximum amount of steps (+1) in sequence, +1 so we can process logic from 1 instead of 0
5 constant slice \ this slice for pin 10 (5A slice)
2 constant r-led
3 constant g-led
4 constant y-led
5 constant b-led
6 constant r-switch
7 constant g-switch
8 constant y-switch
9 constant b-switch
10 constant piezo
-1 constant you-lost
100 constant you-won

variable sequence sequence-size allot \ array that holds the sequence

variable max-steps  \ number of steps for a level.

variable step   \ current step of the sequence

variable speed  \ speed of simon showing the sequence (gets faster every 10 steps)

: add-move ( move --) \ adds the value to sequence array[n]
  sequence + c! ;

: get-move ( index -- move ) \ gets the value from the sequence list at position n
  sequence + c@ ;

: gen-move ( -- n) \ generates random number between 0-4 (4 not incl) by using logic and 3 on random int
  random 3 and ;

: gen-move-seq ( -- ) \ gen. the 31 different random values for the game
  sequence-size 0 ?do gen-move i add-move loop ;

: cs ( -- ) \ clear the stack when something is on there (should never need to clean)
  depth 0 > if depth 0 ?do ." aye!" drop loop then ; 

: set-leds-off ( -- ) \ switch the 4 leds off, this is the starting state of the game
  b-led 1 + r-led ?do 0 i pin! loop ;

: setup ( -- ) \ sets up the pins in such away we can uses maths to calculate step value, led pin and switch from one or the other
  cs
  r-led output-pin 
  g-led output-pin
  y-led output-pin
  b-led output-pin

  r-switch input-pin r-switch pull-up-pin
  g-switch input-pin g-switch pull-up-pin
  y-switch input-pin y-switch pull-up-pin
  b-switch input-pin b-switch pull-up-pin 

  \ pwn setup for sound
  slice bit disable-pwm 
  piezo pwm-pin 
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
  50 ms ;                \ debounce the release, just in case

: poll-keys ( -- ) \ checks for a keypress and timeout after ~1500ms
  0                \ counter for timeout
  begin 
    10 ms          \ delay to add to time out (150 * 10 ms is 1500ms and then exit with -1) doubles as debounce
    1 + dup 150 = if drop -1 exit then 
    0 6 pin@ = if drop r-switch key-is-down exit then 
    0 7 pin@ = if drop g-switch key-is-down exit then 
    0 8 pin@ = if drop y-switch key-is-down exit then 
    0 9 pin@ = if drop b-switch key-is-down exit then 
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
    r-led toggle-pin 
    g-led toggle-pin 
    y-led toggle-pin 
    b-led toggle-pin 
    100 ms
 loop 
 stop-beep ;

: you-beat-the-game ( -- ) \ a little light show
  14 0 ?do 
    r-led toggle-pin 0 play-beep
    y-led toggle-pin 1 play-beep
    100 ms
    g-led toggle-pin 2 play-beep
    b-led toggle-pin 3 play-beep
    100 ms
  loop stop-beep ;

: players-move ( step -- n) \ reads the buttons and compare the value against simon's sequence at that step; 
  0 ?do                     \ returns value at step or -1 when player was wrong
    poll-keys
    dup -1 = if exit then   \ if -1 then timeout accord and return -1 for gameover
    dup i get-move <> if drop -1 exit then \ value from poll-keys didn't match the value from get-move, return -1 for gameover
    drop
  loop 0 ;

: wait-for-level-select ( -- ) \ wait for one of the 4 buttons to be pressed to set level and start game
  begin
    10 ms
    0 6 pin@ = if 10 max-steps ! exit then 
    0 7 pin@ = if 15 max-steps ! exit then 
    0 8 pin@ = if 20 max-steps ! exit then
    0 9 pin@ = if 30 max-steps ! exit then
  again ;

: set-speed ( n -- ) \ n = step, the procedure will update global var speed (order of IFs is important)
  dup  5 < if drop 300 speed ! exit then   \ 1..5 300ms drop removes the step we do not want to return
  dup 24 > if drop 130 speed ! exit then   \ 25..30 130ms (don't dup we don't return value)
  dup 20 > if drop 150 speed ! exit then   \ 21..30 150ms
  dup 15 > if drop 175 speed ! exit then   \ 16..20 175ms
  dup 10 > if drop 200 speed ! exit then   \ 11 .. 15 200ms
  dup  4 > if drop 250 speed ! exit then ; \ 5..10 250ms

: game-loop ( -- n ) \ loop from 1 to sequence-size if you don't get to sequence-size then returns -1 else returns 100
  begin
    step @ 1 + step !     \ go to next step in the sequence (steps==0 when game is reset)
    step @ set-speed      \ set simon's speed depedning on step
    step @ simons-move    \ simon plays the sequence until step
    step @ players-move   \ user repeats simons sequence
    you-lost = if -1 exit then  \ we break out when players-move is -1 (indicating gameover)
    800 ms                \ wait 800 ms and then let simon start next sequence
    step @ max-steps @ =  \ did we reach the whole sequence no? continue TODO: victory light show after until
  until you-won ;         \ 100 is to indicate you beat the whole sequemce

: simon ( -- )            \ SIMON game entry point, loops indefinitely
  setup
  begin
    reset-game
    wait-for-level-select
    1000 ms               \ wait 1 second for player to get ready
    game-loop
    dup you-lost = if drop game-over then
    dup you-won = if drop you-beat-the-game then
    key?
  until ;
