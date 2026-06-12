## Why

`jjay session-open` reopens every previously-spawned workspace by **re-running the exact command that originally launched it** — for an apply spawn, `claude "/opsx:apply <change>" --dangerously-skip-permissions …`. Reopening a dead session therefore starts a **brand-new** Claude conversation and re-runs `/opsx:apply` from zero, discarding the in-flight conversation that was already doing the work. Bean [jjay-tzl3](../../.beans/jjay-tzl3--session-open-should-not-run-claude-opsxapply.md) states the intent directly: *"Ideally they should resume the claude session but leave it from there."*

**Root cause.** `internal/spawn/spawn.go` has a single command template per spawn flow and routes **both** first-launch and reopen through `OpenWindow`/`setupPanes`. The doc comment treats "Spawn and session-open cannot diverge" as a virtue — but that is exactly the bug. Launch and resume are different intents:

```
  FIRST SPAWN  ──► claude "/opsx:apply foo" …     ✓ correct: start the work
  SESSION REOPEN ─► claude "/opsx:apply foo" …    ✗ wrong: re-runs apply from scratch
```

The task names two more needs:
- *"per supported agent we should know what the /resume command is"* — resume syntax is agent-specific, so it belongs in per-agent config alongside the launch command.
- *"a new command is needed: 'open workspace in tmux'"* — reopening **one** workspace's window is a distinct primitive from opening a whole session.

There is no jjay config file today: agent commands are a hardcoded Go const overridable only by the `--agent` flag (ADR-006, which explicitly deferred a config file "until usage patterns emerge" and predicted "that's a future config file problem"). Per-agent *launch + resume* knowledge is that future arriving — a flag cannot naturally carry two commands per agent.

## What Changes

- **Launch ≠ resume (the fix).** An agent is no longer a single command string but a profile `{ launch, resume }`. `OpenWindow` gains an **intent**: spawn flows use `launch`; reopen flows use `resume`. They still share window/pane setup — only the command string diverges.
- **`session-open` reopens via resume.** `reopenDetached` SHALL run each detached workspace's resolved `resume` command instead of re-running `/opsx:apply`. The default `resume` for `claude` is `claude --resume --add-dir {wsdir}` (interactive picker — hands the conversation back to the human to "leave it from there"). The user owns this string; jjay does not decide what resume *means*.
- **New `jjay tmux-open <workspace>` command.** Reopens a single workspace's tmux window using the resume intent, in the current session. `session-open`'s `reopenDetached` SHALL loop over this same primitive, so there is one reopen code path.
- **jjay config file (3-layer, per-field).** Introduce jjay's own config (distinct from `openspec/config.yaml`):
  - global `~/.config/jjay/config.yaml`, project `<repo>/.jjay/config.yaml`, and a Go built-in default.
  - Resolution is **per field**: `project ?? global ?? builtin` for each of `launch`/`resume`. A user who sets only `launch` still inherits the default `resume`.
  - The Go built-in is the single source of truth: `jjay init` materializes it into `<repo>/.jjay/config.yaml`, and runtime falls back to it, so seeded file and fallback cannot drift.

```yaml
# .jjay/config.yaml  (same schema at ~/.config/jjay/config.yaml)
agents:
  claude:
    launch: 'claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}'
    resume: 'claude --resume --add-dir {wsdir}'
```

## Capabilities

### New Capabilities
- `config`: jjay SHALL resolve per-agent `launch`/`resume` command templates from a 3-layer source (project `.jjay/config.yaml` → global `~/.config/jjay/config.yaml` → Go built-in), merged per field, and SHALL seed the project file from the built-in on `init`.

### Modified Capabilities
- `spawn`: spawning SHALL distinguish a **launch** intent (first spawn → `launch` template) from a **resume** intent (reopen → `resume` template); `session-open` and the new `tmux-open` command SHALL reopen workspaces via the resume intent rather than re-running `/opsx:apply`. A new `jjay tmux-open <workspace>` command SHALL reopen a single workspace's window; `session-open` SHALL reopen all detached workspaces by looping that same primitive.

## Impact

- **Code**:
  - New `internal/config` package: types `{Agents map[string]AgentProfile}`, `AgentProfile{Launch, Resume}`, a per-field 3-layer resolver, and the Go built-in default (the current `DefaultAgentCommand` becomes the built-in `claude.launch`).
  - `internal/spawn/spawn.go`: `OpenWindow`/`openWindow`/`setupPanes` gain an intent; the "cannot diverge" comment is inverted. `Spawn`/`SpawnProposal` resolve via config (built-in still wins for tests via `--agent`).
  - `internal/session/session.go`: `reopenDetached` delegates to the new single-workspace reopen primitive (resume intent). **Also bundles two session-open bug fixes surfaced during implementation:** `SessionName` sanitizes `.`/`:` → `_` for tmux target syntax ([jjay-e3bx](../../.beans/jjay-e3bx--session-open-fails-on-dotscolons-in-repo-dir-name.md)); and reopen is scoped to the **target** repo via `status.ListIn(repoRoot, …)` so `session-open <path>` no longer reopens the caller repo's workspaces into the new session ([jjay-02nr](../../.beans/jjay-02nr--session-open-respawns-tmux-windows-in-wrong-tmux-s.md)).
  - `cmd/jjay/main.go`: new `tmux-open` command (with arg completion over reopenable workspaces).
  - `internal/init`: seed `<repo>/.jjay/config.yaml` from the built-in (idempotent, non-destructive — ADR-008 pattern).
  - Command doc + skill assets under `internal/init/assets/` for `tmux-open` and the updated `session-open` behavior.
- **Severity**: HIGH — `session-open` currently destroys in-flight agent conversations on every reopen.
- **Relation**: **supersedes ADR-006** (config-via-flags-not-file) for the agent-command surface — the deferred config file now exists; flags remain valid overrides. Slots a config foundation under backlog tasks [jjay-iex3](../../.beans/jjay-iex3--configuration-dir.md) (global config dir) and [jjay-euup](../../.beans/jjay-euup--settings-in-project-for-spawned-tmux-layout-and-ag.md) (per-project settings), which shift from "build the config mechanism" to "add more fields (tmux layout, panes, agent-to-use)" on top of it.
- **ADRs**: ADR-014 (agent profiles + 3-layer per-field config; launch≠resume on reopen).
- **Beans**: jjay-tzl3 → in-progress, linked here. Bundled fixes: jjay-e3bx (dot/colon session names), jjay-02nr (cross-project reopen leak) → in-progress, linked here.
- **Bean task ref**: [jjay-tzl3](../../.beans/jjay-tzl3--session-open-should-not-run-claude-opsxapply.md)

## Deferred / Out of Scope

- **Global verb rename.** The task hints `tmux-open` is "more informative" and that command names should be refactored holistically later. This change adds the single well-named `tmux-open` command but does **not** rename `session-open`, `spawn apply`, etc. — that is a separate future task/bean.
- **Agents other than `claude`.** The profile abstraction *allows* `codex`/`mistral`/etc., but only `claude` is populated.
- **Richer per-project settings** (tmux layout, extra panes, `nix develop`) — owned by jjay-euup, built on this config foundation.
