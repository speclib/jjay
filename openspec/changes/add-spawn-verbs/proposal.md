## Why

`jjay spawn` does exactly one thing: spawn an agent running `/opsx:apply` on an **existing** change. But the orchestrator (the main session) is often busy — resolving a merge conflict, applying a change — and that's precisely when you want to kick off a *new* exploration or proposal without blocking. Today you can't: explore/propose only happen in the main session.

Bean [jjay-4ulx](../../.beans/jjay-4ulx--more-spawns-spawn-apply-spawn-explore-spawn-propos.md) asks for "more spawns." Exploration crystallized it: there are **two spawn verbs**, not three — *apply* (work an existing change) and *proposal* (create new thinking from a prompt). "Explore" is not a separate verb; it is the **starting mode of a proposal**, selected by a flag/config. The main session is freed to do what only it should — resolve conflicts — while exploration and proposal creation run in their own isolated windows.

## What Changes

- **`jjay spawn <verb>` subcommands:**
  - `jjay spawn apply <change>` — today's behavior, now namespaced. Validates the change exists, isolates it, runs `/opsx:apply`. Workspace/window named **`app-<change>`**.
  - `jjay spawn proposal <prompt> [--mode explore|propose]` — seed a new proposal spawn from a free-text prompt. No change exists yet; the agent creates it. Workspace/window named **`prop-<slug>`**. `--mode` (default configurable) picks the seed command: `/opsx:explore <prompt>` or `/opsx:propose <prompt>`.
  - **No bare alias.** `jjay spawn` SHALL require a verb; bare `jjay spawn <change>` is removed. With a single user and a pre-1.0 tool we take the clean break over carrying two spellings — callers (`/jjay:spawn`, integration tests) are updated to `spawn apply`.
- **Slug identity for proposal spawns** (no AI, no round-trip): plain-code extraction from the prompt — lowercase, strip punctuation/stopwords, keep salient tokens, cap length, add a uniqueness suffix if needed. The slug is **both the immutable handle and the display name** — it is never remapped, even after the agent invents a real change name.
- **Verb-prefixed names for all spawns:** `app-<change>`, `prop-<slug>`. The prefix encodes the kind and makes `status` readable at a glance.
- **All verbs are isolated** — every spawn gets its own jj workspace (no workspace-less mode). Proposal spawns write `openspec/changes/<ai-name>/` inside their own workspace, so they cannot race the main working copy (the divergence class fixed in jjay-30gc / jjay-ug7y). Accepted cost; it also keeps merge the single robust integration path.
- **`jjay status` shows two tables:** *CHANGES* (`app-*`, the openspec lifecycle with MERGED/ARCHIVED/TASKS) and *PROPOSAL SPAWNS* (`prop-*`, prompt-seeded, no change yet).

## Capabilities

### New Capabilities
<!-- None: this extends the existing spawn capability rather than adding a new one. -->

### Modified Capabilities
- `spawn`: add `apply` and `proposal` subcommands; proposal spawns derive a slug identity and don't require a pre-existing openspec change; verb-prefixed workspace/window names.
- `status`: split the table into CHANGES vs PROPOSAL SPAWNS by spawn kind (prefix).
- `completion`: complete the spawn verbs (`apply`, `proposal`); the existing change-name completion moves from the bare `spawn` argument to `spawn apply`'s argument; `spawn proposal` takes a free-text prompt (no candidates).

## Impact

- **Code**: `cmd/jjay/main.go` (subcommand wiring, no bare-arg form), `internal/spawn/spawn.go` (two flows: apply validates-then-isolates; proposal slugs-then-isolates, no `checkOpenspecChange`), a small slug helper (`internal/spawn` or `internal/workspace`), `internal/workspace.WindowName`/`WorkspaceDir` to carry the prefix, `internal/status` for the two-table view.
- **Caller migration**: the `/jjay:spawn` command and the spawn integration tests are updated from bare `spawn <change>` to `spawn apply <change>`.
- **Naming divergence is intentional:** a `prop-<slug>` workspace will contain a change dir named differently (the AI's eventual name). `merge`/`status` must not assume workspace-name == change-name (today they do). Stated so it isn't mistaken for a bug.
- **Out of scope → [jjay-mg00](../../.beans/jjay-mg00--proposal-status-lifecycle-enum-exploredraftreviewa.md):** the proposal status lifecycle enum (explore → draft → ready-for-review → accepted → request-for-changes). That's a state machine whose storage fights ADR-006 and likely shares the agent-emitted-signal mechanism with [jjay-atq6](../../.beans/jjay-atq6--status-agent-column-busyfinished-via-agent-emitted.md).
- **ADRs**: ADR-011 (two spawn flows + code-derived immutable slug identity).
- **Beans**: 4ulx → in-progress, linked here; mg00 holds the deferred lifecycle.
