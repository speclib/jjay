# Tasks: spawn-config

## 1. Configuration structs

- [ ] 1.1 Add `SpawnOptions` struct to `internal/spawn/spawn.go` with `Agent`, `Session`, `WorkspaceRoot` fields
- [ ] 1.2 Add `CleanupOptions` struct to `internal/cleanup/cleanup.go` with `Session`, `WorkspaceRoot` fields
- [ ] 1.3 Update `workspace.WorkspaceDir()` to accept optional root parameter
- [ ] 1.4 Refactor `Spawn(changeName)` to `Spawn(changeName, opts SpawnOptions)`
- [ ] 1.5 Refactor `Cleanup(changeName)` to `Cleanup(changeName, opts CleanupOptions)`

## 2. Fix shell pane working directory (jjay-ps3d)

- [ ] 2.0 Use `tmux new-window -c <wsDir>` to set window starting directory instead of relying on send-keys cd
- [ ] 2.1 Use `tmux split-window -h -c <wsDir>` for right pane instead of send-keys cd
- [ ] 2.2 Remove the `send-keys cd` for the right pane (no longer needed)

## 3. Agent command templating

- [ ] 3.1 Implement `{change}` and `{wsdir}` placeholder substitution in agent command
- [ ] 3.2 Update `setupPanes()` to use the configured agent command instead of hardcoded claude

## 4. Tmux session targeting

- [ ] 4.1 Update all tmux commands in spawn to use session-prefixed targets when `Session` is set
- [ ] 4.2 Update all tmux commands in cleanup to use session-prefixed targets when `Session` is set

## 5. CLI flags

- [ ] 5.1 Add `--agent`, `--session`, `--workspace-root` flags to spawn subcommand in `cmd/jjay/main.go`
- [ ] 5.2 Add `--session`, `--workspace-root` flags to cleanup subcommand
- [ ] 5.3 Wire flags to options structs

## 6. Fake agent

- [ ] 6.1 Create `testdata/fake-agent.sh` that creates `agent-was-here.txt` and exits

## 7. Integration test

- [ ] 7.1 Create `internal/spawn/spawn_integration_test.go` with `//go:build integration` tag
- [ ] 7.2 Implement test setup: temp dir, jj repo, openspec change, tmux session
- [ ] 7.3 Test spawn: verify tmux window, jj workspace, directory, agent marker file
- [ ] 7.4 Test spawn: assert both panes have working directory set to workspace dir (via `tmux display-message -p -t <pane> '#{pane_current_path}'`)
- [ ] 7.5 Test cleanup: verify all resources removed
- [ ] 7.6 Implement teardown: kill tmux session, remove temp dirs

## 8. Verification

- [ ] 8.1 Verify `make test` passes (unit tests, no integration tag)
- [ ] 8.2 Verify `make build` and `make lint` pass
- [ ] 8.3 Run integration test: `go test -tags integration ./internal/spawn/ -v`
