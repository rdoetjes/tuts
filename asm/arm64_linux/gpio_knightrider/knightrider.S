#include "gpio.S"
#include "system.S"

.text
.align 8
.global _start

_start:
	m_gpio_map			//x10 will contain the memory map, x11 the /dev/mem fd

	m_gpio_setDirectionOut gpio2	//set gpio 2 as output
	m_gpio_setDirectionOut gpio3	//set gpio 3 as output
	m_gpio_setDirectionOut gpio4	//set gpio 4 as output
	m_gpio_setDirectionOut gpio5	//set gpio 5 as output
	m_gpio_setDirectionOut gpio6	//set gpio 6 as output
	m_gpio_setDirectionOut gpio17	//set gpio 17 as output
	m_gpio_setDirectionOut gpio22	//set gpio 22 as output
	m_gpio_setDirectionOut gpio27	//set gpio 27 as output

	m_gpio_setDirectionOut gpio23	//set gpio 23 as output (the slector lsb msb)

reset:
	mov x0, #0
	bl gpio_Set16bit
	m_nanosleep

	mov x0 ,#1
1:
	bl gpio_Set16bit
	bl Sleep
	lsl x0, x0, #1
	cmp x0, #0x10000
	bne 1b

2:
	lsr x0, x0, #1
	bl gpio_Set16bit
	bl Sleep
	cmp x0, #1
	bne 2b
	b 1b
	m_exit 0

//this procedure saves us the extra step of 
//saving the arguments x0 and x1 which we use in the main also
//and keeps the code cleaner
Sleep:
	stp x29, x30, [sp, #-16]!
        stp x0, x1, [sp, #-16]!

	m_nanosleep

	ldp x0, x1, [sp], #16
        ldp x29, x30, [sp], #16
	ret

exit:
	m_exit 0


.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 050000000
