.global _start                      //_start is the starting address for MacOS
.align 4                            //align on even boundaries, to satisfy drawin

_start:
    mov W3, #0xFFFFF
_loop:
        //source string
        adrp X0, myString@PAGE          //set X1 to address of the string
        add X0, X0, myString@PAGEOFF    //align to 64 bits    
        //dest string to upper case
        adrp X1, upperCase@PAGE          //set X1 to address of the string
        add X1, X1, upperCase@PAGEOFF    //align to 64 bits
        //length string
        ldr X2, =lenMyString             //set X2 to the number of bytes in the string by using len = . - myString
        bl toUpperCase                   //call out toUppercase procedure

        mov W3, #0xFFFF                  //print the upper case messafge FFFF times
loop:
        adrp X1, upperCase@PAGE          //set X1 to address of the string
        add X1, X1, upperCase@PAGEOFF    //align to 64 bits
        //length string
        ldr X2, =lenMyString             //set X2 to the number of bytes in the string by using len = . - myString
        bl print                         //call the print procedure
        subs W3, W3, #1                  //decrement the W3 (FFFF) counter 
        bne loop                         //if W3 not 0 then repeat printing

        bl exit                          //call exit procedure

toUpperCase:
        ldrb W4, [X0], #1                //load a byte in to W4 and increment X0 offset by 1 (to fetch next byte)
        cmp W4, #'a'                     //if W4 is greate or equal to 'a' char then turn uppercase
        bge _uppercase                      
        b _store                         //if W4 < 'a' char then just store the byte in the new char array; as it not [a-z]

_uppercase:
        cmp W4, #'z'                     //if W4 >z then just store the byte in the new char array; as it's not [a-z]
        bgt _store                       
        subs W4, W4, #32                 //if W4 is a-z then substract 32, that turns it into A-Z (capitols)

_store:
        strb W4, [X1], #1                //store the byte in W4, into X1 pointer, and increment X1 by one, to store next char
        subs X2, X2, #1                  //decrement the length of the string, to read the next char and turn into uppercase
        bne toUpperCase                  //if X2 (length) is not 0 then continue, because we are not at the beginning of the string
        ret                              //return to the caller's next instruction

print:
        mov X0, #1
        mov X16, #4                     //X16 holds the sys call, sycall 4 is write to 
        svc #0x80                       //call software interrupt 80h
        ret                             //return to the caller's next instruction; although this will never happen since we exit
exit:
        mov X16, #1                     //X16 holds syscall, syscall #1 is exit 
        svc #0x80                       //call sofwtare interrupt 80h
        ret

.data
myString:   
    .ascii "Oh boy Apple, this has been released too soon! It's bug ridden!\r\n"
lenMyString = . - myString
.align 2
upperCase:
        .fill lenMyString, 1, 0
