.PHONY: build test lint clean

build:
	go build ./cmd/lx

test:
	go test ./... -v

lint:
	golangci-lint run

clean:
	go clean
