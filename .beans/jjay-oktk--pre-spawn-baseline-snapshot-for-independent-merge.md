---
# jjay-oktk
title: Pre-spawn baseline snapshot for independent merge verification
status: todo
type: feature
priority: normal
created_at: 2026-06-08T14:28:24Z
updated_at: 2026-06-08T14:28:24Z
parent: jjay-5y1a
blocked_by:
    - jjay-q6ko
---

Capture a **baseline snapshot of `main`** at the moment a workspace is spawned, so that `jjay merge`'s post-merge verification has an **independent reference point** instead of trusting the same commit graph the merge itself walks.

## Why

The verification in [jjay-q6ko](jjay-q6ko) (proposal `harden-merge-verification`, ADR-013) added a post-merge smoke test (L1 "main gained something", L2 "every frontier file is present on main"). But L2's notion of "what the workspace touched" is computed from the **work frontier** — the same revset the merge operates on. That makes the failure modes **correlated**: if the merge mis-defines the frontier (exactly the instance-3 orphaned-sibling bug, where the real proposal lived on a divergent commit the revset never reached), then both the merge *and* the L2 capture miss the same files. The smoke test cannot flag a file it never knew existed — it marks its own homework against the flawed source.

A baseline taken **before spawn** breaks that correlation. It records what `main` looked like the moment the workspace was born, independent of any later revset. Merge can then compute the true divergence as `diff(baseline, all workspace heads)` — a definition that does not depend on the frontier revset being correct.

This is also the natural foundation for **L3 content equivalence**, which `harden-merge-verification` explicitly defers: with both endpoints (baseline + workspace tips) recorded, content comparison becomes tractable.

## Scope notes

- Touches **spawn** (`/jjay:spawn` / spawn internals), not just `internal/merge/merge.go` — a different command and lifecycle moment than q6ko.
- Needs a **persistence mechanism** so merge can read the baseline later (a jj ref/bookmark, a recorded op id, or workspace metadata). Choosing this is the main design decision.
- Composes with ADR-013's frontier definition and smoke test rather than replacing them — baseline is an *additional, independent* check.

## TODO (to refine during proposal)

- [ ] Decide baseline persistence mechanism (ref / op id / metadata file)
- [ ] Capture baseline at spawn time
- [ ] merge: compute independent divergence baseline..workspace-heads
- [ ] Feed independent divergence into L2 (catch files the frontier revset misses)
- [ ] Lay groundwork for / implement L3 content equivalence
- [ ] Tests for the orphaned-sibling topology proving the baseline catches what the frontier missed
