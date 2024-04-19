* = $1000

loop:
    clc
    lda bytes+3
    adc #1
    sta bytes+3
    bcc loop
    lda bytes+2
    clc
    adc #1
    sta bytes+2
    bcc loop
    lda bytes+1
    clc
    adc #1
    sta bytes+1
    bcc loop
    lda bytes+0
    clc
    adc #1
    sta bytes+0
    bcc loop

done:
    rts

bytes:
    .byte $0, $0 ,$0, $0
