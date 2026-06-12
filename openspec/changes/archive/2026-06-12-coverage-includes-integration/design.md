## Context

`make coverage` runs `go test -coverprofile=coverage.out ./...` — unit-only, default attribution. Two blind spots (proven in jjay-gwpc):
- No `-tags integration` → integration-tested code counts 0 (merge: 4.8% vs ~82%).
- No `-coverpkg` → `test/integration` (a separate package) exercising `internal/spawn` etc. in-process isn't attributed (spawn: 5.8% vs ~70% with `-coverpkg`).

`test-integration` already depends on `clean-tests` and uses `gotestsum` when present; `coverage` should align with that integration-running shape. `badge` depends on `coverage` and patches the README; `coverage` alone does not.

## Goals / Non-Goals

**Goals:**
- `make coverage` = whole-repo, integration-inclusive coverage (`-tags integration -coverpkg=./...`), debris swept.
- `make coverage-unit` = same `-coverpkg` measurement, no integration, for tmux/jj-free environments.
- Document that `badge` (not `coverage`) updates the README.

**Non-Goals:**
- Wiring CI to install tmux/jj (separate; `coverage-unit` is the bare-CI fallback).
- A curated package subset for the denominator — whole-repo `./...` is the chosen, honest number.
- Changing any test or production code.

## Decisions

- **`coverage` flags:** `go test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...`. Whole-repo `-coverpkg` (chosen over a curated subset) — the true "how much code is tested" number, accepting that thin `cmd/main` is in the denominator.
- **`coverage: clean-tests` dependency.** `-tags integration` spawns real `jjay-test-*` sessions/temp dirs; reuse the existing sweep so coverage runs don't leak (consistent with `test-integration`).
- **`coverage-unit`** mirrors `coverage` minus `-tags integration`: `go test -coverpkg=./... -coverprofile=coverage.out ./...`. Same coverprofile output so `badge` can consume either. No `clean-tests` dep needed (no integration spawns).
- **`badge` unchanged** (already `badge: coverage`); just documented as the README-updating target. (Optional: `badge` could depend on `coverage-unit` for CI, but local/release use the integration `coverage` — leave `badge: coverage`.)
- **`.PHONY`** gains `coverage-unit`.

## Risks / Trade-offs

- **Badge number will drop/shift.** Whole-repo denominator now includes `cmd/main` (thin, hard to unit-test) — the number is *truer* but may look lower than a curated subset. Accepted (honest > flattering); revisit only if misleading.
- **`coverage` now needs tmux+jj.** A developer on a machine without them gets a failure; `coverage-unit` is the documented escape. Mention in README/Makefile help.
- **`gotestsum` vs coverage output.** `test-integration` pipes through `gotestsum`; `coverage` needs the raw coverprofile, so it should call `go test` directly (not via gotestsum) to keep `-coverprofile` simple. Keep coverage's invocation plain `go test`.
- **Coverage run is slower** (full integration suite). Acceptable for a coverage/badge step (not the hot path); `make test` stays unit-only and fast.
