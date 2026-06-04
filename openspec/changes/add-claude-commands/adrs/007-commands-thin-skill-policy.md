# ADR-007: Slash commands are thin binary wrappers; the skill carries policy

**Status**: Proposed

## Context

jjay needs a Claude Code integration so agents (and users) drive the tool the way it was designed: spawn isolated workspaces rather than apply changes in place. Claude Code offers two mechanisms — slash **commands** (explicit, user-typed, `.claude/commands/`) and **skills** (auto-loaded by `description` match, `.claude/skills/`). The question is how to divide responsibility between them, and how much logic to put in the commands.

A second force: spawned workspaces are jj working copies of the same repo, and `.claude/` is checked in. So whatever lives in `.claude/commands/` and `.claude/skills/` is present in *every* spawned workspace too — the integration is self-propagating, but it must therefore behave sensibly in both the orchestrator session and a worker session.

## Options Considered

- **Commands only** — `/jjay:spawn` etc., no skill. Gives buttons, but nothing makes the model *prefer* spawning over manual `/opsx:apply`. The original problem persists when nobody types the command.
- **Skill only** — one skill describing the whole workflow, no commands. Fixes the policy gap, but no ergonomic typeable entry point and no discoverable per-verb surface.
- **Fat commands** — reimplement spawn/merge logic in the command prompts. Duplicates the binary; drifts from it; two sources of truth for the same behavior.
- **Thin commands + policy skill** — commands are minimal wrappers that just invoke `jjay <verb>`; a single skill encodes when/why to use them and the lifecycle. One behavior source (the binary), one policy source (the skill).

## Decision

- **Commands are thin wrappers.** Each `/jjay:<verb>` command's job is to invoke `jjay <verb> <args>` and relay output. No orchestration logic in the prompt — the binary owns behavior (consistent with ADR-006: the binary derives state live).
- **The skill carries policy.** A single `jjay` skill encodes the lifecycle (explore → propose → spawn → status → merge → cleanup) and the rule *implement changes by spawning an isolated workspace, not by `/opsx:apply` in the main session*. Its `description` is written to auto-trigger when the conversation is about implementing/managing changes in this repo.
- **Orchestrator vs worker.** The skill names the distinction: the **orchestrator** session is where `/jjay:spawn` is run; the spawned **worker** session runs `/opsx:apply` inside its workspace. A worker should not recursively spawn. The skill states this so the self-propagated commands don't cause nesting.

## Consequences

- **Positive**: Single source of behavior (binary) and single source of policy (skill); commands can't drift from the binary.
- **Positive**: Self-propagating via checked-in `.claude/` — every workspace has the same commands.
- **Positive**: Fixes the root cause (model choosing manual apply) via the skill, not just symptoms.
- **Negative**: Requires the `jjay` binary on `PATH` in the session; a command in a workspace without the binary fails. Acceptable — jjay is the prerequisite of the whole repo.
- **Negative**: The orchestrator/worker rule is prompt-enforced, not mechanically enforced; a misbehaving worker could still try to spawn. Mitigated by the skill's explicit guidance; a hard guard in the binary is a possible future ADR.
