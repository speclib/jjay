---
# jjay-mg00
title: proposal status lifecycle enum (explore→draft→review→accepted)
status: draft
type: task
priority: normal
created_at: 2026-06-04T23:33:47Z
updated_at: 2026-06-05T10:05:53Z
parent: jjay-5y1a
blocked_by:
    - jjay-atq6
---

Split out of jjay-4ulx (spawn verbs). A `proposal` spawn (prop-<slug>) should carry a lifecycle status, since explore is just the earliest state of a proposal, not a separate thing.

## Status enum
[explore, draft, ready-for-review, accepted, request-for-changes]

- explore: agent is thinking/sketching (seeded by /opsx:explore)
- draft: artifacts exist, not yet reviewed
- ready-for-review: author hands it to the user
- accepted: becomes a real change → can be applied (app-<change>)
- request-for-changes: bounced back to the agent

Promotion explore→propose is a STATUS transition, NOT a rename — the handle prop-<slug> is immutable (decided in 4ulx).

## The hard question: where does status live?
It is NOT derivable from jj+tmux, so it fights ADR-006 (no state file). Likely shares a mechanism with jjay-atq6 (agent emits a status signal jjay reads). Candidates:
- a workspace-local file the agent writes (agent-owned, not jjay state — may be acceptable under ADR-006)
- map onto openspec change status once the change exists
- jjay-managed metadata (needs an ADR weighing it against the no-state-file invariant)

## Surfaces
- `jjay status` second table (proposal spawns) shows this enum per spawn (the two-table view from 4ulx).
- transitions: who sets them? CLI cmd (`jjay proposal <slug> --status ready`)? agent signal? user?

## Relationships
- depends on jjay-4ulx (spawn proposal must exist first)
- likely shares the signal mechanism with jjay-atq6 (agent-emitted busy/finished + agent name)
