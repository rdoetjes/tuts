.text
.align 8
.global _start

.include "common.s"

_start:	
	printString s_bitcount, len_s_bitcount

	ldr x0, =value			//set the number that we want to count the number of set bits from
	bl PrintBin

	printString s_bitcount1, len_s_bitcount1

	bl bitCount			//count the number of set bits in X0, result will be in X0
	bl printUInt			//prints the number of set bits
	
	exit	0			//call exit

//Prints the value in X0 into 64 bit binary format
//Argument X0 contains value to print
PrintBin:
	stp x29, x30, [sp, #-16]!       //store frame pointer and  stack pointer on the stack
        stp x4, x5, [sp, #-16]!         //store x4 and x5 on the stack, so they won't get globbered
        stp x2, x3, [sp, #-16]!         //store x2 and x4 on the stack, so they won't get globbered
        stp x0, x1, [sp, #-16]!         //store x0 and x1 on the stack, so they won't get globbered

	mov w3, #0			//bit counter
	mov x4, #0x8000000000000000	//bit mask (in binary, it a 1 with 63 zeros)
PrintBin_mask:
	tst x0, x4			//test if the bit in X4 is set in x0
	beq PrintBin_zero		//the zero bit is set by the tst opcode (hence printing zero if equal; inverse login)

PrintBin_one:
	printChar '1'
	b PrintBin_counter		//jump to continue the loop

PrintBin_zero:
	printChar '0'

PrintBin_counter:
	lsr x4, x4, #1			//shift the mask over 1 bit to the left
	add w3, w3, #1			//increment the counter
	cmp w3, #64			//if counter hits 64 then we printed the 64 bits (we started at 0)
	bne PrintBin_mask		//if counter not 64, then continue the loop
	
PrintBin_Exit:
        ldp x0, x1, [sp], #16           //restore x0 and x1 so they won't be globbered
        ldp x2, x3, [sp], #16           //restore x2 and x3 so they won't be globbered
        ldp x4, x5, [sp], #16           //restore x4 and x5 so they won't be globbered
        ldp x29, x30, [sp], #16         //restore fp and sp so they won't be globbered
        ret


//We will print the contents of X0 to the screen in hex
PrintHex:
	stp x29, x30, [sp, #-16]!	//store frame pointer and  stack pointer on the stack
	stp x4, x5, [sp, #-16]!		//store x4 and x5 on the stack, so they won't get globbered
	stp x2, x3, [sp, #-16]!         //store x2 and x4 on the stack, so they won't get globbered
	stp x0, x1, [sp, #-16]!         //store x0 and x1 on the stack, so they won't get globbered

	mov w4, #0			//we use this for a leading zero flag (as long as this is 0, no digits will printed, as they'd be leading zeros)
	mov w3, #0			//this is our nibble counter used to determine when we are done
PrintHex_mask:
	and x5, x0, #0xf000000000000000	//mask the top nibble from x0 (x3 is the mask) store result in x5 for shifting down
	lsr x5, x5, #60			//we need the low nibble as an offset index, so we shift the result down 60 bits
	lsl x0, x0, #4			//we shift x0 up 4 bits, so we can process the next nibble 
					//(we are automatically printing from high nibble to low nibble this way)
	add w3, w3, #1			//increment the nibble counter (we stop this procedure when this is 16)

	cmp w4, #1			//if (x4==1) then print
	beq PrintHex_print		//else continue without printing

	cmp w5, #0			//if (w5>0) 
	bgt PrintHex_setFlag		//then set x4 (no more leading zeros flag) to 1
	
	cmp w3, #16                     //has x0 been shifted all the way up, in effect being 0, no? Then continue
        bne PrintHex_mask

        //if x0 argument was set to 0, then we fall through here and we print the one and only 0

PrintHex_setFlag:
	mov w4, #1			//the first digit is no longer 0 so we can print and continue to print digits
	
PrintHex_print:
	ldr x1, =hexArray		//load the approptiate hex value from the array, x1 is the offset in the hex char array
	ldrb w5, [x1, x5] 		//pick the correct hex character to print from the hexArray, x5 has the right offset, because it's shifted down to the 1st nibble

	printChar w5

	cmp w3, #16			//has the lowest nibble, been shifted all the way up (in effect being 0) 
	bne PrintHex_mask		//No? then continue printing nibble

PrintHex_Exit:
	ldp x0, x1, [sp], #16		//restore x0 and x1 so they won't be globbered
	ldp x2, x3, [sp], #16		//restore x2 and x3 so they won't be globbered
	ldp x4, x5, [sp], #16		//restore x4 and x5 so they won't be globbered
	ldp x29, x30, [sp], #16		//restore fp and sp so they won't be globbered
	ret


//Will count the number of set bits in the number in register X0
//Result X0 will contain the number of set bits after calling this routine
bitCount:
	stp x29, x30, [sp, #-16]!	//store frame pointer and  stack pointer on the stack
	stp x1, x2, [sp, #-16]!		//store x1 and x2 on the stack (so they won't be globbered)
	
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
	ldp x1, x2, [sp], #16		//pop x1 and x2 from stack (so they won't be globbered)
	ldp x29, x30, [sp], #16		//pop fp and sp from stack (so they won't be globbered)
	ret				//return

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
	beq printUInt_Zero		//we set the value on the stack to 0

printUInt_Count:
	udiv x2, x0, x7			//divide the value x0 by 10
	msub x3, x2, x7, x0		//obtain the remainder (x3) and the Quotient (x2)
	add x5, x5, #1			//increment the digit counter (x5)
	strb w3, [sp, x5]		//store the digit on the stack as single byte
	mov x0, x2			//copy the Quotient (x2) into x0 which is the new value to divide by 10
	cmp x0, #0			//if the Quotient (x0) is 00 then we found all individual digits
	bne printUInt_Count		//if x0 is not yet zero than there's more digits to extract
	b printUInt_print			//we set all the digits on the stack now we can pop them off and print them

printUInt_Zero:				//this is the exceptional case when x0 is 0 then we need to push this ourselves to the stack
	add x5, x5, #1 			//x5 is not used so still 0, there for we need to offset it by 1 for the sp offset
	strb w0, [sp, x5]		//set the value 0 to the stack, so that it can be printed to the screen

	//using the stacl guarantees that the digits are printed in the right order (from large to smallest_
printUInt_print:
	ldrb w3, [sp, x5]		//pop the last digit from the stack (the biggest value)
	add w3, w3, #48

	printChar w3

	subs x5, x5, #1			//reduce x5 by 1, pointing to the next digit on the stack 
	bne printUInt_print		//if x5 is not 0 then there are still digits on the stack, that should be printed

printUInt_exit:
	add sp, sp, #128		//reclaim the 128 bytes local storage on the stack 
	ldp x0, x1, [sp], #16		//pop x0 and x1 from stack (so they won't be globbered)
	ldp x2, x3, [sp], #16		//pop x2 and x3 from stack (so they won't be globbered)
	ldp x5, x7, [sp], #16		//pop x5 and x7 from stack (so they won't be globbered)
	ldp x29, x30, [sp], #16		//pop fp and sp from stack (so they won't be globbered)
	ret 				//return

.align 8
.data
s_bitcount:	.ascii "We are counting the number of set bits in b"
len_s_bitcount = . - s_bitcount

s_bitcount1: 	.ascii " which there are #"
len_s_bitcount1 = . - s_bitcount1

value = 255

char: 		.byte 0

hexArray:	.ascii "0123456789ABCDEF"
