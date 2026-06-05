---
# jjay-2g7t
title: blocking proposals
status: draft
type: task
priority: normal
created_at: 2026-06-04T19:43:28Z
updated_at: 2026-06-04T20:29:39Z
parent: jjay-5y1a
---

 sometimes build order is important. We need a way to define this and make sure jjay respect this when we spawn one or multiple jobs

## Concrete instance (2026-06-05)

`./jjay merge add-spawn-verbs` failed with a real conflict (correctly surfaced by rebase-before-merge — NOT a tool bug):

```
Error: rebase produced conflicts in workspace "add-spawn-verbs". Resolve them manually...
  README.md                       2-sided conflict
  internal/status/status.go       2-sided conflict (4 conflicts)
  internal/status/status_test.go  2-sided conflict
```

Why: `add-status-merged-column` and `add-spawn-verbs` were both spawned from the same base and BOTH modify `internal/status` (the MERGED/TMUX columns vs the two-table view) and README. `add-status-merged-column` merged first; when `add-spawn-verbs` then rebased onto the new main, the overlapping edits collided (e.g. both touch the `Spawn` struct fields and `Render`/`List`). Genuine semantic conflict — needs a human/agent to COMBINE both changes' intent, not a mechanical resolve.

This is exactly the case this bean anticipates: jjay gave NO warning that the two active spawns overlapped, so they were worked fully in parallel and collided at the second merge.

What a fix could provide (scoping TBD — proposal deferred):
- Overlap DETECTION: warn at spawn/merge time when two active spawns touch the same files ("add-spawn-verbs overlaps add-status-merged-column in internal/status, README").
- Declared BLOCKING / build order: mark change B as depending on A so jjay spawns B based on A (or orders/refuses the merge) — the bean's original ask.

Resolution of the current conflict: left to the add-spawn-verbs worker (resolve the 4 status.go conflicts by combining the MERGED-column work with the two-table work, then re-run merge). Not resolved in the main session.
