.text
.align 8
.global _start

_start:

	ldr X1, =string
	ldr X2, =len
	bl print

	mov X0, #0
	bl exit

//X1 is string ptr, X2 is string lenmg
print:
	mov X0, #1
	mov X8, #64	//write syscall 
	svc #0
	ret

//X0 contains exit code
exit:
	mov X8, #93
	svc #0
	ret

.data
string:
	.ascii "Hello World, ARM64 rules!\n"
len = . - string
