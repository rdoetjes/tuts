.text							//code segment
.align 8						//align to 64 bits
.global _start						//set entry point

_start:
	ldr X1, =string					//load string ptr in X1
	ldr X2, =len					//load the lenghth into X2
	bl print					//call print function

	mov X0, #0					//set X0 (exit status) to 0
	bl exit						//call exit function

//X1 is string ptr, X2 is string lenmg
print:
	mov X0, #1					//set file descriptor for write() to STDOUT
	mov X8, #64	 				//Set write system call
	svc #0						//call syscall (64 => write())
	ret						//return to caller

//X0 contains exit code
exit:
	mov X8, #93					//set exit system call
	svc #0						//call syscall (93 => exit)
	ret						//return to caller (this will never be executed, because we've already exitted)

.data
string:
	.ascii "Hello World, ARM64 rules!\n"
len = . - string
