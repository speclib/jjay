## Why

jjay exists so that an agent's work happens in an **isolated spawned workspace**, not in the main session. But nothing teaches Claude that. With only the `jjay` binary, an agent helping in this repo will reach for raw `/opsx:apply <change>` in the main working copy — exactly the manual, un-isolated flow jjay was built to replace. There is also no typeable entry point: a user must remember the binary's flags instead of an ergonomic `/jjay:spawn <change>`.

Bean [jjay-8xuj](../../.beans/jjay-8xuj--jjay-skill-and-or-command.md) ("jjay skill and/or command", parent [jjay-5y1a](../../.beans/jjay-5y1a--backlog.md) backlog) asks for the Claude Code integration layer. The "and/or" resolves to **both**: slash **commands** as the buttons, and a **skill** as the policy that makes Claude choose spawn-into-workspace over apply-in-place.

## What Changes

- **`/jjay:*` slash commands** mirroring the CLI surface — thin wrappers that run the `jjay` binary:
  - `/jjay:spawn <change>` → `jjay spawn <change>` (workspace + tmux window + agent running `opsx:apply`)
  - `/jjay:status` → `jjay status`
  - `/jjay:merge <change>` → `jjay merge <change>`
  - `/jjay:cleanup <change>` → `jjay cleanup <change>`
  - `/jjay:session-open <path>` → `jjay session-open <path>`
  Committed under `.claude/commands/jjay/`, so they propagate into every spawned workspace automatically (workspaces are jj copies of the repo).
- **`jjay` orchestrator skill** (`.claude/skills/jjay/SKILL.md`) — its `description` triggers auto-load in this repo and encodes the lifecycle policy: *implement a change by spawning an isolated agent workspace (`/jjay:spawn`), not by applying in the main session*. Documents the spawn → status → merge → cleanup loop and the orchestrator-vs-worker session distinction.
- **README** section documenting the `/jjay:*` commands and the skill.

## Capabilities

### New Capabilities
- `claude-commands`: the `/jjay:*` slash-command set and the `jjay` orchestrator skill — the Claude Code integration layer over the jjay binary.

### Modified Capabilities
<!-- None: this layer wraps existing CLI behavior; it does not change spawn/merge/cleanup/status/session-open requirements. -->

## Impact

- **Files**: new `.claude/commands/jjay/{spawn,status,merge,cleanup,session-open}.md`; new `.claude/skills/jjay/SKILL.md`; README update.
- **No Go code** — this is a documentation/prompt layer over the existing binary.
- **Depends on** the `jjay` binary being on `PATH` in the session (commands shell out to it).
- **Interacts with** `workspace-aware-session` (`/jjay:status`, reopen) — commands should cover those surfaces once that change lands; this proposal specifies the wrappers, not the underlying behavior.
- **ADRs**: ADR-007 (commands as thin wrappers; skill as policy).
- **Beans**: 8xuj → in-progress, linked here.
