import board
import busio
from digitalio import Direction
from adafruit_mcp230xx.mcp23017 import MCP23017
i2c = busio.I2C(board.SCL, board.SDA)

mcp1 = MCP23017(i2c, address=0x20)
mcp2 = MCP23017(i2c, address=0x21)

p = 0
while True:
    for i in range (0,16):
        pin = mcp1.get_pin(i)
        pin.direction = Direction.OUTPUT
        if p % 2 == 1:
          pin.value = True
        else:
          pin.value = False

        pin = mcp2.get_pin(i)
        pin.direction = Direction.OUTPUT
        if p % 2 == 1:
          pin.value = False
        else:
          pin.value = True
    p += 1




