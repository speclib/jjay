# Tasks: spawn-config

## 1. Configuration structs

- [ ] 1.1 Add `SpawnOptions` struct to `internal/spawn/spawn.go` with `Agent`, `Session`, `WorkspaceRoot` fields
- [ ] 1.2 Add `CleanupOptions` struct to `internal/cleanup/cleanup.go` with `Session`, `WorkspaceRoot` fields
- [ ] 1.3 Update `workspace.WorkspaceDir()` to accept optional root parameter
- [ ] 1.4 Refactor `Spawn(changeName)` to `Spawn(changeName, opts SpawnOptions)`
- [ ] 1.5 Refactor `Cleanup(changeName)` to `Cleanup(changeName, opts CleanupOptions)`

## 2. Agent command templating

- [ ] 2.1 Implement `{change}` and `{wsdir}` placeholder substitution in agent command
- [ ] 2.2 Update `setupPanes()` to use the configured agent command instead of hardcoded claude

## 3. Tmux session targeting

- [ ] 3.1 Update all tmux commands in spawn to use session-prefixed targets when `Session` is set
- [ ] 3.2 Update all tmux commands in cleanup to use session-prefixed targets when `Session` is set

## 4. CLI flags

- [ ] 4.1 Add `--agent`, `--session`, `--workspace-root` flags to spawn subcommand in `cmd/jjay/main.go`
- [ ] 4.2 Add `--session`, `--workspace-root` flags to cleanup subcommand
- [ ] 4.3 Wire flags to options structs

## 5. Fake agent

- [ ] 5.1 Create `testdata/fake-agent.sh` that creates `agent-was-here.txt` and exits

## 6. Integration test

- [ ] 6.1 Create `internal/spawn/spawn_integration_test.go` with `//go:build integration` tag
- [ ] 6.2 Implement test setup: temp dir, jj repo, openspec change, tmux session
- [ ] 6.3 Test spawn: verify tmux window, jj workspace, directory, agent marker file
- [ ] 6.4 Test cleanup: verify all resources removed
- [ ] 6.5 Implement teardown: kill tmux session, remove temp dirs

## 7. Verification

- [ ] 7.1 Verify `make test` passes (unit tests, no integration tag)
- [ ] 7.2 Verify `make build` and `make lint` pass
- [ ] 7.3 Run integration test: `go test -tags integration ./internal/spawn/ -v`
