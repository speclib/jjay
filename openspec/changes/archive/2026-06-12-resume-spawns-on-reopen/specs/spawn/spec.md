## ADDED Requirements

### Requirement: Reopen resumes the agent instead of re-running apply
When reopening a previously-spawned workspace (via `session-open` or `tmux-open`), jjay SHALL run that workspace's resolved `resume` command, NOT its `launch` command. Reopening SHALL NOT re-run `/opsx:apply` (or `/opsx:explore` / `/opsx:propose`) from scratch. The window SHALL open with the workspace directory as its working directory so the resume command is scoped to that workspace.

#### Scenario: session-open resumes, does not re-apply
- **WHEN** a tmux session died and `jjay session-open <path>` reopens a detached apply spawn `app-foo`
- **THEN** the reopened window runs the resolved `resume` command (default `claude --resume --add-dir {wsdir}`)
- **THEN** the window does NOT run `claude "/opsx:apply foo" …`

#### Scenario: Resume command is user-configured verbatim
- **WHEN** the resolved `resume` for an agent is set by the user (e.g. `claude --continue …`)
- **THEN** jjay substitutes `{change}`/`{prompt}`/`{wsdir}` and runs that command as-is
- **THEN** jjay does not impose its own resume semantics

### Requirement: tmux-open reopens a single workspace
jjay SHALL provide a `jjay tmux-open <workspace>` command that recreates the tmux window and panes for one existing spawned workspace, in the current session (or `--session`), and launches the agent's `resume` command. It SHALL NOT create or modify the jj workspace; it SHALL operate on an already-existing workspace.

#### Scenario: Reopen one workspace
- **WHEN** `jjay tmux-open app-foo` runs and the `app-foo` workspace exists with no live window
- **THEN** a `ws-app-foo` window is created in the target session, opened in the workspace directory
- **THEN** the window runs the agent's resolved `resume` command

### Requirement: session-open reopens via the single-workspace primitive
`jjay session-open` SHALL reopen all detached spawns by looping over the same single-workspace reopen primitive that `tmux-open` exposes, so reopen has exactly one code path. Reopening SHALL remain best-effort: a per-workspace failure SHALL be logged and the remaining workspaces SHALL still be reopened.

#### Scenario: All detached workspaces reopened via resume
- **WHEN** `jjay session-open <path>` runs with multiple detached spawns
- **THEN** each detached workspace is reopened by the shared primitive using its `resume` command
- **THEN** a workspace that already has a live window is skipped (no duplicate)

#### Scenario: Per-workspace reopen failure is non-fatal
- **WHEN** one workspace fails to reopen during `session-open`
- **THEN** the failure is logged and the remaining workspaces are still reopened
- **THEN** `session-open` still succeeds overall

### Requirement: session-open reopens the TARGET repo's workspaces, not the caller's
`jjay session-open <path>` SHALL enumerate and reopen workspaces of the jj repository at `<path>`, NOT of the directory jjay is invoked from. Workspace enumeration SHALL be scoped to the target repo (e.g. `jj -R <path> workspace list`), and reopened workspace directories SHALL resolve under the target repo's workspace root.

#### Scenario: Cross-project isolation (jjay-02nr)
- **WHEN** `jjay session-open ~/other/proj` is run from inside the jjay repo, and jjay has a spawn `app-foo` while `~/other/proj` has a spawn `app-bar`
- **THEN** the new `jjay->proj` session reopens `app-bar` (the target repo's spawn)
- **THEN** it does NOT reopen `app-foo` (the caller repo's spawn)

### Requirement: tmux session names are sanitized for tmux target syntax
The tmux session name derived from a repo path SHALL replace characters tmux reserves in target names (`.` and `:`) with `_`, so the name used to create the session and the name used to target it (switch-client, list-windows) agree.

#### Scenario: Dotted repo dir name (jjay-e3bx)
- **WHEN** `jjay session-open ~/cLinden/mip.rs/` is run
- **THEN** the session is created and targeted as `jjay->mip_rs` (dot normalized to `_`)
- **THEN** switching to the session succeeds (no "can't find pane" error)

## MODIFIED Requirements

### Requirement: Spawn launches the agent via the resolved launch command
Spawning a workspace (`jjay spawn apply <change>` / `jjay spawn proposal <prompt>`) SHALL launch the agent using the resolved **launch** command for the selected agent, substituting `{change}`/`{prompt}`/`{wsdir}`. The launch and reopen paths SHALL share window/pane setup but SHALL diverge on which command is run: spawn uses `launch`, reopen uses `resume`. The `--agent` flag, when given, SHALL remain the highest-priority override of the resolved launch command.

#### Scenario: First spawn uses launch
- **WHEN** `jjay spawn apply foo` runs
- **THEN** the new window runs the resolved `launch` command (default `claude "/opsx:apply foo" --dangerously-skip-permissions --add-dir {wsdir}`)

#### Scenario: Flag overrides resolved launch
- **WHEN** `jjay spawn apply foo --agent '<fake-agent> {change}'` runs
- **THEN** the `--agent` value is used as the launch command, taking priority over config and built-in
