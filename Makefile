VERSION := $(shell cat VERSION)

.PHONY: build test test-spawn test-integration lint coverage coverage-unit badge clean-tests

build:
	go build -ldflags "-X main.version=$(VERSION)" -o jjay ./cmd/jjay

test:
	go test ./...

test-spawn:
	go test -tags integration ./test/integration/ -v -run TestSpawn$$

test-integration: clean-tests
	@if command -v gotestsum >/dev/null 2>&1; then \
		gotestsum --no-color=false --format testname -- -tags integration ./...; \
	else \
		go test -tags integration ./...; \
	fi

clean-tests:
	-tmux list-sessions -F '#{session_name}' 2>/dev/null | grep '^jjay-test-' | xargs -r -n1 tmux kill-session -t
	rm -rf /tmp/jjay-test-* /tmp/jjay-merge-test-*

lint:
	go vet ./...

# coverage runs the WHOLE suite including integration tests, with -coverpkg=./...
# so coverage is attributed across every package regardless of which test package
# exercised it (test/integration drives internal/spawn, internal/cleanup, etc.
# in-process; without -coverpkg those don't count). Requires tmux + jj on PATH
# (it runs the integration suite) and sweeps test debris first via clean-tests.
# Plain `go test` (not gotestsum) so -coverprofile stays simple.
coverage: clean-tests
	go test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$NF}')"

# coverage-unit is the tmux/jj-free fallback (bare CI): same whole-repo -coverpkg
# attribution but WITHOUT -tags integration, so integration-only coverage is
# simply absent and no tmux/jj is required.
coverage-unit:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$NF}')"

badge: coverage
	$(eval PCT := $(shell go tool cover -func=coverage.out | grep total | awk '{print $$NF}' | tr -d '%'))
	$(eval COLOR := $(shell echo "$(PCT)" | awk '{if ($$1 >= 80) print "brightgreen"; else if ($$1 >= 60) print "yellow"; else print "red"}'))
	sed -i'' -e 's|https://img.shields.io/badge/coverage-[^)]*|https://img.shields.io/badge/coverage-$(PCT)%25-$(COLOR)|' README.md
