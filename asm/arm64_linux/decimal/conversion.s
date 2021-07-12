.text
.align 8

//Converts the string in x1 into decimal
//x1 is input
//x2 is length of input
//x0 contains output
atoi:
	stp x29, x30, [sp, #-16]!       //store frame pointer and  stack pointer on the stack
        stp x5, x6, [sp, #-16]!         //push x5 and x7 to stack (so they won't be globbered)
        stp x3, x4, [sp, #-16]!         //push x2 and x3 to stack (so they won't be globbered)
        stp x1, x2, [sp, #-16]!         //push x0 and x1 to stack (so they won't be globbered)

        mov x5, #1                      //multiplier
        mov x0, #0                      //The decimal result that will be returned, is set to 0 to allow base10 operations
        mov x3, #0                      //intermediate value to calulator decimnal part
        mov x6, #10                     //the decimal multiplier

atoi_readchar:
	subs x2, x2, #1			//subtract 1 from the byte counter (because length starts at 1 but is byte offset 0 in buffer)
	ldrb w3, [x1, x2]		//load the byte from the heap (from highest byte on the heap (least significat decimal) to the beginning
	sub x3, x3, #48			//turn char into integer
	mul x3, x3, x5			//multiply the value by it's decimal denominator (1, 10, 100, 1000 etc)
	add x0, x0, x3			//add that decimal position to the result
	mul x5, x5, x6			//mulitply the multiplier by 10, to get the next decimal denominator
	cmp x2, #0			//are we at the beginning of our input string? 
	bne atoi_readchar		//if not continue to read and convert the next character

	ldp x1, x2, [sp], #16           //pop x0 and x1 from stack (so they won't be globbered)
        ldp x3, x4, [sp], #16           //pop x2 and x3 from stack (so they won't be globbered)
        ldp x5, x6, [sp], #16           //pop x5 and x7 from stack (so they won't be globbered)
        ldp x29, x30, [sp], #16         //pop fp and sp from stack (so they won't be globbered)
        ret
.data
