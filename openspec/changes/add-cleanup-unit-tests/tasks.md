## 1. Test scaffolding

- [ ] 1.1 Create `internal/cleanup/cleanup_test.go` (package `cleanup`, white-box) with a temp-dir helper using `os.MkdirTemp`.

## 2. Unit tests

- [ ] 2.1 `removeDirectory`: create a dir under a temp workspace-root, remove it via `removeDirectory(change, root)`, assert it's gone.
- [ ] 2.2 `removeDirectory`: missing-dir branch returns cleanly (no panic, dir absent).
- [ ] 2.3 `tmuxTarget`: assert `window` with no session and `session:window` with a session.
- [ ] 2.4 `killWindow` / `forgetWorkspace`: call for a nonexistent change; assert tolerant (no panic, returns).
- [ ] 2.5 `Cleanup`: all-missing resources with a temp empty `WorkspaceRoot` returns `nil`.

## 3. Verify

- [ ] 3.1 `go test ./internal/cleanup/` passes.
- [ ] 3.2 `make coverage` shows `internal/cleanup` is no longer 0%.
- [ ] 3.3 Full suite (`go test -tags integration ./...`) still green.

## 4. Beans

- [ ] 4.1 Set `jjay-kyx6` to `in-progress`; add `openspec-link` on archive.
