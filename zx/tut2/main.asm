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


; SCROLL THE BUFFER WITH LETTER 
buffer_roll_print:
        ld iy, buf
        ld hl, APOS + 31
        ld de, 32
        ld b,8
        ld a,(color)
        ld c, a
_s4:     
        ld a,0
        sla (iy+0)
        jr nc,_s5
        ld a, c
_s5:     
        ld (hl), a
        add hl,de
        ld (hl), a
        add hl,de
        inc iy
        djnz _s4
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
_p1:     
        ld a, (hl)
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
        djnz _p1
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
