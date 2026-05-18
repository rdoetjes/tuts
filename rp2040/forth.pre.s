@ =============================================================================
@ RP2040 Forth - Subroutine Threaded (STC)
@ =============================================================================
@ This is a native Forth implementation for the RP2040 (Cortex-M0+).
@ Unlike traditional Forth VMs, this system compiles words into native ARM
@ machine code subroutines. Every Forth word is a real ARM function.
@ =============================================================================

.include "rp2040.inc"
.syntax unified
.cpu cortex-m0plus
.thumb

@ -----------------------------------------------------------------------------
@ DICTIONARY MANAGEMENT & STORAGE
@ -----------------------------------------------------------------------------
@ The dictionary is a linked list of "word headers" stored in two locations:
@ 1. Built-in words: Stored in Flash (read-only), defined during assembly.
@ 2. User words: Stored in RAM (starting at _forth_ram_start), created at runtime.
@
@ Each word header follows this structure:
@ [Link (4 bytes)] -> Points to the previous word's header.
@ [Flags (1 byte)] -> Metadata (e.g. if the word is 'Immediate').
@ [Len (1 byte)]   -> Number of characters in the word's name.
@ [Name (N bytes)] -> The actual name string (e.g., "DUP").
@ [Code (Native)]  -> The ARM Thumb instructions for this word.
@
@ Both words (colon definitions) and variables use this same header structure,
@ making them both searchable by the interpreter.

.set _link, 0  @ Tracks the head of the dictionary during assembly

@ Define a new Forth word header
@ Usage: defcode "NAME", NAME_LENGTH, FLAGS
.macro defcode name, len, flags=0
    .align 2
_L\@:
    .word _link            @ Link to previous word
    .set _link, _L\@       @ Update link head to this word
    .byte \flags, \len     @ Metadata
    .ascii "\name"         @ Name string
    .align 2
    .thumb_func
"\name":                   @ The assembly label is the Forth name in quotes
.endm

@ -----------------------------------------------------------------------------
@ ASSEMBLY HELPERS
@ -----------------------------------------------------------------------------

@ boilerplate for starting and ending a word that calls other words
.macro ENTER
    push {lr}
.endm

.macro EXIT
    pop {pc}
.endm

@ Macro for literal numbers (renamed from PUSH to avoid conflict)
.macro LIT val
    ldr r0, =\val
    subs r5, #4
    str r4, [r5]
    mov r4, r0
.endm

@ Quick calls for common words to keep assembly readable
.macro DUP
    bl "DUP"
.endm

.macro DROP
    bl "DROP"
.endm

.macro SWAP
    bl "SWAP"
.endm

.macro OVER
    bl "OVER"
.endm

.macro SLASH
    bl "/"
.endm

@ -----------------------------------------------------------------------------
@ HARDWARE VECTORS & RESET
@ -----------------------------------------------------------------------------
.section .vectors, "ax"
    .word _stack_top       @ Initial SP (Return Stack)
    .word _reset + 1       @ Initial PC (Entry Point)

.section .text
.global _reset
.thumb_func
_reset:
    @ Detect if we are on real RP2040 (Cortex-M0+) or QEMU (Cortex-M3)
    ldr r0, =0xe000ed00    @ CPUID base
    ldr r1, [r0]
    ldr r2, =0x410cc600    @ ARM, Cortex-M0+ (RP2040)
    ldr r3, =0xffffff00
    ands r1, r3
    cmp r1, r2
    bne skip_hw_init       @ Skip GPIO init if in QEMU (M3)
    bl _gpio_init
skip_hw_init:
    ldr r0, =_sidata
    ldr r1, =_sdata
    ldr r2, =_edata
1:  cmp r1, r2
    bhs 2f
    ldr r3, [r0]
    adds r0, #4
    str r3, [r1]
    adds r1, #4
    b 1b

2:  ldr r1, =_sbss
    ldr r2, =_ebss
    movs r3, #0
3:  cmp r1, r2
    bhs 4f
    str r3, [r1]
    adds r1, #4
    b 3b

