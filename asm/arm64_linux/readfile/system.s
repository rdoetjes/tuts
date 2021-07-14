.include "system.macro"

.text
.align

//Exit to operating system, X0 will contain the exit code
//X0 contains exit code
exit:
	mov X8, #93			//exit system call
	svc #0				//call system call (93 => exit)
	ret				//this won't be called because exit will terminate program before we get here

