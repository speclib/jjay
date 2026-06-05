## ADDED Requirements

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
