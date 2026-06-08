# Changelog

## Unreleased

### Added

- **`make clean-tests`** — sweep leaked integration-test debris
  - kills `jjay-test-*` tmux sessions and removes `/tmp/jjay-test-*` + `/tmp/jjay-merge-test-*`; `test-integration` runs it first. Real `jjay->` sessions are never touched. See [proposal](openspec/changes/archive/2026-06-05-add-clean-tests-target/proposal.md).

## 0.3.0 - 2026-06-05

### Added

- **`jjay spawn apply` / `jjay spawn proposal`** — two spawn verbs
  - `spawn apply <change>` is the old spawn, now named `app-<change>`; `spawn proposal <prompt> [--mode explore|propose]` seeds new work in `prop-<slug>` with a code-derived slug (no AI), freeing the orchestrator for conflict work.
  - No bare `spawn <change>` form — `/jjay:spawn` and the spawn integration tests migrate to `spawn apply`. Pre-existing unprefixed workspaces (none expected, single pre-1.0 user) are still read by `status` via a legacy fallback.
  - `status` splits CHANGES (`app-*`) from PROPOSAL SPAWNS (`prop-*`); a proposal's workspace name is never assumed equal to its produced change. See [proposal](openspec/changes/archive/add-spawn-verbs/proposal.md), [ADR-011](openspec/changes/archive/add-spawn-verbs/adrs/011-spawn-verbs-and-slug-identity.md).
- **`jjay init [path]`** — idempotent, non-destructive project bootstrap
  - openspec (via `openspec init --tools claude`), `/jjay:*` commands + skill into `.claude/`, `AGENTS.md`; `--with-jj`/`--with-hooks` opt-in; `--force` to overwrite. New `internal/init`. See [proposal](openspec/changes/archive/2026-06-05-add-init-command/proposal.md), [ADR-008](openspec/changes/archive/2026-06-05-add-init-command/adrs/008-init-idempotent-orchestrator.md).
- **Unit tests for `internal/cleanup`** — the last untested command package
  - covers `removeDirectory`, `tmuxTarget`, and the tolerance branches; 0% → 59.3%. Test-only. See [proposal](openspec/changes/archive/2026-06-05-add-cleanup-unit-tests/proposal.md).
- **Claude Code integration layer** — `/jjay:*` commands + `jjay` skill
  - thin wrappers over the binary; skill steers toward `/jjay:spawn` over `/opsx:apply`, documents lifecycle and orchestrator-vs-worker split. See [proposal](openspec/changes/archive/2026-06-04-add-claude-commands/proposal.md), [ADR-007](openspec/changes/archive/2026-06-04-add-claude-commands/adrs/007-commands-thin-skill-policy.md).
- **Shell completion for change names** — `spawn`/`merge`/`cleanup`, per-verb
  - `spawn` = changes without a workspace; `merge`/`cleanup` = existing workspaces. Reads `openspec list` + `jj workspace list` only. New `internal/completion`, `internal/openspec`, `status.WorkspaceNames()`. See [proposal](openspec/changes/archive/2026-06-05-add-change-completion/proposal.md), [ADR-009](openspec/changes/archive/2026-06-05-add-change-completion/adrs/009-completion-package-depends-on-data-sources.md).
- **`jjay status`** — lists spawned workspaces, task progress, attached/detached
  - paths relative to the main repo root; read-only, derived from `jj workspace list` + `tmux list-windows`; tolerant of a missing tmux server. See [proposal](openspec/changes/archive/2026-06-04-workspace-aware-session/proposal.md).
- **Reopen on `session-open`** — recreates `ws-<change>` windows for spawns
  - best-effort: a per-spawn failure is reported and skipped without aborting.
- **ADR-006** — the jj workspace is the single source of truth
  - a spawn is "open" iff its workspace exists; tmux/agent are recreatable views.

### Changed

- **`jjay status` gains a `MERGED` column; `STATUS` renamed to `TMUX`**
  - **MERGED** (yes/no) flags whether a spawn's work is already on `main`, derived live from jj (`main..<change>@` empty ⇒ merged); a merged-but-not-archived spawn is the "ready to clean up" signal. The tmux-state column is now **TMUX**, freeing "status" for a future agent-state column. New column order `CHANGE WORKSPACE TASKS TMUX MERGED ARCHIVED`; folds `jjay status` into the lifecycle integration test. See [proposal](openspec/changes/add-status-merged-column/proposal.md).
- **Archive blog gate** — `/opsx:archive` auto-creates a missing `blog` artifact
  - written retrospectively; other incomplete artifacts (e.g. `adr`) stay warn-only. See [proposal](openspec/changes/archive/2026-06-04-fix-devlog-archive-gate/proposal.md).

### Fixed

- **`jjay merge` silently dropped files** — rebase workspace onto main first
  - eliminates jj's 3-way silent file picking; conflicts abort explicitly. Adds 6 e2e scenarios. See [proposal](openspec/changes/archive/2026-06-04-rebase-before-merge/proposal.md), [ADR-007](openspec/adrs/007-rebase-before-merge.md).
