MEMORY
{
    FLASH(rx) : ORIGIN = 0x00000000, LENGTH = 2048k
    SRAM(rwx) : ORIGIN = 0x20000000, LENGTH = 256k
}

ENTRY(_reset)

SECTIONS
{
    .text : {
        KEEP(*(.vectors))
        *(.text*)
        *(.rodata*)
        . = ALIGN(4);
    } > FLASH

    _sidata = LOADADDR(.data);

    .data : {
        . = ALIGN(4);
        _sdata = .;
        *(.data*)
        . = ALIGN(4);
        _edata = .;
    } > SRAM AT > FLASH

    .bss : {
        . = ALIGN(4);
        _sbss = .;
        *(.bss*)
        *(COMMON)
        . = ALIGN(4);
        _ebss = .;
    } > SRAM

    /* Pointer to end of SRAM for stack */
    _stack_top = ORIGIN(SRAM) + LENGTH(SRAM);

    /* Forth RAM allocation */
    .forth_mem : {
        . = ALIGN(4);
        _forth_ram_start = .;
    } > SRAM
}
