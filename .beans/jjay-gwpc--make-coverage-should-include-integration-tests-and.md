---
# jjay-gwpc
title: make coverage should include integration tests (and clarify badge vs coverage)
status: todo
type: task
priority: normal
created_at: 2026-06-05T00:39:22Z
updated_at: 2026-06-05T09:36:47Z
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
