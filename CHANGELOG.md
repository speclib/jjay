# Changelog

## Unreleased

### Added

- **`jjay init [path]`** — idempotent, non-destructive project bootstrap that prepares a target project (default: cwd) for orchestration by jjay. It delegates openspec scaffolding to `openspec init <path> --tools claude` (never reimplementing it), installs the `/jjay:*` commands and the `jjay` skill into the target's `.claude/` from embedded copies of this repo's own `.claude/` (a drift test keeps them byte-identical), and writes an `AGENTS.md` of jjay conventions (openspec archive flow, beans tasks, jj usage). jj init (`--with-jj`) and example hooks (`--with-hooks`) are opt-in. Each step detects whether its artifact already exists and skips it; `--yes` accepts creation defaults but only `--force` authorizes overwriting an existing file. Re-running is a no-op on a prepared project and completes only missing steps on a partial one. New `internal/init` package and `initCmd`. See [proposal](openspec/changes/archive/2026-06-05-add-init-command/proposal.md) and [ADR-008](openspec/changes/archive/2026-06-05-add-init-command/adrs/008-init-idempotent-orchestrator.md).
- **Claude Code integration layer** — `/jjay:*` slash commands (`spawn`, `status`, `merge`, `cleanup`, `session-open`) as thin wrappers over the `jjay` binary, plus a `jjay` orchestrator skill (`.claude/skills/jjay/SKILL.md`) that steers Claude to implement changes by spawning an isolated workspace (`/jjay:spawn`) rather than running `/opsx:apply` in the main session. The skill documents the lifecycle (explore → propose → spawn → status → merge → cleanup) and the orchestrator-vs-worker split (a worker applies in place and does not recursively spawn). Committed under `.claude/`, so it propagates into every spawned workspace. Requires the `jjay` binary on `PATH`. See [proposal](openspec/changes/archive/2026-06-04-add-claude-commands/proposal.md) and [ADR-007](openspec/changes/archive/2026-06-04-add-claude-commands/adrs/007-commands-thin-skill-policy.md).
- **Shell completion for change names** — the positional argument of `spawn`, `merge`, and `cleanup` now tab-completes, filtered per verb: `spawn` offers openspec changes **without** a workspace (set-minus of `openspec list` and `jj workspace list`), while `merge`/`cleanup` offer existing spawned workspaces (`default` excluded). Completion reads only `openspec list` + `jj workspace list` — no tmux, no task files — and degrades silently to no candidates on any read error. Wired via cobra `ValidArgsFunction`; install with `jjay completion <shell>`. New `internal/completion`, `internal/openspec` (lifted from `spawn.go`, which now reuses it), and `internal/status.WorkspaceNames()`. See [proposal](openspec/changes/add-change-completion/proposal.md) and [ADR-009](openspec/adrs/009-completion-package-depends-on-data-sources.md).
- **`jjay status`** — lists every spawned jj workspace with task progress (`done/total (pct%)` read from each workspace's `openspec/changes/<change>/tasks.md`) and reports each as **attached** (a `ws-<change>` tmux window exists in the current session) or **detached** (workspace still open, no window). Workspace paths are shown **relative to the main repo root**, resolved correctly even when run from inside a child workspace (via the `.jj/repo` pointer). Read-only; derived live from `jj workspace list` + `tmux list-windows` with no state file, and tolerant of a missing tmux server (all detached). See [proposal](openspec/changes/archive/2026-06-04-workspace-aware-session/proposal.md).
- **Reopen on `session-open`** — after switching to the session, `jjay session-open` recreates a `ws-<change>` window and relaunches the agent for every spawned workspace that lacks one, restoring the tmux view to match the workspaces on disk. Best-effort: a per-spawn failure is reported and skipped without aborting session-open.
- **ADR-006** — records the founding invariant that the jj workspace is the single source of truth (a spawn is "open" iff its workspace exists; tmux windows and the agent are recreatable views), so `status` and reopen share one definition.

### Changed

- **Archive blog gate** — `/opsx:archive` now auto-creates a missing `blog` artifact before archiving (via the "OpenSpec Archive triggers" section in `CLAUDE.md`), instead of only warning. The blog is written retrospectively from `proposal.md` and completed `tasks.md`; other incomplete artifacts (e.g. `adr`) stay warn-only. See [proposal](openspec/changes/archive/2026-06-04-fix-devlog-archive-gate/proposal.md).

### Fixed

- **`jjay merge` silently dropped files** — merge now rebases the workspace onto current main (`jj rebase -b <change>@ -d main`) before creating the merge commit, eliminating jj's 3-way silent file picking that lost task progress, blog posts, and beans. Conflicts are surfaced explicitly and abort the merge. Adds 6 e2e merge scenarios. See [proposal](openspec/changes/archive/2026-06-04-rebase-before-merge/proposal.md) and [ADR-007](openspec/adrs/007-rebase-before-merge.md).
- **`jjay merge` silently dropped main-side work created after a spawn** — work committed in the main working copy *ahead of* the `main` bookmark (new proposals, bean edits) was excluded from the merge and orphaned when the bookmark advanced. Merge now folds ahead-of-bookmark main work into `main` (`latest(main..@ & ~empty())`) before merging, so it survives. Because jj auto-snapshots the working copy, uncommitted main edits are preserved too. Adds the `TestMerge_MainAddsNewFiles` regression scenario. See ADR-010.

## 0.2.1 - 2026-06-03

### Added

- **`make coverage` target** — runs unit tests with `-coverprofile`, generates HTML report, prints coverage percentage. See [proposal](openspec/changes/archive/2026-06-03-add-test-coverage/proposal.md).
- **`make badge` target** — extracts coverage percentage and patches README with a shields.io badge (green ≥80%, yellow ≥60%, red <60%).
- **README coverage badge** — shields.io badge between hero image and heading.

## 0.2.0 - 2026-06-03

### Fixed

- **tmux pane working directories** — spawn now uses `tmux new-window -c` and `split-window -c` to set pane working directories at creation time, replacing the racy `send-keys cd` approach. Both panes reliably start in the workspace directory. See [proposal](openspec/changes/archive/2026-06-03-fix-pane-dirs/proposal.md).
- **Integration test: pane directory assertion** — added `assertPaneDir` to verify both panes report the correct working directory via `tmux display-message #{pane_current_path}`.

### Added

- **Devblog with Kaa persona** — blog artifact in schema, retroactive posts for 2026-06-02 work. Kaa (Eurasian Jay in jujutsu gi) narrates in first person. See [proposal](openspec/changes/archive/2026-06-03-add-devblog/proposal.md).
- **`jjay session-open <path>`** — create and switch to a dedicated tmux session (`jjay-><dirname>`) for a jj repo. Validates jj repo, prevents duplicate sessions. See [proposal](openspec/changes/archive/2026-06-03-session-open/proposal.md).
- **Configurable spawn** — `--agent`, `--session`, `--workspace-root` flags on spawn; `--session`, `--workspace-root` flags on cleanup. Enables custom agents, dedicated tmux sessions, and flexible workspace locations. See [proposal](openspec/changes/archive/2026-06-03-spawn-config/proposal.md).
- **Integration test** — full spawn → cleanup lifecycle test using fake agent, isolated tmux session, and temp jj repo (`go test -tags integration`). See [proposal](openspec/changes/archive/2026-06-03-spawn-config/proposal.md).
- **ADR-006**: Configuration via CLI flags, not config file. See [ADR](openspec/adrs/006-config-via-flags-not-file.md).
- **`jjay merge <change>`** — merge a workspace's work into main: resolve workspace change via jj revset, create merge commit, move bookmark, create fresh change. See [proposal](openspec/changes/archive/2026-06-03-merge-command/proposal.md).

### Added

- **`jjay cleanup <change>`** — tear down spawned workspace: kill tmux window, forget jj workspace, remove directory. Tolerant execution skips missing resources. See [proposal](openspec/changes/archive/2026-06-02-cleanup-command/proposal.md).
- **`internal/workspace` package** — shared `WindowName()` and `WorkspaceDir()` helpers extracted from spawn, used by both spawn and cleanup.

### Changed

- **`internal/spawn`** — refactored to use shared `workspace` package (no behavior change).

### Fixed

- **CRITICAL: workspace isolation** — spawn now runs `jj new` before creating child workspace, then uses `--revision @-`. Prevents data loss when main workspace becomes stale during concurrent agent work. See [proposal](openspec/changes/archive/2026-06-02-fix-workspace-isolation/proposal.md).
- **ADR-005**: Workspace isolation via jj new snapshot. See [ADR](openspec/adrs/005-workspace-isolation-via-snapshot.md).

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
