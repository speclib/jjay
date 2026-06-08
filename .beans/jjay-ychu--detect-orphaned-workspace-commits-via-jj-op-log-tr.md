---
# jjay-ychu
title: detect orphaned workspace commits via jj op-log (true-orphan blind spot)
status: draft
type: task
priority: normal
created_at: 2026-06-08T15:54:55Z
updated_at: 2026-06-08T15:54:55Z
parent: jjay-hjjg
---

Discovered while implementing harden-merge-verification (ADR-013).

## The blind spot
`jjay merge` defines the workspace's work as `ancestors(<change>@) & main.. & ~empty()`. This closes q6ko instances 1 (staleness) and 2 (empty @, work in @-), and the smoke test (L1/L2) catches a frontier-with-work that fails to land. BUT a TRUE orphan — work on a divergent sibling commit that `@` never descended from, with `@` left empty (the real zroyto case, q6ko instance 3) — is UNDETECTABLE by merge:

- jj associates exactly ONE working-copy commit (`@`) per workspace.
- A commit `@` never visited has no association with the workspace name.
- So the frontier is empty → the smoke test has nothing to expect → it cannot flag the miss. The orphan is invisible at every layer reachable from `<change>@`.

Confirmed by spike + a failing test that asserted more than the design can deliver (the test was reframed to the detectable cases; this bean tracks the real gap).

## The only known fix: op-log forensics
Scan `jj op log` for commits CREATED while `@` was in this workspace's directory (i.e. operations originating from the spawned workspace), then treat those as candidate work even if unreachable from `<change>@`. Merge/smoke-test could then warn: "workspace authored commit X, not on main".

## Why it's deferred (not done now)
- Heavy and fragile: op-log format/semantics are jj-version-dependent; parsing them is brittle.
- ADR-013 explicitly scoped op-log discovery as a non-goal for the first pass (verbose-and-honest over clever-and-fragile, pre-1.0).
- The common cases (1, 2, L1/L2) are closed without it. True-orphan + empty-@ is rare (requires deliberately moving @ off the work).

## Acceptance idea
Given a spawned workspace whose work is on a sibling `@` never descended from, `jjay merge` (or a doctor/verify step) detects the orphaned commit and warns with a recovery handle — instead of a silent empty merge.

Relationships: extends jjay-q6ko / ADR-013 (harden-merge-verification). Sibling of the smoke-test mechanism. Possibly belongs with a future `jjay doctor` (jjay-wsuv).
