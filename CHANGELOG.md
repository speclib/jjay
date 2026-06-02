# Changelog

## Unreleased

### Added

- **`jjay cleanup <change>`** — tear down spawned workspace: kill tmux window, forget jj workspace, remove directory. Tolerant execution skips missing resources. See [proposal](openspec/changes/archive/2026-06-02-cleanup-command/proposal.md).
- **`internal/workspace` package** — shared `WindowName()` and `WorkspaceDir()` helpers extracted from spawn, used by both spawn and cleanup.

### Changed

- **`internal/spawn`** — refactored to use shared `workspace` package (no behavior change).

### Added

- **Release process** — VERSION file as single source of truth, goreleaser for multi-platform builds, GitHub Actions workflow on `v*` tags, interactive release script with gum and nix vendorHash auto-update. See [proposal](openspec/changes/archive/2026-06-02-release-process/proposal.md).
- **ADR-004**: VERSION file as single source of truth. See [ADR](openspec/adrs/004-version-single-source-of-truth.md).

### Changed

- **flake.nix** — reads version from VERSION file, devShell includes goreleaser and gum.
- **Makefile** — injects version from VERSION file via ldflags.

### Decisions

- **Go as implementation language** — chosen for team familiarity, CLI ecosystem (cobra, bubbletea), single-binary distribution, and orchestrator-friendly stdlib. See [proposal](openspec/changes/archive/2026-06-02-techstack-go/proposal.md).

### Changed

- **openspec/config.yaml** — filled in project context (Go, cobra, jj/tmux orchestration) and light per-artifact rules. See [proposal](openspec/changes/archive/2026-06-02-config-yaml/proposal.md).

### Added

- **spec-driven-with-adr schema** — forked `spec-driven`, added persistent ADR artifact that generates to `openspec/adrs/`. Includes superseding convention. See [proposal](openspec/changes/archive/2026-06-02-spec-driven-with-adr/proposal.md).
- **ADR-001**: Use Go as implementation language. See [ADR](openspec/adrs/001-use-go.md).
- **ADR-002**: OpenSpec config — project context and light rules. See [ADR](openspec/adrs/002-openspec-config.md).
- **Project scaffold** — Go module, cobra CLI (`jjay version`), Makefile, Nix flake, initial tests. See [proposal](openspec/changes/archive/2026-06-02-project-scaffold/proposal.md).
- **`jjay spawn <change>`** — create jj workspace, tmux window with two-pane layout, launch claude agent. See [proposal](openspec/changes/archive/2026-06-02-spawn-command/proposal.md).
- **ADR-003**: Spawn orchestration — sequential subprocess, no rollback. See [ADR](openspec/adrs/003-spawn-orchestration-pattern.md).
