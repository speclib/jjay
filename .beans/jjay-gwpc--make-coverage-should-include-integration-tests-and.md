---
# jjay-gwpc
title: make coverage should include integration tests (and clarify badge vs coverage)
status: todo
type: task
priority: normal
created_at: 2026-06-05T00:39:22Z
updated_at: 2026-06-12T12:21:58Z
parent: jjay-hjjg
---

initial bean was: jjay-7rol

`make coverage` currently runs unit tests only — `go test -coverprofile ./...` with NO `-tags integration`. So integration-tested code (the whole spawn → cleanup lifecycle, merge scenarios, the status subtest) contributes 0 to the reported coverage %, and the README badge understates real coverage.

Wanted: coverage should include integration tests so the badge reflects unit + integration.

## Caveat (why it's not a one-liner)
Integration tests require tmux AND jj on PATH (`test/integration/helpers_test.go` calls `requireCmd(t, "tmux")`, creates a real jj repo + tmux session). So `go test -tags integration -coverprofile` only works in an environment that has them (local dev, nix devShell) — NOT a bare CI runner. Options to weigh:
- combined coverage target that assumes tmux+jj (devShell/local only); keep a unit-only fallback for CI.
- run integration with coverage in CI by ensuring tmux+jj are installed in the workflow.
- merge two coverprofiles (unit + integration) into one total.

## Also (the confusion that surfaced this)
`make coverage` does NOT touch the README — only `make badge` (which depends on coverage) patches the shields.io badge. Easy to conflate. Either fold badge-patching into coverage, or document clearly that `badge` is the README-updating target.

Found while prepping a release: the user expected `make coverage` to run integration tests and update the README; it does neither on its own.
Related: jjay-znb0 (created the current targets, archived), jjay-7rol (test-integration output readability).



## CRITICAL refinement (2026-06-12): -coverpkg is the other (bigger) half

Two separate measurement blind spots, not one. The `-tags integration` fix alone is INSUFFICIENT:

1. `-tags integration` — includes the integration tests. Fixes packages whose tests live IN-package (e.g. internal/merge: 4.8% → 81.6% with the tag).

2. `-coverpkg=./...` — REQUIRED for packages whose tests live in test/integration (a different package). By default `go test -cover ./internal/spawn/` only attributes coverage from tests inside internal/spawn. The spawn integration tests DO call spawn.Spawn() in-process (test/integration/helpers_test.go:139) — but because they live in `package integration`, that exercise is not attributed back to internal/spawn. So spawn shows 5.8% even WITH the tag. Adding `-coverpkg=./internal/spawn/...` attributes it: spawn jumps to ~70% (setupPanes 80%, tmuxTarget 67%, resolveAgentCommand 100%).

Proven numbers (with `-tags integration -coverpkg=./...`):
- internal/merge:  4.8% → ~82%
- internal/spawn:  5.8% → ~70%
Neither package is under-tested; the default measurement was blind in two ways.

## The actual fix (both flags)
```makefile
coverage:
	go test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$NF}')"
```
- `-tags integration` → runs integration tests (needs tmux+jj on PATH; see caveat above re CI).
- `-coverpkg=./...` → attributes coverage of ALL packages regardless of which test package exercised them — this is what makes test/integration's exercise of spawn/cleanup/status count.

Note: `-coverpkg=./...` also slightly changes the total (counts every package as a denominator, including main/cmd). Decide whether the badge should reflect that whole-repo number or a curated subset.
