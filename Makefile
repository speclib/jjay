VERSION := $(shell cat VERSION)

.PHONY: build test test-spawn test-integration lint

build:
	go build -ldflags "-X main.version=$(VERSION)" -o jjay ./cmd/jjay

test:
	go test ./...

test-spawn:
	go test -tags integration ./test/integration/ -v -run TestSpawn$$

test-integration:
	go test -tags integration -v ./...

lint:
	go vet ./...
