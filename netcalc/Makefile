build:
	go build -o bin/netcalc -ldflags="-s -w" -trimpath

clean:
	rm bin/netcalc

test:
	$(MAKE) build
	go test ./lib -v || exit 1
	go test -v || exit 1
