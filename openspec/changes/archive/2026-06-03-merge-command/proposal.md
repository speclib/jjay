# Proposal: jjay merge command

**Change**: merge-command
**Status**: proposed
**Bean**: [jjay-uxs9 — jjay merge command](../../../.beans/jjay-uxs9--jjay-merge-command.md)

## Why

After a spawned agent finishes, the user must manually run three jj commands to merge the workspace's work into main. This is error-prone — you need to find the right change ID, remember the bookmark dance, and create a fresh change afterwards.

## What Changes

- Add `jjay merge <change-name>` cobra subcommand
- Resolve workspace's working copy via `jj log -r "<change>@"`
- Create merge commit: `jj new main <change>@ -m "merge <change> into main"`
- Move bookmark: `jj bookmark set main -r @`
- Create fresh change: `jj new`
- Precondition checks: workspace exists, workspace has changes
- Package: `internal/merge/merge.go`

## Capabilities

### New Capabilities

- `merge`: Merge a workspace's work into main bookmark

### Modified Capabilities

_(none)_

## Impact

- New: `internal/merge/merge.go`
- Modified: `cmd/jjay/main.go` (add merge subcommand)
- Modified: `README.md` (move merge from Planned to CLI section)
