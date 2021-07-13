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

	m_readFile x11, file, file_len

	m_print file, x0

	m_closeFile x0

	mov x0, #0
	bl exit

.align 8
.data
prompt:		.ascii "enter file name: "
prompt_len = . - prompt

file_name:	.fill 255,1,0
file_name_len = . - file_name

file:	.fill 1024,1,0
file_len = . - file
