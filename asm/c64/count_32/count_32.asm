* = $1000

main:
    lda #0
    sta bytes+1

loop:
    clc
    lda bytes+3
    adc #1
    sta bytes+3
    bcc loop
    clc
    lda bytes+2
    adc #1
    sta bytes+2
    bcc loop
    clc
    lda bytes+1
    adc #1
    sta bytes+1
    bcc loop
    clc
    lda bytes+0
    adc #1
    sta bytes+0
    bcc loop

done:
    rts

bytes:
    .byte $0, $0 ,$0, $0