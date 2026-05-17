@ RP2040 Forth Implementation - Subroutine Threaded Code (STC)
@ ==========================================================
@ Register Mapping:
@ r4 : TOS (Top of Data Stack)
@ r5 : DSP (Data Stack Pointer)
@ r0-r3: Scratch / Arguments
@ lr : Instruction Pointer (STC)

.include "rp2040.inc"

.syntax unified
.cpu cortex-m0plus
.thumb

@ --- Dictionary Header Macro ---
.set _link, 0

.macro defheader name, label, namelen, flags=0
    .align 2
    .word _link
    .set _link, .
    .byte \flags
    .byte \namelen
    .ascii "\name"
    .align 2
    .thumb_func
\label:
.endm

@ --- Primitives Macros ---
.macro NEXT;   bx lr;       .endm
.macro PUSH_TOS; adds r5, #-4; str r4, [r5]; .endm
.macro POP_TOS;  ldr r4, [r5]; adds r5, #4;  .endm

@ --- Vector Table ---
.section .vectors, "ax"
.align 2
    .word _stack_top
    .word _reset + 1

.section .text

.global _reset
.thumb_func
_reset:
    @ 0. Copy data and zero BSS
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
4:
    @ Detect QEMU
    ldr r0, =RESETS_BASE
    ldr r1, [r0, #RESETS_RESET_DONE_OFFSET]
    ldr r2, =_is_qemu
    movs r3, #0
    cmp r1, #0
    beq 10f
    b 11f
10: movs r3, #1
11: strb r3, [r2]

    @ 1. Unreset (Hardware only)
    cmp r3, #0
    bne 6f
    ldr r1, =(RESETS_RESET_IO_BANK0_BITS | RESETS_RESET_PADS_BANK0_BITS | RESETS_RESET_UART0_BITS)
    str r1, [r0, #RESETS_RESET_OFFSET]
    movs r2, #0
    str r2, [r0, #RESETS_RESET_OFFSET]
5:  ldr r2, [r0, #RESETS_RESET_DONE_OFFSET]
    tst r2, r1
    beq 5b
6:
    @ 2. Setup UART (Hardware only)
    ldr r2, =_is_qemu
    ldrb r2, [r2]
    cmp r2, #0
    bne 7f
    ldr r0, =UART0_BASE
    movs r1, #67
    str r1, [r0, #UART_IBRD_OFFSET]
    movs r1, #52
    str r1, [r0, #UART_FBRD_OFFSET]
    movs r1, #0x70
    str r1, [r0, #UART_LCR_H_OFFSET]
    movs r1, #0x01
    str r1, [r0, #UART_CR_OFFSET]
    ldr r0, =IO_BANK0_BASE
    movs r1, #2
    str r1, [r0, #GPIO0_CTRL_OFFSET]
    str r1, [r0, #GPIO1_CTRL_OFFSET]

7:  ldr r0, =CMSDK_UART0_BASE
    movs r1, #3
    str r1, [r0, #CMSDK_UART_CTRL_OFFSET]

    @ Initialize Forth
    ldr r5, =_data_stack_base
    movs r4, #0
    ldr r0, =_latest
    ldr r1, =_link_builtin
    str r1, [r0]
    ldr r0, =_here
    ldr r1, =_forth_ram_start
    str r1, [r0]

    @ Run REPL
    b QUIT

.align 2
.ltorg

@ --- Core Words ---

defheader "DUP", f_dup, 3
    PUSH_TOS
    NEXT

defheader "DROP", f_drop, 4
    POP_TOS
    NEXT

defheader "SWAP", f_swap, 4
    ldr r0, [r5]
    str r4, [r5]
    mov r4, r0
    NEXT

defheader "OVER", f_over, 4
    ldr r0, [r5]
    PUSH_TOS
    mov r4, r0
    NEXT

defheader "+", f_plus, 1
    ldr r0, [r5]
    adds r4, r4, r0
    adds r5, #4
    NEXT

defheader "-", f_minus, 1
    ldr r0, [r5]
    subs r4, r0, r4
    adds r5, #4
    NEXT

defheader "*", f_star, 1
    ldr r0, [r5]
    muls r4, r0, r4
    adds r5, #4
    NEXT

defheader "AND", f_and, 3
    ldr r0, [r5]
    ands r4, r4, r0
    adds r5, #4
    NEXT

defheader "OR", f_or, 2
    ldr r0, [r5]
    orrs r4, r4, r0
    adds r5, #4
    NEXT

defheader "INVERT", f_invert, 6
    mvns r4, r4
    NEXT

defheader "@", f_fetch, 1
    ldr r4, [r4]
    NEXT

defheader "!", f_store, 1
    ldr r0, [r5]
    str r0, [r4]
    POP_TOS
    NEXT

defheader ".", f_dot, 1
    push {lr}
    mov r0, r4
    bl _print_hex
    POP_TOS
    pop {pc}

.thumb_func
_print_hex:
    push {r4, lr}
    mov r2, r0
    movs r3, #8
1:  lsrs r0, r2, #28
    cmp r0, #10
    blt 2f
    adds r0, #7
2:  adds r0, #48
    mov r4, r0
    bl _emit_internal
    lsls r2, r2, #4
    subs r3, #1
    bne 1b
    pop {r4, pc}

defheader "WORDS", f_words, 5
    push {r4, lr}
    ldr r2, =_latest
    ldr r2, [r2]
1:  cmp r2, #0
    beq 3f
    ldrb r1, [r2, #5]
    adds r3, r2, #6
2:  ldrb r4, [r3]
    bl _emit_internal
    adds r3, #1
    subs r1, #1
    bne 2b
    movs r4, #32
    bl _emit_internal
    ldr r2, [r2]
    b 1b
3:  pop {r4, pc}

.thumb_func
_emit_internal:
    push {r0, r1, r2, lr}
    ldr r1, =CMSDK_UART0_BASE
    str r4, [r1, #CMSDK_UART_DATA_OFFSET]
    ldr r1, =_is_qemu
    ldrb r1, [r1]
    cmp r1, #0
    bne 2f
    ldr r1, =UART0_BASE
1:  ldr r0, [r1, #UART_FR_OFFSET]
    movs r2, #0x20
    tst r0, r2
    bne 1b
    str r4, [r1, #UART_DR_OFFSET]
2:  pop {r0, r1, r2, pc}

defheader "KEY", f_key, 3
    PUSH_TOS
1:  ldr r1, =CMSDK_UART0_BASE
    ldr r0, [r1, #CMSDK_UART_STATE_OFFSET]
    movs r2, #0x02
    tst r0, r2
    beq 2f
    ldr r4, [r1, #CMSDK_UART_DATA_OFFSET]
    NEXT
2:  ldr r1, =_is_qemu
    ldrb r1, [r1]
    cmp r1, #0
    bne 1b
    ldr r1, =UART0_BASE
    ldr r0, [r1, #UART_FR_OFFSET]
    movs r2, #0x10
    tst r0, r2
    bne 1b
    ldr r4, [r1, #UART_DR_OFFSET]
    NEXT

defheader "EMIT", f_emit, 4
    push {lr}
    bl _emit_internal
    POP_TOS
    pop {pc}

defheader "HERE", f_here, 4
    PUSH_TOS
    ldr r0, =_here
    ldr r4, [r0]
    NEXT

defheader "LATEST", f_latest, 6
    PUSH_TOS
    ldr r0, =_latest
    ldr r4, [r0]
    NEXT

defheader "STATE", f_state, 5
    PUSH_TOS
    ldr r0, =_state
    ldr r4, [r0]
    NEXT

defheader ",", f_comma, 1
    ldr r0, =_here
    ldr r1, [r0]
    str r4, [r1]
    adds r1, #4
    str r1, [r0]
    POP_TOS
    NEXT

defheader "[", f_lbrace, 1, 1
    ldr r0, =_state
    movs r1, #0
    str r1, [r0]
    NEXT

defheader "]", f_rbrace, 1
    ldr r0, =_state
    movs r1, #1
    str r1, [r0]
    NEXT

defheader ":", f_colon, 1
    push {lr}
    bl _get_word
    ldr r2, =_here
    ldr r3, [r2]
    adds r3, #3
    lsrs r3, #2
    lsls r3, #2
    ldr r7, =_latest
    ldr r1, [r7]
    str r1, [r3]
    str r3, [r7]
    adds r3, #4
    movs r1, #0
    strb r1, [r3]
    strb r0, [r3, #1]
    adds r3, #2
    ldr r1, =_tib
    movs r7, #0
1:  ldrb r6, [r1, r7]
    strb r6, [r3, r7]
    adds r7, #1
    cmp r7, r0
    blt 1b
    add r3, r0
    adds r3, #3
    lsrs r3, #2
    lsls r3, #2
    str r3, [r2]
    ldr r0, =_state
    movs r1, #1
    str r1, [r0]
    pop {pc}

defheader ";", f_semicolon, 1, 1
    ldr r0, =_here
    ldr r1, [r0]
    ldr r2, =0x4770
    strh r2, [r1]
    adds r1, #2
    adds r1, #3
    lsrs r1, #2
    lsls r1, #2
    str r1, [r0]
    ldr r0, =_state
    movs r1, #0
    str r1, [r0]
    NEXT

defheader "IMMEDIATE", f_immediate, 9
    ldr r0, =_latest
    ldr r0, [r0]
    adds r0, #4
    ldrb r1, [r0]
    movs r2, #1
    orrs r1, r2
    strb r1, [r0]
    NEXT

defheader "\\", f_backslash, 1, 1
1:  bl f_key
    mov r0, r4
    POP_TOS
    cmp r0, #10
    beq 2f
    cmp r0, #13
    beq 2f
    b 1b
2:  NEXT

.set _link_builtin, _link

@ --- Internal Functions ---

.thumb_func
_get_word:
    push {lr}
    ldr r2, =_tib
    movs r3, #0
1:  bl f_key
    mov r7, r4
    POP_TOS
    cmp r7, #32
    bls 2f
    strb r7, [r2, r3]
    adds r3, #1
    cmp r3, #31
    blt 1b
2:  cmp r3, #0
    beq 1b
    mov r0, r3
    pop {pc}

.thumb_func
_find:
    push {r4, r5, r6, r7, lr}
    ldr r2, =_latest
    ldr r2, [r2]
1:  cmp r2, #0
    beq _find_fail
    ldrb r3, [r2, #5]
    cmp r3, r1
    bne 4f
    movs r4, #0
2:  adds r7, r2, #6
    ldrb r7, [r7, r4]
    ldrb r6, [r0, r4]
    cmp r7, r6
    bne 4f
    adds r4, #1
    cmp r4, r1
    blt 2b
    ldrb r4, [r2, #4]
    adds r2, #6
    add r2, r1
    adds r2, #3
    lsrs r2, #2
    lsls r2, #2
    adds r2, #1
    mov r0, r2
    mov r1, r4
    pop {r4, r5, r6, r7, pc}
4:  ldr r2, [r2]
    b 1b
_find_fail:
    movs r0, #0
    pop {r4, r5, r6, r7, pc}

.thumb_func
_number:
    ldr r0, =_tib
    movs r1, #0
    movs r2, #0
    ldrb r3, [r0]
    cmp r3, #0
    beq _num_fail
1:  ldrb r3, [r0, r2]
    subs r3, #48
    cmp r3, #9
    bhi _num_fail
    movs r7, #10
    muls r1, r7, r1
    adds r1, r3
    adds r2, #1
    cmp r2, r6
    blt 1b
    mov r0, r1
    movs r1, #1
    bx lr
_num_fail:
    movs r1, #0
    bx lr

.thumb_func
_compile_call:
    ldr r2, =_here
    ldr r3, [r2]
    ldr r1, =0x4801
    strh r1, [r3]
    ldr r1, =0x4780
    strh r1, [r3, #2]
    ldr r1, =0xe002
    strh r1, [r3, #4]
    str r0, [r3, #8]
    adds r3, #12
    str r3, [r2]
    bx lr

.thumb_func
_compile_literal:
    ldr r2, =_here
    ldr r3, [r2]
    ldr r1, =0x3d04
    strh r1, [r3]
    ldr r1, =0x602c
    strh r1, [r3, #2]
    ldr r1, =0x4c00
    strh r1, [r3, #4]
    ldr r1, =0xe001
    strh r1, [r3, #6]
    str r0, [r3, #8]
    adds r3, #12
    str r3, [r2]
    bx lr

.align 2
.ltorg

@ --- REPL ---
.thumb_func
QUIT:
    ldr r5, =_data_stack_base

repl_main:
    ldr r0, =_state
    ldr r0, [r0]
    cmp r0, #0
    beq 1f
    movs r4, #62
    bl _emit_internal
    b 2f
1:  movs r4, #10
    bl _emit_internal
    movs r4, #111
    bl _emit_internal
    movs r4, #107
    bl _emit_internal
    movs r4, #32
    bl _emit_internal
2:
    bl _get_word
    mov r6, r0
    ldr r0, =_tib
    mov r1, r6
    bl _find
    cmp r0, #0
    beq repl_not_word
    mov r7, r0
    mov r3, r1
    ldr r1, =_state
    ldr r1, [r1]
    cmp r1, #0
    beq repl_execute
    movs r1, #1
    tst r3, r1
    bne repl_execute
    mov r0, r7
    bl _compile_call
    b repl_main

repl_execute:
    blx r7
    b repl_main

repl_not_word:
    bl _number
    cmp r1, #0
    beq repl_error
    mov r7, r0
    ldr r1, =_state
    ldr r1, [r1]
    cmp r1, #0
    beq repl_push_num
    mov r0, r7
    bl _compile_literal
    b repl_main

repl_push_num:
    PUSH_TOS
    mov r4, r7
    b repl_main

repl_error:
    ldr r2, =_tib
    movs r3, #0
1:  cmp r3, r6
    bge 2f
    ldrb r4, [r2, r3]
    bl _emit_internal
    adds r3, #1
    b 1b
2:  movs r4, #32
    bl _emit_internal
    movs r4, #63
    bl _emit_internal
    ldr r0, =_state
    movs r1, #0
    str r1, [r0]
    b repl_main

.align 4
.ltorg

@ --- Data ---
.section .data
.align 2
_latest:  .word 0
_state:   .word 0
_here:    .word 0
_is_qemu: .byte 0

.section .bss
.align 4
_data_stack_lim:  .space 512
_data_stack_base:
_return_stack_lim: .space 512
_return_stack_base:
_tib:             .space 64
