@ RP2040 Forth - Subroutine Threaded (STC)
.include "rp2040.inc"
.syntax unified
.cpu cortex-m0plus
.thumb

.set _link, 0
.macro defh name, label, len, flags=0
    .align 2
_link_\label:
    .word _link
    .set _link, _link_\label
    .byte \flags, \len
    .ascii "\name"
    .align 2
    .thumb_func
\label:
.endm

.section .vectors, "ax"
    .word _stack_top
    .word _reset + 1

.section .text
.global _reset
.thumb_func
_reset:
    ldr r0, =_sidata; ldr r1, =_sdata; ldr r2, =_edata
1:  cmp r1, r2; bhs 2f; ldr r3, [r0]; adds r0, #4; str r3, [r1]; adds r1, #4; b 1b
2:  ldr r1, =_sbss; ldr r2, =_ebss; movs r3, #0
3:  cmp r1, r2; bhs 4f; str r3, [r1]; adds r1, #4; b 3b
4:  ldr r0, =CMSDK_UART0_BASE; movs r1, #3; str r1, [r0, #CMSDK_UART_CTRL_OFFSET]
    ldr r5, =_data_stack_base; movs r4, #0
    ldr r0, =_latest; ldr r1, =_link_builtins; str r1, [r0]
    ldr r0, =_here; ldr r1, =_forth_ram_start; str r1, [r0]
    b QUIT

.thumb_func
_emit:
    push {r1, r2, r3, lr}; ldr r1, =CMSDK_UART0_BASE
