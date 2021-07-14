//This macro is more or less an interface to the actual functions
//allowing to have an abstraction between the (sys)call and the
//user

.equ  O_RDONLY, 0
.equ  O_WRONLY, 1
.equ  O_CREAT,  0100
.equ  O_EXCL,   0200
.equ  S_RDWR,   0666

.macro  m_openFile filename, flags
	ldr x1, =\filename
	mov x2, #\flags
	mov x3, #S_RDWR
	bl sys_open
.endm

.macro m_readFile fd, buffer, len	
	mov x0, \fd
	ldr x1, =\buffer
	ldr x2, =\len
	bl sys_read
.endm

.macro m_print string, len
	ldr x1, =\string
	mov x2, \len
	bl print
.endm

.macro m_input buffer, len
	ldr x1, =\buffer
	ldr x2, =\len
	bl sys_input
.endm

.macro m_closeFile fd
	mov x0, \fd
	bl sys_close
.endm
