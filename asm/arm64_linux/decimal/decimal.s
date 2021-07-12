.text
.align 8
.global _start

_start:
	ldr x1, =input
	ldr x2, =input_len
	bl inputSTDIN

	mov x2, x0		//read number of bytes read in inputSTDIN to x2
	ldr x1, =input		
	bl atoi
 	bl printUInt

	mov x0, #0
	bl exit

.include "common.s"
.include "input.s"
.include "conversion.s"

.align 8
.data
input:	.fill 4,1,0
input_len = . - input
