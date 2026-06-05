## ADDED Requirements

### Requirement: Completion handles the spawn verbs
With `jjay spawn` taking verb subcommands, completion SHALL operate at two levels: the verb, then the verb's argument. `jjay spawn <TAB>` SHALL suggest the verbs (`apply`, `proposal`). `jjay spawn proposal <TAB>` takes a free-text prompt and SHALL NOT offer change-name candidates.

#### Scenario: Spawn verb completion
- **WHEN** completion is requested for `jjay spawn <TAB>`
- **THEN** `apply` and `proposal` are offered

#### Scenario: Proposal takes a free-text prompt
- **WHEN** completion is requested for `jjay spawn proposal <TAB>`
- **THEN** no change-name candidates are offered

## MODIFIED Requirements

### Requirement: Spawn completes un-spawned changes
Shell completion for the change-name argument of `jjay spawn apply` SHALL suggest the names of openspec changes that do not currently have a spawned jj workspace (the set of all change names minus the set of existing workspace names). Completion SHALL NOT offer file-path fallback for this argument. (The completion attaches to the `apply` verb's argument; `jjay spawn` no longer has a bare positional form — see `add-spawn-verbs`.)

#### Scenario: Only un-spawned changes suggested
- **WHEN** openspec has changes `add-foo`, `fix-bar` and a workspace exists for `add-foo`, and the shell requests completion for `jjay spawn apply <TAB>`
- **THEN** `fix-bar` is offered
- **THEN** `add-foo` is not offered

#### Scenario: No file fallback
- **WHEN** completion is requested for `jjay spawn apply <TAB>`
- **THEN** the completion directive suppresses default file-name completion
