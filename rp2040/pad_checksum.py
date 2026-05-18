import struct
import sys


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 pad_checksum.py <binary_file>")
        sys.exit(1)

    with open(sys.argv[1], "rb") as f:
        data = bytearray(f.read())

    # RP2040 Bootrom CRC32
    # Polynomial: 0x04C11DB7
    # Initial value: 0xFFFFFFFF
    # Final XOR: 0x00000000
    # Data is NOT reflected, but the result is.

    def crc32_rp2040(data):
        crc = 0xFFFFFFFF
        for byte in data:
            crc ^= byte << 24
            for _ in range(8):
                if crc & 0x80000000:
                    crc = ((crc << 1) ^ 0x04C11DB7) & 0xFFFFFFFF
                else:
                    crc = (crc << 1) & 0xFFFFFFFF
        return crc

    # Compute CRC over the first 252 bytes of boot2
    checksum = crc32_rp2040(data[:252])
    print(f"{checksum:#x}")

    # Write the checksum into the last 4 bytes of the 256-byte boot2 block
    struct.pack_into("<I", data, 252, checksum)
    with open(sys.argv[1], "wb") as f:
        f.write(data)


if __name__ == "__main__":
    main()
