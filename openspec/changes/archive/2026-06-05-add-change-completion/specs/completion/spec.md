## ADDED Requirements

### Requirement: Spawn completes un-spawned changes
Shell completion for the positional argument of `jjay spawn` SHALL suggest the names of openspec changes that do not currently have a spawned jj workspace (the set of all change names minus the set of existing workspace names). Completion SHALL NOT offer file-path fallback for this argument.

#### Scenario: Only un-spawned changes suggested
- **WHEN** openspec has changes `add-foo`, `fix-bar` and a workspace exists for `add-foo`, and the shell requests completion for `jjay spawn <TAB>`
- **THEN** `fix-bar` is offered
- **THEN** `add-foo` is not offered

#### Scenario: No file fallback
- **WHEN** completion is requested for `jjay spawn <TAB>`
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

### Requirement: Completion handles the spawn verbs
Once `jjay spawn` takes verb subcommands (`apply`, `proposal` — see `add-spawn-verbs`), completion SHALL operate at two levels: completing the verb itself, then completing the verb's argument. `jjay spawn <TAB>` SHALL suggest the verbs (`apply`, `proposal`). `jjay spawn apply <TAB>` SHALL suggest the filtered, un-spawned change names (the behavior previously specified for the bare positional argument now applies after the `apply` verb). `jjay spawn proposal <TAB>` takes a free-text prompt and SHALL NOT offer change-name candidates.

> Note: this requirement was recorded after archive, capturing the interaction with `add-spawn-verbs` (which introduced the verbs after this change shipped). It documents the intended behavior; the wiring lives in `add-spawn-verbs`. It was NOT synced into `openspec/specs/completion/spec.md` (archives are frozen) — the live requirement belongs in the `add-spawn-verbs` spec.

#### Scenario: Spawn verb completion
- **WHEN** completion is requested for `jjay spawn <TAB>`
- **THEN** `apply` and `proposal` are offered

#### Scenario: Apply argument completion
- **WHEN** completion is requested for `jjay spawn apply <TAB>`
- **THEN** the filtered un-spawned change names are offered (as for the former bare `spawn <TAB>`)

#### Scenario: Proposal takes a free-text prompt
- **WHEN** completion is requested for `jjay spawn proposal <TAB>`
- **THEN** no change-name candidates are offered (the argument is a free-text prompt)
