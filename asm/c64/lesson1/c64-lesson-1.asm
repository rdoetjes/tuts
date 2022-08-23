.const SCREEN = $0400

BasicUpstart2(main)

* = $01000

main:
  jsr cls         //jump to subroutine cls
  //jsr $FF81     //jump to kernal clear screen subroutine
  rts             //return to basic

cls:
  //A accumulator, X and Y index registers
  lda #32         //load a with value 32 (space)
  ldx #0          //load x with value 0 (is our counter)

cls_loop:         //fill screen memory with 1000 spaces
  //screen size is 40x25 = 1000
  sta SCREEN, X           //store value in a to screen  x
  sta SCREEN + $0100, X   //store value in a to screen + $100 x
  sta SCREEN + $0200, X   //store value in a to screen + $200 x
  sta SCREEN + $02e8, X   //store value in a to screen + $300 x
  dex                     //subtract one from x
  bne cls_loop            // jump to cls_loop when x is not 0

  rts     //return from subroutine (back to basic)