VERSION := $(shell cat VERSION)

.PHONY: build test lint

build:
	go build -ldflags "-X main.version=$(VERSION)" -o jjay ./cmd/jjay

test:
	go test ./...

lint:
	go vet ./...
