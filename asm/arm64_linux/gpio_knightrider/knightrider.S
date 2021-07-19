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
	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	m_gpio_value gpio23, #0
	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	mov x0, #0
	bl sys_gpioByte
	m_nanosleep

	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	m_gpio_value gpio23, #1
	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	mov x0, #0
	bl sys_gpioByte
	m_nanosleep
	mov x3, #0

	mov x0 ,#1
1:
	bl SetGpio16bit
	bl Sleep
	lsl x0, x0, #1
	cmp x0, #0x10000
	bne 1b

2:
	lsr x0, x0, #1
	bl SetGpio16bit
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

//X0 contains the 16 bit value to set
//X1 contains the pointer to the gpio byte (8 consecutive pins in the gpio.s gpio table) in our case gpio_byte_1
//X10 contains the virtua memory pointer (obtained by m_gpio_map
SetGpio16bit:
	stp x29, x30, [sp, #-16]!
	stp x2, x3, [sp, #-16]!
	stp x0, x1, [sp, #-16]!

	mov x3, x0			//save x0 because it will be overridden in the m_gpio_value 

	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	m_gpio_value gpio23, #1		//select lowest byte
	and x0, x3, #255		//bask off the lowest byte and set it in x0 to be set by sys_gpioByte
	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	bl sys_gpioByte			//set the byte currently in x0 to the 8 consecutive gpio pins, were gpio_byte_1 points to

	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	m_gpio_value gpio23, #0		//select most significant byte on the display
	mov x0, x3			//restore the saved copy of the x0 argument, to process top 8 bits
	lsr x0, x0, #8			//move top 8 bits in lower byte
	and x0, x0, #255		//(this and is not required but ther efo consistency with the latter code) mask them off
	ldr x1, =gpio_byte_1		//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
	bl sys_gpioByte			//set the byte currently in x0 to the 8 consecutive gpio pins, were gpio_byte_1 points to

	ldp x0, x1, [sp], #16
	ldp x2, x3, [sp], #16
	ldp x29, x30, [sp], #16
	ret

exit:
	m_exit 0


.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 130000000
welcome:	.asciz "This is a GPIO demo\n"
