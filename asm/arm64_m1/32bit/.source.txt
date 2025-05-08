.global _start                      //_start is the starting address for MacOS

.text
.align 4                            //align on even boundaries, to satisfy drawin

_start:
    mov W0, #0x0
_loop:
    adds W0, W0, #1
    bne _loop

_exit:
    mov X16, #1                     //X16 holds syscall, syscall #1 is exit 
    svc #0x80                       //call sofwtare interrupt 80h
