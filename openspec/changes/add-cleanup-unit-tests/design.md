## Context

`internal/cleanup/cleanup.go` has `Cleanup`, `killWindow`, `forgetWorkspace`, `removeDirectory`, and `tmuxTarget`, but no test file. The other command packages establish the testing patterns to mirror: `merge/merge_test.go` exercises a "not found" parsing branch by calling against a nonexistent resource; `spawn/spawn_test.go` tests pure helpers and env-driven branches without a live server. Cleanup is tolerant by design (every step skips missing resources and never returns an error), which makes its branches straightforward to unit-test.

## Goals / Non-Goals

**Goals:**
- Unit tests for the parts of cleanup that need no live tmux/jj: `removeDirectory`, `tmuxTarget`, the tolerance branches, and `Cleanup`'s all-missing path.
- Match the existing unit-test style (table-free, `t.Run` where it clarifies, temp dirs via `os.MkdirTemp`).

**Non-Goals:**
- Changing any cleanup behavior.
- End-to-end teardown of a real tmux window + jj workspace — that path is already covered by the full-lifecycle integration test (`test/integration`, spawn → cleanup). This change is unit-level.
- A live-tmux/jj integration test for cleanup (the integration suite already drives cleanup).

## Decisions

- **New file `internal/cleanup/cleanup_test.go`**, package `cleanup` (white-box, so it can call unexported `tmuxTarget`/`removeDirectory`/`killWindow`/`forgetWorkspace`).
- **`removeDirectory`:** create a temp dir under a temp workspace-root, call `removeDirectory(change, root)`, assert it's gone; call again (or on a never-created name) to assert the missing-dir branch returns cleanly. Use the `WorkspaceRoot` override so the test controls the path and never touches the real `../<project>-workspaces`.
- **`tmuxTarget`:** direct assertions for the with/without-session cases (pure function).
- **Tolerance branches:** call `killWindow`/`forgetWorkspace` for a nonexistent change; assert they do not panic and return (no error type to check — they print and return), mirroring `merge`'s nonexistent-resource approach. Where tmux/jj are absent in CI, the early `exec.Command(...).Output()` error path is exercised.
- **`Cleanup` all-missing:** call `Cleanup("nonexistent-xyz", CleanupOptions{WorkspaceRoot: tempEmptyRoot})`; assert `err == nil`.

## Risks / Trade-offs

- **Environment sensitivity:** if a tmux server or jj repo *is* present in the test environment, `killWindow`/`forgetWorkspace` take the "found/exists?" branch instead of the error branch. Tests assert tolerance (no error / no panic) rather than a specific branch, so they pass either way. Use a change name guaranteed not to exist.
- **`removeDirectory` path safety:** always pass an explicit `WorkspaceRoot` pointing at a temp dir so a buggy test can never remove a real workspace.
- **Low risk overall:** test-only change; the worst case is a flaky test, mitigated by the tolerance-oriented assertions above.
