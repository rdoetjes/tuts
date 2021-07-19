#ifndef SYSTEM_S
#define SYSTEM_S

.include "system.macro"

.text
.align

//X0 contains a pointer to a sleep value in nano seconds
//x1 contains a poinyer to a sleep value in nano seconds when interrupt is called (ignored)
sys_nanoSleep:
	stp x29, x30, [sp, #-16]!
	stp x0, x1, [sp, #-16]!

	mov x8, 101		//nano sleep syscall 
	svc #0

	ldp x0, x1, [sp], #16
	ldp x29, x30, [sp], #16
	ret

//Exit to operating system, X0 will contain the exit code
//X0 contains exit code
sys_exit:
	mov X8, #93			//exit system call
	svc #0				//call system call (93 => exit)
	ret				//this won't be called because exit will terminate program before we get here
#endif
