.text						//code segment
.align 4					//align to even boundary
.global _start					//tell assembler what the entry point is

_start:						//entry point label
	mov R0, #1				//set write() file descriptor to 1 (STDOUT)
	ldr R1, =string				//load in R0 the address where "string" begins
	ldr R2, =len
	mov R7, #4				//set syscall to 4, which is write() syscall
	svc #0					//call software interrupt 0, which executes syscall

	mov R7, #1				//set syscall to 1, which is exit() syscall
	svc #0					//call software interrupt 0, which executes syscall
	
.data
string:	
	.ascii "Hello World, ARM cpus rule!\n"	//this is our string
len = . - string				//the length of the string, calculated at assemble time (. is current address)
