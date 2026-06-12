## Why

`make coverage` reports a coverage number that badly understates reality, in **two independent ways** (both proven this session — bean [jjay-gwpc](../../.beans/jjay-gwpc--make-coverage-should-include-integration-tests-and.md)):

1. **No `-tags integration`.** Integration-tested code (the whole spawn → cleanup lifecycle, merge scenarios, the status subtest) contributes **0** to the number. `internal/merge` reads as 4.8% but is ~82% with the tag.
2. **No `-coverpkg`.** Packages whose tests live in `test/integration` (a *different* package) get no attribution by default — `go test -cover ./internal/spawn/` only counts tests *inside* `internal/spawn`. The spawn integration tests call `spawn.Spawn()` in-process but from `package integration`, so spawn reads 5.8% even *with* the tag; `-coverpkg=./internal/spawn/...` lifts it to ~70%.

Neither package is under-tested — the measurement was blind twice. The README badge (patched by `make badge`, which depends on `coverage`) inherits the lie.

## What Changes

- **`make coverage` runs integration + whole-repo attribution:** `go test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...`. The badge then reflects unit **and** integration coverage across every package (`cmd/main` included — the honest "how much of the code is tested" number).
- **`make coverage` sweeps its own debris:** it depends on `clean-tests` (like `test-integration` already does), since `-tags integration` spawns real `jjay-test-*` tmux sessions + temp dirs.
- **New `make coverage-unit` fallback:** `go test -coverpkg=./... -coverprofile=coverage.out ./...` (no `-tags integration`) for environments without tmux/jj on PATH (bare CI). The integration `coverage` is the default and the one the badge uses locally.
- **`make badge` is documented as the README-updating target** (it depends on `coverage`); `coverage` itself only computes + prints. This resolves the "I ran coverage but the README didn't change" confusion.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `test-infrastructure`: the `coverage` target SHALL run integration tests with whole-repo `-coverpkg` attribution and sweep test debris; a `coverage-unit` target SHALL provide a tmux/jj-free fallback; the `coverage`/`badge` split (compute vs README-patch) SHALL be documented.

## Impact

- **Files**: `Makefile` (`coverage` gains `-tags integration -coverpkg=./...` + `clean-tests` dep; new `coverage-unit`; `.PHONY` updated). README badge value changes to the true whole-repo number on the next `make badge`. No production-code or test changes.
- **Behavior**: `make coverage` now requires tmux + jj on PATH (it runs the integration suite) — documented; `coverage-unit` is the escape hatch. The badge number will shift (merge/spawn etc. now counted; `cmd/main` now in the denominator).
- **Relation**: builds on `clean-tests` (jjay-zgqx, archived) and the integration targets (jjay-znb0, archived). Pairs with [jjay-bi7i](../../.beans/jjay-bi7i--release-script-should-run-tests-coverage-and-commi.md) (release script runs coverage + commits the badge) — bi7i should call this `coverage`.
- **ADRs**: none — a build-tooling change, no architectural decision. (The `-coverpkg`/denominator choice is recorded in design.md.)
- **Beans**: gwpc → in-progress, linked here.

## Deferred / Out of Scope

- **CI wiring** to install tmux+jj so integration coverage runs in CI — separate concern; `coverage-unit` is the bare-CI fallback for now.
- **Curated-subset denominator** — rejected in favor of the honest whole-repo number; revisit only if `cmd/main` noise proves misleading.
