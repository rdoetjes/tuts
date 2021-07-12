.text
.align 8

inputSTDIN:
	stp x29, x30, [sp, #-16]!
	stp x2, x8, [sp, #-16]!

	mov x0, #0
	mov x8, #63
	svc #0

	subs x0, x0, #1
	ldrb w8, [x1, x0]
	cmp w8, #'\n'
	beq inputSTDIN_exit
	add x0, x0, #1 

inputSTDIN_exit:
	ldp x2, x8, [sp], #16
	ldp x29, x30, [sp], #16
	ret
.data
