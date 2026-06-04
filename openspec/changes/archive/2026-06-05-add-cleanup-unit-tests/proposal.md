## Why

`internal/cleanup` is the only command package with **no test file** — `spawn`, `merge`, `session`, `status`, and `workspace` all have unit tests, but `cleanup.go` has none. Cleanup tears down real resources (kills a tmux window, forgets a jj workspace, removes a directory), so its tolerance logic — "skip missing resources, never fail" — and its path resolution are exactly the kind of behavior a regression could silently break. Bean [jjay-kyx6](../../.beans/jjay-kyx6--cleanup-has-not-unit-tests.md) flags the gap.

## What Changes

- **Add `internal/cleanup/cleanup_test.go`** with unit tests covering the parts of cleanup that are testable without a live tmux server or jj repo:
  - `removeDirectory` — removes an existing dir; tolerates a missing dir; resolves the workspace path via the `workspace` package (honors `WorkspaceRoot`).
  - `tmuxTarget` — pure function: returns `window` with no session, `session:window` with a session.
  - `killWindow` / `forgetWorkspace` — the tolerance branches: when tmux/jj are unavailable or the resource is absent, the function returns without error (mirrors how `merge`/`spawn` exercise their "not found" paths).
  - `Cleanup` — orchestration: runs all three steps and returns nil even when every resource is missing (tolerant by design).
- **No production code changes** — this is a test-only change. `cleanup.go` behavior is unchanged.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `cleanup`: add a unit-test coverage requirement for the cleanup command (no behavior change).

## Impact

- **Code**: new `internal/cleanup/cleanup_test.go` only.
- **Tools**: tests use a temp dir for `removeDirectory`; tolerance tests rely on absent tmux/jj resources (no live server needed), consistent with `merge`/`spawn` unit tests.
- **Coverage**: raises `internal/cleanup` from 0% toward parity with the other command packages; surfaces in `make coverage`.
- **ADRs**: none (test-only; no architectural decision).
- **Beans**: kyx6 → in-progress, linked here.
