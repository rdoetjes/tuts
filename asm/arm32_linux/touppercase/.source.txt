.text						//code segment
.align 4					//align to even boundary
.global _start					//tell assembler what the entry point is

_start:						//entry point label
	ldr R0, =string				//load in R0 the address where "string" begins
	ldr R1, =uppercase			//load in R1 the address where the (empty) uppercase string begins
	ldr R2, =len				//load in R2 the length of the string, named string
	bl toupper				//call the procedure to upper

	mov R3, #0xFFFF				//set R3 to the value 64k, this is our loop counter
loop:					
	ldr R1, =uppercase			//load the begin address for the uppercase string in R1
	ldr R2, =len				//load the length of the string named string in R2
	bl print				//call procedure print

	subs R3, R3, #1				//decrement R3 by one
	bne loop				//if R0 not zero than continue to loop and print the uppercase message

	mov R0, #0				//set exit code to 0 (normal exit)
	bl exit					//exit the program

//toUpper function changes string in R0 to uppercase
//parameters: 	R0 is input string
//		R1 is the output string (must have same length!)
//		R2 is the length of the input string (and output string)
//Clobbers: R3
toupper:					
	push { lr }				//save the link register to jump back to
_toupper:
	ldrb R3, [R0], #1			//load the byte pointed to by R0 in R3 and increment R0 ptr by 1
	cmp R3, #'a'				//compare if char read in R3 is greater or equal to a (97) 
	bge _uppercase				//the turn the value into uppercase uppercase
	b _store				//else store the value in the output string

_uppercase:				
	cmp R3, #'z'				//check if the R3 character is greater than z (122)
	bgt _store				//if so than it's not a char in range [a-z], so just store it in output string
	subs R3, R3, #32			//now R3 has a char in range of [a-z] subtract 32 to make it uppercase

_store:
	strb R3, [R1], #1			//store the bye in R3 to the address pointed in R1 and increment R1 otr by 1
	subs R2, R2, #1				//decrement the string length (R2) by 1
	bne _toupper				//is string length (R2) not 0 then continue to change case 
	pop { pc }				//return 

//print function prints what is in R1 to the STDOUT
//parameters:	R1 is the string you want to write to STDOUT
//Clobbers: R7 and R0
print:
	push { lr }				//save the ling register to jump back to
	mov R0, #1				//set write() file descriptor to 1 (STDOUT)
	mov R7, #4				//set syscall to 4, which is write() syscall
	svc #0					//call software interrupt 0, which executes syscall
	pop { pc }				//return

//exit function, exits to OS and retruns what is in R0
//parameters:	R0 is the exit code returned to OS
//Clobbers: nothing, but who cares at this point ;)
exit:
	push { lr }				//store link register to jump back to
	mov R7, #1				//set syscall to 1, which is exit() syscall
	svc #0					//call software interrupt 0, which executes syscall
	pop { pc }				//return
	
.data
string:	
	.ascii "Hello World, ARM cpus rule!\n"	//this is our string
len = . - string				//the length of the string, calculated at assemble time (. is current address)
.align 2					//align nicely (this is optional)
uppercase:
	.fill len, 1, 0				//the uppercase string is the same length is the precalulated string length
