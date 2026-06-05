---
# jjay-q6ko
title: merging now almost always give a stale error
status: todo
type: bug
priority: critical
created_at: 2026-06-04T22:23:50Z
updated_at: 2026-06-05T11:56:13Z
parent: jjay-hjjg
---

this is how I solve it from the workspace claude session:

me: i tried ./jjay merge from the main ws and this ws has this error: Error: The working copy is stale (not updated since operation f3fe0a226e2c).
claude: The staleness is resolved — but updating revealed a 2-sided conflict in .beans/jjay-ofk7--bootstrap-task.md. Main moved the bean (changed parent to jjay-qltp and bumped updated_at to 2026-06-04T22:08:46Z), while this workspace changed updated_at to 2026-06-05
  and left parent: jjay-5y1a.

## 2nd instance (2026-06-05): empty merge + stale workspace

Merging prop-jjay-7rol: the work was committed in the workspace's @- (mmptmw), while its @ (qtputx) was EMPTY. `jjay merge` builds the merge from <change>@, so the merge commit came out EMPTY ("0 files changed") — looked like "merged but I don't see the change". The work was NOT lost (mmptmw was already on main's line, so it landed), but the empty merge commit and the empty @ are confusing and fragile: in a different topology, merging the empty @ instead of the real work in @- would silently merge nothing.

Two things to fix here:
- merge should operate on the workspace's real HEAD, not blindly on <change>@ when @ is empty (or it should warn/refuse on an empty @, like checkWorkspaceEmpty already detects but only warns).
- staleness: prop-jjay-7rol workspace is stale after the merge (same root as above — merge rewrites from the main session without updating the workspace).


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



## 3rd instance (2026-06-05): merge grabbed the WRONG commit — real proposal orphaned on a sibling branch

Deeper investigation of the SAME prop-jjay-7rol merge revealed the empty-merge was worse than first thought. The workspace had THREE relevant commits:
- mmptmw — bean edits (landed on main)
- qtputx — empty @ (what merge used)
- zroyto — THE ACTUAL PROPOSAL (openspec/changes/make-integration-output-readable/, full artifact set: proposal, design, ADR-012, spec delta, tasks)

zroyto was a SIBLING (parent rqyvlv), not an ancestor of <change>@. So `jjay merge` — building from <change>@ — never saw it. The merge reported success, but the real created proposal was left ORPHANED off-main entirely. Not in @, not in @-, on a divergent branch the merge's revset didn't reach. The user correctly insisted 'we still miss the merged proposal' while tooling kept reporting the merge fine.

NEW failure mode (distinct from instances 1-2):
- Not just 'empty @' — the workspace's real work can live on a DIVERGENT commit that <change>@ does not descend from. `jjay merge <change>` only ever follows <change>@'s line; any sibling work is silently dropped.
- This is the strongest argument yet for rse4 (post-merge smoke test): comparing 'files the workspace touched (across ALL its commits, not just @)' vs 'what landed on main' would have flagged: proposal dir present in workspace history, absent on main.

Implication for the fix: merge must define 'the workspace's work' robustly — likely ALL non-empty commits in the workspace not yet on main (e.g. main..(workspace heads), including divergent ones), not just <change>@. An empty @ with work in @- AND a sibling commit with the real proposal both broke the naive <change>@ assumption in one session.

Recovery used: `jj restore --from zroyto openspec/changes/make-integration-output-readable/` to surface the proposal into the working copy for inspection (not a real merge — landing it properly still needs the .beans/jjay-7rol conflict resolved).
