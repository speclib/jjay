## Context

Claude Code resolves `/foo:bar` to `.claude/commands/foo/bar.md` (a prompt template with frontmatter) and auto-loads skills from `.claude/skills/<name>/SKILL.md` by matching the skill's `description` against the conversation. This repo already ships `/opsx:*` commands and `openspec-*` skills the same way (`.claude/commands/opsx/apply.md`, `.claude/skills/openspec-apply-change/SKILL.md`) — this change adds the `jjay` equivalents. Per ADR-007: commands are thin binary wrappers, the skill carries policy.

## Goals / Non-Goals

**Goals:**
- A `/jjay:<verb>` command per CLI verb (spawn, status, merge, cleanup, session-open) that shells out to `jjay`.
- A single `jjay` skill that makes Claude prefer spawn-into-workspace over apply-in-place, and documents the lifecycle.
- Self-propagation: commands + skill checked into `.claude/`, so they exist in every spawned workspace.

**Non-Goals:**
- Changing any jjay CLI behavior (that's the binary's specs).
- A TUI ([jjay-bc1c](../../.beans/jjay-bc1c--jjay-tui-mode.md)).
- Mechanically preventing a worker from spawning (prompt-enforced for now; ADR-007 notes a hard guard as future work).
- Auto-installing the commands into *other* repos — scope is this repo's `.claude/`.

## Decisions

- **Mirror, don't invent.** Match the existing `/opsx:*` command file shape (frontmatter: `name`, `description`, `category`, `tags`; body = steps). One file per verb under `.claude/commands/jjay/`.
- **Thin body.** Each command body: resolve args (prompt if a required arg is missing, listing candidates via `openspec list --json` for change-name verbs), run `jjay <verb> <args>`, relay output. No reimplementation.
- **Skill `description` is the trigger.** Phrase it so it fires on "implement/work on/manage a change in this repo" — the moment the model would otherwise reach for `/opsx:apply`. The body states the spawn-first rule, the lifecycle, and the orchestrator/worker split.
- **Self-propagation is a feature, handled by the worker rule.** Because `.claude/` is committed, workers inherit the commands. The skill's orchestrator/worker section prevents a worker from re-spawning (which would nest workspaces).
- **Binary on PATH is a precondition.** Commands assume `jjay` resolves; document in README. (Future: the skill could detect absence and point at install.)

## Risks / Trade-offs

- **Skill over-triggering** — too-broad a `description` loads it in unrelated conversations. Mitigate by scoping wording to this repo's change-implementation context; tune after observing.
- **Command/CLI drift if bodies grow logic** — guard by keeping bodies thin (ADR-007); the binary is the only behavior source.
- **Worker recursion** — a worker that ignores the rule could `/jjay:spawn` inside its own workspace. Prompt-enforced now; flagged for a possible binary guard later.
- **PATH dependency in workers** — a spawned workspace must also have `jjay` on PATH; true today since workers run in the same environment, but worth a README note.
