# ADR-012: gotestsum for integration test output

**Status**: Proposed

## Context

`make test-integration` (`go test -tags integration -v ./...`) emits ~340 lines where the real signal is ~20. Two noise sources share one column: `go test -v` machinery and raw subprocess banners from jj/OpenSpec. The banners are the bigger problem and, critically, they are written to `os.Stdout` from the `runIn()` test helper — they never enter Go's `go test -json` event stream. So a JSON-consuming formatter alone cannot fold or reattribute them; the source must also change.

We need: per-scenario headings, subprocess detail nested under its scenario, a run summary, and color — without coupling the suite to a formatter being installed.

## Options Considered

- **gotestsum** — mature `go test -json` wrapper; colored per-test lines, format presets, `DONE N tests in Xs` summary; packaged in nixpkgs. Con: a dev-time tool dependency; can't touch banners by itself.
- **tparse** — similar wrapper, table-style summary. Con: less natural for a stream-then-summarize flow; same banner limitation.
- **Home-grown formatter** — zero dependency. Con: reinvents color/summary/format handling and still can't address banners.
- **Status quo (`go test -v`)** — no new tooling. Con: it is the problem.

## Decision

Adopt **gotestsum** with `--format testname --no-color=false`, added to the Nix devShell. `testname` is chosen over `standard-verbose` because the latter keeps all subprocess body on screen, leaving the output mostly uncolored plain text; `testname` emits one colored `✓/✗ Test` line per test (high color density, fast to scan) and folds subprocess banners away on passing tests, surfacing them only on failure where they are needed. gotestsum 1.13.0 defaults `--no-color` to `true`, so `--no-color=false` is passed to enable color on a TTY. Pair it with a source-level fix: `runIn()` captures each subprocess's combined output and emits it via `t.Logf`, so the text enters the JSON event stream and nests under the correct subtest (and is what shows up under a failing test). The Makefile recipe guards on `command -v gotestsum` and falls back to plain `go test -tags integration ./...` when absent.

## Consequences

- **Positive:** Colored per-scenario results, subprocess detail attributed to its scenario, and a run summary. No production/binary change. Fallback keeps the suite runnable without the tool.
- **Negative / accepted:** Subprocess banners are hidden on green runs (folded by `testname`) and appear only on failure; switchable back to `standard-verbose` with a one-word change if always-on detail is wanted. Adds a dev-time tool (Nix-provided, not a Go module dependency). Subprocess output is logged on completion rather than streamed live — immaterial for these short helper commands.
