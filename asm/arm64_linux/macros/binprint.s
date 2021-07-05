.include "common.s"

.text

.align 8

.global _start

_start:	
	printString s_bitcount, len_s_bitcount

	ldr x0, =value			//set the number that we want to count the number of set bits from
	bl PrintBin

	printString s_bitcount1, len_s_bitcount1

	bl bitCount			//count the number of set bits in X0, result will be in X0
	bl printUInt			//prints the number of set bits
	
	exit	0			//call exit

//Will count the number of set bits in the number in register X0
//Result X0 will contain the number of set bits after calling this routine
bitCount:
	pushPair x29, x30		//store frame pointer and  stack pointer on the stack
	pushPair x1, x2			//store x1 and x2 on the stack (so they won't be globbered)
	
	mov x2, #0			//x2 is used to count the set bits, initialize it to 0
bitCount_Counter:
	cmp x0, #0			//if x0 is 0 then there's no more set bits and we should return
	beq bitCount_Exit		//return when there's no more set bits
	add x2, x2, #1			//increment the bit counter, everytimne we hit the loop a bit is set
	//n = n & n-1
	sub x1, x0, #1			//n &= n-1 this is the Kernighan algorithm, to count bits efficiently
	and x0, x0, x1			//(where x0 = n and x1 = n-1)

	bne bitCount_Counter		//x0 is not 0 yet, so there are still be set bits in the value x0

bitCount_Exit:
	mov x0, x2			//copy number of set bits to x0, as x0 contains our result
	popPair x1, x2		//pop x1 and x2 from stack (so they won't be globbered)
	popPair x29, x30		//pop fp and sp from stack (so they won't be globbered)
	ret				//return

.data
.align 8
s_bitcount:	.ascii "We are counting the number of set bits in 0x"
len_s_bitcount = . - s_bitcount

s_bitcount1: 	.ascii " which there are #"
len_s_bitcount1 = . - s_bitcount1

value = 255
