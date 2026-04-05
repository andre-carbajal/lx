.PHONY: build test lint clean

build:
	go build ./cmd/lx

test:
	go test ./... -v -coverprofile=coverage.out

lint:
	golangci-lint run

clean:
	rm -f coverage.out
	go clean
