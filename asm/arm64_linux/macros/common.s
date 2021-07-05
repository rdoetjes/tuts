
.macro printChar chr
	mov w2, \chr
	ldr x1, =char
	strb w2, [x1]
	mov x2, #1
	bl print
.endm

.macro printString string, length
        ldr x1, =\string
	mov x2, #\length
        bl print
.endm


//Exit to operating system, X0 will contain the exit code
//X0 contains exit code
.macro exit exitCode
	mov x0, \exitCode
        mov X8, #93                     //exit system call
        svc #0                          //call system call (93 => exit)
        ret                             //this won't be called because exit will terminate program before we get here
.endm

//Print to STDOUT the message pointed by X1
//X1 is string ptr 
//X2 is string length of the string
print:
	stp x29, x30, [sp, #-16]!	//store fp and sp on the stack
	stp x2, x8, [sp, #-16]!		//store x2 and x8 on the stack (so they won't be globbered)
	stp x0, x1, [sp, #-16]!		//store x0 and x1 on the stack (so they won't be globbered)

	mov x0, #1			//set X0 to point to standard out
	mov x8, #64			//write syscall 
	svc #0				//call system call (64 => write)

	ldp x0, x1, [sp], #16		//pop x0 and x1 from stack (so they won't be globbered)
	ldp x2, x8, [sp], #16		//pop x2 and x8 from stack (so they won't be globbered)
	ldp x29, x30, [sp], #16		//pop fp and sp from stacl
	ret				//return to caller
