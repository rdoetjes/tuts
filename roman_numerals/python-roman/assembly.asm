.const screenPosStartManual = $0400 + (40*25) - 115

    BasicUpstart2(main)
    *=$1000

main:
    // intialize vic
    jsr $FF81

	// set to 25 line text mode and turn on the screen
	lda #$1B
	sta $D011

	// disable SHIFT-Commodore
	lda #$80
	sta $0291

	// set screen memory ($0400) and charset bitmap offset ($2000)
	lda #$18
	sta $D018

	// set border color
	lda #$00
	sta $D020

	// set background color
	lda #$00
	sta $D021

	// set sprite multicolors
	lda #$02
	sta $d025
	lda #$06
	sta $d026

	// positioning sprites
	lda #$58
	sta $d000	// #0. sprite X low .byte
	lda #$86
	sta $d001	// #0. sprite Y
	lda #$F9
	sta $d002	// #1. sprite X low .byte
	lda #$67
	sta $d003	// #1. sprite Y
	lda #$58
	sta $d004	// #2. sprite X low .byte
	lda #$67
	sta $d005	// #2. sprite Y
	lda #$F9
	sta $d006	// #3. sprite X low .byte
	lda #$86
	sta $d007	// #3. sprite Y
	lda #$49
	sta $d008	// #4. sprite X low .byte
	lda #$77
	sta $d009	// #4. sprite Y
	lda #$EA
	sta $d00A	// #5. sprite X low .byte
	lda #$77
	sta $d00B	// #5. sprite Y
	lda #$06
	sta $d00C	// #6. sprite X low .byte
	lda #$77
	sta $d00D	// #6. sprite Y
	lda #$65
	sta $d00E	// #7. sprite X low .byte
	lda #$77
	sta $d00F	// #7. sprite Y

	// X coordinate high bits
	lda #$40
	sta $d010

	// expand sprites
	lda #$00
	sta $d01d
	lda #$00
	sta $d017

	// set multicolor flags
	lda #$00
	sta $d01c

	// set screen-sprite priority flags
	lda #$00
	sta $d01b

	// set sprite pointers
	lda #$28
	sta $07F8
	lda #$29
	sta $07F9
	lda #$2A
	sta $07FA
	lda #$2B
	sta $07FB
	lda #$2C
	sta $07FC
	lda #$2D
	sta $07FD
	lda #$2E
	sta $07FE
	lda #$2F
	sta $07FF

	// turn on sprites
	lda #$FF
	sta $d015

    //main loop
	jsr printManual

loop:
    ldx #1      //read joystick 1
    jsr drwButton
    ldx #1      //read joystick 1
    jsr drwStick

    ldx #2
    jsr drwButton
    ldx #2
    jsr drwStick

wait:
    lda $d012
    cmp $ff
    bne wait
jsr drwAction

!cont:
	ldy #32
	jsr setAction
	jmp loop

    //read the button of the joystick X holds the joystick, y gets the color value
readButton:
    cpx #1
    beq !joystick1+
    lda $dc00       //read joystick 2
    jmp !readButton+

    !joystick1:
        lda $dc01       //read joystick 1

    !readButton:
        and #$10        //check button press
        beq !redColor+
        ldy #05
        rts

    !redColor:
    	txa
    	pha
    	ldy #4
    	jsr setAction
    	pla
    	tax
    	ldy #02
        rts


//changes the color of the approtiate up, dwn, rght, lft sprite, X holds joy stick 1 or 2 and Y holds color
drwStick:
        stx joynr           //save joystick number
        cpx #1
        beq !joystick1+
        lda $dc00           //read joystick 2
        jmp jProc

    !joystick1:
        lda $dc01

    jProc:
        sta jvalue           //store joystick register value
        ldy #0               //y is the bit counter

    !loop:
        ldx joynr
        cpx #1
        bne !js1+
        lda js2, y          //get the sprite pointer for this y index
		jmp !next+
    !js1:
        lda js1, y          //get the sprite pointer for this y index

	!next:
        sta sprite1+1     //self modifying code for the sprite index
        sta sprite2+1     //self modifying code for the sprite index
        lda jvalue
        and bits, y          //check up
        beq !up+
        ldx #05
    sprite1:
        stx $d000           //00 is overwritten by sprite index
        jmp !next+
    !up:
		jsr setAction
        ldx #02
    sprite2:
        stx $d000           //00 is overwritten by sprite index

    !next:
        iny
        cpy #4
        bne !loop-
    rts