e1: ldr r2, [r1, #CMSDK_UART_STATE_OFFSET]; movs r3, #1; tst r2, r3; bne e1
    str r0, [r1, #CMSDK_UART_DATA_OFFSET]; pop {r1, r2, r3, pc}

.thumb_func
_key:
    push {r1, r2, r3, lr}
    ldr r1, =CMSDK_UART0_BASE
k1: ldr r2, [r1, #CMSDK_UART_STATE_OFFSET]; movs r3, #2; tst r2, r3; beq k1
    ldr r0, [r1, #CMSDK_UART_DATA_OFFSET];
    movs r1, #255; ands r0, r1
    pop {r1, r2, r3, pc}

.thumb_func
_print_dec:
    push {r4, r5, r6, lr}; mov r4, r0; ldr r6, =powers_of_10; movs r5, #0
    cmp r4, #0; bge pd1
    mov r0, r4; mvns r0, r0; adds r0, #1; mov r4, r0
    movs r0, #45; bl _emit
pd1:ldr r2, [r6]; cmp r2, #0; beq pd5; movs r0, #0
pd2:cmp r4, r2; blo pd3; subs r4, r2; adds r0, #1; b pd2
pd3:cmp r0, #0; bne pd4; cmp r5, #0; bne pd4; cmp r2, #1; beq pd4; b pd6
pd4:movs r5, #1; adds r0, #48; bl _emit
pd6:adds r6, #4; b pd1
pd5:pop {r4, r5, r6, pc}

.align 2
powers_of_10:
    .word 1000000000, 100000000, 10000000, 1000000, 100000, 10000, 1000, 100, 10, 1, 0

defh "DUP", f_dup, 3
    subs r5, #4; str r4, [r5]; bx lr
defh "DROP", f_drop, 4
    ldr r4, [r5]; adds r5, #4; bx lr
defh "SWAP", f_swap, 4
    ldr r0, [r5]; str r4, [r5]; mov r4, r0; bx lr
defh "+", f_plus, 1
    ldr r0, [r5]; adds r4, r4, r0; adds r5, #4; bx lr
defh "-", f_minus, 1
    ldr r0, [r5]; subs r4, r0, r4; adds r5, #4; bx lr
defh "*", f_star, 1
    ldr r0, [r5]; muls r4, r0, r4; adds r5, #4; bx lr

defh "=", f_equal, 1
    ldr r0, [r5]; adds r5, #4; cmp r0, r4; bne eq1
    movs r4, #0; mvns r4, r4; bx lr
eq1:movs r4, #0; bx lr

defh "<", f_less, 1
    ldr r0, [r5]; adds r5, #4; cmp r0, r4; bge lt1
    movs r4, #0; mvns r4, r4; bx lr
lt1:movs r4, #0; bx lr

defh ".", f_dot, 1
    push {lr}; mov r0, r4; bl _print_dec; movs r0, #32; bl _emit
    ldr r4, [r5]; adds r5, #4; pop {pc}

defh ".S", f_dot_s, 2
    push {r4, r5, r6, r7, lr}
    movs r0, #60; bl _emit @ '<'
    ldr r6, =_data_stack_base; subs r0, r6, r5; lsrs r0, r0, #2
    bl _print_dec
    movs r0, #62; bl _emit @ '>'
    movs r0, #32; bl _emit
    ldr r6, =_data_stack_base; subs r6, #8
ds1:cmp r6, r5; blo ds2; ldr r0, [r6]; bl _print_dec; movs r0, #32; bl _emit; subs r6, #4; b ds1
ds2:ldr r6, =_data_stack_base; cmp r5, r6; beq ds3; mov r0, r4; bl _print_dec; movs r0, #32; bl _emit
ds3:pop {r4, r5, r6, r7, pc}

defh "WORDS", f_words, 5
    push {r4, lr}; ldr r2, =_latest; ldr r2, [r2]
w1: cmp r2, #0; beq w3; ldrb r1, [r2, #5]; adds r3, r2, #6
w2: ldrb r0, [r3]; bl _emit; adds r3, #1; subs r1, #1; bne w2
    movs r0, #32; bl _emit; ldr r2, [r2]; b w1
w3: pop {r4, pc}

.align 2
.ltorg

defh ":", f_colon, 1
    push {lr}; bl _getw; ldr r2, =_here; ldr r3, [r2]
    adds r3, #3; lsrs r3, #2; lsls r3, #2
    ldr r7, =_latest; ldr r1, [r7]; str r1, [r3]; str r3, [r7]
    adds r3, #4; movs r1, #0; strb r1, [r3]; strb r0, [r3, #1]; adds r3, #2
    ldr r1, =_tib; movs r7, #0
c1: ldrb r6, [r1, r7]; strb r6, [r3, r7]; adds r7, #1; cmp r7, r0; blt c1
    add r3, r0; adds r3, #3; lsrs r3, #2; lsls r3, #2
    ldr r6, =0x46c0b500 @ push {lr}
    str r6, [r3]; adds r3, #4
    str r3, [r2]
    ldr r0, =_state; movs r1, #1; str r1, [r0]; pop {pc}

defh ";", f_semicolon, 1, 1
    ldr r0, =_here; ldr r1, [r0]
    ldr r2, =0x46c0bd00 @ pop {pc}
    str r2, [r1]; adds r1, #4
    adds r1, #3; lsrs r1, #2; lsls r1, #2; str r1, [r0]
    ldr r0, =_state; movs r1, #0; str r1, [r0]; bx lr

defh "IF", f_if, 2, 1
    push {lr}
    ldr r0, =_here; ldr r1, [r0]
    @ align to halfword just in case
    adds r1, #1; lsrs r1, #1; lsls r1, #1

    @ cmp r4, #0
    ldr r2, =0x2c00; strh r2, [r1]; adds r1, #2
    @ ldr r4, [r5]
    ldr r2, =0x682c; strh r2, [r1]; adds r1, #2
    @ adds r5, #4
    ldr r2, =0x3504; strh r2, [r1]; adds r1, #2
    @ beq <forward> - jump if condition is false
    ldr r2, =0xd000; strh r2, [r1]

    subs r5, #4; str r4, [r5]; mov r4, r1
    adds r1, #2; str r1, [r0]
    pop {pc}

defh "THEN", f_then, 4, 1
    push {lr}
    ldr r0, =_here; ldr r1, [r0]
    subs r2, r1, r4; subs r2, #4; asrs r2, #1
    strb r2, [r4]
    ldr r4, [r5]; adds r5, #4
    pop {pc}

defh "ELSE", f_else, 4, 1
    push {lr}
    ldr r0, =_here; ldr r1, [r0]
    adds r1, #1; lsrs r1, #1; lsls r1, #1
    ldr r2, =0xe000; strh r2, [r1] @ b <forward>
    mov r7, r1
    adds r1, #2; str r1, [r0]
    
    ldr r1, [r0]
    subs r2, r1, r4; subs r2, #4; asrs r2, #1
    strb r2, [r4]
    
    mov r4, r7
    pop {pc}

defh "\\", f_bs, 1, 1
b1: bl _key; cmp r0, #10; beq b2; cmp r0, #13; beq b2; b b1
b2: bx lr

.align 2
.ltorg

.set _link_builtins, _link

_getw:
    push {r4, r5, lr}
    ldr r4, =_tib
    movs r5, #0
g1: bl _key
    mov r7, r0
    bl _emit
    cmp r7, #32
    bls g2
    strb r7, [r4, r5]
    adds r5, #1
    b g1
g2: cmp r5, #0
    beq g1
    mov r0, r5
    pop {r4, r5, pc}

_find:
    push {r4, r5, r6, r7, lr}; ldr r2, =_latest; ldr r2, [r2]
f1: cmp r2, #0; beq f3; ldrb r3, [r2, #5]; cmp r3, r1; bne f4; movs r7, #0
f2: adds r0, r2, #6; ldrb r0, [r0, r7]; ldr r4, =_tib; ldrb r4, [r4, r7]
    cmp r0, r4; bne f4; adds r7, #1; cmp r7, r3; blt f2
    ldrb r1, [r2, #4]; adds r2, #6; add r2, r3; adds r2, #3; lsrs r2, #2; lsls r2, #2
    adds r2, #1; mov r0, r2; pop {r4, r5, r6, r7, pc}
f4: ldr r2, [r2]; b f1
f3: movs r0, #0; pop {r4, r5, r6, r7, pc}

_num:
    push {r4, r5, r6, r7, lr}; ldr r0, =_tib; mov r6, r1; movs r1, #0; movs r2, #0
    ldrb r3, [r0, #0]; movs r4, #0; cmp r3, #45; bne n1
    movs r4, #1; adds r2, #1
n1: ldrb r3, [r0, r2]; subs r3, #48; cmp r3, #9; bhi n2
    movs r7, #10; muls r1, r7, r1; adds r1, r3; adds r2, #1; cmp r2, r6; blt n1
    cmp r4, #0; beq n3
    mov r0, r1; mvns r0, r0; adds r0, #1; mov r1, r0
n3: mov r0, r1; movs r1, #1; pop {r4, r5, r6, r7, pc}
n2: movs r1, #0; pop {r4, r5, r6, r7, pc}

.align 2
.ltorg

QUIT:
    movs r0, #10; bl _emit; movs r0, #111; bl _emit; movs r0, #107; bl _emit; movs r0, #32; bl _emit
qloop:
    bl _getw; mov r6, r0; mov r1, r6; bl _find
    cmp r0, #0; beq qnum
    mov r7, r0; mov r2, r1; ldr r1, =_state; ldr r1, [r1]
    cmp r1, #0; beq qexec
    movs r1, #1; tst r2, r1; bne qexec
    ldr r2, =_here; ldr r3, [r2]
    adds r3, #3; lsrs r3, #2; lsls r3, #2
    ldr r1, =0x4801; strh r1, [r3]
    ldr r1, =0x4780; strh r1, [r3, #2]
    ldr r1, =0xe002; strh r1, [r3, #4]
    ldr r1, =0xbf00; strh r1, [r3, #6]
    str r7, [r3, #8]
    adds r3, #12; str r3, [r2]; b qloop
qexec:
    blx r7; ldr r1, =_state; ldr r1, [r1]; cmp r1, #0; beq QUIT; b qloop
qnum:
    mov r1, r6; bl _num; cmp r1, #0; beq qerr
    mov r7, r0; ldr r1, =_state; ldr r1, [r1]
    cmp r1, #0; beq qpush
    ldr r2, =_here; ldr r3, [r2]
    adds r3, #3; lsrs r3, #2; lsls r3, #2
    ldr r1, =0x4801; strh r1, [r3]
    ldr r1, =0x4780; strh r1, [r3, #2]
    ldr r1, =0xe004; strh r1, [r3, #4]
    ldr r1, =0xbf00; strh r1, [r3, #6]
    ldr r0, =_clit; str r0, [r3, #8]
    str r7, [r3, #12]; adds r3, #16; str r3, [r2]; b qloop
qpush:
    subs r5, #4; str r4, [r5]; mov r4, r7; b QUIT
qerr:
    movs r0, #63; bl _emit; b QUIT

.thumb_func
_clit:
    subs r5, #4; str r4, [r5]
    mov r0, lr
    subs r0, #1
    ldr r4, [r0, #8]
    adds r0, #13
    bx r0

.align 2
.ltorg

.section .data
_latest: .word 0
_here: .word 0
_state: .word 0

.section .bss
.align 4
_data_stack_lim: .space 512
_data_stack_base:
_tib: .space 64
