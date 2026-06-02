# Tasks: cleanup-command

## 1. Extract shared helpers

- [ ] 1.1 Create `internal/workspace/workspace.go` with exported `WindowName()` and `WorkspaceDir()` functions
- [ ] 1.2 Update `internal/spawn/spawn.go` to import and use `workspace.WindowName()` and `workspace.WorkspaceDir()`
- [ ] 1.3 Verify `make test`, `make build`, `make lint` pass after refactor

## 2. Cleanup package

- [ ] 2.1 Create `internal/cleanup/cleanup.go` with `Cleanup(changeName string) error` function
- [ ] 2.2 Implement tolerant tmux window kill (skip if not found)
- [ ] 2.3 Implement tolerant jj workspace forget (skip if not found)
- [ ] 2.4 Implement tolerant workspace directory removal (skip if not found)
- [ ] 2.5 Print summary of what was cleaned up vs skipped

## 3. CLI integration

- [ ] 3.1 Add `cleanup` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [ ] 3.2 Wire cleanup subcommand to call `cleanup.Cleanup()`

## 4. Tests

- [ ] 4.1 Add unit tests for cleanup (window name and workspace dir via shared package)
- [ ] 4.2 Verify `make test`, `make build`, `make lint` all pass

## 5. Documentation

- [ ] 5.1 Update README CLI section with `cleanup` command
