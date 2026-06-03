# Tasks: fix-pane-dirs

## 1. Fix tmux pane directories

- [ ] 1.1 Update `createWindow()` to pass `-c wsDir` to `tmux new-window`
- [ ] 1.2 Update `setupPanes()` to pass `-c wsDir` to `tmux split-window`
- [ ] 1.3 Remove `send-keys cd` for right pane
- [ ] 1.4 Remove `cd <wsDir> &&` prefix from agent command (window starts in wsDir)

## 2. Integration test

- [ ] 2.1 Create `internal/spawn/spawn_integration_test.go` with `//go:build integration`
- [ ] 2.2 Implement test setup: temp dir, jj repo init, openspec change dir, tmux session
- [ ] 2.3 Test spawn: verify tmux window, jj workspace, directory exist
- [ ] 2.4 Test spawn: assert both panes have correct working directory via `tmux display-message`
- [ ] 2.5 Test spawn: verify fake agent marker file exists in workspace
- [ ] 2.6 Test cleanup: verify all resources removed
- [ ] 2.7 Implement teardown: kill tmux session, remove temp dirs

## 3. Verification

- [ ] 3.1 Verify `make test` passes (unit tests)
- [ ] 3.2 Verify `make build` and `make lint` pass
- [ ] 3.3 Run `make test-integration` and verify all assertions pass
