BasicUpstart2(main)

        * = $1000         // Program starts here

// --- Timer value (10 ms) ---
.const timer_value = $268C         // ~9852 decimal

main:
        sei                 // Disable interrupts

        lda #$7f
        sta $dc0d           // Disable CIA #1 interrupts
        sta $dd0d           // Disable CIA #2 interrupts

        lda $dc0d           // Acknowledge any pending CIA #1 interrupt
        lda $dd0d           // Acknowledge any pending CIA #2 interrupt

        lda #$01
        sta $d01a           // Enable IRQs from CIA #1

        lda #<irq_handler
        sta $0314           // Set IRQ vector low byte
        lda #>irq_handler
        sta $0315           // Set IRQ vector high byte

        lda #<timer_value
        sta $dc04           // CIA1 Timer A low byte
        lda #>timer_value
        sta $dc05           // CIA1 Timer A high byte

        lda #%00010001      // Start Timer A, continuous mode, count CPU cycles
        sta $dc0e

        lda #%00000001      // Enable Timer A interrupt
        sta $dc0d

        cli                 // Enable interrupts

loop:
        jmp loop            // Endless loop

// --- IRQ handler ---
irq_handler:
        lda $dc0d           // Acknowledge CIA1 interrupt

        // Do something here if needed
        inc $d020           // Flash border color for debug

        jmp $ea31           // Jump to standard IRQ handler (restore registers, etc.)

starttelegram:
.byte %01110010
.byte %00100010

stoptelegram:
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
