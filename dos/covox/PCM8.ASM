.model tiny

.data

.code
	org 100h

main proc
	; Install our own timer interrupt vector
	; we don't care about the old timer vector, since we will reboot
	; it is a demo program after all
	push ds
	mov ax,2508h
	push cs
	pop ds
	mov dx,OFFSET sound
	int 21h
	pop ds

	mov al, 036h
	out 043h, al	;tell PIT we are using channel 0

	;closest to 12000hz 1193182/12000 should be 100, but 80 works better
	mov ax, 80
	out 040h, al
	xor ax, ax
	out 040h, al
	mov dx, 03BCh
	mov di, 0h
_wait:
	jmp _wait

sound:
	push ax
	push bx

	mov al, byte ptr sample[di]
	sub al, 80h
	out dx, al

	cmp di, (endsample-sample)
	jz resetDi

	inc di
	jmp _iexit
resetDi:
	mov di, 0h
_iexit:
	mov al, 020h
	out 020h, al

	pop bx
	pop ax
	iret
main endp

sample:
	include hello3.sam
endsample:

end main
