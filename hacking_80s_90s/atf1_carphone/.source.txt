BasicUpstart2(main)

        * = $1000

number_to_convert:
.text "s1234520717666e"         // this can be overwritten in basic by poking from 4097 (overwrites the 1 after the S)
.text "$"

// --- Constants ---
.const CIA1_TIMER_A_LOW  = $dc04
.const CIA1_TIMER_A_HIGH = $dc05
.const CIA1_ICR          = $dc0d
.const CIA1_CRA          = $dc0e
.const IRQ_VECTOR        = $0314
.const TIMER_10ms = $268C
.const TIMER_START_CONTINUOUS = %00010001       // load value, start timer
//SID1 as in SID_CHANNEL1
.const SID1_FREQ_HI = $d401
.const SID1_FREQ_LO = $d400
.const SID_CHANNEL_A_VOLUME = $d41b
.const SID_VOLUME = $d418
.const SID1_ATTACK_DECAY = $d405
.const SID1_SUSTAIN_RELEASE = $d406
.const SID1_OSCILATOR_TYPE = $d404


// --- Main Program ---
main:
        // setup the cia, sid and reset variables
        jsr setup
        // send preeamble 600ms 2050hz
        jsr ms600_preeamble
        // process the dial
        jsr process_dial
        //reset the variables so we can do it again
        jsr reset_var
        // do it twice like a car phone does, we found out in the past that this increased the reliability
        jsr process_dial
        // turn frequency to 0Hz (silence)
        jsr hz0
        // exit to basic
        jmp ($A002)

reset_var:
        jsr $e544       //clr screen

        lda #$00
        sta timer_600ms_lapsed
        sta byte_count
        sta offset_table
        sta irq_counter
        rts

// clean variables, setup the sid and cia
// configure the CIA to trigger 10ms interrupts, in the interrupt routine we'll increment the irq_counter
// and we check in the main what beep to send and wait during the duration until the irq_counter is incremented again
// in affect giving a 10ms beep, we can process ~8000 instructions in that period so whe are good, we miss mere microseconds of the 10ms
// beep because of parseing the number. So we are well withing range. In more timing critical situations we could change the timer value to get that longer
// duration -- but microseconds on millisecond lengths are not that important.
// But defenitely something you need to be aware of when dealing with realtime systems.
// Hence this can't be done in BASIC.
setup:
        sei
       
        jsr reset_var
        jsr setup_sid
        jsr hz2070

        lda #<dial_irq
        sta IRQ_VECTOR
        lda #>dial_irq
        sta IRQ_VECTOR+1

        // Configure and start the 10ms CONTINUOUS timer
        lda #<TIMER_10ms
        sta CIA1_TIMER_A_LOW
        lda #>TIMER_10ms
        sta CIA1_TIMER_A_HIGH

        lda #%10000001
        sta CIA1_ICR 

        lda #TIMER_START_CONTINUOUS
        sta CIA1_CRA

        cli
        rts

// wait for 600ms and then turn off the frequency
ms600_preeamble:
        lda timer_600ms_lapsed
        cmp #1  
        bne ms600_preeamble

        // 600ms clear preeamble send
        lda #$0
        jsr hz0
        rts

// parse the number from the number_to_convert string
// fetch the offset from the table that contains the binary bytes to send
// send each bit as the correct frequency for 5ms
process_dial:
        ldx byte_count
        // read letter from string and check if we are at the end of the string
        lda number_to_convert,x
        sta $0400,x

        sec
        sbc #'0'         // Normalize ASCII digit chars to 0-9 (subtract 48)
        cmp #10
        bcc !digit+      // If 0–9, go handle digit
        clc
        adc #'0'         // Restore A

        cmp #'s'
        beq !set_s+
        cmp #'e'
        beq !set_e+
        cmp #'c'
        beq !set_c+
        jmp !no_match+

!digit:
        clc
        asl              // Multiply index by 2 (each entry = 2 bytes)
        adc #6
        tay
        sty offset_table
        lda table, y
        jmp !done+

!set_s:
        lda #0
        sta offset_table
        lda table+0
        jmp !done+

!set_e:
        lda #2
        sta offset_table
        lda table+2
        jmp !done+

!set_c:
        lda #4
        sta offset_table
        lda table+4
        jmp !done+

!no_match:
    // No valid character matched. Do nothing or handle error
        rts

!done:

        //process bytes (send the 8 bits)
!:
        inc offset_table
        jsr send_8bits

        ldx offset_table
        lda table,x
        jsr send_8bits

        inc byte_count
        jmp process_dial
!end:
        rts

// Takes whatever is in A and rols the bit out in the carry and based on that varry value
// sounds either 1950Hz for a ! and 2070Hz for a 0
// it will wait unul irq_counter changes (indicating 5ms has passed) to process the next bit      
send_8bits:
        ldy #8
!rol:
        rol
        bcc !_0+
        jsr hz1950
        jmp !_5ms+
!_0:
        jsr hz2070

!_5ms:    
        //wait for the next 5 ms interrupt to change the irq_counter. so we can send the next beep.
        //We could potentially do something else here
        //We just spin wait
        ldx irq_counter
        cpx irq_counter
        beq *-3
        dey
        bne !rol-
!end:
        rts

// --- Preamble IRQ (Handles 10ms interrupts) ---
dial_irq:
        pha
        tya
        pha
        txa
        pha

        inc irq_counter
        lda irq_counter
        cmp #60
        bne !+

        lda #1  
        sta timer_600ms_lapsed
        lda #$ff
        sta irq_counter
!:
        lda $dc0d

        pla
        tax
        pla
        tay
        pla
        jmp $ea31

// Configures SID1 to play to use triangle wave and no attack and no decay, sustain to full release no
setup_sid:
        lda #$0f // set volume full
        sta SID_VOLUME         
        
        lda #$00        // immediate attack, no decay
        sta SID1_ATTACK_DECAY   

        lda #$f0        // full sustain, no release
        sta SID1_SUSTAIN_RELEASE 
        
        lda #$11        //use triangle wave it's closes to sine
        sta SID1_OSCILATOR_TYPE
        rts

// set the frequency to 0Hz
hz0:
        pha

        lda #$00
        sta SID1_FREQ_LO
        lda #$00
        sta SID1_FREQ_HI

        pla
        rts

// set the frequency to 2070Hz
hz2070:
        pha

        lda #$ea
        sta SID1_FREQ_LO
        lda #$89
        sta SID1_FREQ_HI

        pla
        rts

// set the frequency to 1950Hz
hz1950:
        pha
        
        lda #$a0
        sta SID1_FREQ_LO
        lda #$81
        sta SID1_FREQ_HI

        pla
        rts

// --- Variables ---
irq_counter: .byte 0
timer_600ms_lapsed: .byte 0
byte_count: .byte 0
offset_table: .byte 0

table:
starttelegram:  //S in the string
.byte %01110010
.byte %00100010

stoptelegram:   //E in the string
.byte %01110100
.byte %00100001

canceltelegram: //C in the string
.byte %01110101
.byte %01010101

c0:
.byte %01110110
.byte %00000011

c1:
.byte %01110101
.byte %00000101

c2:
.byte %01110100
.byte %10001001

c3:
.byte %01110100
.byte %01010001

c4:
.byte %01110011
.byte %00000110

c5:
.byte %01110010
.byte %10001010

c6:
.byte %01110010
.byte %01010010

c7:
.byte %01110001
.byte %10001100

c8:
.byte %01110001
.byte %01010100

c9:
.byte %01110000
.byte %11011000