.text
.align 8
.global _start

_start:
	ldr X1, =string
	ldr X2, =len
	bl print

	mov X0, #0
	bl exit

print:
	mov X0, #1
	mov X8, #64
	svc #0
	ret

exit:
	mov X8, #93
	svc #0
	ret
	
.data
string:	
	.ascii "Hello, World! ARM CPUs rule!\r\n"
len = . - string
