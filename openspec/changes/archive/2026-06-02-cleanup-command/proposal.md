# Proposal: jjay cleanup command

**Change**: cleanup-command
**Status**: proposed
**Bean**: [jjay-uypj — jjay cleanup command](../../../.beans/jjay-uypj--jjay-cleanup-command.md)

## Why

After a spawned agent finishes (or fails), the user must manually run three commands to tear down the workspace: `jj workspace forget`, `rm -rf`, `tmux kill-window`. This is tedious and error-prone, especially when managing multiple parallel agents.

## What Changes

- Add `jjay cleanup <change-name>` cobra subcommand
- Kill tmux window, forget jj workspace, remove workspace directory
- Tolerant execution: skip missing pieces, clean up what exists
- Extract shared helpers (`windowName`, `workspaceDir`) from spawn into `internal/workspace/`
- Update spawn to use the shared package

## Capabilities

### New Capabilities

- `cleanup`: Tear down workspace + tmux window + directory for a given change name

### Modified Capabilities

- `spawn`: Refactor shared helpers to `internal/workspace/` (no behavior change)

## Impact

- New files: `internal/cleanup/cleanup.go`, `internal/workspace/workspace.go`
- Modified files: `cmd/jjay/main.go` (add cleanup subcommand), `internal/spawn/spawn.go` (use shared workspace package)
- No new external dependencies
