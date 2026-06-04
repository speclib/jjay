---
# jjay-ug7y
title: 'CRITICAL: merge deletes active sibling change-dirs created in main after spawn snapshot'
status: todo
type: bug
priority: critical
created_at: 2026-06-04T21:20:38Z
updated_at: 2026-06-04T21:20:38Z
parent: jjay-qmpg
---

Regression / gap in the rebase-before-merge fix (jjay-30gc). That fix preserved files MODIFIED on both sides, but NEW directories created in main AFTER the spawn snapshot are still silently deleted at merge time.

## Reproduced this session (2026-06-04)
1. Spawned `workspace-aware-session`; it snapshotted main.
2. AFTER the snapshot, three new active changes were created in main: `add-claude-commands`, `add-init-command`, `add-change-completion` (each a new `openspec/changes/<name>/` dir).
3. `jjay merge workspace-aware-session` (+ archive) ran.
4. Result: `add-init-command` and `add-change-completion` were DELETED from main. `add-claude-commands` survived only because it was re-touched after the snapshot.

Recovered both from orphan jj commit 60f1a51 via `git checkout 60f1a51 -- <paths>`. No permanent loss, but silent and dangerous.

## Why 30gc's fix didn't catch it
`rebase-before-merge` rebases `<change>@` onto main before the merge commit, which should fold in main's changes. The escape here is NEW top-level dirs added to main after the workspace's base revision — they did not survive. Either the rebase used a stale main, or the merge's tree resolution dropped main-only additions. NEEDS INVESTIGATION (op log around the merge: jj op log; compare main@ tree before/after).

## Test gap
internal/merge/merge_integration_test.go has TestMerge_WorkspaceAddsNewFiles (workspace adds files, main moves) but NO mirror test for "MAIN adds new files/dirs after snapshot → must survive merge". Add scenario:
- TestMerge_MainAddsNewFiles: main creates a new dir/file after the workspace base; workspace commits unrelated work; after merge BOTH the main-only dir AND the workspace work exist.
This must be implemented in a change proposal (explore mode cannot write code).

## Relationship
Sibling/regression of jjay-30gc (fixed via archive/2026-06-04-rebase-before-merge). Same class (silent merge data loss), different trigger (main-side additions vs both-sides modification).