drwAction:
    ldx #39
	ldy #05

!loop:
    lda actions, x
	clc
	adc #112
    sta $400+(40*15), x
	lda #05
	sta $d800+(40*15), x
    dex
    bne !loop-

    rts

//STore the action in register Y into the actions list
setAction:
	ldx #0

!loop:				//move charters on the line to the left
	lda actions+1,x
    sta actions,x
    inx
    cpx #39
    bne !loop-

    sty actions+39	//store newest action on the end

	rts

    //draws the buttons on the screen, X holds joy stick 1 or 2 and Y holds color
	//TODO: update with slef modifying code to reduce size
drwButton:
	jsr readButton

    cpx #1
    beq !button1+

    !button2:
        //draw the 4 chars that make up the circle
        lda #$4f
        sta $0400 + 394

        lda #$50
        sta $0400 + 395

        lda #$4c
        sta $0400 + 434

        lda #$7a
        sta $0400 + 435

        //store the color in colorram
        sty $d800 + 394
        sty $d800 + 395
        sty $d800 + 434
        sty $d800 + 435

        jmp !end+

    !button1:
        lda #$4f
        sta $0400 + 374

        lda #$50
        sta $0400 + 375

        lda #$4c
        sta $0400 + 414

        lda #$7a
        sta $0400 + 415

        //store the color in colorram
        sty $d800 + 374
        sty $d800 + 375
        sty $d800 + 414
        sty $d800 + 415

    !end:
        rts

printManual:
        //print the manual text
        ldx #0
    !:
        lda mesg, x
        cmp #$ff
        beq !+
        inx
        sta screenPosStartManual, x
        bne !-

    !:
		lda #$05
		ldx #$ff
	!t:						//color lower 255 .bytes green
		sta $db00, x
		dex
		bne !t-
        rts

 //variables
joynr:   	.byte $00
jvalue:  	.byte $00
bits:    	.byte $01, $02, $04, $08
js2:     	.byte $29, $27, $2b, $2e
js1:     	.byte $28, $2a, $2c, $2d
actions: 	.fill 40, 32
mesg:		.text @"joystick tester phonax (c)2021\$ff"

// Sprite bitmaps 8 x 64 .bytes
*=$0A00
// sprite #0
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $3F, $FF, $F0, $3F, $FF, $F0, $1F, $FF, $E0
	.byte $0F, $FF, $C0, $07, $FF, $80, $03, $FF, $00, $01, $FE, $00, $00, $FC, $00, $00, $78, $00, $00, $30, $00
	.byte 0

// sprite #1
	.byte $00, $30, $00, $00, $78, $00, $00, $FC, $00, $01, $FE, $00, $03, $FF, $00, $07, $FF, $80, $0F, $FF, $C0
	.byte $1F, $FF, $E0, $3F, $FF, $F0, $3F, $FF, $F0, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte 0

// sprite #2
	.byte $00, $30, $00, $00, $78, $00, $00, $FC, $00, $01, $FE, $00, $03, $FF, $00, $07, $FF, $80, $0F, $FF, $C0
	.byte $1F, $FF, $E0, $3F, $FF, $F0, $3F, $FF, $F0, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte 0

// sprite #3
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00
	.byte $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $03, $FF, $00, $3F, $FF, $F0, $3F, $FF, $F0, $1F, $FF, $E0
	.byte $0F, $FF, $C0, $07, $FF, $80, $03, $FF, $00, $01, $FE, $00, $00, $FC, $00, $00, $78, $00, $00, $30, $00
	.byte 0

