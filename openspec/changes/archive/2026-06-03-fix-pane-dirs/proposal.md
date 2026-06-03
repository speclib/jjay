# Proposal: Fix pane working directories and add integration test

**Change**: fix-pane-dirs
**Status**: proposed
**Bean**: [jjay-ps3d — tmux shell pane not in correct workspace dir](../../../.beans/jjay-ps3d--tmux-shell-pane-not-in-correct-workspace-dir.md)

## Why

The shell pane (right side) after `jjay spawn` ends up in the wrong directory. The current code uses `tmux send-keys "cd <wsDir>"` which races with shell initialization — fish/bash may not be ready when the `cd` arrives, so it gets lost.

This was supposed to be fixed in spawn-config but was missed due to a merge conflict that dropped the changes.

Additionally, no integration test exists yet to catch this kind of regression.

## What Changes

- Fix `createWindow()`: add `-c <wsDir>` to set window starting directory
- Fix `setupPanes()`: add `-c <wsDir>` to `split-window`, remove `send-keys cd` for right pane
- Remove `cd <wsDir> &&` prefix from agent command (window already starts in workspace dir)
- Add integration test with `//go:build integration` tag covering:
  - Full spawn → verify → cleanup → verify lifecycle
  - Assert both panes have correct working directory via `tmux display-message #{pane_current_path}`
- Fake agent script already exists at `testdata/fake-agent.sh`
- `test-integration` Makefile target already exists

## Capabilities

### Modified Capabilities

- `spawn`: Fix tmux pane working directories using `-c` flag

### New Capabilities

- `integration-test`: Full spawn/cleanup lifecycle test with pane directory assertion

## Impact

- Modified: `internal/spawn/spawn.go` (createWindow, setupPanes)
- New: `internal/spawn/spawn_integration_test.go`
- Uses existing: `testdata/fake-agent.sh`, `Makefile` test-integration target
