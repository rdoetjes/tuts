.text
.align 8
.global _start

_start:

	ldr X0, =string
	ldr X1, =uppercase
	ldr X2, =len
	bl toupper

	ldr X1, =uppercase
	ldr X2, =len
	bl print

	mov X0, #0
	bl exit

toupper:
	ldrb w3, [X0], #1
	cmp w3, #'a'
	bge _toupper
	b _store

_toupper:
	cmp w3, #'z'
	bgt _store
	subs X3, X3, #32

_store:
	strb w3, [X1], #1
	subs X2, X2, #1
	bne toupper
	ret

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
	.ascii "Hello world, ARM64 is the way to go!\n"
len = . - string
uppercase:
	.fill len, 1, 0