// sprite #4
	.byte $00, $00, $00, $00, $C0, $00, $01, $C0, $00, $03, $C0, $00, $07, $C0, $00, $0F, $FF, $F8, $1F, $FF, $F8
	.byte $3F, $FF, $F8, $7F, $FF, $F8, $FF, $FF, $F8, $FF, $FF, $F8, $7F, $FF, $F8, $3F, $FF, $F8, $1F, $FF, $F8
	.byte $0F, $FF, $F8, $07, $C0, $00, $03, $C0, $00, $01, $C0, $00, $00, $C0, $00, $00, $00, $00, $00, $00, $00
	.byte 0

// sprite #5
	.byte $00, $00, $00, $00, $C0, $00, $01, $C0, $00, $03, $C0, $00, $07, $C0, $00, $0F, $FF, $F8, $1F, $FF, $F8
	.byte $3F, $FF, $F8, $7F, $FF, $F8, $FF, $FF, $F8, $FF, $FF, $F8, $7F, $FF, $F8, $3F, $FF, $F8, $1F, $FF, $F8
	.byte $0F, $FF, $F8, $07, $C0, $00, $03, $C0, $00, $01, $C0, $00, $00, $C0, $00, $00, $00, $00, $00, $00, $00
	.byte 0

// sprite #6
	.byte $00, $00, $00, $00, $03, $00, $00, $03, $80, $00, $03, $C0, $00, $03, $E0, $1F, $FF, $F0, $1F, $FF, $F8
	.byte $1F, $FF, $FC, $1F, $FF, $FE, $1F, $FF, $FF, $1F, $FF, $FF, $1F, $FF, $FE, $1F, $FF, $FC, $1F, $FF, $F8
	.byte $1F, $FF, $F0, $00, $03, $E0, $00, $03, $C0, $00, $03, $80, $00, $03, $00, $00, $00, $00, $00, $00, $00
	.byte 0

// sprite #7
	.byte $00, $00, $00, $00, $03, $00, $00, $03, $80, $00, $03, $C0, $00, $03, $E0, $1F, $FF, $F0, $1F, $FF, $F8
	.byte $1F, $FF, $FC, $1F, $FF, $FE, $1F, $FF, $FF, $1F, $FF, $FF, $1F, $FF, $FE, $1F, $FF, $FC, $1F, $FF, $F8
	.byte $1F, $FF, $F0, $00, $03, $E0, $00, $03, $C0, $00, $03, $80, $00, $03, $00, $00, $00, $00, $00, $00, $00
	.byte 0


