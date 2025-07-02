BasicUpstart2(main)

        * = $1000         // Program starts here

// --- Constants ---
.const CIA1_TIMER_A_LOW  = $dc04
.const CIA1_TIMER_A_HIGH = $dc05
.const CIA1_ICR          = $dc0d
.const CIA1_CRA          = $dc0e
.const CIA1_CRB          = $dc0f
.const VIC_IRQ_ACK       = $d019
.const CIA_IRQ_ENABLE    = $d01a
.const IRQ_VECTOR        = $0314
.const VIC_IRQ_STAT      = $d019
.const VIC_IRQ_MASK      = $d01a

.const TIMER_66ms = $FFFF  // Actually ~66.5ms on PAL
.const TIMER_10ms = $268C  // ~10ms on PAL

// --- Main Program ---
// --- Main Program ---
main:
        sei                 // Disable interrupts during setup

        jsr setup_sid
        jsr hz2070          // Turn clear channel for transmit preamble

        // --- Take over the system IRQ ---
        // 1. Disable all interrupt sources first
        lda #$7f
        sta VIC_IRQ_STAT    // Acknowledge any pending VIC IRQs
        lda #%00000000      // Disable ALL VIC interrupt sources (Raster, etc.)
        sta VIC_IRQ_MASK    // VIC_IRQ_MASK is $d01a

        // 2. Point the CPU's IRQ vector to our handler
        lda #<preeamble_irq
        sta IRQ_VECTOR
        lda #>preeamble_irq
        sta IRQ_VECTOR+1

        // 3. Configure and start OUR desired interrupt source (CIA Timer A)
        lda #<TIMER_66ms
        sta CIA1_TIMER_A_LOW
        lda #>TIMER_66ms
        sta CIA1_TIMER_A_HIGH

        // 4. Set the CIA's INTERNAL interrupt mask to allow Timer A to fire
        //    Bit 7=1 means "Set Mask", Bit 0=1 means "for Timer A"
        lda #%10000001
        sta CIA1_ICR        // This is the correct way to enable the CIA source

        // 5. Start the timer
        lda #%00010001      // Load start value into Timer A and start timer
        sta CIA1_CRA
        
        // 6. Enable interrupts globally
        cli

loop:
        // Main loop polls the flag set by the IRQ
        lda timer_600ms_lapsed
        beq loop            // If flag is 0, keep looping

        // --- 600ms has passed ---
        sei                 // Disable interrupts to safely change timer
        
        jsr hz0             // Turn beep off
        inc $d020           // Flash border to show we got here
        
        // Reset the flag for the next time
        lda #0
        sta timer_600ms_lapsed
        sta irq_counter
        
        // Here you would re-configure for the 10ms chirps used for dial sequence
        lda #<TIMER_10ms
        sta CIA1_TIMER_A_LOW
        lda #>TIMER_10ms
        sta CIA1_TIMER_A_HIGH

        // Point to the next handler
        lda #<dial_number_irq
        sta IRQ_VECTOR
        lda #>dial_number_irq
        sta IRQ_VECTOR+1
        
        // Restart the timer
        lda #%00010001      // Load start value into Timer A and start timer
        sta CIA1_CRA
        
        cli                 // Re-enable interrupts
        jmp *               // Infinite loop to end the program

// --- Interrupt Service Routine ---
preeamble_irq:
        // --- Standard IRQ entry ---
        pha
        txa
        pha
        tya
        pha

        // --- Acknowledge the interrupt ---
        // It's crucial to do this first. We only care about CIA1 Timer A.
        lda CIA1_ICR

        // --- IRQ Logic ---
        // This is the only code that should be in the IRQ
        inc irq_counter
        lda irq_counter
        cmp #9              // Have 9 interrupts (approx 600ms) passed?
        bne !+         // If not, exit

        // If yes, set the flag for the main loop and reset the counter
        lda #1
        sta timer_600ms_lapsed
        
!:        // --- Standard IRQ exit ---
        pla
        tay
        pla
        tax
        pla
        rti

dial_number_irq:
         // --- Standard IRQ entry ---
        pha
        txa
        pha
        tya
        pha

        // --- Acknowledge the interrupt ---
        // It's crucial to do this first. We only care about CIA1 Timer A.
        lda CIA1_ICR


!:        // --- Standard IRQ exit ---
        pla
        tay
        pla
        tax
        pla
        rti

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

hz0:
        // Set frequency (example: 0 Hz) PAL   
        lda #$00       // Low byte of frequency
        sta $d400
        lda #$00       // High byte of frequency
        sta $d401
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

irq_counter:
.byte 0

timer_600ms_lapsed:
.byte 0

number_to_dial:
.text "S1234520717666E"