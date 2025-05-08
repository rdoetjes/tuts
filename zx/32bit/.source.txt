org $8000

  ld hl, #$0
  ld de, #$0

loop:
  inc hl
  ld a, h
  or l
  jr nz, loop
  inc de
  ld a, d
  or e
  jr nz, loop
  ret

end $8000

