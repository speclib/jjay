---
# jjay-q6ko
title: merging now almost always give a stale error
status: todo
type: bug
priority: normal
created_at: 2026-06-04T22:23:50Z
updated_at: 2026-06-04T23:58:31Z
---

this is how I solve it from the workspace claude session:

me: i tried ./jjay merge from the main ws and this ws has this error: Error: The working copy is stale (not updated since operation f3fe0a226e2c).
claude: The staleness is resolved — but updating revealed a 2-sided conflict in .beans/jjay-ofk7--bootstrap-task.md. Main moved the bean (changed parent to jjay-qltp and bumped updated_at to 2026-06-04T22:08:46Z), while this workspace changed updated_at to 2026-06-05
  and left parent: jjay-5y1a.


the conflict are real but the staleness is something we can propably prevent

## Investigation (2026-06-05, reproduced)

Reproduced merging `add-status-merged-column` with `./jjay merge`: the merge SUCCEEDED, but afterwards `jj status` inside `../jjay-workspaces/add-status-merged-column` errors with `The working copy is stale (not updated since operation 003da90efd67)`.

Root cause: `internal/merge/merge.go` rebases the workspace's commit (`jj rebase -b <change>@ -d main`) and builds the merge commit FROM THE MAIN session. That rewrites `<change>@` out from under the spawned workspace, but nothing updates that workspace's working-copy pointer — so the workspace is left pointing at a superseded operation = stale. It happens on essentially every merge because every merge rewrites `<change>@`.

Two things are tangled in the original report:
- (1) STALENESS — structural, preventable. This is the bug.
- (2) the `.beans/jjay-ofk7` conflict — a REAL content conflict (main changed parent→qltp, workspace kept 5y1a), not staleness. Just a genuine divergence to resolve manually.

Likely fix directions (needs a proposal):
- After a successful merge, run `jj workspace update-stale` for the merged workspace (it's about to be cleaned up anyway — or the user may inspect it).
- OR, since the workspace is normally torn down post-merge, have merge/cleanup forget the workspace so staleness never surfaces; if merge is expected to leave the workspace alive, update-stale it.
- Either way: merge should not leave a spawned workspace in a stale state it can't recover from without manual `jj workspace update-stale`.

Status note: surfaced while merging add-status-merged-column (which itself merged fine — the ug7y "fold ahead-of-bookmark work" step ran correctly).
