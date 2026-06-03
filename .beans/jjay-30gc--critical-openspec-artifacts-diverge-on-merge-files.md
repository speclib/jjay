---
# jjay-30gc
title: 'CRITICAL: openspec artifacts diverge on merge — files lost/overwritten'
status: in-progress
type: bug
priority: critical
created_at: 2026-06-03T19:13:26Z
updated_at: 2026-06-03T19:41:12Z
parent: jjay-qmpg
---

When openspec change artifacts (tasks.md, blog posts, etc.) are created in the main workspace before spawn, then the spawned workspace modifies them, both sides have changes to the same files. At merge time jj silently resolves by picking one side — losing work from the other.

Confirmed cases:
- tasks.md checkbox progress lost (merge-command)
- blog posts dropped (session-open, fix-pane-dirs)
- This very bean was deleted by a merge

Root cause: same file modified in both branches (proposal creation vs task completion/blog writing).
Impact: archived changes may show unchecked tasks, blog posts disappear, beans get deleted.
Needs investigation for the right fix — this is the most critical workflow issue.
