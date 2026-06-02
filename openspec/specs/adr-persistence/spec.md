### Requirement: ADR artifact in schema
The `spec-driven-with-adr` schema SHALL include an `adr` artifact that generates files at `openspec/adrs/<number>-<slug>.md`. The artifact SHALL sit alongside `design` in the dependency graph (requires `proposal`, unlocks `tasks`).

#### Scenario: Schema includes ADR artifact
- **WHEN** a change is created using the `spec-driven-with-adr` schema
- **THEN** the artifact graph includes `adr` alongside `proposal`, `specs`, `design`, and `tasks`

#### Scenario: ADR artifact depends on proposal
- **WHEN** `openspec status` is checked for a new change
- **THEN** the `adr` artifact shows `proposal` as a dependency

### Requirement: ADR template structure
Each ADR file SHALL follow a standard template with sections: Title, Status, Context, Options Considered, Decision, and Consequences.

#### Scenario: ADR file follows template
- **WHEN** an ADR is created for a change
- **THEN** the file contains all required sections in order

### Requirement: ADRs persist outside change lifecycle
ADR files SHALL be written to `openspec/adrs/` which is outside the change directory. ADRs SHALL NOT be moved or deleted when a change is archived.

#### Scenario: ADR survives archival
- **WHEN** a change with an ADR is archived via `openspec archive`
- **THEN** the ADR file remains at `openspec/adrs/<number>-<slug>.md`
- **THEN** the archived change directory does NOT contain a copy of the ADR

### Requirement: ADR numbering
ADRs SHALL be numbered sequentially with zero-padded three-digit prefixes (e.g., `001`, `002`). The next number SHALL be determined by scanning existing files in `openspec/adrs/`.

#### Scenario: First ADR in empty project
- **WHEN** no ADRs exist in `openspec/adrs/`
- **THEN** the new ADR is numbered `001`

#### Scenario: Subsequent ADR
- **WHEN** `openspec/adrs/` contains `001-use-go.md` and `002-openspec-config.md`
- **THEN** the next ADR is numbered `003`

### Requirement: Schema config switch
The project `openspec/config.yaml` SHALL use `schema: spec-driven-with-adr` after this change is applied.

#### Scenario: Config references new schema
- **WHEN** the config is read after implementation
- **THEN** the `schema` field is `spec-driven-with-adr`
