.text							//code segment
.align 8						//align to 64 bits
.global _start						//set entry point

_start:
	//save current break address
	mov X0, #0
	bl sbrk
	ldr X1, =curbrk
	str X0, [X1]
	
	//request new brk address to house =len bytes of data
	ldr X1, =curbrk
	ldr X0, =len
	add X0, X0, X1
	bl sbrk
	ldr X1, =newbrk
	str X0, [X1]
	cmp X0, #-1
	beq oom

	//setup for copying from hard coded data to dynamic allocated memory 
	ldr X0, =string
	ldr X1, =newbrk
	ldr X2, =len
	//copy the data 
_copy:
	ldrb W3, [X0], #1
	sub W3, W3, #1
	strb W3, [X1], #1
	subs X2, X2, #1
	bne _copy

	//print the orig data
	ldr X1, =string
	ldr X2, =len
	bl print

	//print the copied data
	ldr X1, =newbrk
	ldr X2, =len
	bl print

	//exit
	mov X0, #0					//set X0 (exit status) to 0
	bl exit						//call exit function

oom:
	ldr X1, =oomstr
	ldr X2, =lenoom
	bl print
	mov X0, #-1
	bl exit

//X1 is string ptr, X2 is string lenmg
print:
	mov X0, #1					//set file descriptor for write() to STDOUT
	mov X8, #64	 				//Set write system call
	svc #0						//call syscall (64 => write())
	ret						//return to caller

sbrk:
	mov X8,# 214	
	svc #0
	ret

//X0 contains exit code
exit:
	mov X8, #93					//set exit system call
	svc #0						//call syscall (93 => exit)
	ret						//return to caller (this will never be executed, because we've already exitted)

.data
//string that is our source string
string:	
	.ascii "asdasdsijfiljropijrpiofjkijgpeoirdsjfposridjrjfoierjgpoidfjskerjfpiorjgdofijv;reoigjwoe;irqjf;oiejgr;oirwtejg;bcdefghijkmlnopqrstuvwxy\n"
len = . - string

//out of memory message
oomstr:
	.ascii "Out of memory error\n"
lenoom = . - oomstr

//dynamic allocated memory results in a page ptr that we will set in newbrk
.align 8
curbrk = .  
newbrk = .   
