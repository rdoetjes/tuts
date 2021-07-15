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
loop_up:
	bl sys_gpioByte
	mov x5, x0			//nano changes x0, so we safe it
	m_nanosleep
	mov x0, x5			//restore x0 after nano call

	lsl x0, x0, #1
	cmp x0, #128
	bne loop_up
loop_down:
	bl sys_gpioByte
        mov x5, x0                      //nano changes x0, so we safe it
        m_nanosleep
        mov x0, x5                      //restore x0 after nano call

        lsr x0, x0, #1
        cmp x0, #1
	bne loop_down
	b loop_up

exit:
	m_exit 0

.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 070000000
