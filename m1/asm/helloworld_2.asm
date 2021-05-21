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
        ldr X2, =lenMyString                    //set X2 to the number of bytes in the string by using len = . - myString
        bl toUpperCase

        mov W3, #0xFFFF
loop:
        adrp X1, upperCase@PAGE          //set X1 to address of the string
        add X1, X1, upperCase@PAGEOFF    //align to 64 bits
        //length string
        ldr X2, =lenMyString                    //set X2 to the number of bytes in the string by using len = . - myString
        bl print
        subs W3, W3, #1
        bne loop

        bl exit

toUpperCase:
        ldrb W4, [X0], #1
        cmp W4, #'a'
        bge _uppercase
        b _store

_uppercase:
        cmp W4, #'z'
        bgt _store
        subs W4, W4, #32

_store:
        strb W4, [X1], #1
        subs X2, X2, #1
        bne toUpperCase
        ret

print:
        mov X0, #1
        mov X16, #4                     //X16 holds the sys call, sycall 4 is write to 
        svc #0x80                       //call software interrupt 80h
        ret
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