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
loop:
        lda timer_600ms_lapsed
        cmp #1  
        bne loop

        // 600ms clear preeamble send
        inc  $d020
        lda #$0
        jsr hz0

        // process the dial
        
stop_loop:
        jmp stop_loop

setup:
        sei
        lda #$00
        sta timer_600ms_lapsed

        jsr setup_sid
        jsr hz2070

        lda #<dial_irq
        sta IRQ_VECTOR
        lda #>dial_irq
        sta IRQ_VECTOR+1

        // Configure and start the 66ms CONTINUOUS timer
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

// --- Preamble IRQ (Handles 66ms interrupts) ---
dial_irq:
        inc irq_counter
        inc $0400
        lda irq_counter
        cmp #60
        bne !+

        lda #1  
        sta timer_600ms_lapsed
!:
        lda $dc0d
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
        lda #$ea
        sta $d400
        lda #$89
        sta $d401
        rts

// --- Variables ---
irq_counter: .byte 0
timer_600ms_lapsed: .byte 0

starttelegram:  //S in the string
.byte %01110010
.byte %00100010

stoptelegram:   //E in the string
.byte %01110100
.byte %00100001

canceltelegram:
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

number_to_dial:
.text "S1234520717666E"