## ADDED Requirements

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
