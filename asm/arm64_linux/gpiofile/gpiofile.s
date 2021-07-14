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

	mov x5, #0
loop:
	
	m_gpio_value p2 x5
	m_gpio_value p3 x5
	m_gpio_value p4 x5
	m_gpio_value p5 x5
	m_gpio_value p6 x5
	m_gpio_value p17 x5
	m_gpio_value p22 x5
	m_gpio_value p27 x5
	eor x5, x5, #1

	m_nanosleep
	b loop

exit:
	m_exit 0

.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 100000000
