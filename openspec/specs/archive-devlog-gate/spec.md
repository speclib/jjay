## Purpose

Define how `/opsx:archive` gates incomplete artifacts: the blog artifact is auto-created at archive time, while other incomplete artifacts (e.g. adr) remain warn-only with user confirmation.

## Requirements

### Requirement: Archive creates missing blog artifact
The `/opsx:archive` skill SHALL create the blog artifact automatically if it is not `done` at archive time, using `openspec instructions blog` for template and context.

#### Scenario: Blog missing at archive time
- **WHEN** `/opsx:archive` is invoked for a change
- **THEN** `openspec status --change "<name>" --json` is checked
- **THEN** if the `blog` artifact status is not `done`, it is created before proceeding

#### Scenario: Blog already exists
- **WHEN** `/opsx:archive` is invoked for a change that already has a completed blog artifact
- **THEN** the blog creation step is skipped
- **THEN** archive proceeds normally

### Requirement: Blog content reflects implementation
The auto-created blog post SHALL be written from a retrospective perspective, referencing what was actually built (from tasks.md and proposal.md), not just what was planned.

#### Scenario: Blog references completed work
- **WHEN** the blog artifact is auto-created at archive time
- **THEN** the blog post reads the proposal and completed tasks for context
- **THEN** the narrative reflects the actual outcome of the change

### Requirement: Other incomplete artifacts remain soft-gated
Incomplete artifacts other than `blog` (e.g., `adr`) SHALL continue to produce a warning with user confirmation, not automatic creation.

#### Scenario: ADR missing at archive time
- **WHEN** `/opsx:archive` is invoked and the `adr` artifact is not `done`
- **THEN** a warning is displayed listing `adr` as incomplete
- **THEN** the user is prompted to confirm proceeding without it
