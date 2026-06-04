## Why

`jjay status` shows whether a spawn's tmux window is attached, its task progress, and whether the change is archived — but not the one thing you most want before tearing a spawn down: **has its work landed on `main` yet?** Bean [jjay-nnjz](../../.beans/jjay-nnjz--more-status-columns-view-configurable.md) asks for more columns; this change is the **easy, ADR-006-pure half** of it. The harder columns (agent name, busy/finished runtime state, configurable views) are split out to [jjay-atq6](../../.beans/jjay-atq6--status-agent-column-busyfinished-via-agent-emitted.md) because they require persisted/emitted state rather than live derivation.

A **MERGED** column is the most-wanted addition: it flags the "done, ready to clean up" state — work is on `main` but the workspace still exists (and may not be archived yet). It composes cleanly with the planned "after merge, ask to cleanup" bean ([jjay-0yko](../../.beans/jjay-0yko--after-merge-ask-to-cleanup-automatically.md)): MERGED gives the visibility, 0yko the action.

## What Changes

- **Rename the `STATUS` column to `TMUX`.** Today's `STATUS` column literally reports the tmux window state (attached/detached); calling it `TMUX` is accurate and frees the word "status" for the future agent-state column (jjay-atq6).
- **Add a `MERGED` column** (yes/no): is the spawn's work already on `main`? Derived live from jj — leading candidate revset `main..<change>@` (empty ⇒ merged) — with no state file, consistent with ADR-006.
- **Keep the existing `ARCHIVED` column.** MERGED (work landed on main) and ARCHIVED (openspec change archived) are distinct: a change can be **merged but not yet archived** — exactly the "ready to archive/cleanup" signal. New column layout: `CHANGE  WORKSPACE  TASKS  TMUX  MERGED  ARCHIVED`.
- **Add `jjay status` to the integration test** (folds in bean [jjay-y5yn](../../.beans/jjay-y5yn--add-cmd-status-to-the-integration-test.md)). `status` currently has unit tests but no end-to-end coverage. The lifecycle test gains a `status` subtest (between spawn and cleanup) asserting the column set and the MERGED value — which is also the most honest way to validate MERGED's edge cases (unmerged spawn, empty workspace) against a real jj repo.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `status`: rename the tmux-state column to `TMUX`; add a `MERGED` column derived live from jj.

## Impact

- **Code**: `internal/status/status.go` — add a `Merged bool` to the `Spawn` struct, compute it during `List`/`join` via a jj revset, and update `Render` (header + row). No new package; no new external tool (jj already used).
- **Tests**: `test/integration/full_lifecycle_test.go` gains a `status` subtest asserting the column set + MERGED; unit tests for the merged derivation in `internal/status`.
- **Out of scope** (→ jjay-atq6): `AGENT` column, `STATUS` busy/finished, view-configurable columns. These need an agent-emitted signal / persisted state and their own design.
- **ADRs**: none — MERGED is a live jj derivation, squarely inside ADR-006; no new architectural decision.
- **Beans**: nnjz → in-progress (easy half) and y5yn → in-progress (status integration test), both linked here; atq6 holds the deferred hard half.
