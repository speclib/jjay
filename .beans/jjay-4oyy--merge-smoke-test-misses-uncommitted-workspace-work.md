---
# jjay-4oyy
title: merge smoke test misses UNCOMMITTED workspace work (empty frontier = false pass)
status: in-progress
type: bug
priority: high
created_at: 2026-06-12T12:26:16Z
updated_at: 2026-06-12T16:00:00Z
parent: jjay-hjjg
---

A 4th merge failure mode, beyond the three in jjay-q6ko / harden-merge-verification (ADR-013). The verification-gated merge STILL reported "verified" while landing nothing.

## Reproduced (2026-06-12)
`./jjay merge prop-no-respawn-claude-apply` printed:
  Warning: workspace "..." has no changes in its working copy.
  ... workspace forgotten (work is on main).
  Rebased and merged ... into main (verified).
But the proposal `resume-spawns-on-reopen` was NOT on main.

## Root cause
The workspace had a full proposal `openspec/changes/resume-spawns-on-reopen/` (proposal, design, adrs, specs, tasks) but it was UNCOMMITTED — sitting in the working copy on disk. All the workspace's jj commits (txyyun, txnvoz) were empty.

merge's frontier = `ancestors(<change>@) & main.. & ~empty()` only sees COMMITTED non-empty work. Uncommitted files are invisible to it → frontier empty → workFrontierFiles = {} → smokeTest sees len(expectedFiles)==0 → "nothing to prove" → PASS → forgets the workspace → reports "verified".

So the smoke test's L1 ("did the workspace have work?") is answered from committed commits only; it cannot see uncommitted work, so an empty-frontier-but-dirty-working-copy workspace passes falsely.

## The signal was there
merge ALREADY printed "Warning: workspace has no changes in its working copy" (checkWorkspaceEmpty). That warning + forget-on-success is the dangerous combo: we warn that there's nothing, then forget the workspace anyway — discarding the (uncommitted) work's only home from jj's view. Recovery worked ONLY because `jj workspace forget` leaves the directory on disk; the proposal was copied back from ../jjay-workspaces/prop-no-respawn-claude-apply/.

## Fix directions (needs proposal)
- Before merge, SNAPSHOT the workspace's working copy (jj auto-snapshots on a jj command run IN that workspace dir — but merge runs from the main session, so the spawned workspace's dirty @ is never snapshotted). merge should force a snapshot of the target workspace (e.g. `jj -R <wsdir> status`/`debug snapshot`, or `jj workspace update-stale`) BEFORE computing the frontier, so uncommitted work becomes a real commit and enters the frontier.
- OR: if checkWorkspaceEmpty warns AND the frontier is empty, do NOT forget the workspace and do NOT claim "verified" — treat empty-frontier as unproven when the workspace dir has uncommitted changes.
- Tie-in: this is exactly why jjay-oktk (pre-spawn baseline snapshot) and the "snapshot before acting" idea matter.

## Relation
Extends ADR-013 (harden-merge-verification, archived 2026-06-12). Sibling of jjay-ychu (orphan-sibling detection). The common thread across q6ko + ychu + this: merge must define "the workspace's work" to include uncommitted + sibling state, not just ancestors(<change>@)'s committed commits.
