#include <iostream>
#include <pigpio.h>
#include <unistd.h>
#include "i2c_wrapper.h"

void scanLeds(int dev1, int dev2, unsigned int sleepTimeOut){
	unsigned int i = 1;
	while(i != 0x80000000){
		i = i << 1;
		set32BitValue(dev1, dev2, i);
		usleep(sleepTimeOut);
	}
	
	i = 0x80000000;		//set i back because it's been shifted over 
	while(i >= 1){
		i = i >> 1;
		set32BitValue(dev1, dev2, i);
		usleep(sleepTimeOut);
	}
}

int main(){
	gpioInitialise();

	int m1 = i2cOpen(1, 0x20, 0);
	int m2 = i2cOpen(1, 0x21, 0);

	if (m1 < 0 || m2 < 0) {
    std::cout << "Could not open device!" << std::endl;
		exit(1);
	}

	setAllPortsToOutput(m1, m2);

  while(1){
		scanLeds(m1, m2, 10000);
	}

	i2cClose(m1);
	i2cClose(m2);

	return 0;
}
