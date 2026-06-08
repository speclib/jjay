VERSION := $(shell cat VERSION)

.PHONY: build test test-spawn test-integration lint coverage badge clean-tests

build:
	go build -ldflags "-X main.version=$(VERSION)" -o jjay ./cmd/jjay

test:
	go test ./...

test-spawn:
	go test -tags integration ./test/integration/ -v -run TestSpawn$$

test-integration: clean-tests
	go test -tags integration -v ./...

clean-tests:
	-tmux list-sessions -F '#{session_name}' 2>/dev/null | grep '^jjay-test-' | xargs -r -n1 tmux kill-session -t
	rm -rf /tmp/jjay-test-* /tmp/jjay-merge-test-*

lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$NF}')"

badge: coverage
	$(eval PCT := $(shell go tool cover -func=coverage.out | grep total | awk '{print $$NF}' | tr -d '%'))
	$(eval COLOR := $(shell echo "$(PCT)" | awk '{if ($$1 >= 80) print "brightgreen"; else if ($$1 >= 60) print "yellow"; else print "red"}'))
	sed -i'' -e 's|https://img.shields.io/badge/coverage-[^)]*|https://img.shields.io/badge/coverage-$(PCT)%25-$(COLOR)|' README.md
