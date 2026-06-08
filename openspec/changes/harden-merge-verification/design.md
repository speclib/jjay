## Context

`jjay merge` (`internal/merge/merge.go`) is single-threaded on one revset, `changeName + "@"`. It is used by the rebase (`jj rebase -b <change>@ -d main`), the merge commit (`jj new main <change>@`), and every check (`checkConflicts`, `checkWorkspaceEmpty`). This assumes `<change>@` *is* the workspace's work. Bean jjay-q6ko reproduced three ways that breaks:

1. **Staleness** — merge rewrites `<change>@` from the main session; the spawned workspace's working-copy pointer is never updated → `The working copy is stale`. On ~every merge.
2. **Empty `@`** — work was in `@-`, `@` empty; `jj new main <change>@` produced a 0-file merge commit reported as success.
3. **Orphaned sibling** — the real proposal lived on a divergent commit (`zroyto`, parent `rqyvlv`) that `<change>@` did not descend from; the merge never saw it and silently landed nothing.

ADR-010 (`advanceMainToHead`) and `rebase-before-merge` (jjay-30gc) already protect *main-side* work. Today merge proves only a **negative** — "no conflicts". Conflict-free ≠ merged (jjay-rse4).

The existing integration suite (`TestMerge_MultipleWorkspaceCommits`) passes today but only exercises a **linear** workspace chain, which `<change>@` reaches. There is no test for the **empty-`@`** or **divergent-sibling** topology — exactly the cases that broke.

## Goals / Non-Goals

**Goals:**
- No silent empty/wrong merge: define "the work" robustly (all non-empty commits in `main..(workspace heads)`), so siblings and `@-` work are included.
- Prove the work landed (smoke test L1+L2) before declaring success.
- Gate the workspace lifecycle on that proof: forget on success (staleness impossible), keep + loudly warn on failure (non-destructive recovery state).
- Integration tests that **fail on today's code and pass after** for each reproduced q6ko instance.

**Non-Goals:**
- Reverting ADR-010 / `rebase-before-merge` — this composes with them.
- L3 content equivalence now (documented follow-up; the rebase-combination caveat in rse4 needs care).
- Auto-rollback / self-healing (rse4 defers; pre-1.0 stance is loud + non-destructive).
- Resolving genuine content conflicts (e.g. the `.beans/jjay-ofk7` divergence in q6ko instance 1) — that is a real conflict to resolve by hand, explicitly out of scope.

## Decisions

- **Work frontier, not `<change>@` (ADR-013).** Define `WORK = main..(heads(<workspace commits>)) & ~empty()`. This includes `@`, `@-`, and divergent siblings. The rebase and merge operate over this frontier. Enumerating "a workspace's commits" precisely is the one genuinely tricky revset (see Risks); the implementation task must validate it end-to-end against the empty-`@` and sibling topologies, not by unit reasoning.
- **Empty-`@`-with-work-in-`@-` needs no special case.** `main..heads` already includes `@-`, so instance 2 falls out of the frontier definition. State this explicitly in tests.
- **Pre-merge op snapshot.** Capture `jj op log` head before any rewrite (verified `jj op restore <id>` exists). Surface it in the failure message as the recovery handle.
- **Smoke test L1+L2 (rse4).** Capture the changed-file set across `WORK` before merge. After merge: L1 — if the workspace had work but main gained nothing → FAIL. L2 — every captured file must be present on main → FAIL on any miss. Verbose by default (pre-1.0); failure message lists expected vs missing files + the op id.
- **Verification-gated lifecycle.** On L1+L2 pass → `jj workspace forget <change>` (reuse the `cleanup.forgetWorkspace` pattern), so there is no live pointer to go stale → instance 1 is structurally closed. On fail → keep the workspace intact, do not move the bookmark beyond what the merge already did, emit the loud warning, exit non-zero. Staleness now only ever appears in the failure branch, where the workspace is deliberately kept for inspection.
- **Test mirrors the q6ko instances.** Add one integration test per reproduced instance (empty-`@`, sibling, staleness/forget-on-success) plus an unproven-keeps-workspace test and an L1 empty-land detector. Two new helpers: `isStale(wsDir)` (runs `jj status` inside the spawned workspace dir, checks for the stale string) and a way to build a **non-linear** workspace DAG (`jj new <earlier>` to fork a sibling).

## Risks / Trade-offs

- **The `heads()` revset is the crux.** Getting "all of a workspace's commits" right across linear, empty-`@`, and divergent topologies is the main risk. Must be validated against the real spawn lifecycle (ADR-010's lesson: validate end-to-end, not by unit reasoning).
- **Forget-on-success changes lifecycle.** The workspace dir is no longer left around after a clean merge (the jj workspace is forgotten; directory cleanup remains the separate `cleanup` step). This is the intended fix but is a visible behavior change — document it. Composes with the existing `cleanup` command and the fd0z/l80s/uypj "archive+cleanup after merge" beans.
- **L2 false positives.** A file legitimately removed/renamed by the workspace could trip "missing on main". Scope L2 to **added/modified** files captured from `WORK`, not deletions; note this limitation.
- **Interaction with `advanceMainToHead`.** The frontier is computed relative to `main` *after* the ahead-of-bookmark fold; ordering matters. Keep `advanceMainToHead` first, then compute `WORK` against the advanced bookmark.
- **Recovery guidance.** Until merged, advise that orphaned work is recoverable via `jj op log` / `jj restore --from <commit>` (as used to recover instance 3's `zroyto`).
