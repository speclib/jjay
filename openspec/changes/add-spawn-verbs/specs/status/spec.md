## ADDED Requirements

### Requirement: Status separates change spawns from proposal spawns
`jjay status` SHALL group spawns by kind, derived from the name prefix: apply spawns (`app-*`, which track an openspec change) and proposal spawns (`prop-*`, prompt-seeded with no change yet) SHALL be shown in two distinct tables. Change-only columns (MERGED, ARCHIVED, TASKS) apply to the CHANGES table; proposal spawns SHALL NOT be forced into change-shaped columns that are meaningless for them.

#### Scenario: Two tables shown
- **WHEN** both an `app-add-foo` and a `prop-dark-mode` spawn exist and `jjay status` runs
- **THEN** `app-add-foo` appears under a CHANGES table
- **THEN** `prop-dark-mode` appears under a PROPOSAL SPAWNS table

#### Scenario: Only one kind present
- **WHEN** only proposal spawns exist
- **THEN** the PROPOSAL SPAWNS table is shown and the CHANGES table is empty or omitted
