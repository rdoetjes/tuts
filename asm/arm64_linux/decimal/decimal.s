.text
.align 8
.global _start

_start:
	ldr x1, =input
	ldr x2, =input_len
	bl inputSTDIN

	mov x2, x0
	bl atoi
	bl printUInt

	mov x0, #0
	bl exit

.include "common.s"
.include "input.s"
.include "conversion.s"

.data
input:	.fill 20,1,0
input_len = . - input
