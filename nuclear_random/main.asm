// The idea is that we have a free running timer that is incremented by the CPU clock. 
// When a decay happens the joystick button is pressed by a transistor connected between fire button and ground.
// As the button is pressed the LSB of the timer B is read and shifted into a byte when 8 bits are filled we have
// a truly random byte generated
// we actually generate 255 bytes in total, so we can do some analysis on the distribution of the bytes
// those random bytes start at $2000 (8192)
// Since we use a free running timer that we read when a decay happens, we avoid system entropy.
.const CONT_TIMER_B = $DC0F 
.const LSB_TIMER_B = $DC06 
.const CURRENT_SCANLINE = $D012
.const JOYSTICK_2 = $DC00
.const SCREEN_MEM = $0400

BasicUpstart2(main)

*=$1000
main:
    jsr setup
    jsr start_timer_b

loop:
    jsr print_bit_pattern

    // the detector has a 10ms pulse so we wait max 40ms for the next decay, so we don't flood the whole byte with same bit 
    jsr wait_for_raster_line    // give a 20 ms wait
    jsr wait_for_raster_line    // give a 20 ms wait

    jsr wait_for_decay_pulse    // wait till the isotope decays

    jsr mask_lowest_bit         //mask the lowest bit of timer B and if 1 shift without add if 0 shift with add of 1
    beq !shift_add_1+           // if 0 shift left and add one

    jsr shift_add_nothing       // if 1 shift left and don't add one
    beq !back_to_basic+         // if sub routine set Z flag then exit to basic

    jmp loop                    // read the next pulse

!shift_add_1:   
    jsr shift_add_1             // shift left and add 1 to the intermediate result in random_result byte that offset x points to
    beq !back_to_basic+         // if subroutine set Z flag then exit to basic

    jmp loop                    // read next bit

//Exit to basic gracefully
!back_to_basic:
    sei
    lda     #$37
    sta     $01             
    cli
    jsr     $ff8a           
    jsr     $ff81           
    jsr     $ff84           
    rts


//
//SUBROUTINES
//

//sets up the counters
setup:
    lda CONT_TIMER_B       // load the configuration of TIMER B	   
    and #%1000101   // configure to count clock cycles
    ora #%0000000   // start the timer
    sta CONT_TIMER_B       // store the timer config and start
    ldx #$00        // offset ponter into random_bytes array 
    ldy #$00        // ofset pointer to bit in the random_byte pointed to by random_bytes+x
    rts

// start the timer B again
start_timer_b:
    ora #%00000001      //start timer B again
    sta CONT_TIMER_B           //set timer B register
    rts
 
// this waits for raster line ff to give a bit of a delay
wait_for_raster_line:
    !wait:
    lda CURRENT_SCANLINE // read current rasterline
    cmp $ff              // if current rasterline is not ff than wait
    bne !wait-
    rts

// wait for the joystick button (which is triggered by the geiger counter when a isotope decays)
wait_for_decay_pulse:
   !await_pulse:
    //poll joystick button, a button press is a tick from the geiger counter
    lda JOYSTICK_2      // read joystick 2
    and #$10            // check button pressed
    bne !await_pulse-   // no new decay measured than wait
    rts

// mask the lowerst bit of timer b LSB value
mask_lowest_bit:
    lda LSB_TIMER_B           // read lsb of timer B
    sta SCREEN_MEM+1          // show the low byte timer B on screen
    and #$1                   // mask lowest bit to see if the timer B count is odd or eveb
    rts

// shift the value in random_bytes+x left one position
shift_add_nothing:
    lda random_bytes,X  // load the current random byte (X is the offset into the 255 random bytes)
    asl                 // if even just shift result left without adding one to the result
    jsr prep_for_next_bit // set everything up for the next bit (or next byte)
    rts

// shidt the value in random_bytes+x left one position and set the lowest bit to 1
shift_add_1:
    lda random_bytes,X  // load the current random byte (X is the offset into the 255 random bytes)
    asl                 // shift the temporary random result left
    ora #$01            // add one to the result
    jsr prep_for_next_bit   // set everything up for the next bit (or next byte)
    rts

//prepare the offset to y for the next bit, if y is 8 then set to 0 and increment x pointer
prep_for_next_bit:
    sta random_bytes,X  //store the temporary random byte back into the array   
    iny                 //increment the bit counter
    cpy #$08            //if bit counter is 8 then reset bit counter and increment array counter to fill next byte
    beq !reset_y+       //if y==8 then reset
    rts
!reset_y:
    jsr reset_y         // reset y to 0
    rts

reset_y:
    lda random_bytes,X  // get the whole random byte from offset x
    sta SCREEN_MEM           // show the random byte on screen
    ldy #$00            // reset the bit counter to 0
    inx                 // since y was 8, it means we should also increment x to point to the next random byte
    rts

print_bit_pattern:  
    // store a and y to the stack
    pha
    tya
    pha

    ldy #%000000001     // load bit pattern to #%000000001
    sty bit_pattern     // store bit pattern 
    ror bit_pattern     // set bit pattern to the one in the carry, so the next ror it will be 100000000
    ldy #8              // 8 bits to print

    !loop:
    ror bit_pattern     // shift bit to mask, to the right
    lda random_bytes, X // load the random byte from the array
    and bit_pattern     // mask the random byte with the bit pattern
    beq !draw_zero+     //  if the result is 0 then draw a zero

    lda #49             // draw a 1
    jmp !draw_bit+      // draw the actual 1 bit to the screen

    !draw_zero: 
    lda #48             // draw a 0

    !draw_bit:
    sta SCREEN_MEM+2,y         // draw the bit to the screen starting from 402 with offset y
    dey                 // decrement the bit counter
    bne !loop-          // if bit counter is not 0 then loop again

    // restore a and y from the stack
    pla                 
    tay
    pla
    rts

*=$2000 "random bytes"
// array that holds the 255 random bytes that were generated
random_bytes: .fill 255,00     
bit_pattern: .byte 0