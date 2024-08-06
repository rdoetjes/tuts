.const ledScreen = $0400 + ((40*25)/2) - (16/2) //starting position of the input area
.const screenPosStartManual = $0400 + (40*25) - 120 //starting position of the manual on the screen


BasicUpstart2(main)      
main:
    ldx #00
    stx ledPos
    stx ptr

    jsr init_userport
   
    lda #01
    jsr clrScreenAndSetTextColor
    jsr drawScreen

    forever:
        jsr keyb
        jsr loop
        jmp forever

    rts

    /*
        clrScreenAndSetTextColor

        Clear the screen
        And set color ram to color set to whatever value is in A register

        X contains the offset and will countdown from ff (then before the check X is decremented frin 00 to ff) to 00
    */
    clrScreenAndSetTextColor:
        ldx #$00
        tay
    !:
        tya
        sta $d800, x    
        sta $d900, x
        sta $da00, x
        sta $db00, x

        lda #$20
        sta $0400, x    
        sta $0500, x
        sta $0600, x
        sta $0700, x

        dex
        bne !-
        rts

    /*
        printManual

        Prints the instructions starting on the screen location pointed to by screenPosStartManual const 

        X is used as the offset counter from the memory loctaion on the screen pointed
        to by screenPosStartManual
    */
    printManual:
        //print the manual text
        ldx #0
    !:
        lda manualText, x
        cmp #$ff
        beq !+
        inx
        sta screenPosStartManual, x 
        bne !-
        
    !:
        rts

    /*
        printBorder

        Prints a border around the LED stack area
        This is done less efficient but quick and clear
    */
    printBorder:
        //top corners
        lda #$70            //left upper corner char
        sta ledScreen-240-2

        lda #$6e            //right upper corner char
        sta ledScreen-240+17

        //top bar
        ldx #0
    !:
        lda #$43            //horizontal bar (stripe) connecting the two top corner characters
        sta ledScreen-240-1, x
        inx 
        cpx #18
        bne !-

        //vertical bars unrolled loop
        lda #$42            //vertical bar (stripes) connecting to the left and right corner characters
        sta ledScreen-200-2
        sta ledScreen-200+17

        sta ledScreen-160-2
        sta ledScreen-160+17

        sta ledScreen-120-2
        sta ledScreen-120+17

        sta ledScreen-80-2
        sta ledScreen-80+17

        sta ledScreen-40-2
        sta ledScreen-40+17

        sta ledScreen-2
        sta ledScreen+17

        //bottom corners
        lda #$6d            //left bottom corner character
        sta ledScreen+40-2

        lda #$7d            //right bottom corner character
        sta ledScreen+40+17

        //bottom bar
        ldx #0
    !:
        lda #$43            //horizontal bar (stripe) connecting the two bottom corner characters
        sta ledScreen+40-1, x
        inx 
        cpx #18
        bne !-

        rts

    /*
        drawScreen

        Makes the screen and screen border black.
        
        Draws the manual on screen
        
        and prints the border around the inpit area

        Draws 16 dots '.' on the screen to where ledSreen constant is pointing 
        These dots serve as the input field, illustrating how many more values
        one needs to input before it is a valid 16 bit data block, that can be pushed
        into our ledValue array
    */
    drawScreen:
        //set screen to black
        ldx #0
        stx $d020
        stx $d021
        //reset ledPos inidcator to 0 (staring from bit 0 with input)
        stx ledPos
        
        //Write red . on bit as to indicate the user will be filling this bit in
        lda #2
        sta ledScreen+$d400,x
        lda #46
        sta ledScreen,x
        inx 

        //Write red arrow to indicate which line will be editted
        lda #2
        sta ledScreen-1+$d400

        lda #62
        sta ledScreen-1

        //set the rest of the 16 bits to a . to show how many bits are left to fill in
    !:    
        lda #46
        sta ledScreen, x
        lda #1
        sta ledScreen+$d400, x
        inx
        cpx #16
        bne !-

        //call the draw manual 
        jsr printManual
        //call the draw border
        jsr printBorder

        rts

    /*
        init_userport

        Sets user port to parrallel and 8 bit output
    */
    init_userport:
        //set to write mode
        lda #$ff
        sta $dd03
        lda #$01
        sta $dd01
        rts

    /*
        printValue

        Reads the byte from the ledValue array, that is 10 bytes below the current size.
        Because there are 5 lines visiable on screen and each line value is represented by
        two bytes in the ledValue array.
        
        This is used in combination with insertLedValueTopOfList which is called
        when deleting a value from the array. And the ledValue that is
        offscreen (if there's more than 5 values) will then be printed and be shown.
    */
    printValue:
        lda size
        sbc #10     //get the byte from ledValue that is 10 bytes back from the bottom of the ledValue array (indicated by size const how many bytes the array currently is) 
        tax

        ldy #0

        //first byte to print
        lda #1
    !:
        lda pow2, y
        and ledValue, x
        cmp pow2, y
        beq !one+
        lda #48
        sta ledScreen-200,y
        jmp !next+
    !one:
        lda #49
        sta ledScreen-200,y
    !next:
        iny 
        cpy #8
        bne !-

        //second byte to print
        lda size
        sbc #9
        tax

        ldy #0

        lda #1
    !:
        lda pow2, y
        and ledValue, x
        cmp pow2, y
        beq !one+
        lda #48
        sta ledScreen-200+8,y
        jmp !next+
    !one:
        lda #49
        sta ledScreen-200+8,y
    !next:
        iny 
        cpy #8
        bne !-

        rts

    /*
        insertLedValueTopOfList

        Get the ledValue entry based on size - 5 if size<5 then insert empty line.
        This is used in combination with printValue to actually print the binary
        representation.
    */
    insertLedValueTopOfList:
        ldx size
        cpx #11
        bcs !++

        //if the size is zero than don't do anything.
        cpx #0
        beq !end+

        //we are on the bottom 5 entries on the list, so we can't add a new entry at the top, sow we put empty line on top 
        //decreate the list size by 2 bytes
        dex
        dex
        stx size  
        
        //raw empty chars on the top line of the list
        ldx #16
        lda #32
    !: 
        sta ledScreen-200-1, x
        dex
        bne !-
        jmp !end+

        //We have an entry left in the list so we push it on the top of the screen
    !:
        ldx size
        dex
        dex
        stx size
        jsr printValue

    !end:
        rts

    /*
        moveListDown

        Move down and add the led value of ptr-5 ot the top row.
        It relies upon insertLedValueTopOfList to show the next entry
        in the array, that is off screen.
    */
    moveListDown:
        //no need to move is size is 0, just redraw drawScreen (input area)
        ldx size
        cpx #0
        beq !zero+

        ldx #16
    !:
        dex

        lda ledScreen-40, x
        sta ledScreen, x

        lda ledScreen-80, x
        sta ledScreen-40, x

        lda ledScreen-120, x
        sta ledScreen-80, x

        lda ledScreen-160, x
        sta ledScreen-120, x

        lda ledScreen-200, x
        sta ledScreen-160, x

        cpx #0
        bne !-

        //pop either the previous entry from the list and show it or create empty line to be copied
        jsr insertLedValueTopOfList

    !zero:
        jsr drawScreen

    !end:
        rts

    /*
        moveLineUp

        Copy the line and move it up one, and do this for the other lines as well.
        Resulting in the binary representation on screen of the ledValue in the array,
        to move upwards (and remove thge top entry from the screen when there's 5 entries in the list)
    */
    moveLineUp:
        ldx #29
    !:
        dex

        lda ledScreen-160, x
        sta ledScreen-200, x

        lda ledScreen-120, x
        sta ledScreen-160, x

        lda ledScreen-80, x
        sta ledScreen-120, x

        lda ledScreen-40, x
        sta ledScreen-80, x

        lda ledScreen, x
        sta ledScreen-40, x
        cpx #0
        bne !-

        lda #0
        sta ledPos

    !end:
        rts

    /*
        storeLed
        Calculates the hexadecimal value of the 16 bits displayed on the screen.
        We iteratre through the screen positions (no seperate bytes allocated for that) testing each
        character position to see if it's '1' or '0'. 
        We shift the bit over for every character and if it's 1 we will "OR 1" in to the A that holds
        eventually, after shifting, the 8 bit value.

        We do this in two steps/llops one for the high 8 bits and the exact same loop copied 
        (not efficient, but works and fine for now... Hmmm refactor I will says Yoda)
        to calculate the next byte.

        This new string of the binary ledValue at the offset of pos.
    */
    storeLed:
        ldx #0
        lda #0
    !:
        asl
        ldy ledScreen, x
        cpy #48
        beq !notOne+
        ora #1
    !notOne:
        inx
        cpx #8
        bne !-
        //first byte calculated
        ldy size
        sta ledValue, y

        ldx #8
        lda #0
    !:
        asl
        ldy ledScreen, x
        cpy #48
        beq !notOne+
        ora #1
    !notOne:
        inx
        cpx #16
        bne !-

        //second byte calculated
        ldy size
        iny
        sta ledValue, y
        iny
        sty size
        
        jsr moveLineUp
        jsr drawScreen
        rts

    /*
        ledToggle
        Switches one of the 16 dots on the screen (which simulate the LED BAR) to '1' or '0' , corresponding to the value in A pointed, 
        which contains the value 49 (char 1) or 48 (char 0) .
        
        ledPos contains the bit number that we are filling in.

        The next bit we will fill will be marked by turning that dot red.
    */
    ledToggle:
        ldx ledPos
        sta ledScreen, x

        //store the new position offset
        inx
        stx ledPos
        
        //when we have 16 entries, we will store the two bytes in memory, os the can be shown by LED bar
        cpx #16
        beq storeLed

        //Write a red . for the next position to be filled by user
        ldx ledPos
        lda #46
        sta ledScreen, x
        lda #2
        sta ledScreen+$d400, x

        rts

    /*
        ledPosLeft

        Perform a backspace, so you can overwrite the previous bit(s).

        We decrement ledPos and draw a . on the old offset pointed
        to by ledPos
    */
    ledPosLeft:
    !:
        jsr $F142   //extra deboune of backspace just to be sure
        bne !-

        ldx ledPos
        beq !+
        
        //draw . on current position
        lda #46
        sta ledScreen, x
        lda #1
        sta ledScreen+$d400, x
        
        //dec x and write ?
        dex
        stx ledPos
        lda #46
        sta ledScreen, x
        lda #2
        sta ledScreen+$d400, x

    !:      
        rts    

    /*
        resetList

        resets the size to 0
        and sets the first ledValue to 0 (all LED's off)
        So you can start a new LED sequence

        It also redraws the whole screen
    */
    resetList:
        ldx #0
        stx size
        stx ledValue
        stx ledValue+1
        
        lda #01
        jsr clrScreenAndSetTextColor

        jsr drawScreen
        rts


    /*
        keyb

        read from keyboard using KERAL calls: ff9f and ffe4
    */
    keyb:
        jsr $ff9f
        clc         //clearning the carry to force a debounce
        jsr $ffe4   //get char
        cli         //clear interrupt to debounce

        cmp #43  //+ speed up
        beq incSpeed

        cmp #45  //- slow down
        beq decSpeed

        cmp #61  //= normal speed
        beq normSpeed

        cmp #48  //0 set LED on position OFF
        beq ledToggle

        cmp #49  //1 set LED on position ON
        beq ledToggle

        cmp #20  //del key, reduce ledPos by one
        beq ledPosLeft
 
        cmp #82  //r key, will reset the list
        beq resetList

        cmp #145  //cursor down key, will reset the list
        beq jMoveListDown
        jmp end
    jMoveListDown:
        jsr moveListDown

    end:
        rts

    /*
        incSpeed

        increase the userPort display speed
    */
    incSpeed:
        dec delayP
        lda delayP
        rts

    /*
        decSpeed

        slow the userPort display speed
    */
    decSpeed:
        inc delayP
        lda delayP
        rts
    
    /*
        reset to normal speed
    */
    normSpeed:
        lda #$08
        sta delayP
        rts    

    /*
        delay

        slow down whole propgran, we may want to use a timer
    */
    delay:
        ldx #$80
    w:
        dex
        bne w
        rts


    /*
        loop

        loop through the array ledValue load the data for each array element
        and output that on the userport.

        The data is stored as two bytes values, for the
        LED bar's 16 LEDs
        First the lowbyte is written and then the high byte
    */
    loop:
        ldx #0
        //set low led bar (requires PB2 to be high)
        lda $dd00
        and #%11111011
        sta $dd00
        ldx ptr
        lda (ledValue),x
        sta $dd01

        ldy delayP                                             
    wait1:  
        jsr delay
        dey
        bne wait1
        
        //set high led bar (requires PB2 to be low)
        lda $dd00
        ora #%00000100
        sta $dd00
        ldx ptr
        inx
        lda (ledValue),x
        sta $dd01
        
        ldy delayP
    wait2:  
        jsr delay
        dey
        bne wait2
        inc ptr
        inc ptr
        ldx ptr
        cpx size
        bcs reset

        //continue 
        rts
    reset:
        ldx #$00
        stx ptr
        rts


/*
    variable memory
*/
keyPressed:         //which key was pressed is stored in here
        .byte 00    
delayP:             //the delay time (large is slower)
        .byte $08   
ptr:                //to which element in the array we are currently pointing
        .byte 00    
ledPos:             //shows where on the input, for a new array entrey we are (0-15)
        .byte 00    
pow2:               //this prevents us from implementing shifting logic, that would've taken as much space and be slower
        .byte 128
        .byte 64
        .byte 32
        .byte 16
        .byte 08
        .byte 04
        .byte 02
        .byte 01
    
manualText:         //this is the manual text that is shown on the screen
    .text @"speed: + faster, - slower, = normal     "
    .text @"r reset, up remove last entry           "
    .text @"1 bit on, 0 bit off, del remove bit\$ff"

/*
    The array with the LED values
    size contains the number of bytes currently used in this 512 byte array
*/
*=$1000
size:           //shows how many bytes are in
    .byte 00
ledValue:       //so you can hold 256 positions, each having two values (therefore 512 bytes reserved) 
   .fill 512, 00