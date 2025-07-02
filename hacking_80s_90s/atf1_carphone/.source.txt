BasicUpstart2(main)

        * = $1000

// --- Constants ---
.const CIA1_TIMER_A_LOW  = $dc04
.const CIA1_TIMER_A_HIGH = $dc05
.const CIA1_ICR          = $dc0d
.const CIA1_CRA          = $dc0e
.const IRQ_VECTOR        = $0314
.const TIMER_10ms = $268C
.const TIMER_START_CONTINUOUS = %00010001

// --- Main Program ---
main:
        jsr setup
       
        // Wait until the 600ms preamble is finished
        lda timer_600ms_lapsed

        // send preeamble 600ms 2050hz
        jsr ms600_preeamble

        // process the dial
        jsr process_dial

stop_loop:
        jsr hz0
        jmp ($A002)

setup:
        sei
        lda #$00
        sta timer_600ms_lapsed
        sta byte_count
        sta offset_table
        sta irq_counter

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

ms600_preeamble:
        lda timer_600ms_lapsed
        cmp #1  
        bne ms600_preeamble

        // 600ms clear preeamble send
        lda #$0
        jsr hz0
        rts

process_dial:
        ldx byte_count
        // read letter from string and check if we are at the end of the string
        lda number_to_convert,x
        sta $0400,x

        sec
        sbc #'0'         // Normalize ASCII digit chars to 0-9
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

        //process bits
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
        //wait 5 ms (set by irq) to next bit
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
        inc $d020

        pla
        tax
        pla
        tay
        pla
        jmp $ea31

// --- SID Routines ---
setup_sid:
        lda #$0f
        sta $d418
        lda #$00
        sta $d405
        lda #$f0
        sta $d406
        lda #$11
        sta $d404
        rts

hz0:
        lda #$00
        sta $d400
        lda #$00
        sta $d401
        rts

hz2070:
        pha

        lda #$ea
        sta $d400
        lda #$89
        sta $d401

        pla
        rts

hz1950:
        pha
        
        lda #$a0
        sta $d400
        lda #$81
        sta $d401

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

number_to_convert:
.text "s1234520717666e"
.text "$"