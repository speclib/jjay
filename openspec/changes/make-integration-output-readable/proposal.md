## Why

`make test-integration` runs `go test -tags integration -v ./...`, which produces ~340 lines of output where the genuine signal — which scenarios passed, how long they took, what failed — is buried under raw subprocess banners (jj's "Working copy (@) now at…", OpenSpec's setup ASCII and "Getting started" blurb) interleaved with `go test -v` machinery. There are no visual boundaries between scenarios, no run summary, and no color. A passing run is indistinguishable at a glance from a failing one.

Bean: [jjay-7rol](../../../.beans/jjay-7rol--make-test-integration-output-more-readable.md)

## What Changes

- Adopt **gotestsum** as the integration-test output formatter. The `test-integration` Makefile target pipes through `gotestsum --format standard-verbose -- -tags integration ./...`, gaining per-scenario PASS/FAIL lines with **color** and a `DONE N tests in Xs` **stats** footer. `standard-verbose` is chosen deliberately so subprocess banners remain visible on every run (not only on failure).
- Add `gotestsum` to the Nix flake devShell `buildInputs` so it is available in the dev environment.
- Route subprocess banners through the test logger: `runIn()` in `test/integration/helpers_test.go` captures each subprocess's combined stdout/stderr and emits it via `t.Logf` instead of streaming raw to `os.Stdout`. This makes jj/OpenSpec output **nest under the correct subtest heading** (the "clean headings vs body" fix) rather than floating free.
- Preserve a **graceful fallback**: plain `go test -tags integration ./...` (no gotestsum) MUST still run and report correctly, so the suite is not hard-coupled to the formatter being installed.

The four bean bullets map to: cleanup → banner routing via `t.Log`; clean headings vs body → gotestsum scenario lines + nested logs; stats → gotestsum `DONE` footer; colors → gotestsum built-in color.

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `test-infrastructure`: The Makefile target requirement is amended so `make test-integration` formats output via gotestsum; a new requirement is added covering readable integration-test output (per-scenario headings, nested subprocess detail, a run summary, color) and the gotestsum-absent fallback.

## Impact

- **Makefile**: `test-integration` target changes its command (and a `test-integration` target should be declared in `.PHONY`).
- **flake.nix**: devShell `buildInputs` gains `gotestsum`.
- **test/integration/helpers_test.go**: `runIn()` switches from `cmd.Stdout/Stderr = os.Stdout/os.Stderr` to capturing combined output and logging it via `t.Logf`.
- **Dependencies**: adds `gotestsum` as a dev-time tool (Nix-provided; not a Go module dependency, not shipped in the binary).
- No production code, no CLI behavior, no user-facing binary change.
