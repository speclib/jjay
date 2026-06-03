# Tasks: spawn-config

## 1. Configuration structs

- [x] 1.1 Add `SpawnOptions` struct to `internal/spawn/spawn.go` with `Agent`, `Session`, `WorkspaceRoot` fields
- [x] 1.2 Add `CleanupOptions` struct to `internal/cleanup/cleanup.go` with `Session`, `WorkspaceRoot` fields
- [x] 1.3 Update `workspace.WorkspaceDir()` to accept optional root parameter
- [x] 1.4 Refactor `Spawn(changeName)` to `Spawn(changeName, opts SpawnOptions)`
- [x] 1.5 Refactor `Cleanup(changeName)` to `Cleanup(changeName, opts CleanupOptions)`

## 2. Agent command templating

- [x] 2.1 Implement `{change}` and `{wsdir}` placeholder substitution in agent command
- [x] 2.2 Update `setupPanes()` to use the configured agent command instead of hardcoded claude

## 3. Tmux session targeting

- [x] 3.1 Update all tmux commands in spawn to use session-prefixed targets when `Session` is set
- [x] 3.2 Update all tmux commands in cleanup to use session-prefixed targets when `Session` is set

## 4. CLI flags

- [x] 4.1 Add `--agent`, `--session`, `--workspace-root` flags to spawn subcommand in `cmd/jjay/main.go`
- [x] 4.2 Add `--session`, `--workspace-root` flags to cleanup subcommand
- [x] 4.3 Wire flags to options structs

## 5. Fake agent

- [x] 5.1 Create `testdata/fake-agent.sh` that creates `agent-was-here.txt` and exits

## 6. Integration test

- [x] 6.1 Create `test/integration/full_lifecycle_test.go` with `//go:build integration` tag
- [x] 6.2 Implement test setup: temp dir, jj repo, openspec change, tmux session
- [x] 6.3 Test spawn: verify tmux window, jj workspace, directory, agent marker file
- [x] 6.4 Test cleanup: verify all resources removed
- [x] 6.5 Implement teardown: kill tmux session, remove temp dirs

## 7. Verification

- [x] 7.1 Verify `make test` passes (unit tests, no integration tag)
- [x] 7.2 Verify `make build` and `make lint` pass
- [x] 7.3 Run integration test: `go test -tags integration ./test/integration/ -v`
