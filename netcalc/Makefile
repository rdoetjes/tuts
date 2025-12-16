build:
	go build -o bin/netcalc -ldflags="-s -w" -trimpath

clean:
	rm bin/netcalc

test:
	$(MAKE) build
	go test ./... -v -count=1|| exit 1

install:
	if [ -f ./bin/netcalc ]; then \
		cp bin/netcalc /usr/local/bin/; \
	else \
		echo "Run make before running make install"; \
	fi
