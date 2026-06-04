## Context

A project must be prepared before jjay can orchestrate it: openspec (schema + `config.yaml`), the `/jjay:*` Claude integration, `AGENTS.md`, optional jj, optional hooks. Today that is a manual checklist. `jjay init` automates it. Bean [jjay-ofk7](../../.beans/jjay-ofk7--bootstrap-task.md) lists the steps; this design follows ADR-008 (idempotent, non-destructive orchestrator) and consumes the command/skill content authored in `add-claude-commands` per ADR-007 (install into the *target* project's `.claude/`).

External initializers to delegate to:
- `openspec init [path] --tools claude` (non-interactive via `--tools`/`--force`/`--profile`) — already the tool that scaffolds `openspec/`.
- `jj` own init — optional.

Existing patterns to mirror: cobra subcommand registration in `cmd/jjay/main.go`; the precondition-check style and "delegate to external tools, no rollback" stance of `internal/spawn` (ADR-003).

## Goals / Non-Goals

**Goals:**
- One `jjay init [path]` command that bootstraps openspec + jjay Claude integration + AGENTS.md, with jj/hooks optional.
- Idempotent and non-destructive: safe to re-run, never clobbers user files without `--force`.
- Non-interactive flags for automation.
- Embed the `/jjay:*` command + skill files as assets written into the target's `.claude/`.

**Non-Goals:**
- Authoring the command/skill *content* — that is `add-claude-commands`; init only installs it.
- Reimplementing `openspec init` or jj init.
- Spawning agents (no tmux involvement in init).
- Migrating/upgrading an existing openspec schema (init ensures presence, not migration).

## Decisions

- **New `internal/init` package + `initCmd`** in `cmd/jjay/main.go`, registered alongside spawn/merge/cleanup/status. `init [path]` with `path` defaulting to cwd.
- **Step pipeline, each step = detect → act.** openspec, claude-integration, agents-md are core; jj, hooks are opt-in. Each step checks for its artifact and skips if present-and-valid (idempotency, ADR-008). Order: openspec → claude integration → AGENTS.md → (jj) → (hooks).
- **Embed templates with `go:embed`.** The `/jjay:*` command files and `jjay/SKILL.md` are embedded from the repo's `.claude/` (the canonical copies from `add-claude-commands`) so init writes the exact same content into a target's `.claude/`. Single source — no drift between dogfooded and installed copies.
- **Delegate openspec.** Shell out to `openspec init <path> --tools claude` (+ `--force` only under jjay's `--force`); then ensure `config.yaml` exists, seeding from the schema template and prompting for project context (skipped under `--yes`).
- **Non-destructive contract.** `--yes` authorizes *creating* missing files with defaults; overwriting an existing user file requires `--force`. Each step reports created / skipped / would-overwrite.
- **No rollback (ADR-003 consistency).** A failed step is reported; completed steps remain; re-running resumes via idempotency.

## Risks / Trade-offs

- **Template drift** between the embedded assets and the live `.claude/` content — mitigated by embedding the actual repo files rather than a separate copy; a test can assert the embed matches `.claude/`.
- **`openspec init` interface changes** (flag names, tool list) could break the delegation — pin behavior with an integration test invoking the real binary; surface a clear error if `openspec` is absent.
- **"Valid" detection per step is fuzzy** — e.g. an `openspec/` dir that exists but is misconfigured. Init checks presence, not deep validity; deep migration is a non-goal. Document that init ensures presence.
- **Sequencing dependency on `add-claude-commands`** — init embeds assets that change must produce first. If init is built earlier, the embed target won't exist. Build `add-claude-commands` first.
