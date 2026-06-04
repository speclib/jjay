---
# jjay-atq6
title: 'status: agent column + busy/finished via agent-emitted signal (view-configurable)'
status: draft
type: task
priority: normal
created_at: 2026-06-04T22:44:00Z
updated_at: 2026-06-04T23:55:54Z
parent: jjay-5y1a
---

The "hard half" of jjay-nnjz, split out because it fights ADR-006 (everything derived live from jj+tmux, no state file) and needs real design.

## Columns wanted
- AGENT: claude / codex / etc. — only knowable while ATTACHED (read tmux pane command); a detached spawn has no pane. Knowing it when detached would require persisting what spawn launched (state file → conflicts with ADR-006). Needs a decision.
- STATUS: busy / finished — genuinely ephemeral runtime state. No clean signal exists today: tmux pane activity is crude (a claude at an input prompt looks the same as mid-thought); tasks.md all-checked ≠ agent stopped.

## Chosen direction (user's idea)
Agents EMIT a status signal that jjay reads — a small beacon (e.g. a file the agent writes in its workspace, or a tmux/title signal) that jjay can poll. This is a NEW mechanism, not a derivation, so it deserves its own ADR weighing it against the no-state-file invariant (a workspace-local, agent-owned status file may be acceptable since it's not jjay state — it's the agent reporting itself).

## Also here
- view-configurable columns: `jjay status --columns=...` (or named views). Opt-in columns let sometimes-unknowable fields (AGENT when detached) be hidden rather than shown as "unknown".

## Relationship
- Splits from jjay-nnjz. The EASY half (rename STATUS→TMUX, add MERGED column) is its own proposal (add-status-merged-column / archive link TBD).
- AGENT/STATUS depend on the agent-signal mechanism; build that first, then the columns + the --columns view flag.
