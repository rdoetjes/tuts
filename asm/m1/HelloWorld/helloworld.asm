.global _start                      //_start is the starting address for MacOS
.align 4                            //align on even boundaries, to satisfy drawin

_start:
    mov W3, #0xFFFFF
_loop:
    mov X0, #1                      //set file descriptor to STDOUT (1)
    adrp X1, myString@PAGE          //set X1 to address of the string
    add X1, X1, myString@PAGEOFF    //align to 64 bits
    ldr X2, =lenMyString                    //set X2 to the number of bytes in the string by using len = . - myString
    mov X16, #4                     //X16 holds the sys call, sycall 4 is write to 
    svc #0x80                       //call software interrupt 80h
    subs W3, W3, #1
    bne _loop

_exit:
    mov X0, #0                      //exit return code
    mov X16, #1                     //X16 holds syscall, syscall #1 is exit 
    svc #0x80                       //call sofwtare interrupt 80h

.align 4
myString:   
    .ascii "Oh boy Apple, this has been released too soon! It's bug ridden!\r\n"
lenMyString = . - myString
