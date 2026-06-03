## ADDED Requirements

### Requirement: Session-open creates tmux session
The `jjay session-open <path>` command SHALL create a new tmux session named `jjay-><dirname>` where `<dirname>` is the basename of the given path. The session's working directory SHALL be set to the given path.

#### Scenario: Session created
- **WHEN** `jjay session-open ~/projects/myapp` is executed
- **THEN** a tmux session named `jjay->myapp` is created
- **THEN** the session's working directory is `~/projects/myapp`

### Requirement: Session-open switches to the session
After creating the session, the command SHALL switch the current tmux client to the new session.

#### Scenario: Client switches
- **WHEN** `jjay session-open ~/projects/myapp` completes
- **THEN** the user is in the `jjay->myapp` tmux session

### Requirement: Session-open requires jj repo
The command SHALL verify the given path contains a jj repository (`.jj/` directory exists).

#### Scenario: Valid jj repo
- **WHEN** `jjay session-open ~/projects/myapp` is executed and `~/projects/myapp/.jj/` exists
- **THEN** session creation proceeds

#### Scenario: Not a jj repo
- **WHEN** `jjay session-open /tmp/notarepo` is executed and no `.jj/` directory exists
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the path is not a jj repository

### Requirement: Session-open rejects duplicate sessions
The command SHALL check if a session with that name already exists.

#### Scenario: Session already exists
- **WHEN** `jjay session-open ~/projects/myapp` is executed and `jjay->myapp` session exists
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the session already exists

### Requirement: Session-open requires path argument
The command SHALL require exactly one argument: the path to the repo.

#### Scenario: No argument
- **WHEN** `jjay session-open` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code
