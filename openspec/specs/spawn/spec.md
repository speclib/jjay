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
The command SHALL create a new tmux window in the current session named `ws-<change>`.

#### Scenario: Window created
- **WHEN** `jjay spawn feat-payments` is executed
- **THEN** a tmux window named `ws-feat-payments` is created in the current tmux session

#### Scenario: Window name already taken
- **WHEN** `jjay spawn feat-payments` is executed and a tmux window named `ws-feat-payments` already exists
- **THEN** jjay exits with a non-zero exit code
- **THEN** an error message is printed indicating the window name is taken

### Requirement: Spawn creates two-pane layout
The tmux window SHALL be split into two panes. The left pane SHALL run the claude agent with `/opsx:apply <change>` and `--add-dir` pointing to the workspace directory (to bypass the trust dialog). The right pane SHALL be a shell cd'd to the workspace directory.

#### Scenario: Two-pane layout
- **WHEN** `jjay spawn feat-payments` is executed successfully in project `jjay`
- **THEN** the tmux window has two panes side by side
- **THEN** the left pane runs `claude "/opsx:apply feat-payments" --dangerously-skip-permissions --add-dir ../jjay-workspaces/feat-payments`
- **THEN** the right pane is a shell with working directory `../jjay-workspaces/feat-payments`

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

### Requirement: Configurable agent command
The `jjay spawn` command SHALL accept an `--agent` flag that specifies the full command to run in the agent pane. The default SHALL be `claude "/opsx:apply <change>" --dangerously-skip-permissions --add-dir <wsDir>`. The `<change>` and `<wsDir>` placeholders SHALL be substituted at runtime.

#### Scenario: Default agent
- **WHEN** `jjay spawn feat-payments` is executed without `--agent`
- **THEN** the agent pane runs the default claude command

#### Scenario: Custom agent
- **WHEN** `jjay spawn --agent "./fake-agent.sh {change}" feat-payments` is executed
- **THEN** the agent pane runs `./fake-agent.sh feat-payments`

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
