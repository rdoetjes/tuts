.text
.align 8
.global _start

_start:
	ldr x1, =s_bitcount
	ldr x2, =len_s_bitcount
	bl print

	ldr x0, =value
	bl printNr

	ldr x1, =s_bitcount1
	ldr x2, =len_s_bitcount1
	bl print

	bl bitCount
	bl printNr

	mov X0, #0			//set the exit code to 0
	bl exit				//call exit

//Will count the number of set bits in the number in register X0
//Result X0 will contain the number of set bits after calling this routine
bitCount:
	stp x29, x30, [sp, #-16]!
	stp x1, x2, [sp, #-16]!
	
	mov x2, #0
bitCount_Counter:
	cmp x0, #0
	beq bitCount_Exit
	add x2, x2, #1
	sub x1, x0, #1
	and x0, x0, x1
	bne bitCount_Counter

bitCount_Exit:
	mov x0, x2
	ldp x1, x2, [sp], #16
	ldp x29, x30, [sp], #16
	ret

//Print value in X0 as an unisgned int to screen
printNr:
	stp x29, x30, [sp, #-16]!
	stp x5, x7, [sp, #-16]!
	stp x2, x3, [sp, #-16]!
	stp x0, x1, [sp, #-16]!

	mov x7, #10
	mov x5, #0
	sub sp, sp, #128
printNr_Count:
	udiv x2, x0, x7
	msub x3, x2, x7, x0
	add x5, x5, #1
	strb w3, [sp, x5]
	mov x0, x2
	cmp x0, #0
	bne printNr_Count

printNr_print:
	cmp x5, #00
	beq printNr_exit

	ldrb w3, [sp, x5]
	ldr X1,=char
	add w3, w3, 48
	strb w3, [x1]
	mov x2,#1 
	bl print
	subs x5, x5, #1
	bne printNr_print

printNr_exit:
	add sp, sp, #128
	ldp x0, x1, [sp], #16
	ldp x2, x3, [sp], #16
	ldp x5, x7, [sp], #16
	ldp x29, x30, [sp], #16
	ret 

//Print to STDOUT the message pointed by X1
//X1 is string ptr 
//X2 is string length of the string
print:
	stp x29, x30, [sp, #-16]!
	stp x2, x8, [sp, #-16]!
	stp x0, x1, [sp, #-16]!

	mov X0, #1			//set X0 to point to standard out
	mov X8, #64			//write syscall 
	svc #0				//call system call (64 => write)

	ldp x0, x1, [sp], #16
	ldp x2, x8, [sp], #16
	ldp x29, x30, [sp], #16
	ret				//return to caller

//Exit to operating system, X0 will contain the exit code
//X0 contains exit code
exit:
	mov X8, #93			//exit system call
	svc #0				//call system call (93 => exit)
	ret

.data
s_bitcount:	.ascii "We are counting the number of set bits in #"
len_s_bitcount = . - s_bitcount

s_bitcount1: 	.ascii " which there are #"
len_s_bitcount1 = . - s_bitcount1

value = 65535

char: 		.byte 0
