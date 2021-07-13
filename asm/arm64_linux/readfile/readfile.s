.include "common.s"
.include "fileio.s"

.text
.align 8
.global _start

_start:
	m_print prompt, prompt_len

	m_input file_name, file_name_len

	m_openFile file_name, O_RDONLY
	mov x11, x0		//keep copy of file desc
	cmp x0, #0
	blt open_file_error

more_data:
	m_readFile x11, file, file_len
	cmp  x0, #0
	beq close_file
	m_print file, x0
	b more_data

open_file_error:
	m_print error, error_len

close_file:
	m_closeFile x0

	mov x0, #0
	bl exit
.align 8
.data
prompt:		.ascii "enter file name: "
prompt_len = . - prompt

error:		.ascii "Couldn't open file!\n"
error_len = . - error

file_name:	.fill 255,1,0
file_name_len = . - file_name

file:		.fill 255,1,0
file_len = . - file
