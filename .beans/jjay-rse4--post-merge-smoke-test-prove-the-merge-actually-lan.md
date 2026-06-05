---
# jjay-rse4
title: 'post-merge smoke test: prove the merge actually landed (changed-file + content check)'
status: draft
type: task
priority: high
created_at: 2026-06-05T10:01:37Z
updated_at: 2026-06-05T11:12:19Z
parent: jjay-qltp
---

Merge-without-AI is the HEART of jjay. Today `jjay merge` proves only a negative — "no conflicts". It does NOT prove the work actually landed. The prop-jjay-7rol merge was conflict-free but produced an EMPTY merge commit (work was in @-, @ was empty); "merged" was reported, nothing visibly changed. Conflict-free != merged.

## Goal
After a merge, jjay runs an in-app smoke test that proves the workspace's work is now on main — extra proof beyond conflict-free.

## Proof levels (build toward L3)
- L1 non-empty: if <change> had changes but the merge added nothing to main → FAIL. (catches the prop-jjay-7rol empty-merge bug)
- L2 changed-file set: capture files the workspace changed (diff base..<change>@) BEFORE merge; after merge assert each appears in main's new state. (catches partial drops — the ug7y class)
- L3 content equivalence (target): each changed file's content on main matches the workspace version. Useful now. CAVEAT: after `rebase -b <change>@ -d main`, a changed file on main may legitimately be a COMBINATION of both sides, not byte-identical — the check must allow legitimate rebase merges, not flag them as failures.

## Pre-1.0 stance (user)
- Be VERBOSE by default under 1.0. Later make it optional via a flag / project / global setting.
- On failure: do NOTHING destructive. Emit a LOUD warning with enough structured information for an AI (or human) to fix the situation — which files were expected, which are missing/divergent, and the recovery handle.
- NOT self-healing yet. Self-healing comes later, once we understand all failure modes. For now: maximize observability and learning about edge cases.

## Tension-1: pre-merge snapshot (feasible)
Capture the jj operation id BEFORE the rebase (`jj op log` gives it; verified `jj op restore <id>` exists). The smoke-test failure message can then hand the user/AI an exact recovery point ("to undo: jj op restore <preMergeOpId>"). This is the "memory-snapshot before rebase" the user asked for — cheap, jj-native.

## Tension-2: toward rock-solid / self-healing
End goal: merge so solid it self-heals once failure modes are catalogued. Each smoke-test failure observed in the wild should be recorded as a new edge case (extend this bean or q6ko). Treat the smoke test as a learning instrument first, a guard second.

## Failure handling decision (for now)
Loud warning + recovery hint (incl. the pre-merge op id). Bookmark stays where it is; no auto-rollback yet. Auto-rollback / refuse-to-leave-fresh-change are future options once modes are understood.

## Relationships
- Sibling of jjay-q6ko (prevent/diagnose bad merges; staleness) — q6ko now logs 2 concrete instances.
- Sibling of jjay-l80s (archive after merge).
- Feeds the "make merge rock solid" theme; merge is the no-AI core capability.



## Refinement (2026-06-05): capture work across ALL workspace commits, not just <change>@

A concrete case (prop-jjay-7rol, logged in q6ko 3rd instance) showed the workspace's real proposal lived on a SIBLING commit (zroyto) that <change>@ did not descend from. So a smoke test that captures 'files changed' via `diff base..<change>@` would MISS it entirely — same blind spot as the merge itself.

Therefore the BEFORE-merge capture must be the union over all the workspace's relevant commits, e.g. files touched in `main..(all heads reachable in the workspace)` — including divergent/sibling commits — not just `base..<change>@`. The smoke test's whole value is catching exactly the work the merge's <change>@ revset fails to follow; if it uses the same revset, it inherits the same blindness.

This also sharpens the failure message: 'workspace history touches openspec/changes/X and Y; main gained neither' is the signal that makes an orphaned-sibling proposal obvious.
