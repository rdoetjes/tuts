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
.const INTERRUPT_VECTOR_LO = $0318
.const INTERRUPT_VECTOR_HI = $0319

BasicUpstart2(main)

*=$1000
main:   
    jsr setup
    jsr start_timer_b
    jmp *
//
//SUBROUTINES
//

//sets up the counters
setup:
    jsr $e544           // clear screen for next number
    lda CONT_TIMER_B       // load the configuration of TIMER B	   
    and #%1000101   // configure to count clock cycles
    ora #%0000000   // start the timer
    sta CONT_TIMER_B       // store the timer config and start
    ldy #$00        // ofset pointer to bit in the random_byte pointed to by random_bytes+x

    sei                          //stop interrupts
    lda #<irq                   // setup low byte irq handler
    sta INTERRUPT_VECTOR_LO     // store the low byte of the interrupt handler

    lda #>irq                   // setup high byte irq handler
    sta INTERRUPT_VECTOR_HI     // store the high byte of the interrupt handler

    lda #%10010000               // enable NMI interrupts for FLAG and VBLANK
    sta $DD0D                    // store the interrupt flags
    cli
    rts

//interrupt handler
irq:
    //save registers
    pha
    tay
    pha
    tax
    pha

    ldy y                       //load y with the offset into the random_byte pointed to by random_bytes+x
    jsr mask_lowest_bit         //mask the lowest bit of timer B and if 1 shift without add if 0 shift with add of 1
    bne !shift_add_1+           // if 0 shift left and add one

    jsr shift_add_nothing       // if 1 shift left and don't add one
    jmp !exit+

    !shift_add_1:   
    jsr shift_add_1             // shift left and add 1 to the intermediate result in random_result byte that offset x points to

!exit:
    //restore registers
    pla
    tax
    pla
    tay
    pla
    bit $dd0d                   // ack interrupt so we can receive a new one
    rti


// start the timer B again
start_timer_b:
    ora #%00000001      //start timer B again
    sta CONT_TIMER_B           //set timer B register
    rts
 
// mask the lowerst bit of timer b LSB value
mask_lowest_bit:
    lda LSB_TIMER_B           // read lsb of timer B
    and #$1                   // mask lowest bit to see if the timer B count is odd or eveb
    rts

// shift the value in random_bytes+x left one position
shift_add_nothing:
    lda random_byte       // load the current random byte (X is the offset into the 255 random bytes)
    asl                   // if even just shift result left without adding one to the result
    jsr prep_for_next_bit // set everything up for the next bit (or next byte)
    rts

// shidt the value in random_bytes+x left one position and set the lowest bit to 1
shift_add_1:
    lda random_byte     // load the current random byte 
    asl                 // shift the temporary random result left
    ora #$01            // add one to the result
    jsr prep_for_next_bit   // set everything up for the next bit (or next byte)
    rts

//prepare the offset to y for the next bit, if y is 8 then set to 0 and increment x pointer
prep_for_next_bit:
    sta random_byte     //store the temporary random byte 
    iny                 //increment the bit counter
    cpy #$08            //if bit counter is 8 then reset bit counter and increment array counter to fill next byte
    beq !reset_y+       //if y==8 then reset
    sty y               // store the offset into the bit counter to the y register
    rts

    !reset_y:
    jsr reset_y         // reset y to 0
    rts

reset_y:
    // print the decimal format by loading random byte in X and setting A to 0
    ldx random_byte
    lda #$0
    jsr $bdcd           // print XA as integer
    jsr $ab3b           // print space

    ldy #$00            // reset the bit counter to 0
    sty random_byte     //store the temporary random byte
    sty y               // store the offset into the bit counter to the y register
    rts

*=$2000 "random bytes"
// array that holds the 255 random bytes that were generated
random_byte: .byte 0
y: .byte 0