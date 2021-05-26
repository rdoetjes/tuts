.text
.align 4
.global _start

_start:
	ldr R0, =string
	ldr R1, =uppercase
	ldr R2, =len
	bl toupper

	mov R3, #0xFFFF
loop:
	ldr R1, =uppercase
	ldr R2, =len
	bl print

	subs R3, R3, #1
	bne loop
	bl exit

toupper:
	push { lr }
_toupper:
	ldrb R3, [R0], #1
	cmp R3, #'a'
	bge _uppercase
	b _store

_uppercase:
	cmp R3, #'z'
	bgt _store
	subs R3, R3, #32

_store:
	strb R3, [R1], #1
	subs R2, R2, #1
	bne _toupper
	pop { pc }

print:
	push { lr }
	mov R0, #1
	mov R7, #4
	svc #0
	pop { pc }

exit:
	push { lr }
	mov R0, #0
	mov R7, #1
	svc #0
	pop { pc }
	
.data
string:	
	.ascii "Hello World, ARM cpu's rule\n"
len = . - string
.align 2
uppercase:
	.fill len, 1, 0
