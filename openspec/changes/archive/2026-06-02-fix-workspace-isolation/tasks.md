# Tasks: fix-workspace-isolation

## 1. Fix spawn sequence

- [x] 1.1 Add `jj new` step after precondition checks but before workspace creation in `internal/spawn/spawn.go`
- [x] 1.2 Change `--revision @` to `--revision @-` in `createWorkspace()`
- [x] 1.3 Update success message to inform user that main workspace is now on a fresh change

## 2. Tests

- [x] 2.1 Update existing tests if affected by the new step
- [x] 2.2 Verify `make test`, `make build`, `make lint` all pass

## 3. Documentation

- [x] 3.1 Update spawn spec in main specs to reflect new behavior
