## Context

jjay has a Makefile with `build`, `test`, `test-integration`, and `lint` targets. There is no coverage measurement or reporting. The project uses standard Go tooling (`go test`, `go tool cover`).

## Goals / Non-Goals

**Goals:**
- Measure unit test code coverage with a single `make coverage` command
- Generate an HTML coverage report for local browsing
- Print the coverage percentage to stdout
- Provide `make badge` to patch the README with a shields.io coverage badge
- Keep generated coverage artifacts out of version control

**Non-Goals:**
- CI integration or automated badge updates (local-only for now)
- Integration test coverage (unit tests only — integration tests spawn real processes and would distort numbers)
- Coverage thresholds or enforcement

## Decisions

### Coverage scope: unit tests only
Run `go test -coverprofile=coverage.out ./...` without `-tags integration`. Integration tests shell out to jj/tmux and measure orchestration glue, not meaningful branch coverage. Unit tests give a cleaner signal.

### Badge mechanism: shields.io static badge URL
Use `https://img.shields.io/badge/coverage-XX%25-COLOR` where COLOR is green/yellow/red based on thresholds (≥80 green, ≥60 yellow, <60 red). The percentage is extracted from `go tool cover -func` total line via `grep` + `awk`, then `sed` patches the README in-place. No external tools required beyond coreutils.

### Subprocess invocations

`make coverage`:
1. `go test -coverprofile=coverage.out ./...` → produces `coverage.out`
2. `go tool cover -html=coverage.out -o coverage.html` → produces `coverage.html`
3. `go tool cover -func=coverage.out | grep total | awk '{print $$NF}'` → prints percentage (e.g. `74.2%`)

`make badge`:
1. Runs `make coverage` first (dependency)
2. Extracts percentage from `coverage.out` using same `go tool cover -func` pipeline
3. Determines color based on numeric value
4. Uses `sed -i` to replace the badge image URL in `README.md`

### Badge placement
The badge goes on the line immediately after the closing `</p>` of the hero image block, before the `# jjay` heading.

## Risks / Trade-offs

- **`sed -i` portability**: GNU sed and BSD sed differ on `-i` syntax. Using `sed -i'' -e` works on both. Alternatively, use a temp file approach. Since the bean says local-only and the project targets Linux/macOS, we'll use the portable form.
- **Badge goes stale**: Since it's manually updated, the badge will only be as current as the last `make badge` run. Acceptable for now; CI can automate later.