- **`jjay merge` dropped post-spawn main work** — fold ahead-of-bookmark work
  - merge folds `latest(main..@ & ~empty())` into main first; jj auto-snapshots preserve uncommitted edits. Adds `TestMerge_MainAddsNewFiles`. See ADR-010.

## 0.2.1 - 2026-06-03

### Added

- **`make coverage` target** — unit tests with coverprofile + HTML + percentage
  - See [proposal](openspec/changes/archive/2026-06-03-add-test-coverage/proposal.md).
- **`make badge` target** — patches README with a shields.io coverage badge
- **README coverage badge** — between hero image and heading

## 0.2.0 - 2026-06-03

### Fixed

- **tmux pane working directories** — spawn sets pane dirs at creation time
  - uses `new-window -c` / `split-window -c`, replacing the racy `send-keys cd`; both panes start in the workspace dir. See [proposal](openspec/changes/archive/2026-06-03-fix-pane-dirs/proposal.md).
- **Pane directory assertion** — `assertPaneDir` checks both panes in the test

### Added

- **Devblog with Kaa persona** — blog artifact in schema, retroactive posts
  - Kaa (Eurasian Jay in jujutsu gi) narrates in first person. See [proposal](openspec/changes/archive/2026-06-03-add-devblog/proposal.md).
- **`jjay session-open <path>`** — create and switch to a tmux session
  - validates jj repo, prevents duplicate sessions. See [proposal](openspec/changes/archive/2026-06-03-session-open/proposal.md).
- **Configurable spawn** — `--agent`, `--session`, `--workspace-root` flags
  - plus `--session`/`--workspace-root` on cleanup. See [proposal](openspec/changes/archive/2026-06-03-spawn-config/proposal.md).
- **Integration test** — full spawn → cleanup lifecycle with fake agent
  - isolated tmux session + temp jj repo (`go test -tags integration`).
- **ADR-006** — configuration via CLI flags, not a config file
  - See [ADR](openspec/adrs/006-config-via-flags-not-file.md).
- **`jjay merge <change>`** — merge a workspace's work into main
  - resolve via jj revset, create merge commit, move bookmark, fresh change. See [proposal](openspec/changes/archive/2026-06-03-merge-command/proposal.md).
- **`jjay cleanup <change>`** — tear down workspace, window, and directory
  - tolerant execution skips missing resources. See [proposal](openspec/changes/archive/2026-06-02-cleanup-command/proposal.md).
- **`internal/workspace` package** — shared `WindowName()` / `WorkspaceDir()`

### Changed

- **`internal/spawn`** — refactored to use shared `workspace` package

### Fixed

- **CRITICAL: workspace isolation** — spawn runs `jj new` before child workspace
  - then uses `--revision @-`; prevents data loss when main goes stale. See [proposal](openspec/changes/archive/2026-06-02-fix-workspace-isolation/proposal.md).
- **ADR-005** — workspace isolation via jj new snapshot
  - See [ADR](openspec/adrs/005-workspace-isolation-via-snapshot.md).

### Added

- **Release process** — VERSION single source of truth, goreleaser, GH Actions
  - interactive release script with gum + nix vendorHash auto-update. See [proposal](openspec/changes/archive/2026-06-02-release-process/proposal.md).
- **ADR-004** — VERSION file as single source of truth
  - See [ADR](openspec/adrs/004-version-single-source-of-truth.md).

### Changed

- **flake.nix** — reads version from VERSION file; devShell adds goreleaser, gum
- **Makefile** — injects version from VERSION file via ldflags

### Decisions

- **Go as implementation language** — cobra/bubbletea, single binary
  - See [proposal](openspec/changes/archive/2026-06-02-techstack-go/proposal.md).

### Changed

- **openspec/config.yaml** — project context + light per-artifact rules
  - See [proposal](openspec/changes/archive/2026-06-02-config-yaml/proposal.md).

### Added

- **spec-driven-with-adr schema** — forked `spec-driven` + persistent ADR
  - generates to `openspec/adrs/`, includes superseding convention. See [proposal](openspec/changes/archive/2026-06-02-spec-driven-with-adr/proposal.md).
- **ADR-001** — use Go as implementation language
  - See [ADR](openspec/adrs/001-use-go.md).
- **ADR-002** — OpenSpec config: project context and light rules
  - See [ADR](openspec/adrs/002-openspec-config.md).
- **Project scaffold** — Go module, cobra CLI, Makefile, Nix flake
  - See [proposal](openspec/changes/archive/2026-06-02-project-scaffold/proposal.md).
- **`jjay spawn <change>`** — create jj workspace, two-pane tmux window, agent
  - See [proposal](openspec/changes/archive/2026-06-02-spawn-command/proposal.md).
- **ADR-003** — spawn orchestration: sequential subprocess, no rollback
  - See [ADR](openspec/adrs/003-spawn-orchestration-pattern.md).
