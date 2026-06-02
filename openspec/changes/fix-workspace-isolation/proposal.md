# Proposal: Fix workspace isolation to prevent staleness data loss

**Change**: fix-workspace-isolation
**Status**: proposed
**Bean**: [jjay-w88o — CRITICAL: fix workspace staleness](../../../.beans/jjay-w88o--critical-fix-workspace-staleness-spawn-must-isolat.md)

## Why

During testing, `jjay spawn` created a workspace with `--revision @`. The spawned agent worked in the new workspace, creating jj operations. When the user returned to the main workspace and ran `jj workspace update-stale`, uncommitted work was lost and required manual recovery via `jj restore`.

The root cause: the main workspace's `@` had uncommitted changes (the openspec change directory, code edits, etc.). When the release-process workspace created operations, the main workspace became stale. Reconciling the stale working copy against the new repo state caused the loss.

## What Changes

- Before creating the child workspace, run `jj new` in the main workspace to snapshot all uncommitted work into `@-` and start a fresh empty `@`
- Create the child workspace with `--revision @-` (the snapshot with all files) instead of `--revision @` (the new empty change)
- This ensures the main workspace's `@` has nothing to lose if it becomes stale

## Capabilities

### Modified Capabilities

- `spawn`: Change workspace creation sequence to isolate working copies safely

## Impact

- Modified: `internal/spawn/spawn.go` — add `jj new` step, change `--revision @` to `--revision @-`
- No new files, no new dependencies
- **User-visible change**: after `jjay spawn`, the main workspace will be on a new empty change. The previous work is in `@-` (committed, safe).
