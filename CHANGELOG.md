# Changelog

## Unreleased

### Decisions

- **Go as implementation language** — chosen for team familiarity, CLI ecosystem (cobra, bubbletea), single-binary distribution, and orchestrator-friendly stdlib. See [proposal](openspec/changes/archive/2026-06-02-techstack-go/proposal.md).

### Changed

- **openspec/config.yaml** — filled in project context (Go, cobra, jj/tmux orchestration) and light per-artifact rules. See [proposal](openspec/changes/archive/2026-06-02-config-yaml/proposal.md).

### Added

- **spec-driven-with-adr schema** — forked `spec-driven`, added persistent ADR artifact that generates to `openspec/adrs/`. Includes superseding convention. See [proposal](openspec/changes/archive/2026-06-02-spec-driven-with-adr/proposal.md).
- **ADR-001**: Use Go as implementation language. See [ADR](openspec/adrs/001-use-go.md).
- **ADR-002**: OpenSpec config — project context and light rules. See [ADR](openspec/adrs/002-openspec-config.md).
