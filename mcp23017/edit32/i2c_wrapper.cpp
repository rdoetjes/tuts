#include "i2c_wrapper.h"

int setGpioOut(int dev, char port){
  if (port == 'A' || port == 'a') return i2cWriteWordData(dev, 0x00, 0x00);
  if (port == 'B' || port == 'b') return i2cWriteWordData(dev, 0x01, 0x00);
  return -1;
}
/*
 * Sets the 4 ports on the board to the 32 bit value set in value
 * dev1 needs to point to the opened i2c for address 20
 * dev2 needs to point to the opened i2c for address 21
 */
void set32BitValue(int dev1, int dev2, unsigned int value){
  char byte = value;
  i2cWriteWordData(dev1, 0x14, byte);

  byte = value >> 8;
  i2cWriteWordData(dev1, 0x15, byte);

  byte = value >> 16;
  i2cWriteWordData(dev2, 0x14, byte);

  byte = value >> 24;
  i2cWriteWordData(dev2, 0x15, byte);
}

void setAllPortsToOutput(int dev1, int dev2){
  //set the two ports of each of the MCP23017 to be outputs
  setGpioOut(dev1, 'A');
  setGpioOut(dev1, 'B');
  setGpioOut(dev2, 'A');
  setGpioOut(dev2, 'B');
}

