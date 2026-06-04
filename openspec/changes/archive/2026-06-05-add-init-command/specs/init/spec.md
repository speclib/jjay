## ADDED Requirements

### Requirement: Init prepares a target project
The `jjay init [path]` command SHALL prepare the project at `path` (default: current directory) for orchestration by jjay. It SHALL initialize openspec, install the jjay Claude integration, and write an `AGENTS.md`, with jj and hooks as optional steps.

#### Scenario: Bare project initialized
- **WHEN** `jjay init` is run in a project with none of the jjay scaffolding
- **THEN** openspec is initialized for the project
- **THEN** the `/jjay:*` commands and the `jjay` skill are installed under the project's `.claude/`
- **THEN** an `AGENTS.md` documenting jjay conventions exists
- **THEN** jjay exits zero

### Requirement: Init is idempotent
Re-running `jjay init` on an already-prepared project SHALL be a no-op for steps whose artifacts already exist and are valid, and SHALL complete steps that are missing. It SHALL NOT fail merely because some artifacts already exist.

#### Scenario: Re-run on prepared project
- **WHEN** `jjay init` is run a second time on a project it already prepared
- **THEN** existing artifacts are left unchanged
- **THEN** jjay reports the project is already initialized (or completes any newly-added step) and exits zero

#### Scenario: Partially-initialized project completed
- **WHEN** a project has openspec but no `.claude/commands/jjay/`, and `jjay init` is run
- **THEN** the missing `/jjay:*` commands and skill are installed
- **THEN** the existing openspec setup is left unchanged

### Requirement: Init is non-destructive
Init SHALL NOT overwrite an existing user file (e.g. `config.yaml`, `AGENTS.md`, a customized command) without explicit confirmation. The `--yes` flag SHALL accept creation defaults but SHALL NOT authorize overwriting existing files; overwriting SHALL require `--force`.

#### Scenario: Existing AGENTS.md preserved
- **WHEN** the project already has an `AGENTS.md` and `jjay init --yes` is run
- **THEN** the existing `AGENTS.md` is not overwritten
- **THEN** init reports it was left in place

#### Scenario: Force overwrites
- **WHEN** `jjay init --force` is run and an `AGENTS.md` exists
- **THEN** the file is overwritten with the jjay template

### Requirement: Init delegates to canonical initializers
Init SHALL initialize openspec by invoking `openspec init` (selecting the `claude` tool and the project schema) rather than reimplementing openspec scaffolding, and SHALL initialize jj via jj's own commands.

#### Scenario: openspec initialized via openspec
- **WHEN** `jjay init` initializes openspec
- **THEN** it invokes `openspec init` for the target project
- **THEN** the resulting `openspec/` directory is configured for the `claude` tool

### Requirement: Optional steps are opt-in
jj initialization and hook scaffolding SHALL be optional and off by default, enabled via flags (e.g. `--with-jj`). Steps SHALL be individually skippable for non-interactive runs.

#### Scenario: jj initialized only when requested
- **WHEN** `jjay init` is run without `--with-jj` in a non-jj directory
- **THEN** no jj repository is created

#### Scenario: jj initialized when requested
- **WHEN** `jjay init --with-jj` is run in a non-jj directory
- **THEN** a jj repository is initialized

### Requirement: Init supports non-interactive use
The command SHALL provide flags to run without prompts (`--yes` and per-step flags), mirroring `openspec init`'s non-interactive options.

#### Scenario: Non-interactive run
- **WHEN** `jjay init --yes` is run
- **THEN** init proceeds using defaults without interactive prompts
- **THEN** jjay exits zero
