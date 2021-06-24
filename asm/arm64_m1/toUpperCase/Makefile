TOOLPATH = $(shell dirname $(shell xcodebuild -find clang))

helloworld: helloworld.o
	$(TOOLPATH)/ld -o helloworld helloworld.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64

helloworld.o: helloworld_2.asm
	$(TOOLPATH)/as -o helloworld.o helloworld_2.asm

clean:
	rm *.o
	rm helloworld