### Requirement: Spawn creates jj workspace
The `jjay spawn <change>` command SHALL first run `jj new` in the main workspace to snapshot uncommitted work, then create a jj workspace named `<change>` at `../<project-name>-workspaces/<change>` via `jj workspace add --name <change> --revision @- <path>`. The `@-` revision contains the snapshotted files. The parent directory SHALL be created if it does not exist.

#### Scenario: Workspace created with isolation
- **WHEN** `jjay spawn feat-payments` is executed in project `jjay`
- **THEN** `jj new` is run first in the main workspace
- **THEN** `jj workspace add --name feat-payments --revision @- ../jjay-workspaces/feat-payments` is run
- **THEN** the child workspace contains all files from the snapshot
- **THEN** the main workspace's `@` is a fresh empty change (nothing to lose)

#### Scenario: Main workspace safe during concurrent work
- **WHEN** a spawned agent creates jj operations in the child workspace
- **THEN** the main workspace may become stale
- **THEN** running `jj workspace update-stale` in the main workspace does NOT lose work (because `@` is empty, all prior work is in `@-`)

#### Scenario: Workspace already exists
- **WHEN** `jjay spawn feat-payments` is executed and jj workspace `feat-payments` already exists
- **THEN** jjay exits with a non-zero exit code
- **THEN** `jj new` has NOT been run (preconditions checked before any mutations)
- **THEN** an error message is printed indicating the workspace already exists

### Requirement: Spawn creates tmux window
The command SHALL create a new tmux window in the target session named `ws-<change>` with its starting directory set to the workspace directory via `tmux new-window -c <wsDir>`.

#### Scenario: Window starts in workspace dir
- **WHEN** `jjay spawn feat-payments` is executed
- **THEN** the tmux window's starting directory is the workspace directory
- **THEN** verified via `tmux display-message -p -t <window> '#{pane_current_path}'`

#### Scenario: Window name already taken
- **WHEN** `jjay spawn feat-payments` is executed and a tmux window named `ws-feat-payments` already exists
- **THEN** jjay exits with a non-zero exit code
- **THEN** an error message is printed indicating the window name is taken

### Requirement: Spawn creates two-pane layout
The tmux window SHALL be split into two panes via `tmux split-window -h -c <wsDir>`. Both panes SHALL have their working directory set to the workspace directory at creation time. The left pane SHALL run the agent command. The right pane SHALL be a plain shell. No `send-keys cd` SHALL be used for setting working directories.

#### Scenario: Both panes in workspace dir
- **WHEN** `jjay spawn feat-payments` is executed successfully
- **THEN** both pane 0 and pane 1 have working directory set to the workspace directory
- **THEN** verified via `tmux display-message -p -t <pane> '#{pane_current_path}'`

### Requirement: Spawn requires openspec change
The command SHALL verify that an openspec change with the given name exists before proceeding. It uses `openspec list --json` to check.

#### Scenario: Change exists
- **WHEN** `jjay spawn feat-payments` is executed and openspec change `feat-payments` exists
- **THEN** spawn proceeds normally

#### Scenario: Change does not exist
- **WHEN** `jjay spawn feat-payments` is executed and no openspec change `feat-payments` exists
- **THEN** jjay exits with a non-zero exit code
- **THEN** an error message is printed indicating the change does not exist

### Requirement: Spawn requires tmux session
The command SHALL verify it is running inside a tmux session before proceeding.

#### Scenario: Inside tmux
- **WHEN** `jjay spawn` is executed inside a tmux session
- **THEN** spawn proceeds normally

#### Scenario: Outside tmux
- **WHEN** `jjay spawn` is executed outside a tmux session (no TMUX env var)
- **THEN** jjay exits with a non-zero exit code
- **THEN** an error message indicates jjay must be run inside tmux

### Requirement: Spawn requires change name argument
The command SHALL require exactly one argument: the change name.

#### Scenario: No argument
- **WHEN** `jjay spawn` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code

### Requirement: Spawn launches the agent via the resolved launch command
Spawning a workspace (`jjay spawn apply <change>` / `jjay spawn proposal <prompt>`) SHALL launch the agent using the resolved **launch** command for the selected agent, substituting `{change}`/`{prompt}`/`{wsdir}`. The launch and reopen paths SHALL share window/pane setup but SHALL diverge on which command is run: spawn uses `launch`, reopen uses `resume`. The `--agent` flag, when given, SHALL remain the highest-priority override of the resolved launch command.

#### Scenario: First spawn uses launch
- **WHEN** `jjay spawn apply foo` runs
- **THEN** the new window runs the resolved `launch` command (default `claude "/opsx:apply foo" --dangerously-skip-permissions --add-dir {wsdir}`)

#### Scenario: Flag overrides resolved launch
- **WHEN** `jjay spawn apply foo --agent '<fake-agent> {change}'` runs
- **THEN** the `--agent` value is used as the launch command, taking priority over config and built-in

### Requirement: Configurable tmux session
The `jjay spawn` and `jjay cleanup` commands SHALL accept a `--session` flag that specifies which tmux session to target. The default SHALL be the current tmux session. When set, all tmux commands (new-window, send-keys, split-window, kill-window, list-windows) SHALL target the specified session.

#### Scenario: Default session
- **WHEN** `jjay spawn feat-payments` is executed without `--session`
- **THEN** tmux commands target the current session

