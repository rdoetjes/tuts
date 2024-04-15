TOOLPATH = $(shell dirname $(shell xcodebuild -find clang))

32bit: 32bit.o
	$(TOOLPATH)/ld -o 32bit 32bit.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64

32bit.o: 32bit.asm
	$(TOOLPATH)/as -o 32bit.o 32bit.asm

clean:
	rm *.o
	rm 32bit
