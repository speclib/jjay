## Purpose

Shell completion for the `jjay` CLI verbs that take a change/workspace name (`spawn`, `merge`, `cleanup`), so that pressing TAB suggests the right names from openspec changes and existing jj workspaces, fast and without side effects.

## Requirements

### Requirement: Completion handles the spawn verbs
With `jjay spawn` taking verb subcommands, completion SHALL operate at two levels: the verb, then the verb's argument. `jjay spawn <TAB>` SHALL suggest the verbs (`apply`, `proposal`). `jjay spawn proposal <TAB>` takes a free-text prompt and SHALL NOT offer change-name candidates.

#### Scenario: Spawn verb completion
- **WHEN** completion is requested for `jjay spawn <TAB>`
- **THEN** `apply` and `proposal` are offered

#### Scenario: Proposal takes a free-text prompt
- **WHEN** completion is requested for `jjay spawn proposal <TAB>`
- **THEN** no change-name candidates are offered

### Requirement: Spawn completes un-spawned changes
Shell completion for the change-name argument of `jjay spawn apply` SHALL suggest the names of openspec changes that do not currently have a spawned jj workspace (the set of all change names minus the set of existing workspace names). Completion SHALL NOT offer file-path fallback for this argument. (The completion attaches to the `apply` verb's argument; `jjay spawn` no longer has a bare positional form — see `add-spawn-verbs`.)

#### Scenario: Only un-spawned changes suggested
- **WHEN** openspec has changes `add-foo`, `fix-bar` and a workspace exists for `add-foo`, and the shell requests completion for `jjay spawn apply <TAB>`
- **THEN** `fix-bar` is offered
- **THEN** `add-foo` is not offered

#### Scenario: No file fallback
- **WHEN** completion is requested for `jjay spawn apply <TAB>`
- **THEN** the completion directive suppresses default file-name completion

### Requirement: Merge completes existing workspaces
Shell completion for the positional argument of `jjay merge` SHALL suggest the names of existing spawned jj workspaces.

#### Scenario: Mergeable workspaces suggested
- **WHEN** workspaces exist for `add-foo` and `fix-bar` and the shell requests completion for `jjay merge <TAB>`
- **THEN** `add-foo` and `fix-bar` are offered

#### Scenario: No workspaces
- **WHEN** no spawned workspaces exist and completion is requested for `jjay merge <TAB>`
- **THEN** no change-name candidates are offered

### Requirement: Cleanup completes existing workspaces
Shell completion for the positional argument of `jjay cleanup` SHALL suggest the names of existing spawned jj workspaces.

#### Scenario: Cleanable workspaces suggested
- **WHEN** workspaces exist for `add-foo` and `fix-bar` and the shell requests completion for `jjay cleanup <TAB>`
- **THEN** `add-foo` and `fix-bar` are offered

### Requirement: Completion excludes the main workspace
Completion candidates for merge and cleanup SHALL exclude the default (main) jj workspace, which is not a spawn.

#### Scenario: Default workspace not offered
- **WHEN** `jj workspace list` includes the `default` workspace and completion is requested for `jjay merge <TAB>`
- **THEN** `default` is not offered

### Requirement: Completion is fast and side-effect free
Completion SHALL derive candidates from name-only reads (`openspec list` and `jj workspace list`) without probing tmux or reading task files, and SHALL persist nothing.

#### Scenario: No tmux probing on completion
- **WHEN** completion candidates are computed for any of the three verbs
- **THEN** no tmux command is invoked

### Requirement: Completion degrades gracefully
If a candidate source cannot be read (e.g. `openspec` or `jj` unavailable, or not in a repo), completion SHALL return no candidates rather than erroring or blocking the shell.

#### Scenario: Source unavailable
- **WHEN** `openspec list` fails and completion is requested for `jjay spawn <TAB>`
- **THEN** completion returns no candidates and does not emit an error to the shell
