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
        hz0()
        // exit to basic
        jmp ($A002)

// set the variables to the startup values and clear the screen
// this is required if you want to run the problem again
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
        hz2070()

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

// wait for 600ms by checking to see if timer_600ms_lapsed is set to 1
ms600_preeamble:
        lda timer_600ms_lapsed
        cmp #1  
        bne ms600_preeamble
        rts

// parse the number from the number_to_convert string
// fetch the offset from the table that contains the binary bytes to send
// send each bit as the correct frequency for 10ms
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
        adc #6           // add 6 because the s,e,c are 6 bytes long and before the digits
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
        inc offset_table        // increment the offset table for the next byte
        send_8bits()

        ldx offset_table
        lda table,x
        send_8bits()

        inc byte_count
        jmp process_dial
!end:
        rts

// Macro saves a bit of time over jsr (we want to get as close to 10ms as possible
// eventhough those 8 are so cycles are not that important).
// Takes whatever is in A and rols the bit out in the carry and based on that varry value
// sounds either 1950Hz for a ! and 2070Hz for a 0
// it will wait unul irq_counter changes (indicating 10ms has passed) to process the next bit      
.macro send_8bits(){
        ldy #8
!rol:
        rol
        bcc !_0+
        hz1950()
        jmp !_10ms+
!_0:
        hz2070()

!_10ms:    
        //wait for the next 5 ms interrupt to change the irq_counter. so we can send the next beep.
        //We could potentially do something else here
        //We just spin wait
        ldx irq_counter
        cpx irq_counter
        beq *-3
        dey
        bne !rol-
!end:
}

// When IRQ is triggered we increment the irq_counter (the program spin locks on these variables)
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

// We use macros to set the frequency to 0Hz, 1950Hz, 2070Hz this way we can save a few cycles not jsring the function
// set the frequency to 0Hz
.macro hz0(){
        pha
        lda #$00
        sta SID1_FREQ_LO
        lda #$00
        sta SID1_FREQ_HI
        pla
}

// set the frequency to 2070Hz
.macro hz2070(){
        pha
        lda #$ea
        sta SID1_FREQ_LO
        lda #$89
        sta SID1_FREQ_HI
        pla
}

// set the frequency to 1950Hz
.macro hz1950(){
        pha
        lda #$a0
        sta SID1_FREQ_LO
        lda #$81
        sta SID1_FREQ_HI
        pla
}

// --- Variables ---
irq_counter: .byte 0
timer_600ms_lapsed: .byte 0
byte_count: .byte 0
offset_table: .byte 0

/*
  classic DTMF
  sid has max f of 4000hz, mapped from 0Hz to 40kHz in 16 bits
  4000/65536=0,061035Hz per step

            1209Hz=$4d60, 1336=$5581, 1477=$5e87, 1633=6883
697=$2c9b           1           2           3           A
770=$3147           4           5           6           B
852=$3687           7           8           9           C
941=$3c39           *           0           #           D
*/
dtmf:   .byte '0', $81, $55, $39, $3c      //0
        .byte '1', $60, $4d, $9b, $2c      //1
        .byte '2', $81, $55, $9b, $2c      //2
        .byte '3', $87, $5e, $9b, $2c      //3           
        .byte '4', $60, $4d, $47, $31      //4
        .byte '5', $81, $55, $47, $31      //5
        .byte '6', $87, $5e, $47, $31      //6
        .byte '7', $60, $4d, $87, $36      //7
        .byte '8', $81, $55, $87, $36      //8
        .byte '9', $87, $5e, $87, $36      //9
        .byte 'a', $83, $68, $9b, $2c      //A
        .byte 'b', $83, $68, $47, $31      //B
        .byte 'c', $83, $68, $87, $36      //C
        .byte 'd', $83, $68, $39, $3c      //D
        .byte '*', $60, $4d, $39, $3c      //*
        .byte '#', $60, $4d, $87, $5e      //#

/*
CCITT 5

DIGIT	FREQ1	SID_DEC	                SID_HEX	FREQ2	SID_DEC	                SID HEX	mS	
0	    1300	21299,2545260916	    5333	1700	27852,8713033505	    6CCC	55	
1	    700	    11468,8293602032	    2CCC	900	    14745,6377488326	    3999	55	
2   	700	    11468,8293602032	    2CCC	1100	18022,4461374621	    4666	55	
3	    900	    14745,6377488326	    3999	1100	18022,4461374621	    4666	55	
4	    700	    11468,8293602032	    2CCC	1300	21299,2545260916	    5333	55
5	    900	    14745,6377488326	    3999	1300	21299,2545260916	    5333	55
6	    1100	18022,4461374621	    4666	1300	21299,2545260916	    5333	55
7	    700	    11468,8293602032	    2CCC	1500	24576,0629147211	    6000	55
8	    900	    14745,6377488326	    3999	1500	24576,0629147211	    6000	55
9	    1100	18022,4461374621	    4666	1500	24576,0629147211	    6000	55
Code11	700	    11468,8293602032	    2CCC	1700	27852,8713033505	    6CCC	55
Code12	900	    14745,6377488326	    3999	1700	27852,8713033505	    6CCC	55
KP1	    1100	18022,4461374621	    4666	1700	27852,8713033505	    6CCC	100
KP2	    1300	21299,2545260916	    5333	1700	27852,8713033505	    6CCC	100
ST	    1500	24576,0629147211	    6000	1700	27852,8713033505	    6CCC	55
*/


ccitt5: .byte '0', $33, $53, $cc, $6c  //0
        .byte '1', $cc, $2c, $99, $39  //1
        
        .byte '2', $cc, $2c, $66, $46  //2
        .byte '3', $99, $39, $66, $46  //3
        .byte '4', $cc, $2c, $33, $53  //4
        .byte '5', $99, $39, $33, $53  //5
        .byte '6', $66, $46, $33, $53  //6
        .byte '7', $cc, $2c, $00, $60  //7
        .byte '8', $99, $39, $00, $60  //8
        .byte '9', $66, $46, $00, $60  //9
        .byte 'k', $66, $46, $cc, $6c  //kp1
        .byte 'K', $00, $60, $cc, $6c  //kp2
        .byte 'c', $cc, $2c, $cc, $6c  //code11
        .byte 'C', $99, $39, $cc, $6c  //code12