# Tasks: spawn-command

## 1. Spawn package

- [x] 1.1 Create `internal/spawn/spawn.go` with `Spawn(changeName string) error` function
- [x] 1.2 Implement precondition checks: tmux session, openspec change exists, workspace doesn't exist, window name not taken
- [x] 1.3 Implement workspace creation via `jj workspace add ./<change>`
- [x] 1.4 Implement tmux window creation via `tmux new-window -n "ws:<change>"`
- [x] 1.5 Implement two-pane layout: agent in left pane, shell in right pane

## 2. CLI integration

- [x] 2.1 Add `spawn` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [x] 2.2 Wire spawn subcommand to call `spawn.Spawn()`

## 3. Tests

- [x] 3.1 Add unit tests for precondition check logic (tmux env var detection, argument validation)
- [x] 3.2 Verify `make test`, `make build`, `make lint` all pass

## 4. Documentation

- [x] 4.1 Update README CLI preview section with `spawn` command details
