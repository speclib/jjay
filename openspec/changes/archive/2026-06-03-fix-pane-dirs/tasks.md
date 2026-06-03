# Tasks: fix-pane-dirs

## 1. Fix tmux pane directories

- [x] 1.1 Update `createWindow()` to pass `-c wsDir` to `tmux new-window`
- [x] 1.2 Update `setupPanes()` to pass `-c wsDir` to `tmux split-window`
- [x] 1.3 Remove `send-keys cd` for right pane
- [x] 1.4 Remove `cd <wsDir> &&` prefix from agent command (window starts in wsDir)

## 2. Integration test

- [x] 2.1 Create `internal/spawn/spawn_integration_test.go` with `//go:build integration`
- [x] 2.2 Implement test setup: temp dir, jj repo init, openspec change dir, tmux session
- [x] 2.3 Test spawn: verify tmux window, jj workspace, directory exist
- [x] 2.4 Test spawn: assert both panes have correct working directory via `tmux display-message`
- [x] 2.5 Test spawn: verify fake agent marker file exists in workspace
- [x] 2.6 Test cleanup: verify all resources removed
- [x] 2.7 Implement teardown: kill tmux session, remove temp dirs

## 3. Verification

- [x] 3.1 Verify `make test` passes (unit tests)
- [x] 3.2 Verify `make build` and `make lint` pass
- [x] 3.3 Run `make test-integration` and verify all assertions pass
