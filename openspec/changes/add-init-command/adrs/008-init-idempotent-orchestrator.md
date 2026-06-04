# ADR-008: init is an idempotent, non-destructive orchestrator of external initializers

**Status**: Proposed

## Context

`jjay init` bootstraps a project for orchestration: openspec, the `/jjay:*` Claude integration, `AGENTS.md`, optional jj, optional hooks. Several of these have their own initializers (`openspec init`, `jj git init`) that jjay must not reimplement. init will also be run more than once on the same project (re-run after upgrading jjay, or to add a step skipped earlier), and on projects that already have *some* of the pieces. The question is how init treats existing state and how much it owns versus delegates.

This mirrors `jjay spawn`'s philosophy (ADR-003): check preconditions, delegate to the real tools (jj, tmux), don't roll back. And it inherits ADR-007's decision that the `/jjay:*` command + skill content authored in `add-claude-commands` installs into a *target* project's `.claude/`.

## Options Considered

- **Overwrite everything** — init writes all files unconditionally. Simple, but destroys a user's customized `config.yaml`/`AGENTS.md` on re-run. Unacceptable.
- **Fail if anything exists** — init refuses on a partially-initialized project. Safe but useless for the common "add the step I skipped" / "re-run after upgrade" cases.
- **Reimplement initializers** — jjay writes openspec/jj scaffolding itself. Duplicates and drifts from `openspec init` / `jj`.
- **Idempotent orchestrator** — init delegates each step to the canonical tool, treats existing artifacts as satisfied (skip, don't clobber), and only prompts/overwrites with explicit confirmation. Re-runnable; composes existing tools.

## Decision

`jjay init` is an **idempotent, non-destructive orchestrator**:

- **Delegate, don't reimplement.** openspec → `openspec init` (with `--tools claude` and the schema); jj → `jj`'s own init. jjay owns only the jjay-specific assets (the `/jjay:*` commands, the `jjay` skill, `AGENTS.md` content).
- **Idempotent.** Each step detects whether its artifact already exists and treats present-and-valid as done. Re-running init on a prepared project is a no-op (modulo new steps).
- **Non-destructive.** Existing user files (`config.yaml`, `AGENTS.md`, customized commands) are never overwritten without explicit confirmation; `--yes` accepts *creation* defaults but does not authorize *clobbering*. A separate `--force` would be required to overwrite.
- **Per-step, skippable.** Optional steps (jj, hooks) are opt-in; any step can be skipped via flags for non-interactive runs (mirroring `openspec init`'s `--tools`/`--force`).
- **No rollback.** Consistent with ADR-003: if a step fails mid-sequence, init reports it and leaves completed steps in place; re-running resumes.

## Consequences

- **Positive**: Safe to re-run — upgrades and incremental setup just work.
- **Positive**: One source of truth per concern (openspec owns openspec init; jjay owns the jjay assets).
- **Positive**: Non-interactive flags make init usable in `jjay init` automation and CI.
- **Negative**: "Detect if already done" logic per step is more code than blind overwrite; each step needs an existence/validity check.
- **Negative**: Partial-failure leaves a half-initialized project (no rollback). Mitigated by idempotency — re-running completes it.
- **Negative**: Couples init's release to the command/skill template content; the embedded assets must be kept in sync with `add-claude-commands`. Mitigated by embedding the same files, not a copy.
