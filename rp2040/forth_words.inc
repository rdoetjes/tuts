
defcode "2DUP", 4
    ENTER
bl "OVER"
bl "OVER"
    EXIT

defcode "TEST-WORD", 9
    ENTER
    LIT 12
    LIT 13
bl "*"
    EXIT

defcode "BLINKY", 6
    ENTER
bl "BEGIN"
    LIT 300
bl "MS_SLEEP"
    LIT 25
    LIT 0
bl "SET_PIN_TO"
    LIT 300
bl "MS_SLEEP"
    LIT 25
    LIT 1
bl "SET_PIN_TO"
bl "UNTIL"
    EXIT