#### Scenario: Specific session
- **WHEN** `jjay spawn --session work feat-payments` is executed
- **THEN** tmux commands target session `work`

### Requirement: Configurable workspace root
The `jjay spawn` and `jjay cleanup` commands SHALL accept a `--workspace-root` flag that overrides the default workspace root directory. The default SHALL be `../<project-name>-workspaces`.

#### Scenario: Default root
- **WHEN** `jjay spawn feat-payments` is executed in project `jjay` without `--workspace-root`
- **THEN** workspace is created at `../jjay-workspaces/feat-payments`

#### Scenario: Custom root
- **WHEN** `jjay spawn --workspace-root /tmp/test-workspaces feat-payments` is executed
- **THEN** workspace is created at `/tmp/test-workspaces/feat-payments`

### Requirement: Integration test for spawn/cleanup lifecycle
A Go integration test (build tag `integration`) SHALL test the full lifecycle: spawn with fake agent → verify resources and pane directories → cleanup → verify all resources removed.

#### Scenario: Panes are in correct working directory
- **WHEN** the integration test inspects panes after spawn
- **THEN** both panes report the workspace directory as their current path

#### Scenario: Spawn creates all resources
- **WHEN** the integration test runs spawn with a fake agent
- **THEN** a jj workspace exists, tmux window exists, workspace directory exists, agent marker file exists

#### Scenario: Cleanup removes all resources
- **WHEN** the integration test runs cleanup after spawn
- **THEN** tmux window, jj workspace, and workspace directory are all gone

### Requirement: Spawn supports apply and proposal verbs
`jjay spawn` SHALL provide two verb subcommands and SHALL require one. `jjay spawn apply <change>` SHALL behave as the existing spawn (validate the change exists, isolate it, run `/opsx:apply`). `jjay spawn proposal <prompt>` SHALL spawn an agent seeded from a free-text prompt to create new work, without requiring a pre-existing openspec change. There is no bare `jjay spawn <change>` form — invoking `spawn` without a verb SHALL print usage and exit non-zero.

#### Scenario: Apply verb
- **WHEN** `jjay spawn apply add-foo` is run and `add-foo` is an existing openspec change
- **THEN** a workspace and window are created and an agent runs `/opsx:apply add-foo`

#### Scenario: Verb is required
- **WHEN** `jjay spawn add-foo` is run (no verb)
- **THEN** cobra prints usage and exits non-zero
- **THEN** no workspace or window is created

#### Scenario: Proposal verb needs no existing change
- **WHEN** `jjay spawn proposal "dark mode for settings"` is run
- **THEN** a workspace and window are created without requiring any existing openspec change
- **THEN** an agent runs seeded by the prompt

### Requirement: Proposal mode selects the seed command
`jjay spawn proposal` SHALL accept a mode (flag, with a configurable default) selecting whether the agent is seeded with `/opsx:explore` or `/opsx:propose`. Explore is a mode of a proposal spawn, not a separate verb.

#### Scenario: Explore mode
- **WHEN** `jjay spawn proposal "dark mode" --mode explore` is run
- **THEN** the agent is launched with `/opsx:explore` on the prompt

#### Scenario: Propose mode
- **WHEN** `jjay spawn proposal "dark mode" --mode propose` is run
- **THEN** the agent is launched with `/opsx:propose` on the prompt

### Requirement: Spawn names are verb-prefixed
Spawn workspace and window names SHALL be verb-prefixed: `app-<change>` for apply spawns and `prop-<slug>` for proposal spawns. The window name follows the existing `ws-` convention applied to the prefixed name.

#### Scenario: Apply prefix
- **WHEN** `jjay spawn apply add-foo` runs
- **THEN** the workspace is named `app-add-foo` and the window `ws-app-add-foo`

#### Scenario: Proposal prefix
- **WHEN** a proposal spawn derives slug `dark-mode`
- **THEN** the workspace is named `prop-dark-mode` and the window `ws-prop-dark-mode`

### Requirement: Proposal spawns derive a slug identity
A proposal spawn SHALL derive its identity from the prompt by deterministic code (no AI call): lowercase, strip punctuation and stopwords, keep salient tokens, cap length, and append a uniqueness suffix if the slug collides with an existing workspace or window. The derived slug SHALL be the immutable handle and the display name; it SHALL NOT be renamed after the agent creates a differently-named openspec change.

#### Scenario: Slug derived from prompt
- **WHEN** `jjay spawn proposal "add dark mode to the settings page"` runs
- **THEN** a short human-readable slug (e.g. `dark-mode-settings`) is derived without any AI call
- **THEN** the workspace/window use `prop-<slug>` and that name does not change later

#### Scenario: Slug collision
- **WHEN** a derived slug matches an existing spawn's name
- **THEN** a uniqueness suffix is appended so the new spawn's name is distinct

### Requirement: Proposal spawns are isolated and may produce a differently-named change
A proposal spawn SHALL run in its own jj workspace. The openspec change the agent creates inside that workspace MAY have a name different from the workspace slug. Commands that operate on spawns SHALL NOT assume the workspace name equals the openspec change name.

#### Scenario: Workspace name differs from produced change
- **WHEN** a `prop-dark-mode` proposal spawn's agent creates an openspec change named `add-dark-mode`
- **THEN** the workspace remains `prop-dark-mode`
- **THEN** the work is not lost or mis-keyed by merge/status due to the name difference

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
