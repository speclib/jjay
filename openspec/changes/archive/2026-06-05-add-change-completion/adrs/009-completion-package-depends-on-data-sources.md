# ADR-009: Completion depends on data-source packages, not command packages

**Status**: Accepted

## Context

Adding per-verb change-name completion to `spawn`/`merge`/`cleanup` requires two candidate sets: openspec change names (for `spawn`) and spawned-workspace names (for `merge`/`cleanup`, and as the subtrahend in `spawn`'s set-minus). Both already have readers: `spawn.go` parses `openspec list --json` (privately), and `internal/status.List()` (from `workspace-aware-session`) enumerates spawned workspaces — but `List()` also probes tmux and reads tasks.md for the rich status table.

The naive wiring — have a completion package import `spawn` (for the openspec parsing) and `status`/`merge`/`cleanup` — would couple completion to every command package, force `spawn.go`'s private types public, and create a cross-package web that the explicitly-anticipated future CLI redesign would turn into spaghetti. A second concern: TAB completion has a tighter latency budget than a status table — calling `List()` (tmux + file I/O) on every keystroke is wasteful.

## Options Considered

- **Completion imports command packages** — reach into `spawn`/`merge`/`cleanup` for their lookups. Cross-package coupling; private types leak; cycles likely; brittle under CLI change.
- **Duplicate the readers in completion** — re-parse `openspec list` and `jj workspace list` inline. No coupling, but two+ copies drift; violates single-source.
- **Reuse status.List() directly** — one source, but pays tmux + tasks.md cost on every TAB; couples completion to the rich-table shape.
- **Completion over lean data-source readers** — extract name-only readers into data packages; completion composes them; commands bind to completion functions. One-way deps, fast, redesign-safe.

## Decision

- **Dependency direction is one-way:** `cmd → completion → {openspec, status} → exec`. The `internal/completion` package imports **only data-source packages** (`internal/openspec`, `internal/status`), never command packages (`spawn`, `merge`, `cleanup`). Commands bind to completion via `ValidArgsFunction`, not the reverse.
- **Lean name-only readers.** Add `internal/openspec.ChangeNames()` (lifted from `spawn.go`'s existing parsing; `spawn.go` then reuses it — no behavior change) and `internal/status.WorkspaceNames()` (`jj workspace list` only — no tmux, no tasks.md). Completion uses these, not the heavy `List()`, to honor the TAB latency budget. `List()` is unchanged and continues to back the status table.
- **Named-apart verb functions.** `completion.Spawnable`, `Mergeable`, `Cleanable` are distinct exported functions even though `Mergeable == Cleanable` today (both = spawned workspaces). Naming the *intent* per verb lets candidate sets diverge later (e.g. cleanup also targeting archived changes) without renaming or re-plumbing. Shared bodies may delegate to a small helper.
- **Filtering semantics.** `Spawnable = ChangeNames() \ WorkspaceNames()` (set minus); `Mergeable = Cleanable = WorkspaceNames()`. Completion is advisory — the commands' own precondition checks remain authoritative.

## Consequences

- **Positive**: A future CLI redesign moves only the `ValidArgsFunction` bindings in `main.go`; the completion package and its data sources are untouched.
- **Positive**: Fast TAB — completion never probes tmux or reads task files.
- **Positive**: Single source per fact (`ChangeNames`, `WorkspaceNames`); `spawn.go` shares the lifted reader rather than owning a private copy.
- **Negative**: Hard build-order gate — `internal/status` (with `WorkspaceNames`) must exist on main first, i.e. `workspace-aware-session` must merge before this change builds.
- **Negative**: Two identical function bodies (`Mergeable`/`Cleanable`) today — accepted as cheap, intent-revealing insurance against future divergence.
- **Negative**: `Spawnable`'s set-minus calls two readers (openspec + jj) per TAB; still fast (two subprocess reads, no tmux/file walks), but more than a single call.
