# Proposal: Rebase before merge to prevent silent file drops

**Change**: rebase-before-merge
**Status**: proposed
**Bean**: [jjay-30gc — CRITICAL: openspec artifacts diverge on merge](../../../.beans/jjay-30gc--critical-openspec-artifacts-diverge-on-merge-files.md)

## Why

`jjay merge` currently creates a merge commit directly via `jj new main <workspace>@`. When both main and the workspace have modified the same files (or when the workspace adds new files that main doesn't know about), jj's 3-way merge silently picks one side — dropping files, losing task checkbox progress, and even deleting beans.

Confirmed losses:
- tasks.md checkbox progress (merge-command)
- Blog posts (session-open, fix-pane-dirs)
- The merge conflict bean itself was deleted by a merge

## What Changes

- Before creating the merge commit, rebase the workspace branch onto current main via `jj rebase -b <workspace>@ -d main`
- After rebase, the workspace already includes all of main's changes — the merge becomes trivial (no 3-way file picking)
- If rebase surfaces real conflicts, they appear as explicit jj conflict markers — the user can resolve before proceeding
- Add 6 e2e test scenarios covering all merge cases

## Capabilities

### Modified Capabilities

- `merge`: Rebase workspace onto main before merging, eliminating silent file drops

### New Capabilities

- `merge-e2e-tests`: Integration tests for merge covering clean merges, divergent files, new files, conflicts, empty workspaces, multi-commit branches

## Impact

- Modified: `internal/merge/merge.go` (add rebase step, check for conflicts after rebase)
- New: `internal/merge/merge_integration_test.go` (6 e2e test scenarios)
- No new dependencies
