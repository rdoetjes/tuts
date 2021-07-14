//The macro's serve as a sort of an interface that 
//adds an abstracts for handling of the arguments for us.
//This way we can use either the sys_* functions and setup the
//arguments ourselves. Or we can use the convenient (legible)
//method, by using macros
.include "fileio.macro"

.text

.align 8

//x0 filename
//x1 flags
//x2 mode
//RETURN X0 contains file descriptor
sys_open:
	stp x29, x30, [sp, #-16]!

	mov x8, #56		//open syscall
	svc #0

	ldp x29, x30, [sp], #16
	ret

sys_read:
	stp x29, x30, [sp, #-16]!

	mov x8, #63
	svc #0

	ldp x29, x30, [sp], #16
        ret

//x0 file descriptor
//RETRUN x0 error code
sys_close:
  	stp x29, x30, [sp, #-16]!

        mov x8, #57             //close syscall
        svc #0

        ldp x29, x30, [sp], #16
        ret

//x1 contains pointer to string
//x2 is the length of the target string
//x0 contains the nunber of byte sread
sys_input:
	stp x29, x30, [sp, #-16]!
	stp x2, x8, [sp, #-16]!

	mov x0, #0
	mov x8, #63
	svc #0
	
	subs x0, x0, #1
	ldrb w8, [x1, x0]
	cmp w8, #'\n'
	beq sys_input_clean
	add x0, x0, #1
	b sys_input_exit

sys_input_clean:
	mov w8, #0
	strb w8, [x1, x0]

sys_input_exit:
	ldp x2, x8, [sp], #16
	ldp x29, x30, [sp], #16
	ret
	
.data

