.MODEL tiny

.STACK 100h

PALETTE_INDEX EQU 3c8h
PALETTE_DATA EQU 3c9h
INPUT_STATUS EQU 3dah
HRETRACE EQU 01h
VRETRACE EQU 08h

.DATA
	bpos	db 00
	sine	db 128,131,134,137,140,143,146,149
		db 152,156,159,162,165,168,171,174
		db 176,179,182,185,188,191,193,196
		db 199,201,204,206,209,211,213,216
		db 218,220,222,224,226,228,230,232
		db 234,235,237,239,240,242,243,244
		db 246,247,248,249,250,251,251,252
		db 253,253,254,254,254,255,255,255
		db 255,255,255,255,254,254,253,253
		db 252,252,251,250,249,248,247,246
		db 245,244,242,241,239,238,236,235
		db 233,231,229,227,225,223,221,219
		db 217,215,212,210,207,205,202,200
		db 197,195,192,189,186,184,181,178
		db 175,172,169,166,163,160,157,154
		db 151,148,145,142,138,135,132,129
		db 126,123,120,117,113,110,107,104
		db 101,98,95,92,89,86,83,80
		db 77,74,71,69,66,63,60,58
		db 55,53,50,48,45,43,40,38
		db 36,34,32,30,28,26,24,22
		db 20,19,17,16,14,13,11,10
		db 9,8,7,6,5,4,3,3
		db 2,2,1,1,0,0,0,0
		db 0,0,0,1,1,1,2,2
		db 3,4,4,5,6,7,8,9
		db 11,12,13,15,16,18,20,21
		db 23,25,27,29,31,33,35,37
		db 39,42,44,46,49,51,54,56
		db 59,62,64,67,70,73,76,79
		db 81,84,87,90,93,96,99,103
		db 106,109,112,115,118,121,124,128
	endSine	db 00
	cop1Pos dw 00
	cop2Pos	dw 00
	cop3Pos	dw 00
	r	db 00
	g	db 00
	b	db 00

.CODE
	;used for keypress exit
	include STDFUNC.ASM

MAIN PROC
	;initialize the program
	mov ax,SEG _DATA
	mov ds, ax

	mov ax, 13h
	int 10h

_loop:
	;sine table offset pointer
	mov di, 00h
__loop:
	;read the sine from the sine list
	mov cx, word ptr [sine+di]  ;set bar 2 to offset 0 in the sine lost
	mov [cop2Pos], cx           ;store the bar begin pos in cop2Pos

	mov si, di
	add si, 20          ;set bar 1 to offset 20 in the sine list
	and si, 00ffh       ;keep only lowest 8 bits (bar pos is 0-255)
	mov cx, word ptr [sine+si] ;read the bar begin pos from the sine list
	mov [cop1Pos], cx   ;store the bar begin pos in var cop1Pos

	add si, 20          ;set bar 3 to offset 20 from cop1Pos in sine list
	and si, 00ffh       ;keep only the lowest 8 bits (bar pos is 0-255)
	mov cx, word ptr [sine+si] ;store the bar beging pos from sine list
	mov [cop3Pos], cx   ;stire the bar begin pos in var cop3Pos

	;wait for the vertical retrace, then we can redraw the next frame
	call vertRetrace

	mov cx, 0h
_lines:
	mov [r], 0
	mov [g], 0
	mov [b], 0

	mov bx, [cop1Pos]
	call drawBar
	mov [r], al

	mov bx,[cop2Pos]
	call drawBar
	mov [b], al

	mov bx, [cop3Pos]
	call drawBar
	mov [g], al

	;set the color as long as we are on the current line
	call setBarColor

	inc cx	;increment line count
	cmp cx, 350
	jnz _lines
	sti

	;if key pressed then exit
	call KeyPressExit

	;Increment DI to point to the next entry in the sine list
	inc di
	cmp di, (endSine-sine)
	jnz _h3
	jmp _loop
_h3:
	jmp __loop
MAIN ENDP

;
; @Desc:	Wait for a vertical retrace of the electron beam
; @Input:	None
;
; @Globbers: 	dx, al
; @Returns:	None

vertRetrace PROC
	;wait for a vertical retrace
	cli
	mov dx, INPUT_STATUS

_v1:	in al, dx
	and al, VRETRACE
	jnz _v1

_v2:	in al, dx
	and al, VRETRACE
	jz _v2
	sti

	ret
vertRetrace ENDP

;
; @Desc:	Wait for a horizontal retrace of the electron beam
; @Input:	None
;
; @Globbers: 	dx, al
; @Returns:	None
horizRetrace PROC
	;wait for horizontal retrace
	cli
	mov dx, INPUT_STATUS

_h1:	in al, dx
	and al, HRETRACE
	jnz _h1

_h2:	in al, dx
	and al, HRETRACE
	jz _h2

	ret
horizRetrace ENDP

;
; @Desc:	Draw a shaded copper bar
; @Input:	BX start rasterline of the copper bar
;    		CX current rasterline
; @Globbers: 	AX
; @Returns:	Returns the shade of the bar in AL
drawBar PROC
	and bx, 00ffh		;only keep low byte (sine goes from 0-255)
	cmp cx, bx		;is rasterline somewhere between our bar?
	jae bar         	;yes? draw the bar line
	mov al,00h              ;no? draw the black background
	ret
  bar:
	add bx, 30		;draw a bar 30 lines wide
	cmp cx, bx
	jg colOff               ;is scanline beyond the bar, then draw black
	jz barLastLine          ;last line has a brighter tint but not ffh
				;suggesting light from below the bar
	mov al, 0e0h
	sub bx,20               ;top shade brighter shade of the bar
	cmp cx, bx
	push cx
	jg subC
	sub cx, bx              ;shade darker based of width
	add al, cl              ;1byte move result to cl
	jmp short dcd
  subC:
	sub cx, bx
	sub al, cl
  dcd:
	pop cx                  ;restore line counter
	ret
  barLastLine:
	mov al,0d0h
	ret
  colOff:
	mov al,00h
	ret
drawBar ENDP

setBarColor PROC
	;now set the color of the scan line to black
	mov dx, PALETTE_INDEX
	mov al, 00h
	out dx, al

	;this is the current scanline so set the color not to black
	mov al, byte ptr [r]
	mov dx, PALETTE_DATA
	out dx, al

	mov al, byte ptr [g]
	mov dx, PALETTE_DATA
	out dx, al

	mov al, byte ptr [b]
	mov dx, PALETTE_DATA
	out dx, al
	ret

SetBarColor ENDP
END MAIN
