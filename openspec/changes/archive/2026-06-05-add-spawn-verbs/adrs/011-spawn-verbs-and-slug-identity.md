# ADR-011: Two spawn flows; code-derived immutable slug identity for proposal spawns

**Status**: Accepted

## Context

`jjay spawn` assumes the spawn target already exists: its first step is `checkOpenspecChange`, and the change name is the primary key for the jj workspace (`--name`), the directory, and the tmux window (`ws-<change>`). The whole lifecycle (`cleanup`, `status`, `merge`) keys off that one name.

Bean jjay-4ulx wants to also spawn *new* work — exploration and proposal creation — so the main session stays free for conflict resolution. Those invert the assumption: there is **no change name at spawn time**. A proposal spawn's change is invented later by the agent (and an exploration may never produce one). Yet a spawn needs a stable identity *immediately* for its workspace, directory, and window. Waiting on an AI to name it is unacceptable latency for what should be an instant action.

Exploration also collapsed three verbs into two: "explore" is the earliest *mode* of a proposal, not a separate verb — so promotion explore→propose must not change identity.

## Options Considered

**Identity source for a proposal spawn:**
- **AI-generated name at spawn time** — a round-trip just to name a window; too slow for an instant action.
- **User-supplied name** — defeats the "spawn fast when busy" purpose.
- **Code-derived slug from the prompt** — strip stopwords, take salient tokens, dedupe, cap length, add uniqueness suffix. Microseconds, deterministic, no AI. Chosen.

**Identity stability across promotion:**
- **Rename handle when the AI names the change** — breaks every keyed lookup (workspace name, dir, cleanup/status); reintroduces the remap complexity.
- **Immutable handle = display name; never remap** — the slug is the identity forever; the agent's eventual change name lives *inside* the workspace, not as the handle. Chosen.

**Command shape:**
- **Three verbs (apply/explore/propose)** — but explore promoting to propose would force a handle rename. Rejected.
- **Two verbs (apply / proposal), explore is a `--mode` of proposal** — promotion is a status change, not a rename. Chosen.

## Decision

- **`spawn` has two flows behind subcommands.** `spawn apply <change>` validates the change exists, then isolates (today's path). `spawn proposal <prompt> [--mode]` derives a slug, isolates, and launches `/opsx:explore` or `/opsx:propose` — it does **not** call `checkOpenspecChange` (there is no change yet).
- **Code-derived slug, no AI.** A plain function turns the prompt into a short human-readable slug; it is the immutable handle **and** the display name. No phase-2 remap.
- **Verb-prefixed names:** `app-<change>`, `prop-<slug>`. Explore and propose share the `prop-` prefix (same identity, different seed mode) so promotion never renames.
- **Workspace name ≠ change name, by design, for proposals.** A `prop-<slug>` workspace contains a change dir the agent names differently. `merge`/`status` must stop assuming the two are equal.
- **All spawns isolated.** No workspace-less mode; proposal spawns write inside their own workspace so they cannot race the main working copy.

## Consequences

- **Positive**: Instant spawn (no AI round-trip); stable identity that survives promotion; readable `status` via prefixes; main session freed for conflict work.
- **Positive**: Forces merge to stay robust (every proposal eventually merges) — aligned with the jjay-ug7y/30gc hardening.
- **Negative**: `spawn` is now two divergent flows, and several call sites that assumed workspace-name == change-name must be generalized.
- **Negative**: Slug collisions/awkward slugs are possible; mitigated by a uniqueness suffix and accepting that the slug is a handle, not prose.
- **Negative**: Proposal spawns appear as jj workspaces with no openspec change — `status` must handle the "kind" distinction (two tables). The richer per-proposal lifecycle status is deferred to jjay-mg00 (it needs a state mechanism that fights ADR-006).
