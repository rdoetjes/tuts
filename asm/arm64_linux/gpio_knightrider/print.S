#ifndef PRINT_S
#define PRINT_S
#include "print.macro"

.text
.align 8

//x0 contains the number to print
//Print value in X0 as an unisgned int to screen
printUInt:
	stp x29, x30, [sp, #-16]!	//store frame pointer and  stack pointer on the stack
	stp x5, x7, [sp, #-16]!		//push x5 and x7 to stack (so they won't be globbered)
	stp x2, x3, [sp, #-16]!		//push x2 and x3 to stack (so they won't be globbered)
	stp x0, x1, [sp, #-16]!		//push x0 and x1 to stack (so they won't be globbered)

	mov x7, #10			//x7 will contain the divider (10) used in udiv and msub
	mov x5, #0			//x5 counts the number of digits stored on stackl
	sub sp, sp, #128		//move stack pointer down 128 bytes, so we have space to store the to print digits
	
	cmp x0, #0			//if x0=0 then the division algorith will not work
	beq 2f		//we set the value on the stack to 0

1:
	udiv x2, x0, x7			//divide the value x0 by 10
	msub x3, x2, x7, x0		//obtain the remainder (x3) and the Quotient (x2)
	add x5, x5, #1			//increment the digit counter (x5)
	strb w3, [sp, x5]		//store the digit on the stack as single byte
	mov x0, x2			//copy the Quotient (x2) into x0 which is the new value to divide by 10
	cmp x0, #0			//if the Quotient (x0) is 00 then we found all individual digits
	bne 1b				//if x0 is not yet zero than there's more digits to extract
	b 3f				//we set all the digits on the stack now we can pop them off and print them

2:				//this is the exceptional case when x0 is 0 then we need to push this ourselves to the stack
	add x5, x5, #1 			//x5 is not used so still 0, there for we need to offset it by 1 for the sp offset
	strb w0, [sp, x5]		//set the value 0 to the stack, so that it can be printed to the screen

	//using the stacl guarantees that the digits are printed in the right order (from large to smallest_
3:
	ldrb w3, [sp, x5]		//pop the last digit from the stack (the biggest value)
	ldr X1,=printUInt_char		//set X1 to the char variable address, so we can store the char there later on
	add w3, w3, 48			//add 48 to the number, turning it into an ASCII char 0-9
	strb w3, [x1]			//store the ASCII char in the char variable (pointed to by X1)
	mov x2,#1 			//set the length to write() to 1
	bl print			//call the writer() system call wrapper 
	subs x5, x5, #1			//reduce x5 by 1, pointing to the next digit on the stack 
	bne 3b				//if x5 is not 0 then there are still digits on the stack, that should be printed

	add sp, sp, #128		//reclaim the 128 bytes local storage on the stack 
	ldp x0, x1, [sp], #16		//pop x0 and x1 from stack (so they won't be globbered)
	ldp x2, x3, [sp], #16		//pop x2 and x3 from stack (so they won't be globbered)
	ldp x5, x7, [sp], #16		//pop x5 and x7 from stack (so they won't be globbered)
	ldp x29, x30, [sp], #16		//pop fp and sp from stack (so they won't be globbered)
	ret 				//return


//Print to STDOUT the message pointed by X1
//X1 is string ptr 
printString:
	stp x29, x30, [sp, #-16]!	//store fp and sp on the stack
	stp x2, x8, [sp, #-16]!		//store x2 and x8 on the stack (so they won't be globbered)
	stp x0, x1, [sp, #-16]!		//store x0 and x1 on the stack (so they won't be globbered)

	mov x2, #0
1:
	ldrb w8, [x1, x2]
	add x2, x2, #1
	cmp w8, #0
	bne 1b

	mov x0, #1			//set X0 to point to standard out
	mov x8, #64			//write syscall 
	svc #0				//call system call (64 => write)

	ldp x0, x1, [sp], #16		//pop x0 and x1 from stack (so they won't be globbered)
	ldp x2, x8, [sp], #16		//pop x2 and x8 from stack (so they won't be globbered)
	ldp x29, x30, [sp], #16		//pop fp and sp from stacl
	ret				//return to caller


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

.align 8
.data
printUInt_char:	.byte 0

#endif
