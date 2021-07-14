.include "system.s"
.include "print.s"
.include "fileio.s"

//This is the macro which functions as an interface
.include "gpio.macro"

.equ	PROT_RD, 1
.equ	PROT_WR, 2
.equ	MAP_SHARED, 1
.equ 	BitSetOffset, 28 
.equ 	BitClearOffset, 40

.text
.align 8

//Expects no arguments
//RETURN X11 file descriptor mmap, 
//	 X10 virtual address
sys_gpioMap:
	stp x29, x30, [sp, #-16]!
	stp x5, x8, [sp, #-16]!
	stp x3, x4, [sp, #-16]!
	stp x1, x2, [sp, #-16]!
	

        //try opening /dev/mem file, error handling is required because
        //the file /dev/mem can only be read by root equivalent users
        m_openfile sys_gpioMap_dev_mem, O_RDWR+O_EXCL       	//open memory map file descriptor to map virtuak memory
        mov x11, x0                            		    	//save memory map file descriptor
        cmp x11, #0						//if file descriptor is <0 then it;s an error and we will report and exit
        ble 1f

	//mmap() (technically mmap2()) to obtain virtual address pointer of 4K ram
	mov x0, #0						//set addr to NULL, so that OS choices page aligned address
	mov x1, #4096						//GPIO has a 4K memory mapped IO block
	mov x2, #(PROT_RD+PROT_WR)				//we need a virtual address that is both read and write
	mov x3, #MAP_SHARED					//memory map is shared
	mov x4, x11						//set the mem dev filedescriptor to argument 5 (x4)
	ldr x5, =sys_gpioMap_gpio_addr				//which address to map intp virtual memory
	ldr x5, [x5]						//load physical addres in memory
	mov x8, #222						//mmap() syscall
	svc #0

	mov x10, x0						//set the mmap pointer to x10 return
	cmp x0, #0						//check virtual memory pointer, if negative then theres an error
	ble 2f							//error

        b 3f
1:
        m_printString sys_gpioMap_error_openfile		//write error message to screem
	m_exit -1						//we are forced to stop execution, we exit with -1 (error)

2:
        m_printString sys_gpioMap_error_mmap			//write error message to screen
	m_exit -1						//we are forced to stop execution, we exit with -1 (error)

3:
	ldp x1, x2, [sp], #16
	ldp x3, x4, [sp], #16
	ldp x5, x8, [sp], #16
	ldp x29, x30, [sp], #16
	ret

//the m_gpio_out macro will setup the X2 and X3
//X2 needs to be set to point to the required pin in the pin table in the data segment
//X3 needs to be set to point to the required pin in the pin table in the data segment
//X10 contains the virtual memory addess (obtained by sys_gpioMap
//RETURNS: void
sys_gpioSetDirectionOut:
	stp x29, x30, [sp, #-16]!
        stp x5, x8, [sp, #-16]!
        stp x3, x4, [sp, #-16]!
        stp x1, x2, [sp, #-16]!

	ldr w2, [x2]
	ldr w1, [x10, x2]					//get the address of the register from the table

	add x3, x3, #4						//offset to get amount to shift from table
	ldr w3, [x3]						//amount to shift

	mov x0, #0b111						//mask to 3 bits of status register
	lsl x0, x0, x3						//get the bit for this pin to set
	bic x1, x1, x0						//clear those 3 bits in x0
	mov x0, #1						//one bit to set the register to outpu, that needs to be shifted in place next
	lsl x0, x0, x3						//shit the bit that sets output
	orr x1, x1, x0						//set the gpio bit that indicates output
	str w1, [x10, x2]					//save the regsiter and now the pin is in output
	
	ldp x1, x2, [sp], #16
        ldp x3, x4, [sp], #16
        ldp x5, x8, [sp], #16
        ldp x29, x30, [sp], #16
	ret

//the m_gpio_out macro will setup the X2 and X3
//X0 is 1 or 0 (1 is on and 0 is off)
//X1 needs to be set to point to the required pin in the pin table in the data segment
//X10 contains the virtual memory addess (obtained by sys_gpioMap
sys_gpioValue:	
	stp x29, x30, [sp, #-16]!
        stp x2, x3, [sp, #-16]!
        stp x0, x1, [sp, #-16]!
	
	mov x2, x10						//set x2 to gpio regs address
	cmp x0, #0						//if x0 is 0 than set bit to off
	beq 1f
	add x2, x2, #BitSetOffset				//set bit to off by adding the bit offset
	b 2f
1:
	add x2, x2, #BitClearOffset				//sets the bit off
	
2:
	mov x0, #1
	add x1, x1, #8
	ldr w1, [x1]
	lsl x0, x0, x1
	str w0, [x2]

	ldp x0, x1, [sp], #16
	ldp x2, x3, [sp], #16
        ldp x29, x30, [sp], #16
        ret

	

.data
sys_gpioMap_error_openfile:	.asciz "ERROR: Could not open /dev/mem, please try to run with sudo\n"

sys_gpioMap_error_mmap:		.asciz "ERROR: mmap() failed\n"

sys_gpioMap_dev_mem:		.asciz "/dev/mem"

sys_gpioMap_gpio_addr:		.dword 0xfe200000		//physical address of GPIO 

p2:				.word 0
				.word 6
				.word 2

p3:				.word 0
				.word 9
				.word 3

p4:				.word 0
				.word 12
				.word 4

p5:				.word 0
				.word 15
				.word 5

p6:                             .word 0
                                .word 18
                                .word 6

p17:
				.word 4
				.word 21
				.word 17
p22:
				.word 8
				.word 6
				.word 22
p27:
				.word 8
				.word 21
				.word 27
