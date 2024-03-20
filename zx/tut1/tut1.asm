org $8000

    call CLS
    ld hl, HELLOWORLD
    call PRINTS

    halt   

;print string pointed to by hl
PRINTS  
    push af
_NEXT
    ld a, (hl)
    cp 0
    jr z, _EXIT
    rst $10
    inc hl
    jp _NEXT
_EXIT
    pop af
    ret

; Clear the screen to black.
CLS
    push af
    ld a,71         ; White ink (7), black paper (0), bright (64).
    ld (23693),a    ; Set our screen colours.
    xor a           ; Load accumulator with zero.
    call 8859       ; Set permanent border colours.
    call 3503       ; Clear the screen, open channel 2.
    pop af
    ret

HELLOWORLD: 
    defb "Hello, world"
    defb 13, 0

end $8000
