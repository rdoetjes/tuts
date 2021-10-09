.DATA
	PRN dw 00h
;
;	@Desc: Seed the LGC pseudo random generator
;       @Input Arguments: NONE
;	@Return: NONE
;	@Relies on: Variable PRN word
;
SeedRnd PROC
	xor ax, ax	;clear AX
	int 1ah		;get timer ticks IN CX:DX
	add [PRN], dx   ;we are only interrested in the lower bytes
	ret
SeedRnd ENDP

;
;	@Desc: Generate random number
;       @Input Arguments: NONE
;	@Return: DX and AX contain random numner
;	@Relies on: Variable PRN word
;
Random PROC
	mov ax, 25173	;LCG multiplier
	mul word ptr [PRN]	;DX:AX = LCG multiplr * seed
	add ax, 13849	;add lcg increment value
	mov [PRN], ax
	ret
Random ENDP

;
;	@Desc: Return back to 80x25 16 colors text mode
; 	@Input Arguments: NONE
;	@Return: NONE
;
TextMode PROC
	mov ax, 03h	;video mode 25x80 16 colours
	int 10h         ;set it
	ret
TextMode ENDP

;
;	@Desc: Check if key is pressed, if key pressed exit the program
;	@Input Arguments: NONE
;	@Return: NONE
;
KeyPressExit PROC
	mov ah,06h      ;check if keyboard was pressed
	mov dl,0ffh     ;
	int 21h         ;check the keypress (00 is pressed)
	cmp al,00h      ;not keypressed then goto mainloop else exit
	jnz EXIT
	ret
EXIT:
	call TextMode
	mov ax, 4c00h
	int 21h
	ret
KeyPressExit ENDP
