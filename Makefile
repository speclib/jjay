.PHONY: build test lint

build:
	go build -o jjay ./cmd/jjay

test:
	go test ./...

lint:
	go vet ./...
