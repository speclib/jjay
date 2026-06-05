## Context

`internal/spawn.Spawn(changeName, opts)` runs a fixed sequence: `checkTmuxSession` → `checkOpenspecChange` → `checkWorkspaceNotExists` → `checkWindowNotExists` → `WorkspaceDir` → snapshot → `createWorkspace(--name changeName)` → `OpenWindow(ws-changeName)` → `setupPanes` (agent = `DefaultAgentCommand` with `/opsx:apply {change}`). The change name is the single identity for workspace/dir/window, and `cleanup`/`status`/`merge` all key off it. ADR-011 introduces a second flow (proposal) that has no change at spawn time, plus verb-prefixed, slug-based identity.

## Goals / Non-Goals

**Goals:**
- `jjay spawn apply <change>` and `jjay spawn proposal <prompt> [--mode]`; a verb is required (no bare form).
- Code-derived, immutable slug identity for proposal spawns; verb prefixes `app-`/`prop-`.
- All spawns isolated; proposal spawns may produce a differently-named change.
- `jjay status` two-table view keyed on the prefix.

**Non-Goals (→ jjay-mg00):**
- Proposal status lifecycle enum (explore/draft/review/accepted/request-for-changes) and its storage/transitions.
- AI-generated names (slug is plain code).
- Renaming a spawn after the agent names its change (handle is immutable).

## Decisions

- **Two flows, one command, verb required.** Add cobra subcommands `apply` and `proposal` under `spawn`; `spawn` with no verb prints usage and exits non-zero (no bare-arg form — clean break, single pre-1.0 user). The apply flow is today's `Spawn` with the name prefixed `app-`. The proposal flow skips `checkOpenspecChange`, derives a slug, prefixes it `prop-`, and sets the agent command to `/opsx:explore` or `/opsx:propose` per `--mode`.
- **Slug helper (plain code).** A small function: lowercase → strip punctuation → split → drop a stopword set → keep first N salient tokens → join with `-` → cap length → if the resulting `prop-<slug>` collides with an existing workspace/window, append `-2`, `-3`, … Deterministic, no AI, microseconds. Lives in `internal/spawn` (or `internal/workspace` next to the naming helpers).
- **Prefix flows through naming.** `workspace.WindowName`/`WorkspaceDir` already take a name; pass the already-prefixed name (`app-<change>` / `prop-<slug>`). Avoid baking the prefix into those helpers so they stay generic.
- **Agent command per verb.** `DefaultAgentCommand` stays the apply template. Add proposal templates (explore/propose) with a `{prompt}` placeholder; `--mode` (default from config) selects. Reuse `resolveAgentCommand` with the new placeholder.
- **Decouple workspace-name from change-name.** Audit `merge`/`status`/`cleanup` for the assumption `workspace == change`. Apply spawns keep it (via the `app-` mapping); proposal spawns break it deliberately. Where a command needs the produced change, it must read it from the workspace, not infer it from the workspace name.
- **status kind from prefix.** `status.List` already enumerates workspaces; classify each by `app-`/`prop-` prefix and render two tables. Proposal rows omit MERGED/ARCHIVED/TASKS (or show `—`).

## Risks / Trade-offs

- **`spawn` complexity doubles** — two flows with different preconditions. Mitigate by sharing the isolate→window→panes tail and branching only on validate-vs-slug + agent template.
- **Workspace==change assumption is load-bearing today** — missing a call site could mis-key a proposal spawn in merge/status. The audit (task) is the real work; integration coverage should include a proposal spawn whose produced change name differs from its slug.
- **Slug quality** — naive extraction can yield awkward or colliding slugs. Accept (it's a handle, not prose); uniqueness suffix handles collisions; users see it but it's short and readable.
- **`--add-dir {wsdir}` semantics for proposal** — the apply template adds the workspace dir; confirm the explore/propose templates point the agent at the right working dir so it writes inside the isolated workspace.
- **Caller migration (no alias)** — removing the bare form means `/jjay:spawn` and the spawn integration tests must be updated to `spawn apply <change>`; they break loudly if missed (good — a clean break surfaces every caller). Also verify the `app-` prefix rename doesn't break cleanup/status lookups that expected the unprefixed name.
