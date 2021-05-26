.text
.align 4
.global _start

_start:
	mov X3, #0xFF
loop:
	ldr X1, =string
	ldr X2, =len
	bl print
	subs X3, X3, #1
	bne loop
	bl exit

print:
	push { lr }
	mov X0, #1
	mov X7, #4
	svc #0
	pop { pc }

exit:
	push { lr }
	mov X0, #0
	mov X7, #1
	svc #0
	pop { lr }
	
.data
string:	
	.ascii "Hello World!\r\n"
len = . - string
