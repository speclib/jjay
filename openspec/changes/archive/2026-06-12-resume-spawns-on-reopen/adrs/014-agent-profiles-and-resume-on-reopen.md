# ADR-014: Agent profiles + 3-layer config — resume on reopen, not re-apply

**Status**: Accepted

**Supersedes**: ADR-006 (Configuration via CLI flags, not config file) for the agent-command surface.

## Context

`jjay session-open` reopens detached workspaces by re-running the original launch command (`internal/session/session.go` → `spawn.OpenWindow` → `setupPanes` → `send-keys claude "/opsx:apply <change>" …`). Reopening a dead session therefore starts a fresh Claude conversation and re-runs `/opsx:apply` from scratch, discarding the in-flight work. Bean jjay-tzl3: *"Ideally they should resume the claude session but leave it from there"*, *"per supported agent we should know what the /resume command is"*, and *"a new command is needed: 'open workspace in tmux'"*.

The root cause is structural: a single command template per flow, and a deliberate invariant that *"Spawn and session-open cannot diverge"*. Launch and resume are different intents; collapsing them is the bug.

ADR-006 chose CLI flags over a config file and explicitly deferred the file "until usage patterns emerge", predicting repetition would be "a future config file problem". Per-agent **resume** knowledge is that pattern: a flag cannot carry two commands (launch + resume) per agent, and resume syntax is agent-specific.

## Options Considered

- **Status quo (re-run launch on reopen).** Re-runs `/opsx:apply` from zero, destroying the conversation. Rejected (the bug).
- **Hardcode `--continue` on reopen (Go const).** Fixes the behavior but bakes resume semantics into jjay; can't vary per agent or per user. Rejected — the task wants resume to be agent-configurable.
- **Single config layer (project file only).** Simpler, but no user-global defaults and no clean seam for jjay-iex3's global config dir. Rejected in favor of layering.
- **Agent profiles `{launch, resume}` + 3-layer per-field config (chosen).** Reopen runs the resolved `resume`; spawn runs `launch`; the user owns both strings; one reopen primitive (`tmux-open`) that `session-open` loops. Closes all three asks and lays the config foundation jjay-iex3/jjay-euup extend.

## Decision

1. **Agent = profile `{Launch, Resume}`** (command templates over `{change}`/`{prompt}`/`{wsdir}`), replacing the single agent string.
2. **Intent on the window opener.** Spawn flows pass **launch**; reopen flows (`tmux-open`, `session-open`) pass **resume**. Launch and reopen share window/pane setup but diverge on the command. The "cannot diverge" invariant is inverted.
3. **3-layer, per-field config.** `resolve(agent, field) = project ?? global ?? builtin`, per field, empty = unset; over an ordered source list so layers can be added later. Files: `<repo>/.jjay/config.yaml`, `~/.config/jjay/config.yaml` — both distinct from `openspec/config.yaml`.
4. **Built-in is the single source of truth.** The current `DefaultAgentCommand` becomes `builtin["claude"].Launch`; `Resume` defaults to `claude --resume --add-dir {wsdir}`. `jjay init` materializes the built-in into `<repo>/.jjay/config.yaml` (idempotent, non-destructive — ADR-008). Seeded file and runtime fallback share the const, so they cannot drift.
5. **`tmux-open <workspace>` is the single reopen primitive.** `session-open` loops it. The `--agent` flag remains the highest-priority launch override (integration tests unaffected).

## Consequences

- **Positive**: Reopen resumes instead of re-applying — the core fix. Resume is user/agent-configurable. One reopen code path. A real config foundation; ADR-006's deferred file now exists.
- **Positive**: `--resume` default hands control back to the human ("leave it from there"); a user wanting hands-off resume sets `resume: claude --continue …`.
- **Negative**: Reverses ADR-006's "no config file" — more precedence surface (mitigated by per-field fallback and an ordered-source resolver). Recorded as a supersede.
- **Negative**: Default `--resume` is interactive — `session-open` with N detached spawns opens N windows each awaiting a human pick (intended, but a visible change from the hands-off launch flow).
- **Risk**: `--resume` scoping rides on the window cwd (`-c wsDir`); must be verified end-to-end, not assumed.
- **Deferred**: holistic verb rename; non-`claude` agents; richer per-project settings (jjay-euup, built on this foundation).
