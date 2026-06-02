### Requirement: Spawn creates jj workspace
The `jjay spawn <change>` command SHALL create a jj workspace named `<change>` at `../<project-name>-workspaces/<change>` via `jj workspace add --name <change> --revision @ <path>`. The parent directory SHALL be created if it does not exist. Using `--revision @` ensures the workspace includes uncommitted files from the current working copy (e.g., the active openspec change).

#### Scenario: Workspace created successfully
- **WHEN** `jjay spawn feat-payments` is executed in project `jjay`
- **THEN** `jj workspace add --name feat-payments --revision @ ../jjay-workspaces/feat-payments` is run
- **THEN** the directory `../jjay-workspaces/feat-payments` exists as a jj workspace
- **THEN** the workspace contains uncommitted files from the main working copy

#### Scenario: Workspace already exists
- **WHEN** `jjay spawn feat-payments` is executed and jj workspace `feat-payments` already exists
- **THEN** jjay exits with a non-zero exit code
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
