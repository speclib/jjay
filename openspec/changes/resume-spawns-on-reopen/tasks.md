## 1. Config package (3-layer, per-field)

- [ ] 1.1 New `internal/config` package: types `Config{Agents map[string]AgentProfile}` and `AgentProfile{Launch, Resume string}`.
- [ ] 1.2 Go built-in default: `builtin["claude"] = {Launch: <current spawn.DefaultAgentCommand>, Resume: "claude --resume --add-dir {wsdir}"}`. This is the single source of truth.
- [ ] 1.3 YAML load for a single file path (prefer an existing dep, e.g. `gopkg.in/yaml.v3`, before adding one). Missing file ⇒ empty config, no error.
- [ ] 1.4 Per-field resolver over an **ordered list of sources**: `resolveField(agent, field, sources...)` returning the first non-empty value, falling through to the built-in. Empty string = unset.
- [ ] 1.5 Wire the standard source order: project `<repo>/.jjay/config.yaml` → global `~/.config/jjay/config.yaml` → built-in. Project path located from the main repo root (not the spawned workspace).
- [ ] 1.6 Unit tests: project-overrides-one-field, no-files-uses-builtin, present-but-blank-falls-back, ordered-source insertion.

## 2. init seeds the project config

- [ ] 2.1 `internal/init` writes `<repo>/.jjay/config.yaml` from the built-in `agents` block (serialize the same const used for fallback).
- [ ] 2.2 Idempotent + non-destructive: do not overwrite an existing `.jjay/config.yaml` without `--force` (ADR-008 pattern); report the action like the existing seeded files.
- [ ] 2.3 Confirm jjay never writes its `agents` block into `openspec/config.yaml`.
- [ ] 2.4 Tests: seeds when absent; leaves existing untouched without `--force`.

## 3. Launch ≠ resume intent in spawn

- [ ] 3.1 Introduce an `Intent` (launch | resume) into the window-open path; replace `SpawnOptions.Agent string` usage with a resolved `AgentProfile`.
- [ ] 3.2 `openWindow`/`setupPanes` select `profile.Launch` vs `profile.Resume` by intent. Invert the "cannot diverge" comment to "share setup, diverge on command".
- [ ] 3.3 `Spawn`/`SpawnProposal` resolve the profile via `internal/config` (launch intent). `--agent` flag still overrides the launch command at highest priority.
- [ ] 3.4 `resolveAgentCommand` continues to substitute `{change}`/`{prompt}`/`{wsdir}` for whichever template the intent selected.

## 4. tmux-open command + session-open delegation

- [ ] 4.1 Extract a single-workspace reopen primitive (`OpenWindow` with resume intent) — recreates window/panes, runs the resolved `resume`, does not touch the jj workspace.
- [ ] 4.2 New `jjay tmux-open <workspace>` command in `cmd/jjay/main.go` (current session or `--session`); arg completion over reopenable workspaces.
- [ ] 4.3 `session.reopenDetached` loops the same primitive per detached spawn (resume intent); keep the best-effort, non-fatal, skip-attached behavior.
- [ ] 4.4 Tests: `reopenDetached` uses resume not launch; skips attached; per-workspace failure non-fatal; `tmux-open` reopens one window in the workspace dir.

## 5. Docs & assets

- [ ] 5.1 Add `internal/init/assets/commands/jjay/tmux-open.md` (thin wrapper over `jjay tmux-open <workspace>`).
- [ ] 5.2 Update `session-open.md` to describe resume-on-reopen (no longer re-runs `/opsx:apply`).
- [ ] 5.3 Update `SKILL.md` lifecycle/vocab if it references reopen behavior.
- [ ] 5.4 Document `.jjay/config.yaml` (and the global path) in README/AGENTS.md as jjay's own config, distinct from `openspec/config.yaml`.

## 6. Manual verification

- [ ] 6.1 **Verify `--resume` scoping rides on cwd.** Spawn an apply workspace, let Claude start, kill the tmux session, run `jjay session-open`; confirm the reopened window's `claude --resume` picker is scoped to *that* workspace's sessions (window opened with `-c wsDir`). Record the result.
- [ ] 6.2 Confirm a workspace with no prior session still opens cleanly (empty picker / Claude's own message), never silently re-running apply.
- [ ] 6.3 `session-open` with multiple detached spawns opens one window per spawn, each on its resume command.

## 7. Testing — full suite green

- [ ] 7.1 `go test ./...` green (unit: config resolver, init seeding, spawn intent, reopen path).
- [ ] 7.2 `go test -tags integration ./...` green (reopen/session lifecycle unaffected; `--agent` override still works).

## 8. Beans & changelog

- [ ] 8.1 jjay-tzl3 → `in-progress`, linked to this change. (Done at proposal creation.)
- [ ] 8.2 Note for jjay-iex3 / jjay-euup: config foundation now exists; they extend it (more fields), not build it.
- [ ] 8.3 CHANGELOG updated on archive (headline ≤ 80 chars; detail on sub-bullets).
