# Proposal: jjay spawn command

**Change**: spawn-command
**Status**: proposed
**Bean**: [jjay-iwbt — spawn task](../../../.beans/jjay-iwbt--spawn-task.md)

## Why

The core value of jjay is replacing the repetitive manual workflow of creating jj workspaces, tmux windows, and launching AI agents. `jjay spawn` is the first and most important command — it does all three in one shot.

## What Changes

- Add `jjay spawn <change-name>` cobra subcommand
- Create jj workspace in current project directory
- Create tmux window in current session with two-pane layout
- Launch claude agent in left pane, shell in right pane
- Validate preconditions and fail with clear errors

## Capabilities

### New Capabilities

- `spawn`: Create workspace + tmux window + launch agent from an openspec change name

### Modified Capabilities

_(none)_

## Impact

- New files: `internal/spawn/spawn.go` (spawn logic), `cmd/jjay/main.go` (add spawn subcommand)
- External tool dependencies: `jj`, `tmux`, `claude` (runtime)
- Requires: running inside a tmux session, existing openspec change
