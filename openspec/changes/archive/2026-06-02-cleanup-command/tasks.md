# Tasks: cleanup-command

## 1. Extract shared helpers

- [x] 1.1 Create `internal/workspace/workspace.go` with exported `WindowName()` and `WorkspaceDir()` functions
- [x] 1.2 Update `internal/spawn/spawn.go` to import and use `workspace.WindowName()` and `workspace.WorkspaceDir()`
- [x] 1.3 Verify `make test`, `make build`, `make lint` pass after refactor

## 2. Cleanup package

- [x] 2.1 Create `internal/cleanup/cleanup.go` with `Cleanup(changeName string) error` function
- [x] 2.2 Implement tolerant tmux window kill (skip if not found)
- [x] 2.3 Implement tolerant jj workspace forget (skip if not found)
- [x] 2.4 Implement tolerant workspace directory removal (skip if not found)
- [x] 2.5 Print summary of what was cleaned up vs skipped

## 3. CLI integration

- [x] 3.1 Add `cleanup` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [x] 3.2 Wire cleanup subcommand to call `cleanup.Cleanup()`

## 4. Tests

- [x] 4.1 Add unit tests for cleanup (window name and workspace dir via shared package)
- [x] 4.2 Verify `make test`, `make build`, `make lint` all pass

## 5. Documentation

- [x] 5.1 Update README CLI section with `cleanup` command