// Character bitmap definitions 2k
*=$2000
	.byte	$3C, $66, $6E, $6E, $60, $62, $3C, $00
	.byte	$18, $3C, $66, $7E, $66, $66, $66, $00
	.byte	$7C, $66, $66, $7C, $66, $66, $7C, $00
	.byte	$3C, $66, $60, $60, $60, $66, $3C, $00
	.byte	$78, $6C, $66, $66, $66, $6C, $78, $00
	.byte	$7E, $60, $60, $78, $60, $60, $7E, $00
	.byte	$7E, $60, $60, $78, $60, $60, $60, $00
	.byte	$3C, $66, $60, $6E, $66, $66, $3C, $00
	.byte	$66, $66, $66, $7E, $66, $66, $66, $00
	.byte	$3C, $18, $18, $18, $18, $18, $3C, $00
	.byte	$1E, $0C, $0C, $0C, $0C, $6C, $38, $00
	.byte	$66, $6C, $78, $70, $78, $6C, $66, $00
	.byte	$60, $60, $60, $60, $60, $60, $7E, $00
	.byte	$63, $77, $7F, $6B, $63, $63, $63, $00
	.byte	$66, $76, $7E, $7E, $6E, $66, $66, $00
	.byte	$3C, $66, $66, $66, $66, $66, $3C, $00
	.byte	$7C, $66, $66, $7C, $60, $60, $60, $00
	.byte	$3C, $66, $66, $66, $66, $3C, $0E, $00
	.byte	$7C, $66, $66, $7C, $78, $6C, $66, $00
	.byte	$3C, $66, $60, $3C, $06, $66, $3C, $00
	.byte	$7E, $18, $18, $18, $18, $18, $18, $00
	.byte	$66, $66, $66, $66, $66, $66, $3C, $00
	.byte	$66, $66, $66, $66, $66, $3C, $18, $00
	.byte	$63, $63, $63, $6B, $7F, $77, $63, $00
	.byte	$66, $66, $3C, $18, $3C, $66, $66, $00
	.byte	$66, $66, $66, $3C, $18, $18, $18, $00
	.byte	$7E, $06, $0C, $18, $30, $60, $7E, $00
	.byte	$3C, $30, $30, $30, $30, $30, $3C, $00
	.byte	$0C, $12, $30, $7C, $30, $62, $FC, $00
	.byte	$3C, $0C, $0C, $0C, $0C, $0C, $3C, $00
	.byte	$00, $18, $3C, $7E, $18, $18, $18, $18
	.byte	$00, $10, $30, $7F, $7F, $30, $10, $00
	.byte	$00, $00, $00, $00, $00, $00, $00, $00
	.byte	$18, $18, $18, $18, $00, $00, $18, $00
	.byte	$66, $66, $66, $00, $00, $00, $00, $00
	.byte	$66, $66, $FF, $66, $FF, $66, $66, $00
	.byte	$18, $3E, $60, $3C, $06, $7C, $18, $00
	.byte	$62, $66, $0C, $18, $30, $66, $46, $00
	.byte	$3C, $66, $3C, $38, $67, $66, $3F, $00
	.byte	$06, $0C, $18, $00, $00, $00, $00, $00
	.byte	$0C, $18, $30, $30, $30, $18, $0C, $00
	.byte	$30, $18, $0C, $0C, $0C, $18, $30, $00
	.byte	$00, $66, $3C, $FF, $3C, $66, $00, $00
	.byte	$00, $18, $18, $7E, $18, $18, $00, $00
	.byte	$00, $00, $00, $00, $00, $18, $18, $30
	.byte	$00, $00, $00, $7E, $00, $00, $00, $00
	.byte	$00, $00, $00, $00, $00, $18, $18, $00
	.byte	$00, $03, $06, $0C, $18, $30, $60, $00
	.byte	$3C, $66, $6E, $76, $66, $66, $3C, $00
	.byte	$18, $18, $38, $18, $18, $18, $7E, $00
	.byte	$3C, $66, $06, $0C, $30, $60, $7E, $00
	.byte	$3C, $66, $06, $1C, $06, $66, $3C, $00
	.byte	$06, $0E, $1E, $66, $7F, $06, $06, $00
	.byte	$7E, $60, $7C, $06, $06, $66, $3C, $00
	.byte	$3C, $66, $60, $7C, $66, $66, $3C, $00
	.byte	$7E, $66, $0C, $18, $18, $18, $18, $00
	.byte	$3C, $66, $66, $3C, $66, $66, $3C, $00
	.byte	$3C, $66, $66, $3E, $06, $66, $3C, $00
	.byte	$00, $00, $18, $00, $00, $18, $00, $00
	.byte	$00, $00, $18, $00, $00, $18, $18, $30
	.byte	$0E, $18, $30, $60, $30, $18, $0E, $00
	.byte	$00, $00, $7E, $00, $7E, $00, $00, $00
	.byte	$70, $18, $0C, $06, $0C, $18, $70, $00
	.byte	$3C, $66, $06, $0C, $18, $00, $18, $00
	.byte	$00, $00, $00, $FF, $FF, $00, $00, $00
	.byte	$08, $1C, $3E, $7F, $7F, $1C, $3E, $00
	.byte	$18, $18, $18, $18, $18, $18, $18, $18
	.byte	$00, $00, $00, $FF, $FF, $00, $00, $00
	.byte	$00, $00, $FF, $FF, $00, $00, $00, $00
	.byte	$00, $FF, $FF, $00, $00, $00, $00, $00
	.byte	$00, $00, $00, $00, $FF, $FF, $00, $00
	.byte	$30, $30, $30, $30, $30, $30, $30, $30
	.byte	$0C, $0C, $0C, $0C, $0C, $0C, $0C, $0C
	.byte	$00, $00, $00, $E0, $F0, $38, $18, $18
	.byte	$18, $18, $1C, $0F, $07, $00, $00, $00
	.byte	$18, $18, $38, $F0, $E0, $00, $00, $00
	.byte	$FF, $FF, $7F, $7F, $3F, $1F, $0F, $03
	.byte	$C0, $E0, $70, $38, $1C, $0E, $07, $03
	.byte	$03, $07, $0E, $1C, $38, $70, $E0, $C0
	.byte	$03, $0F, $1F, $3F, $7F, $7F, $FF, $FF
	.byte	$C0, $F0, $F8, $FC, $FE, $FE, $FF, $FF
	.byte	$00, $3C, $7E, $7E, $7E, $7E, $3C, $00
	.byte	$3C, $7E, $E3, $EF, $E7, $EF, $6E, $3C
	.byte	$36, $7F, $7F, $7F, $3E, $1C, $08, $00
	.byte	$60, $60, $60, $60, $60, $60, $60, $60
	.byte	$00, $00, $00, $07, $0F, $1C, $18, $18
	.byte	$C3, $E7, $7E, $3C, $3C, $7E, $E7, $C3
	.byte	$00, $3C, $7E, $66, $66, $7E, $3C, $00
	.byte	$18, $18, $66, $66, $18, $18, $3C, $00
	.byte	$06, $06, $06, $06, $06, $06, $06, $06
	.byte	$08, $1C, $3E, $7F, $3E, $1C, $08, $00
	.byte	$18, $18, $18, $FF, $FF, $18, $18, $18
	.byte	$C0, $C0, $30, $30, $C0, $C0, $30, $30
	.byte	$18, $18, $18, $18, $18, $18, $18, $18
	.byte	$00, $00, $03, $3E, $76, $36, $36, $00
	.byte	$FF, $7F, $3F, $1F, $0F, $07, $03, $01
	.byte	$00, $00, $00, $00, $00, $00, $00, $00
	.byte	$F0, $F0, $F0, $F0, $F0, $F0, $F0, $F0
	.byte	$00, $00, $00, $00, $FF, $FF, $FF, $FF
	.byte	$FF, $00, $00, $00, $00, $00, $00, $00
	.byte	$00, $00, $00, $00, $00, $00, $00, $FF
	.byte	$08, $0C, $0E, $FF, $FF, $0E, $0C, $08
	.byte	$CC, $CC, $33, $33, $CC, $CC, $33, $33
	.byte	$03, $03, $03, $03, $03, $03, $03, $03
	.byte	$00, $00, $00, $00, $CC, $CC, $33, $33
	.byte	$FF, $FE, $FC, $F8, $F0, $E0, $C0, $80
	.byte	$03, $03, $03, $03, $03, $03, $03, $03
	.byte	$18, $18, $18, $1F, $1F, $18, $18, $18
	.byte	$00, $00, $00, $00, $0F, $0F, $0F, $0F
	.byte	$18, $18, $18, $1F, $1F, $00, $00, $00
	.byte	$00, $00, $00, $F8, $F8, $18, $18, $18
	.byte	$18, $3C, $7E, $FF, $18, $18, $18, $18
	.byte	$18, $3C, $7E, $FF, $18, $18, $18, $18
	.byte	$18, $18, $18, $18, $FF, $7E, $3C, $18
	.byte	$10, $30, $70, $FF, $FF, $70, $30, $10
	.byte	$08, $0C, $0E, $FF, $FF, $0E, $0C, $08
	.byte	$3C, $7E, $E3, $EF, $E7, $EF, $6E, $3C
	.byte	$E0, $E0, $E0, $E0, $E0, $E0, $E0, $E0
	.byte	$07, $07, $07, $07, $07, $07, $07, $07
	.byte	$FF, $FF, $00, $00, $00, $00, $00, $00
	.byte	$10, $30, $70, $FF, $FF, $70, $30, $10
	.byte	$18, $18, $18, $18, $FF, $7E, $3C, $18
	.byte	$FF, $FF, $FE, $FE, $FC, $F8, $F0, $C0
	.byte	$00, $00, $00, $00, $F0, $F0, $F0, $F0
	.byte	$0F, $0F, $0F, $0F, $00, $00, $00, $00
	.byte	$18, $18, $18, $F8, $F8, $00, $00, $00
	.byte	$F0, $F0, $F0, $F0, $00, $00, $00, $00
	.byte	$F0, $F0, $F0, $F0, $0F, $0F, $0F, $0F
	.byte	$C3, $99, $91, $91, $9F, $99, $C3, $FF
	.byte	$E7, $C3, $99, $81, $99, $99, $99, $FF
	.byte	$83, $99, $99, $83, $99, $99, $83, $FF
	.byte	$C3, $99, $9F, $9F, $9F, $99, $C3, $FF
	.byte	$87, $93, $99, $99, $99, $93, $87, $FF
	.byte	$81, $9F, $9F, $87, $9F, $9F, $81, $FF
	.byte	$81, $9F, $9F, $87, $9F, $9F, $9F, $FF
	.byte	$C3, $99, $9F, $91, $99, $99, $C3, $FF
	.byte	$99, $99, $99, $81, $99, $99, $99, $FF
	.byte	$C3, $E7, $E7, $E7, $E7, $E7, $C3, $FF
	.byte	$E1, $F3, $F3, $F3, $F3, $93, $C7, $FF
	.byte	$99, $93, $87, $8F, $87, $93, $99, $FF
	.byte	$9F, $9F, $9F, $9F, $9F, $9F, $81, $FF
	.byte	$9C, $88, $80, $94, $9C, $9C, $9C, $FF
	.byte	$99, $89, $81, $81, $91, $99, $99, $FF
	.byte	$C3, $99, $99, $99, $99, $99, $C3, $FF
	.byte	$00, $00, $00, $00, $00, $00, $00, $00
	.byte	$C3, $99, $99, $99, $99, $C3, $F1, $FF
	.byte	$83, $99, $99, $83, $87, $93, $99, $FF
	.byte	$C3, $99, $9F, $C3, $F9, $99, $C3, $FF
	.byte	$81, $E7, $E7, $E7, $E7, $E7, $E7, $FF
	.byte	$99, $99, $99, $99, $99, $99, $C3, $FF
	.byte	$99, $99, $99, $99, $99, $C3, $E7, $FF
	.byte	$9C, $9C, $9C, $94, $80, $88, $9C, $FF
	.byte	$99, $99, $C3, $E7, $C3, $99, $99, $FF
	.byte	$99, $99, $99, $C3, $E7, $E7, $E7, $FF
	.byte	$81, $F9, $F3, $E7, $CF, $9F, $81, $FF
	.byte	$C3, $CF, $CF, $CF, $CF, $CF, $C3, $FF
	.byte	$F3, $ED, $CF, $83, $CF, $9D, $03, $FF
	.byte	$C3, $F3, $F3, $F3, $F3, $F3, $C3, $FF
	.byte	$FF, $E7, $C3, $81, $E7, $E7, $E7, $E7
	.byte	$FF, $EF, $CF, $80, $80, $CF, $EF, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $FF, $FF, $FF
	.byte	$E7, $E7, $E7, $E7, $FF, $FF, $E7, $FF
	.byte	$99, $99, $99, $FF, $FF, $FF, $FF, $FF
	.byte	$99, $99, $00, $99, $00, $99, $99, $FF
	.byte	$E7, $C1, $9F, $C3, $F9, $83, $E7, $FF
	.byte	$9D, $99, $F3, $E7, $CF, $99, $B9, $FF
	.byte	$C3, $99, $C3, $C7, $98, $99, $C0, $FF
	.byte	$F9, $F3, $E7, $FF, $FF, $FF, $FF, $FF
	.byte	$F3, $E7, $CF, $CF, $CF, $E7, $F3, $FF
	.byte	$CF, $E7, $F3, $F3, $F3, $E7, $CF, $FF
	.byte	$FF, $99, $C3, $00, $C3, $99, $FF, $FF
	.byte	$FF, $E7, $E7, $81, $E7, $E7, $FF, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $E7, $E7, $CF
	.byte	$FF, $FF, $FF, $81, $FF, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $E7, $E7, $FF
	.byte	$FF, $FC, $F9, $F3, $E7, $CF, $9F, $FF
	.byte	$C3, $99, $91, $89, $99, $99, $C3, $FF
	.byte	$E7, $E7, $C7, $E7, $E7, $E7, $81, $FF
	.byte	$C3, $99, $F9, $F3, $CF, $9F, $81, $FF
	.byte	$C3, $99, $F9, $E3, $F9, $99, $C3, $FF
	.byte	$F9, $F1, $E1, $99, $80, $F9, $F9, $FF
	.byte	$81, $9F, $83, $F9, $F9, $99, $C3, $FF
	.byte	$C3, $99, $9F, $83, $99, $99, $C3, $FF
	.byte	$81, $99, $F3, $E7, $E7, $E7, $E7, $FF
	.byte	$C3, $99, $99, $C3, $99, $99, $C3, $FF
	.byte	$C3, $99, $99, $C1, $F9, $99, $C3, $FF
	.byte	$FF, $FF, $E7, $FF, $FF, $E7, $FF, $FF
	.byte	$FF, $FF, $E7, $FF, $FF, $E7, $E7, $CF
	.byte	$F1, $E7, $CF, $9F, $CF, $E7, $F1, $FF
	.byte	$FF, $FF, $81, $FF, $81, $FF, $FF, $FF
	.byte	$8F, $E7, $F3, $F9, $F3, $E7, $8F, $FF
	.byte	$C3, $99, $F9, $F3, $E7, $FF, $E7, $FF
	.byte	$FF, $FF, $FF, $00, $00, $FF, $FF, $FF
	.byte	$F7, $E3, $C1, $80, $80, $E3, $C1, $FF
	.byte	$E7, $E7, $E7, $E7, $E7, $E7, $E7, $E7
	.byte	$FF, $FF, $FF, $00, $00, $FF, $FF, $FF
	.byte	$FF, $FF, $00, $00, $FF, $FF, $FF, $FF
	.byte	$FF, $00, $00, $FF, $FF, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $FF, $00, $00, $FF, $FF
	.byte	$CF, $CF, $CF, $CF, $CF, $CF, $CF, $CF
	.byte	$F3, $F3, $F3, $F3, $F3, $F3, $F3, $F3
	.byte	$FF, $FF, $FF, $1F, $0F, $C7, $E7, $E7
	.byte	$E7, $E7, $E3, $F0, $F8, $FF, $FF, $FF
	.byte	$E7, $E7, $C7, $0F, $1F, $FF, $FF, $FF
	.byte	$3F, $3F, $3F, $3F, $3F, $3F, $00, $00
	.byte	$3F, $1F, $8F, $C7, $E3, $F1, $F8, $FC
	.byte	$FC, $F8, $F1, $E3, $C7, $8F, $1F, $3F
	.byte	$00, $00, $3F, $3F, $3F, $3F, $3F, $3F
	.byte	$00, $00, $FC, $FC, $FC, $FC, $FC, $FC
	.byte	$FF, $C3, $81, $81, $81, $81, $C3, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $00, $00, $FF
	.byte	$C9, $80, $80, $80, $C1, $E3, $F7, $FF
	.byte	$9F, $9F, $9F, $9F, $9F, $9F, $9F, $9F
	.byte	$FF, $FF, $FF, $F8, $F0, $E3, $E7, $E7
	.byte	$3C, $18, $81, $C3, $C3, $81, $18, $3C
	.byte	$FF, $C3, $81, $99, $99, $81, $C3, $FF
	.byte	$E7, $E7, $99, $99, $E7, $E7, $C3, $FF
	.byte	$F9, $F9, $F9, $F9, $F9, $F9, $F9, $F9
	.byte	$F7, $E3, $C1, $80, $C1, $E3, $F7, $FF
	.byte	$E7, $E7, $E7, $00, $00, $E7, $E7, $E7
	.byte	$3F, $3F, $CF, $CF, $3F, $3F, $CF, $CF
	.byte	$E7, $E7, $E7, $E7, $E7, $E7, $E7, $E7
	.byte	$FF, $FF, $FC, $C1, $89, $C9, $C9, $FF
	.byte	$00, $80, $C0, $E0, $F0, $F8, $FC, $FE
	.byte	$FF, $FF, $FF, $FF, $FF, $FF, $FF, $FF
	.byte	$0F, $0F, $0F, $0F, $0F, $0F, $0F, $0F
	.byte	$FF, $FF, $FF, $FF, $00, $00, $00, $00
	.byte	$00, $FF, $FF, $FF, $FF, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $FF, $FF, $00
	.byte	$3F, $3F, $3F, $3F, $3F, $3F, $3F, $3F
	.byte	$33, $33, $CC, $CC, $33, $33, $CC, $CC
	.byte	$FC, $FC, $FC, $FC, $FC, $FC, $FC, $FC
	.byte	$FF, $FF, $FF, $FF, $33, $33, $CC, $CC
	.byte	$00, $01, $03, $07, $0F, $1F, $3F, $7F
	.byte	$FC, $FC, $FC, $FC, $FC, $FC, $FC, $FC
	.byte	$E7, $E7, $E7, $E0, $E0, $E7, $E7, $E7
	.byte	$FF, $FF, $FF, $FF, $F0, $F0, $F0, $F0
	.byte	$E7, $E7, $E7, $E0, $E0, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $07, $07, $E7, $E7, $E7
	.byte	$FF, $FF, $FF, $FF, $FF, $FF, $00, $00
	.byte	$FF, $FF, $FF, $E0, $E0, $E7, $E7, $E7
	.byte	$E7, $E7, $E7, $00, $00, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $00, $00, $E7, $E7, $E7
	.byte	$E7, $E7, $E7, $07, $07, $E7, $E7, $E7
	.byte	$3F, $3F, $3F, $3F, $3F, $3F, $3F, $3F
	.byte	$1F, $1F, $1F, $1F, $1F, $1F, $1F, $1F
	.byte	$F8, $F8, $F8, $F8, $F8, $F8, $F8, $F8
	.byte	$00, $00, $FF, $FF, $FF, $FF, $FF, $FF
	.byte	$00, $00, $00, $FF, $FF, $FF, $FF, $FF
	.byte	$FF, $FF, $FF, $FF, $FF, $00, $00, $00
	.byte	$FC, $FC, $FC, $FC, $FC, $FC, $00, $00
	.byte	$FF, $FF, $FF, $FF, $0F, $0F, $0F, $0F
	.byte	$F0, $F0, $F0, $F0, $FF, $FF, $FF, $FF
	.byte	$E7, $E7, $E7, $07, $07, $FF, $FF, $FF
	.byte	$0F, $0F, $0F, $0F, $FF, $FF, $FF, $FF
	.byte	$0F, $0F, $0F, $0F, $F0, $F0, $F0, $F0
