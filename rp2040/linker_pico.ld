MEMORY
{
    FLASH(rx)  : ORIGIN = 0x10000000, LENGTH = 2048k
    SRAM(rwx)  : ORIGIN = 0x20000000, LENGTH = 256k
}

ENTRY(_reset)

SECTIONS
{
    .boot2 ORIGIN(FLASH) : {
        KEEP(*(.boot2))
    } > FLASH

    ASSERT(SIZEOF(.boot2) == 256, "boot2 must be 256 bytes")

    .text : {
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

    .bss (NOLOAD) : {
        . = ALIGN(4);
        _sbss = .;
        *(.bss*)
        *(COMMON)
        . = ALIGN(4);
        _ebss = .;
    } > SRAM

    .forth_mem (NOLOAD) : {
        . = ALIGN(4);
        _forth_ram_start = .;
        . += 200K;
        _forth_ram_end = .;
    } > SRAM

    _stack_top = ORIGIN(SRAM) + LENGTH(SRAM);
}
