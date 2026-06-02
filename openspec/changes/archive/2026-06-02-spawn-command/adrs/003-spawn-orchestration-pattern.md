# ADR-003: Spawn orchestration — sequential subprocess, no rollback

**Status**: Proposed

## Context

`jjay spawn` creates three things in sequence: a jj workspace, a tmux window, and panes running an agent + shell. Each step depends on the previous one succeeding. The question is how to handle failures mid-sequence.

## Options Considered

- **Sequential with rollback** — undo completed steps on failure. Complex: need to track what was created, handle rollback failures.
- **Sequential with no rollback** — fail and leave partial state. Simple: user can clean up manually. Consistent with how jj and tmux already work (they don't roll back each other's state).
- **Check-then-act** — validate all preconditions upfront, then execute. Reduces but doesn't eliminate race conditions.

## Decision

Check all preconditions first (tmux session, openspec change, workspace doesn't exist, window name free), then execute sequentially with no rollback. If a step fails after workspace creation, the user cleans up manually.

## Consequences

- **Positive**: Simple implementation, no rollback state machine
- **Positive**: Precondition checks catch most failures before any side effects
- **Negative**: Rare mid-sequence failures leave partial state (workspace exists but no tmux window, etc.)
- **Negative**: User must know `jj workspace forget` and `tmux kill-window` to clean up — but jjay cleanup command will handle this later
