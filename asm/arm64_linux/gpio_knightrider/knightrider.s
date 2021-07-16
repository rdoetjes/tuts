.include "gpio.s"

.text
.align 8
.global _start

_start:
	m_gpio_map			//x10 will contain the memory map, x11 the /dev/mem fd

	m_gpio_setDirectionOut p2	//set pin 2 as output
	m_gpio_setDirectionOut p3	//set pin 3 as output
	m_gpio_setDirectionOut p4	//set pin 4 as output
	m_gpio_setDirectionOut p5	//set pin 5 as output
	m_gpio_setDirectionOut p6	//set pin 6 as output
	m_gpio_setDirectionOut p17	//set pin 17 as output
	m_gpio_setDirectionOut p22	//set pin 22 as output
	m_gpio_setDirectionOut p27	//set pin 27 as output
reset:
	mov x0, #1			//start with lsb set to 1
	ldr x1, =pins			//set x1 to the beginning of the 8 consecutive pins table (in gpio.s)
					//those 8 consecutive pins in that table make up a GPIO byte
loop_up:
	bl sys_gpioByte			//set x0 onto the gpio "byte"

	m_nanosleep

	lsl x0, x0, #1			//shift byte to the left, to light up next light
	cmp x0, #128			//are we at 128? if not continue otherwise shift left
	bne loop_up			//if less than 128 continue shifting left
loop_down:
	bl sys_gpioByte			//set x0 for to gpio

	m_nanosleep

        lsr x0, x0, #1			//shoft x0 to the right to light up previous led
        cmp x0, #1			//are we at #1? if so shift up, otherwise continue shifting down
	bne loop_down			//not 1 continue to shift down
	b loop_up			//x0 == 1 shift it back up

exit:
	m_exit 0

.align 8
.data
timespecsec:	.dword 0
timespecnano:	.dword 070000000
