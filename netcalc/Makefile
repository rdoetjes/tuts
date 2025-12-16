build:
	go build -o bin/netcalc -ldflags="-s -w" -trimpath

clean:
	rm bin/netcalc

test:
	$(MAKE) build
	go test ./... -v -count=1|| exit 1
