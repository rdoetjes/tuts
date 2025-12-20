build:
	go build -o bin/netcalc -ldflags="-s -w" -trimpath
	(cd netcalc_api && go build -o ../bin/netcalc_api -ldflags="-s -w" -trimpath)

clean:
	rm -f bin/netcalc
	rm -f bin/netcalc_api

test:
	$(MAKE) build
	go test ./... -v -count=1 || exit 1

install:
ifeq ($(OS), Windows_NT)
	if exist bin\netcalc.exe (
		copy bin\netcalc.exe C:\Windows\System32\
	) else (
		$(MAKE) build
		copy bin\netcalc.exe C:\Windows\System32\
	)
else
	if [ -f ./bin/netcalc ]; then \
		cp bin/netcalc /usr/local/bin/; \
	else \
		$(MAKE) build; \
		cp bin/netcalc /usr/local/bin/; \
	fi
endif