4:  ldr r0, =CMSDK_UART0_BASE
    movs r1, #3
    str r1, [r0, #CMSDK_UART_CTRL_OFFSET]

    ldr r5, =_data_stack_base
    movs r4, #0

    ldr r0, =_latest
    ldr r1, =_builtins_latest
    ldr r1, [r1]
    str r1, [r0]

    ldr r0, =_here
    ldr r1, =_forth_ram_start
    str r1, [r0]

    ldr r0, =QUIT
    bx r0

@ -----------------------------------------------------------------------------
@ SERIAL I/O SUBROUTINES
@ -----------------------------------------------------------------------------

@ Write one character in r0 to the serial port
.thumb_func
_emit:
    push {r1, r2, r3, lr}
    ldr r1, =CMSDK_UART0_BASE
e1: ldr r2, [r1, #CMSDK_UART_STATE_OFFSET]
    movs r3, #1
    tst r2, r3
    bne e1
    str r0, [r1, #CMSDK_UART_DATA_OFFSET]
    pop {r1, r2, r3, pc}

.thumb_func
_key:
    push {r1, r2, lr}
    ldr r1, =CMSDK_UART0_BASE
k1: ldr r2, [r1, #CMSDK_UART_STATE_OFFSET]
    movs r3, #2
    tst r2, r3
    beq k1
    ldr r0, [r1, #CMSDK_UART_DATA_OFFSET]
    movs r1, #255
    ands r0, r1
    pop {r1, r2, pc}

@ -----------------------------------------------------------------------------
@ NUMERIC OUTPUT
@ -----------------------------------------------------------------------------

@ Print the number in r0 as decimal text
.thumb_func
_print_dec:
    push {r4, r5, r6, lr}
    mov r4, r0
    ldr r6, =powers_of_10
    movs r5, #0
    cmp r4, #0
    bge pd1
    mov r0, r4
    rsbs r0, r0, #0
    mov r4, r0
    movs r0, #45
    bl _emit
pd1:ldr r2, [r6]
    cmp r2, #0
    beq pd5
    movs r0, #0
pd2:cmp r4, r2
    blo pd3
    subs r4, r2
    adds r0, #1
    b pd2
pd3:cmp r0, #0
    bne pd4
    cmp r5, #0
    bne pd4
    cmp r2, #1
    beq pd4
    b pd6
pd4:movs r5, #1
    adds r0, #48
    bl _emit
pd6:adds r6, #4
    b pd1
pd5:pop {r4, r5, r6, pc}

.align 2
powers_of_10: .word 1000000000, 100000000, 10000000, 1000000, 100000, 10000, 1000, 100, 10, 1, 0

@ -----------------------------------------------------------------------------
@ FORTH PRIMITIVES (Core assembly words)
@ -----------------------------------------------------------------------------

defcode "DUP", 3
    subs r5, #4
    str r4, [r5]
    bx lr

defcode "DROP", 4
    ldr r4, [r5]
    adds r5, #4
    bx lr

defcode "SWAP", 4
    ldr r0, [r5]
    str r4, [r5]
    mov r4, r0
    bx lr

defcode "OVER", 4
    ldr r0, [r5]
    subs r5, #4
    str r4, [r5]
    mov r4, r0
    bx lr

defcode "+", 1
    ldr r0, [r5]
    adds r4, r4, r0
    adds r5, #4
    bx lr

defcode "-", 1
    ldr r0, [r5]
    subs r4, r0, r4
    adds r5, #4
    bx lr

defcode "*", 1
    ldr r0, [r5]
    muls r4, r0, r4
    adds r5, #4
    bx lr

defcode "/", 1
    push {lr}
    ldr r0, [r5]
    mov r1, r4
    bl _idiv
    mov r4, r0
    adds r5, #4
    pop {pc}

@ -----------------------------------------------------------------------------
@ HIGH-LEVEL DEFINITIONS
@ -----------------------------------------------------------------------------

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

@ -----------------------------------------------------------------------------
@ MEMORY AND LOGIC PRIMITIVES
@ -----------------------------------------------------------------------------

defcode "2DROP", 5
    ldr r4, [r5, #4]
    adds r5, #8
    bx lr

defcode "=", 1
    ldr r0, [r5]
    adds r5, #4
    cmp r0, r4
    bne 1f
    movs r4, #0
    mvns r4, r4
    bx lr
1:  movs r4, #0
    bx lr

defcode "<", 1
    ldr r0, [r5]
    adds r5, #4
    cmp r0, r4
    bge 1f
    movs r4, #0
    mvns r4, r4
    bx lr
1:  movs r4, #0
    bx lr

@ --- Word: SET_PIN_TO ( pin state -- ) ---
@ Sets a GPIO pin to High (1) or Low (0).
@ This handles the SIO bitmasks and Output Enable automatically.
defcode "SET_PIN_TO", 10
    ldr r0, [r5]           @ Get pin number from stack
    movs r1, #1
    lsls r1, r1, r0        @ r1 = bitmask (1 << pin)

    ldr r2, =SIO_BASE
    @ Ensure the pin is enabled for output
    str r1, [r2, #SIO_GPIO_OE_SET_OFFSET]

    @ Check if we want to set it High or Low
    cmp r4, #0
    beq 1f                 @ If state is 0, jump to clear logic
    @ Set pin High
    str r1, [r2, #SIO_GPIO_OUT_SET_OFFSET]
    b 2f
1:  @ Set pin Low
    str r1, [r2, #SIO_GPIO_OUT_CLR_OFFSET]
2:  ldr r4, [r5, #4]       @ Drop pin and state from stack, get new TOS
    adds r5, #8
    bx lr

@ --- Word: MS_SLEEP ( ms -- ) ---
@ Wait for the specified number of milliseconds using the hardware timer.
defcode "MS_SLEEP", 8
    push {lr}
    @ QEMU safety check
    ldr r0, =0xe000ed00
    ldr r1, [r0]
    ldr r2, =0x410cc600    @ M0+
    ldr r3, =0xffffff00
    ands r1, r3
    cmp r1, r2
    bne 2f                 @ Skip if in QEMU (no timer)

    mov r0, r4             @ Get ms from TOS
    cmp r0, #0
    ble 2f                 @ Return immediately if ms <= 0
    ldr r1, =1000
    muls r0, r1, r0        @ Convert ms to microseconds
    ldr r2, =TIMER_BASE
    ldr r1, [r2, #TIMER_TIMELW_OFFSET] @ Get start time (low 32 bits)
1:  ldr r3, [r2, #TIMER_TIMELW_OFFSET]
    subs r3, r3, r1        @ Calculate elapsed time (handles 32-bit wrap)
    cmp r3, r0
    blo 1b                 @ Loop until elapsed >= target
2:  ldr r4, [r5]           @ Drop ms from stack, get new TOS
    adds r5, #4
    pop {pc}

@ --- Word: GET_MILIS ( -- ms ) ---
@ Push the number of milliseconds since boot onto the stack.
defcode "GET_MILIS", 9
    push {lr}
    subs r5, #4
    str r4, [r5]           @ Push current TOS to RAM

    @ QEMU safety check
    ldr r0, =0xe000ed00
    ldr r1, [r0]
    ldr r2, =0x410cc600    @ M0+
    ldr r3, =0xffffff00
    ands r1, r3
    cmp r1, r2
    beq 1f
    movs r4, #0            @ Return 0 in QEMU
    pop {pc}

1:  ldr r2, =TIMER_BASE
    ldr r0, [r2, #TIMER_TIMELW_OFFSET] @ Get current microseconds
    ldr r1, =1000
    bl _udiv               @ Divide by 1000 to get milliseconds
    mov r4, r0             @ Put result in TOS
    pop {pc}

defcode "@", 1
    ldr r4, [r4]
    bx lr

defcode "C@", 2
    ldrb r4, [r4]
    bx lr

defcode "!", 1
    ldr r0, [r5]
    str r0, [r4]
    ldr r4, [r5, #4]
    adds r5, #8
    bx lr

defcode "C!", 2
    ldr r0, [r5]
    strb r0, [r4]
    ldr r4, [r5, #4]
    adds r5, #8
    bx lr

defcode ".", 1
    push {lr}
    mov r0, r4
    bl _print_dec
    movs r0, #32
    bl _emit
    ldr r4, [r5]
    adds r5, #4
    pop {pc}

@ Display the current stack state
defcode ".S", 2
    push {r4, r5, r6, r7, lr}
    movs r0, #60
    bl _emit
    ldr r6, =_data_stack_base
    subs r0, r6, r5
    lsrs r0, r0, #2
    bl _print_dec
    movs r0, #62
    bl _emit
    movs r0, #32
    bl _emit
    ldr r6, =_data_stack_base
    subs r6, #8
ds1:cmp r6, r5
    blo ds2
    ldr r0, [r6]
    bl _print_dec
    movs r0, #32
    bl _emit
    subs r6, #4
    b ds1
ds2:ldr r6, =_data_stack_base
    cmp r5, r6
    beq ds3
    mov r0, r4
    bl _print_dec
    movs r0, #32
    bl _emit
ds3:pop {r4, r5, r6, r7, pc}

@ List all words in the dictionary
defcode "WORDS", 5
    push {r4, lr}
    ldr r2, =_latest
    ldr r2, [r2]
w1: cmp r2, #0
    beq w3
    ldrb r1, [r2, #5]
    adds r3, r2, #6
w2: ldrb r0, [r3]
    bl _emit
    adds r3, #1
    subs r1, #1
    bne w2
    movs r0, #32
    bl _emit
    ldr r2, [r2]
    b w1
w3: pop {r4, pc}

.align 2
.ltorg

@ -----------------------------------------------------------------------------
@ COMPILER WORDS
@ -----------------------------------------------------------------------------

@ Start compiling a new word: : NAME ... ;
@ --- Word: \ ---
@ Immediate word that skips the rest of the line.
defcode "\\", 1, 1
1:  bl _key
    cmp r0, #10
    beq 2f
    cmp r0, #13
    beq 2f
    cmp r0, #0
    beq 2f
    b 1b
2:  bx lr

defcode ":", 1
    push {lr}
    bl _getw
    ldr r2, =_latest
    ldr r3, =_here
    ldr r1, [r3]
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    ldr r7, [r2]
    str r7, [r1]
    str r1, [r2]
    adds r1, #4
    movs r7, #0
    strb r7, [r1]
    strb r0, [r1, #1]
    adds r1, #2
    ldr r6, =_tib
    movs r7, #0
1:  ldrb r2, [r6, r7]
    strb r2, [r1, r7]
    adds r7, #1
    cmp r7, r0
    blt 1b
    add r1, r0
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    ldr r6, =0x46c0b500
    str r6, [r1]
    adds r1, #4
    str r1, [r3]
    ldr r0, =_state
    movs r1, #1
    str r1, [r0]
    pop {pc}

defcode ";", 1, 1
    ldr r0, =_here
    ldr r1, [r0]
    ldr r2, =0x46c0bd00
    str r2, [r1]
    adds r1, #4
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    str r1, [r0]
    ldr r0, =_state
    movs r1, #0
    str r1, [r0]
    bx lr

defcode "VARIABLE", 8
    push {lr}
    bl _getw
    ldr r2, =_latest
    ldr r3, =_here
    ldr r1, [r3]
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    ldr r7, [r2]
    str r7, [r1]
    str r1, [r2]
    adds r1, #4
    movs r7, #0
    strb r7, [r1]
    strb r0, [r1, #1]
    adds r1, #2
    ldr r6, =_tib
    movs r7, #0
1:  ldrb r2, [r6, r7]
    strb r2, [r1, r7]
    adds r7, #1
    cmp r7, r0
    blt 1b
    add r1, r0
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    ldr r6, =0x3d04
    strh r6, [r1]
    ldr r6, =0x602c
    strh r6, [r1, #2]
    ldr r6, =0xa400
    strh r6, [r1, #4]
    ldr r6, =0x4770
    strh r6, [r1, #6]
    movs r6, #0
    str r6, [r1, #8]
    adds r1, #12
    str r1, [r3]
    pop {pc}

defcode "HERE", 4
    subs r5, #4
    str r4, [r5]
    ldr r0, =_here
    ldr r4, [r0]
    bx lr

defcode "ALLOT", 5
    ldr r0, =_here
    ldr r1, [r0]
    add r1, r4
    str r1, [r0]
    ldr r4, [r5]
    adds r5, #4
    bx lr

defcode ",", 1
    ldr r0, =_here
    ldr r1, [r0]
    str r4, [r1]
    adds r1, #4
    str r1, [r0]
    ldr r4, [r5]
    adds r5, #4
    bx lr

defcode "ALIGN", 5
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    str r1, [r0]
    bx lr

@ -----------------------------------------------------------------------------
@ CONTROL FLOW (IF, ELSE, THEN, BEGIN, UNTIL, DO, LOOP)
@ -----------------------------------------------------------------------------
@ These words compile ARM branch instructions into RAM code at runtime.

defcode "IF", 2, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1
    lsrs r1, #1
    lsls r1, #1
    ldr r2, =0x4620
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x682c
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x3504
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x2800
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xd100
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xe000
    strh r2, [r1]
    subs r5, #4
    str r4, [r5]
    mov r4, r1
    adds r1, #2
    str r1, [r0]
    pop {pc}

defcode "THEN", 4, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    subs r2, r1, r4
    subs r2, #4
    asrs r2, #1
    ldr r3, =0x7ff
    ands r2, r3
    ldr r3, =0xe000
    orrs r2, r3
    strh r2, [r4]
    ldr r4, [r5]
    adds r5, #4
    pop {pc}

defcode "ELSE", 4, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1
    lsrs r1, #1
    lsls r1, #1
    ldr r2, =0xe000
    strh r2, [r1]
    mov r7, r1
    adds r1, #2
    str r1, [r0]
    ldr r1, [r0]
    subs r2, r1, r4
    subs r2, #4
    asrs r2, #1
    ldr r3, =0x7ff
    ands r2, r3
    ldr r3, =0xe000
    orrs r2, r3
    strh r2, [r4]
    mov r4, r7
    pop {pc}

defcode "BEGIN", 5, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1     @ align next instruction
    lsrs r1, #1     @ Make sure it ends on an
    lsls r1, #1     @ even number 2,4,6 etc
    str r1, [r0]
    subs r5, #4
    str r4, [r5]
    mov r4, r1
    pop {pc}

defcode "UNTIL", 5, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1
    lsrs r1, #1
    lsls r1, #1
    ldr r2, =0x4620 @ mov r0, r4
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x682c @ ldr r4, [r5, #0]
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x3504 @ adds r5, #4
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x2800 @ cmp r0, #0
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xd100 @ bne <label>
    strh r2, [r1]
    adds r1, #2

    @ CALCULATE relative branch offset
    @ r4 - r1 → distance from current position to loop start
    @ -4 → adjust for pipeline / branch position
    @ asrs #1 → convert byte offset → halfword offset (Thumb branches count halfwords)
    @ & 0x7ff → keep 11-bit offset field
    @ | 0xe000 → encode as BNE instruction
    subs r2, r4, r1
    subs r2, #4
    asrs r2, #1
    ldr r3, =0x7ff
    ands r2, r3
    ldr r3, =0xe000
    orrs r2, r3

    strh r2, [r1]
    adds r1, #2
    str r1, [r0]
    ldr r4, [r5]
    adds r5, #4
    pop {pc}

defcode "DO", 2, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1
    lsrs r1, #1
    lsls r1, #1
    ldr r2, =0x6828
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xb411
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x3504
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x682c
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x3504
    strh r2, [r1]
    adds r1, #2
    str r1, [r0]
    subs r5, #4
    str r4, [r5]
    mov r4, r1
    pop {pc}

defcode "LOOP", 4, 1
    push {lr}
    ldr r0, =_here
    ldr r1, [r0]
    adds r1, #1
    lsrs r1, #1
    lsls r1, #1
    ldr r2, =0xbc03
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x3101
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0x4281
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xda01
    strh r2, [r1]
    adds r1, #2
    ldr r2, =0xb403
    strh r2, [r1]
    adds r1, #2
    subs r2, r4, r1
    subs r2, #4
    asrs r2, #1
    ldr r3, =0x7ff
    ands r2, r3
    ldr r3, =0xe000
    orrs r2, r3
    strh r2, [r1]
    adds r1, #2
    str r1, [r0]
    ldr r4, [r5]
    adds r5, #4
    pop {pc}

defcode "I", 1
    subs r5, #4
    str r4, [r5]
    ldr r4, [sp, #4]
    bx lr

.align 2
_builtins_latest: .word _link

@ -----------------------------------------------------------------------------
@ INTERNAL KERNEL ROUTINES
@ -----------------------------------------------------------------------------

@ Read a word from the terminal into _tib
.thumb_func
_getw:
    push {r4, r5, lr}
    ldr r4, =_tib
    movs r5, #0
g1: bl _key
    mov r7, r0
    push {r0, r1, r2, r3, r4, r5}
    mov r0, r7
    bl _emit
    pop {r0, r1, r2, r3, r4, r5}
    cmp r7, #13
    beq g_nl
    cmp r7, #10
    beq g_nl
    cmp r7, #32
    bls g2
    strb r7, [r4, r5]
    adds r5, #1
    b g1
g_nl:
    ldr r0, =_nl_flag
    movs r1, #1
    strb r1, [r0]
g2: cmp r5, #0
    beq g3
    mov r0, r5
    pop {r4, r5, pc}
g3: mov r0, r5
    pop {r4, r5, pc}

.thumb_func
_find:
    push {r4, r5, r6, r7, lr}
    ldr r2, =_latest
    ldr r2, [r2]
f1: cmp r2, #0
    beq f3
    ldrb r3, [r2, #5]
    cmp r3, r1
    bne f4
    movs r7, #0
f2: adds r0, r2, #6
    ldrb r0, [r0, r7]
    ldr r4, =_tib
    ldrb r4, [r4, r7]
    cmp r0, r4
    bne f4
    adds r7, #1
    cmp r7, r3
    blt f2
    ldrb r1, [r2, #4]
    adds r2, #6
    add r2, r3
    adds r2, #3
    lsrs r2, #2
    lsls r2, #2
    adds r2, #1
    mov r0, r2
    pop {r4, r5, r6, r7, pc}
f4: ldr r2, [r2]
    b f1
f3: movs r0, #0
    pop {r4, r5, r6, r7, pc}

.thumb_func
_num:
    push {r4, r5, r6, r7, lr}
    ldr r0, =_tib
    mov r6, r1
    movs r1, #0
    movs r2, #0
    ldrb r3, [r0, #0]
    movs r4, #0
    cmp r3, #45
    bne n1
    movs r4, #1
    adds r2, #1
n1: ldrb r3, [r0, r2]
    subs r3, #48
    cmp r3, #9
    bhi n2
    movs r7, #10
    muls r1, r7, r1
    adds r1, r3
    adds r2, #1
    cmp r2, r6
    blt n1
    cmp r4, #0
    beq n3
    mov r0, r1
    rsbs r0, r0, #0
    mov r1, r0
n3: mov r0, r1
    movs r1, #1
    pop {r4, r5, r6, r7, pc}
n2: movs r1, #0
    pop {r4, r5, r6, r7, pc}

.thumb_func
_idiv:
    push {r4, lr}
    movs r4, #0
    cmp r0, #0
    bge 1f
    rsbs r0, r0, #0
    adds r4, #1
1:  cmp r1, #0
    bge 2f
    rsbs r1, r1, #0
    adds r4, #1
2:  bl _udiv
    lsrs r4, r4, #1
    bcc 3f
    rsbs r0, r0, #0
3:  pop {r4, pc}

.thumb_func
_udiv:
    cmp r1, #0
    beq 4f
    movs r2, #0
    movs r3, #1
1:  lsls r1, #1
    lsls r3, #1
    cmp r1, r0
    bls 1b
2:  lsrs r1, #1
    lsrs r3, #1
    cmp r0, r1
    blo 3f
    subs r0, r0, r1
    orrs r2, r2, r3
3:  cmp r3, #1
    bhi 2b
    mov r0, r2
    bx lr
4:  movs r0, #0
    bx lr

.align 2
.ltorg

@ -----------------------------------------------------------------------------
@ FORTH MAIN LOOP (QUIT)
@ -----------------------------------------------------------------------------
.thumb_func
QUIT:
    movs r0, #10
    bl _emit
    movs r0, #111
    bl _emit
    movs r0, #107
    bl _emit
    movs r0, #32
    bl _emit
    ldr r0, =_nl_flag
    movs r1, #0
    strb r1, [r0]
qloop:
    bl _getw
    mov r6, r0
    cmp r6, #0
    beq q_prompt_check
    mov r1, r6
    bl _find
    cmp r0, #0
    beq qnum
    mov r7, r0
    mov r2, r1
    ldr r1, =_state
    ldr r1, [r1]
    cmp r1, #0
    beq qexec
    movs r1, #1
    tst r2, r1
    bne qexec
    ldr r2, =_here
    ldr r3, [r2]
    adds r3, #3
    lsrs r3, #2
    lsls r3, #2
    ldr r1, =0x4801
    strh r1, [r3]
    ldr r1, =0x4780
    strh r1, [r3, #2]
    ldr r1, =0xe002
    strh r1, [r3, #4]
    ldr r1, =0xbf00
    strh r1, [r3, #6]
    str r7, [r3, #8]
    adds r3, #12
    str r3, [r2]
    b q_prompt_check
qexec:
    blx r7
    ldr r1, =_state
    ldr r1, [r1]
    cmp r1, #0
    beq QUIT
    b q_prompt_check
qnum:
    mov r1, r6
    bl _num
    cmp r1, #0
    beq qerr
    mov r7, r0
    ldr r1, =_state
    ldr r1, [r1]
    cmp r1, #0
    beq qpush
    ldr r2, =_here
    ldr r3, [r2]
    adds r3, #3
    lsrs r3, #2
    lsls r3, #2
    ldr r1, =0x4801
    strh r1, [r3]
    ldr r1, =0x4780
    strh r1, [r3, #2]
    ldr r1, =0xe004
    strh r1, [r3, #4]
    ldr r1, =0xbf00
    strh r1, [r3, #6]
    ldr r0, =_clit
    str r0, [r3, #8]
    str r7, [r3, #12]
    adds r3, #16
    str r3, [r2]
    b q_prompt_check
qpush:
    subs r5, #4
    str r4, [r5]
    mov r4, r7
    b q_prompt_check
qerr:
    movs r0, #63
    bl _emit
    b QUIT
q_prompt_check:
    ldr r0, =_nl_flag
    ldrb r1, [r0]
    cmp r1, #0
    beq qloop
    movs r1, #0
    strb r1, [r0]
    movs r0, #10
    bl _emit
    movs r0, #111
    bl _emit
    movs r0, #107
    bl _emit
    movs r0, #32
    bl _emit
    b qloop

.thumb_func
_gpio_init:
    push {lr}
    @ 1. Unreset IO_BANK0 and PADS_BANK0
    ldr r1, =RESETS_BASE
    ldr r0, =(RESETS_RESET_IO_BANK0_BITS | RESETS_RESET_PADS_BANK0_BITS)
    @ Use atomic clear alias to clear the reset bits
    ldr r2, =ATOMIC_CLR_OFFSET
    add r2, r1
    str r0, [r2, #RESETS_RESET_OFFSET]
1:  ldr r2, [r1, #RESETS_RESET_DONE_OFFSET]
    tst r2, r0
    beq 1b
    @ 2. Set all GPIO pins (0-29) to function 5 (SIO)
    ldr r1, =IO_BANK0_BASE
    adds r1, #4            @ Point to GPIO0_CTRL (STATUS is at +0, CTRL at +4)
    movs r0, #5            @ Function 5 is SIO
    movs r2, #30           @ Total number of GPIO pins
2:  str r0, [r1]           @ Set this pin to SIO
    adds r1, #8            @ Advance to next GPIO_CTRL (skip STATUS)
    subs r2, #1
    bne 2b
    pop {pc}

.thumb_func
_clit:
    subs r5, #4
    str r4, [r5]
    mov r0, lr
    subs r0, #1
    ldr r4, [r0, #8]
    adds r0, #13
    bx r0

.align 2
.ltorg

@ -----------------------------------------------------------------------------
@ DATA AND MEMORY AREAS
@ -----------------------------------------------------------------------------
.section .data
_latest: .word 0
_here:   .word 0
_state:  .word 0

.section .bss
.align 4
_data_stack_lim: .space 512
_data_stack_base:
_tib:            .space 64
_nl_flag:        .byte 0
