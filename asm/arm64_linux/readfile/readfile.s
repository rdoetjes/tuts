.include "system.s"
.include "fileio.s"
.include "print.s"

.text
.align 8
.global _start

_start:
	m_printString prompt

	m_input file_name, file_name_len

	m_openFile file_name, O_RDONLY
	mov x11, x0		//keep copy of file desc
	mov x12, 0		//total bytes read
	cmp x0, #0
	blt open_file_error

more_data:
	m_readFile x11, file, file_len	
	add x12, x12, x0
	cmp  x0, #0
	beq close_file
	m_print file, x0
	b more_data

open_file_error:
	m_printString error

close_file:
	m_printNr x12
	m_printString bytes_read

	m_closeFile x0

	m_exit 0

.align 8
.data
prompt:		.asciz "enter file name: "

error:		.asciz "Couldn't open file!\n"

bytes_read:	.asciz  " bytes read\n"

file_name:	.fill 255,1,0
file_name_len = . - file_name

file:		.fill 255,1,0
file_len = . - file
