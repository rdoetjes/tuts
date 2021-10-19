.model tiny

.code
	org 100h
main proc
	mov dx, 3BCh
	mov al, 0h
sound:
	out dx, al
	inc al
	;call keyPressExit
	jmp sound

	include stdfunc.asm
main endp
end main