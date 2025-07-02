BasicUpstart2(main)

        * = $1000         // Program starts here

// --- Timer value (10 ms) ---
.const ms10_timer = $268C         // ~9852 decimal
.const ms600_timer = $FFFF         // and we need 9 of these

main:
        jsr hz2070
        jsr setup_sid
        // jsr setup_timer
        // jsr setup_timer_600ms
        // jsr start_timer
        // jsr setup_timer_10ms
        // jsr start_timer
        jmp *

setup_timer:
        sei                 // Disable interrupts

        lda #$7f
        sta $dc0d           // Disable CIA #1 interrupts
        sta $dd0d           // Disable CIA #2 interrupts

        lda $dc0d           // Acknowledge any pending CIA #1 interrupt
        lda $dd0d           // Acknowledge any pending CIA #2 interrupt

        lda #$01
        sta $d01a           // Enable IRQs from CIA #1

        rts

setup_timer_10ms:
        lda #<ms10_timer
        sta $dc04           // CIA1 Timer A low byte
        lda #>ms10_timer
        sta $dc05           // CIA1 Timer A high byte

        lda #<ms10irq
        sta $0314           // Set IRQ vector low byte
        lda #>ms10irq
        sta $0315           // Set IRQ vector high byte
        rts

setup_timer_600ms:
        lda #<ms600irq
        sta $0314           // Set IRQ vector low byte
        lda #>ms600irq
        sta $0315           // Set IRQ vector high byte
        rts

start_timer:
        lda #%00010001      // Start Timer A, continuous mode, count CPU cycles
        sta $dc0e

        lda #%00000001      // Enable Timer A interrupt
        sta $dc0d

        cli                 // Enable interrupts
        rts

setup_sid:
        // Set volume
        lda #$0f       // max volume
        sta $d418      // store volume

        lda #$00      //attack decay 0
        sta $d405      

        lda #$f0      //sustain release 0
        sta $d406
        
        // Set control register for pulse wave and gate
        lda #$11       // triangle
        sta $d404      // Control register for voice 1

        rts

hz1950:
        // Set frequency (example: 1950 Hz) PAL
        lda #$a0       // Low byte of frequency 
        sta $d400
        lda #$81       // High byte of frequency 
        sta $d401
        rts

hz2070:
        // Set frequency (example: 2070 Hz) PAL
        lda #$ea       // Low byte of frequency
        sta $d400
        lda #$89       // High byte of frequency 
        sta $d401
        rts

// --- IRQ handler 600ms ---
// we acknowledge the interrupt 9 times, that would lead to a 598 ms delay (close enough)
ms600irq:
        irq_counter:
                .byte 0

        irq_handler:
                pha
                txa
                pha
                tya
                pha

                lda $dc0d           // Acknowledge Timer A

                inc irq_counter
                lda irq_counter
                cmp #9
                bne skip_action

                // --- 600 ms elapsed clear counter for next time ---
                lda #0
                sta irq_counter

                // Set frequency (example: 0 Hz)
                lda #$00       // Low byte of frequency 
                sta $D400
                lda #$00       // High byte of frequency 
                sta $D401

        skip_action:
                pla
                tay
                pla
                tax
                pla
                rti

// --- IRQ handler 10ms ---
ms10irq:
        lda $dc0d           // Acknowledge CIA1 interrupt

        // Do something here if needed
        inc $d020           // Flash border color for debug

        rti

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