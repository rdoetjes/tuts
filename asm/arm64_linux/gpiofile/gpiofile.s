.include "gpio.s"

.text
.align 8
.global _start

_start:
	m_gpio_map

	m_gpio_setDirectionOut p2
	m_gpio_setDirectionOut p3
	m_gpio_setDirectionOut p4
	m_gpio_setDirectionOut p5
	m_gpio_setDirectionOut p6
	m_gpio_setDirectionOut p17
	m_gpio_setDirectionOut p22
	m_gpio_setDirectionOut p27
reset:
	mov x0, #1
loop:
	bl sys_gpioByte
loop1:
	mov x5, x0
	m_nanosleep
	mov x0, x5

	lsl x0, x0, #1
	cmp x0, #256
	beq reset
	b loop

exit:
	m_exit 0

.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 600000000
