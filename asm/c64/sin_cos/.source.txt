  BasicUpstart2(start)

start:
    lda #$93           // Clear screen
    jsr $FFD2

    // Fill screen with PETSCII full block (char 160)
    lda #46
    ldx #0
fill_screen:
    sta $0400, x       // Screen RAM (first 256 bytes)
    sta $0500, x       // Screen RAM (next 256)
    sta $0600, x       // Screen RAM (next 256)
    sta $06E8, x       // Screen RAM (last 232 bytes)
    inx
    bne fill_screen

main_loop:
    inc $d020
    ldx #00
do_plasma:
    lda frame_count
    adc row_count
    tay
    lda sine_table_16, y
    sta $d800,x

    iny
    lda sine_table_16, y
    sta $d900,x

    sta $da00,x
    iny
    lda sine_table_16, y

    sta $db00,x
    iny
    lda sine_table_16, y

    stx row_count
    inx
    cpx #$00
    beq exit_plasma
    dec $d020

exit_plasma:

wait_for_vblank:
    lda $d012
    cmp #$ff
    bne wait_for_vblank
    inc frame_count
    jmp main_loop      // Infinite loop


sine_table_16:
    .byte $8, $9, $A, $B, $C, $D, $D, $E, $E, $F, $F, $F, $F, $F, $F, $F
    .byte $E, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9, $9, $8, $8, $7
    .byte $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1, $1, $0, $0, $F
    .byte $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9, $9, $8, $8
    .byte $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1, $1, $0, $0
    .byte $F, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9, $9, $8
    .byte $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1, $1, $0
    .byte $0, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9, $9, $8
    .byte $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1, $1, $0
    .byte $0, $F, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9, $9
    .byte $8, $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1, $1
    .byte $0, $0, $F, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A, $9
    .byte $9, $8, $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2, $1
    .byte $1, $0, $0, $F, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A, $A
    .byte $9, $9, $8, $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2, $2
    .byte $1, $1, $0, $0, $F, $F, $F, $E, $E, $D, $D, $C, $C, $B, $B, $A
    .byte $A, $9, $9, $8, $8, $7, $7, $6, $6, $5, $5, $4, $4, $3, $3, $2
    .byte $2, $1, $1, $0

sine_table_255:
    .byte $80, $82, $85, $87, $8A, $8C, $8F, $91, $93, $96, $98, $9A, $9D, $9F, $A1, $A3, $A5
    .byte $A7, $A9, $AB, $AD, $AF, $B1, $B2, $B4, $B6, $B7, $B9, $BA, $BC, $BD, $BF, $C0
    .byte $C2, $C3, $C4, $C5, $C7, $C8, $C9, $CA, $CB, $CC, $CD, $CE, $CF, $D0, $D1, $D2
    .byte $D3, $D4, $D5, $D6, $D7, $D8, $D9, $DA, $DB, $DC, $DD, $DE, $DF, $E0, $E1, $E2
    .byte $E3, $E4, $E5, $E6, $E7, $E8, $E9, $EA, $EB, $EC, $ED, $EE, $EF, $F0, $F1, $F2
    .byte $F3, $F4, $F5, $F6, $F7, $F8, $F9, $FA, $FB, $FC, $FD, $FE, $FF, $FF, $FE, $FD
    .byte $FC, $FB, $FA, $F9, $F8, $F7, $F6, $F5, $F4, $F3, $F2, $F1, $F0, $EF, $EE, $ED
    .byte $EC, $EB, $EA, $E9, $E8, $E7, $E6, $E5, $E4, $E3, $E2, $E1, $E0, $DF, $DE, $DD
    .byte $DC, $DB, $DA, $D9, $D8, $D7, $D6, $D5, $D4, $D3, $D2, $D1, $D0, $CF, $CE, $CD
    .byte $CC, $CB, $CA, $C9, $C8, $C7, $C6, $C5, $C4, $C3, $C2, $C1, $C0, $BF, $BE, $BD
    .byte $BC, $BB, $BA, $B9, $B8, $B7, $B6, $B5, $B4, $B3, $B2, $B1, $B0, $AF, $AE, $AD
    .byte $AC, $AB, $AA, $A9, $A8, $A7, $A6, $A5, $A4, $A3, $A2, $A1, $A0, $9F, $9E, $9D
    .byte $9C, $9B, $9A, $99, $98, $97, $96, $95, $94, $93, $92, $91, $90, $8F, $8E, $8D
    .byte $8C, $8B, $8A, $89, $88, $87, $86, $85, $84, $83, $82, $81, $80, $7F, $7E, $7D
    .byte $7C, $7B, $7A, $79, $78, $77, $76, $75, $74, $73, $72, $71, $70, $6F, $6E, $6D
    .byte $6C, $6B, $6A, $69, $68, $67, $66, $65, $64, $63, $62, $61, $60, $5F, $5E, $5D
    .byte $5C, $5B, $5A, $59, $58, $57, $56, $55, $54, $53, $52, $51, $50, $4F, $4E, $4D
    .byte $4C, $4B, $4A, $49, $48, $47, $46, $45, $44, $43, $42, $41, $40, $3F, $3E, $3D
    .byte $3C, $3B, $3A, $39, $38, $37, $36, $35, $34, $33, $32, $31, $30, $2F, $2E, $2D
    .byte $2C, $2B, $2A, $29, $28, $27, $26, $25, $24, $23, $22, $21, $20, $1F, $1E, $1D
    .byte $1C, $1B, $1A, $19, $18, $17, $16, $15, $14, $13, $12, $11, $10, $0F, $0E, $0D
    .byte $0C, $0B, $0A, $09, $08, $07, $06, $05, $04, $03, $02, $01, $00

frame_count:   .byte 0   // Animation frame counter
framee_count_2: .byte 0   // Animation frame counter 2
row_count:     .byte 0   // Row counter
