## Context

`make test-integration` runs `go test -tags integration -v ./...`. The output is ~340 lines, of which the actionable signal (scenario pass/fail, durations, failures) is roughly 20 lines. Two distinct noise sources compete in the same column:

1. **`go test -v` machinery** — `=== RUN`, `=== CONT`, per-package `PASS`/`ok` lines.
2. **Raw subprocess banners** — jj's "Working copy (@) now at…", "Created workspace in…", and OpenSpec's setup ASCII + "Getting started"/"Feedback" blurbs.

The second source is the dominant offender, and crucially it does **not** flow through Go's test event stream: `runIn()` in `test/integration/helpers_test.go` wires subprocess `Stdout`/`Stderr` straight to `os.Stdout`/`os.Stderr` (lines ~182–183). Any formatter that consumes only `go test -json` therefore cannot reattribute or fold that text — it bypasses the JSON events entirely.

This is the key architectural fact driving the design: a pure formatter is necessary but not sufficient. The banners must be redirected at their source (the test helper) for headings/body nesting to work.

## Goals / Non-Goals

**Goals:**
- Per-scenario, color-coded PASS/FAIL lines (clean headings).
- Subprocess detail nested under the scenario that produced it (body, not free-floating).
- A run summary with test count and total duration (stats).
- The integration suite still runs without the formatter installed (graceful fallback).

**Non-Goals:**
- Changing any production/CLI behavior or the shipped binary.
- Suppressing subprocess banners on green runs — the chosen format keeps them visible every run (explicit user preference).
- Reworking `make test` (unit tests) — only the integration target changes.
- Adding a Go module dependency — gotestsum is a dev-time tool, Nix-provided.

## Decisions

### Decision: Use gotestsum as the formatter
**Why:** It is a mature, widely-used wrapper over `go test -json` that gives colored per-test lines, multiple format presets, and a `DONE N tests in Xs` summary footer for free. It is packaged in nixpkgs, so it drops cleanly into the existing devShell `buildInputs` alongside `go`, `gopls`, `goreleaser`, `gum`.

**Alternatives considered:**
- *tparse* — similar, but table-oriented summary; less natural for the "stream the run, then summarize" feel here.
- *Home-grown formatter over `go test -json`* — zero dependency, but reinvents color/summary/format handling and still can't touch the banner noise. Not worth the code to own.
- *Status quo (`go test -v`)* — rejected; that's the problem.

### Decision: Use `--format standard-verbose`, not `testname`
**Why:** gotestsum's quieter formats (`testname`, `pkgname`) hide passing-test body output — they are quiet-on-green. The user explicitly wants subprocess banners visible on every run. `standard-verbose` preserves full per-test output while still adding color and the summary footer. This is a deliberate trade of brevity for always-on forensic detail.

**Trade-off:** Green runs remain long. Accepted per the stated preference; if it later grates, switching to `--format testname` is a one-word Makefile change.

### Decision: Route subprocess banners through `t.Logf` in `runIn()`
**Why:** This is what makes "headings vs body" actually work. Capturing each subprocess's combined output and emitting it via `t.Logf` (instead of `os.Stdout`) means:
- The text enters the `go test -json` event stream, so gotestsum sees it and renders it under the right test.
- It nests beneath the correct subtest (`spawn` / `status` / `cleanup`), giving the scenario→detail hierarchy.

Implementation shape: `runIn()` switches from `cmd.Stdout = os.Stdout; cmd.Stderr = os.Stderr` to `out, err := cmd.CombinedOutput()` and `t.Logf("%s %v:\n%s", name, args, out)`. The helper already takes `t *testing.T`, so no signature change is needed.

**Trade-off:** Output is logged after the subprocess completes rather than streamed live. For these short-lived helper commands this is immaterial, and it buys correct attribution.

### Decision: Graceful fallback via Makefile recipe guard
**Why:** The suite must not hard-fail when gotestsum is absent (e.g., a contributor outside the Nix shell). The `test-integration` recipe detects gotestsum on `PATH` and falls back to plain `go test -tags integration ./...`. A shell `command -v gotestsum` guard in the recipe keeps this self-contained in the Makefile.

## Risks / Trade-offs

- **gotestsum not installed outside Nix shell** → fallback recipe runs plain `go test`; documented in the Makefile recipe.
- **standard-verbose stays long on green** → accepted per user preference; trivially switchable later.
- **Banners logged post-completion, not streamed** → negligible for short helper commands; gains correct nesting.
- **CI color/TTY** → gotestsum auto-detects non-TTY and degrades to plain text; no special handling needed.

## Migration Plan

1. Add `gotestsum` to `flake.nix` devShell `buildInputs`.
2. Update the `test-integration` Makefile recipe (with `.PHONY` already covering it) to use gotestsum-with-fallback.
3. Change `runIn()` to capture combined output and emit via `t.Logf`.
4. Verify: run `make test-integration` inside the dev shell (formatted) and with gotestsum off `PATH` (fallback).

No rollback complexity — the change is confined to dev tooling and a test helper.

## Open Questions

- None. Decisions on formatter, format preset, banner routing, and fallback are settled.
