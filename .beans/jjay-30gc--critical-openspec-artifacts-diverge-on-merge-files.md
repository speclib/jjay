---
# jjay-30gc
title: 'CRITICAL: openspec artifacts diverge on merge — files lost/overwritten'
status: completed
type: bug
priority: critical
created_at: 2026-06-03T19:13:26Z
updated_at: 2026-06-04T00:00:00Z
parent: jjay-qmpg
openspec-link: openspec/changes/archive/2026-06-04-rebase-before-merge
---

When openspec change artifacts (tasks.md, blog posts, etc.) are created in the main workspace before spawn, then the spawned workspace modifies them, both sides have changes to the same files. At merge time jj silently resolves by picking one side — losing work from the other.

Confirmed cases:
- tasks.md checkbox progress lost (merge-command)
- blog posts dropped (session-open, fix-pane-dirs)
- This very bean was deleted by a merge

Root cause: same file modified in both branches (proposal creation vs task completion/blog writing).
Impact: archived changes may show unchecked tasks, blog posts disappear, beans get deleted.
Needs investigation for the right fix — this is the most critical workflow issue.

## Summary of Changes

Fixed via the `rebase-before-merge` change (archived 2026-06-04). `jjay merge` now runs `jj rebase -b <change>@ -d main` before creating the merge commit, so the workspace already includes all of main's changes and the merge no longer triggers jj's 3-way silent file picking. Real conflicts are surfaced explicitly and abort the merge for manual resolution. Added ADR-007 documenting the decision and 6 e2e integration test scenarios (clean merge, divergent files, same-file conflict, new files, empty workspace, multi-commit). The `merge` spec was synced to reflect the rebase-first behavior plus the new file-preservation and e2e-test requirements.
