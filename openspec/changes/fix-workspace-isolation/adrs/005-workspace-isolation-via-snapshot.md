# ADR-005: Workspace isolation via jj new snapshot before spawn

**Status**: Proposed

## Context

`jjay spawn` creates a child jj workspace for an AI agent. The main workspace may have uncommitted changes. When the agent works in the child workspace, the main workspace becomes stale. Running `jj workspace update-stale` can lose uncommitted work from the main workspace's `@`.

This was discovered when a user lost work during testing on 2026-06-02 and had to manually recover via `jj restore`.

## Options Considered

- **`jj new` before spawn** — snapshot uncommitted work into `@-`, create child from `@-`. Main workspace's `@` is empty and safe. Simple, uses standard jj workflow.
- **`jj new` in child after creation** — child briefly shares `@` with main, reintroducing the original race condition.
- **Require user to commit before spawn** — adds friction, easy to forget, doesn't protect against mistakes.
- **Copy files instead of jj workspace** — loses jj integration entirely, defeats the purpose.

## Decision

Run `jj new` in the main workspace before creating the child workspace. The child workspace uses `--revision @-` to get the snapshot. The main workspace ends up on a fresh empty `@` — nothing to lose if staleness occurs.

## Consequences

- **Positive**: Uncommitted work is always safe — it's in `@-`, committed via jj snapshot
- **Positive**: No user discipline required — spawn handles it automatically
- **Positive**: Standard jj workflow — `jj new` is idiomatic
- **Negative**: After spawn, the main workspace is on an empty `@`. User must know to `jj edit @-` if they want to continue working on the previous change. Worth documenting.
