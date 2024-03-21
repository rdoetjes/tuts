        org $8000

; CONSTANTS
SCR     equ $4000
SCRSZ   equ $1800 
ATTR    equ SCR + SCRSZ
ATTRSZ  equ 32*24
APOS    equ ATTR + (24-16)/2*32

start:  call black_screen

l03:    ld ix, text
l02:    ld a,(ix+0) 
        or a
        jr z, l03
        inc ix

        call print_char

        ld b,8    
l01:    push bc

        halt          ; wait till next irq
        halt          ; wait till next irq

        call screen_scroll
        call buffer_roll_print

        pop bc
        djnz l01

        jr l02

; INCHL - function to calculate the address of the byte in the 
;         next line on the ZX screen
; 
; ZX screen addressing explained:
; +-----------------------------------------------+
; |     HIGH BYTE       |   |   LOW BYTE          |
; 0F 0E 0D 0C 0B 0A 09 08   07 06 05 04 03 02 01 00
; -------------------------------------------------
;  0  1  0 P1 P0 L2 L1 L0   Y2 Y1 Y0 X4 X3 X2 X1 X0
;          |---| |------|   |------| |------------|
;            |       |          |           |
;            |       |          |           +-----| X0..X4 : BYTE IN THE CURRENT LINE 0..31 
;            |       |          +-----------------| Y0..Y2 : CHARACTER POSITION IN THE SCREEN PART 0..7
;            |       +----------------------------| L0..L2 : LINE IN THE CHARACTER 0..7
;            +------------------------------------| P0/P1  : PART OF THE SCREEN
inchl:  inc h
        ld a,h
        and 7
        ret nz
        ld a,h
        sub 8
        ld h,a
        ld a,l
        add a,32
        ld l,a
        ret nc
        ld a,h
        add a,8
        ld h,a
        ret

; SCROLL THE BUFFER WITH LETTER 
buffer_roll_print:
        ld iy, buf
        ld hl, APOS + 31
        ld de, 32
        ld b,8
        ld a,(color)
        ld c, a
s4:     ld a,0
        sla (iy+0)
        jr nc,s5
        ld a, c
s5:     ld (hl), a
        add hl,de
        ld (hl), a
        add hl,de
        inc iy
        djnz s4
        ret

; SCROLL ATTRIBUTES MEMORY 
screen_scroll:
        ld hl, APOS + 1
        ld de, APOS
        ld bc, 16*32-1
        ldir

        ret

print_char:
        ld h, 0
        ld l, a
        add hl, hl      ; calculate charter position in the default char list
        add hl, hl
        add hl, hl
        ld de, 15360    ; default CHARACTER start
        add hl, de      ; calculate offset into charter list
        ld de, buf
        ld b,8
p1:     ld a, (hl)
        push bc
        ld c,a
        srl a
        or c
        srl a
        or c
        ld (de),a
        inc hl
        inc de
        pop bc
        djnz p1
        ret

black_screen:
        ; clear the screen
        ld hl, SCR
        ld de, SCR + 1
        ld bc, SCRSZ - 1
        ld (hl), l
        ldir

        ; clear attributes memory
        ld hl, ATTR
        ld de, ATTR + 1
        ld bc, ATTRSZ - 1
        ld a, %00000111
        ld (hl), a 
        ldir

        ; set border to black
        xor a
        out ($fe), a

        ret

color:  defb %00100001
buf:    defs 8

text:   defb "RETRO IS THE NEW BLACK, BY WOLFGANG KIERDORF!!!    ",0

        end $8000
; vim: ft=z8a
