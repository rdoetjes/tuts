.macro printChar chr
        mov w2, \chr
        ldr x1, =common_char
        strb w2, [x1]
        mov x2, #1
        bl print
.endm

.macro printString string, length
        ldr x1, =\string
        mov x2, #\length
        bl print
.endm

.macro pushPair r1, r2
        stp \r1, \r2, [sp, #-16]!
.endm

.macro popPair r1, r2
        ldp \r1, \r2, [sp], #16
.endm

//Exit to operating system, X0 will contain the exit code
//X0 contains exit code
.macro exit exitCode
        mov x0, \exitCode
        mov X8, #93                     //exit system call
        svc #0                          //call system call (93 => exit)
        ret                             //this won't be called because exit will terminate program before we get here
.endm
