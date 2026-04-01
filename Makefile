build:
	CGO_ENABLED=0 go build -o crossroads ./cmd/crossroads/

run: build
	./crossroads

test:
	go test ./...

clean:
	rm -f crossroads

.PHONY: build run test clean
