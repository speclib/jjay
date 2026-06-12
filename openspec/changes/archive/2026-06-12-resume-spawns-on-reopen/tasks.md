## 1. Config package (3-layer, per-field)

- [x] 1.1 New `internal/config` package: types `Config{Agents map[string]AgentProfile}` and `AgentProfile{Launch, Resume string}`.
- [x] 1.2 Go built-in default: `builtin["claude"] = {Launch: <current spawn.DefaultAgentCommand>, Resume: "claude --resume --add-dir {wsdir}"}`. Single source of truth (`config.Builtin()`).
- [x] 1.3 YAML load via `go.yaml.in/yaml/v3` (already vendored transitively; promoted to a direct dep). Missing file ⇒ empty config, no error.
- [x] 1.4 Per-field resolver over an **ordered list of sources** (`resolveField` / `ResolveProfile`): first non-empty, empty = unset, built-in as final fallback.
- [x] 1.5 Standard order wired in `Resolve(agent, repoRoot)`: project `<repo>/.jjay/config.yaml` → global `~/.config/jjay/config.yaml` → built-in. Project path from the main repo root.
- [x] 1.6 Unit tests: no-files-uses-builtin, project-overrides-one-field, ordered-source precedence, present-but-blank-falls-back, missing-file, YAML parse, Builtin-is-copy.

## 2. init seeds the project config

- [x] 2.1 `internal/init` `stepJjayConfig` writes `<repo>/.jjay/config.yaml` by marshaling `config.Builtin()` (same const as the runtime fallback).
- [x] 2.2 Idempotent + non-destructive via the existing `writeFile` (ADR-008): not overwritten without `--force`; reports created/skipped.
- [x] 2.3 jjay writes its `agents` block only to `.jjay/config.yaml`, never `openspec/config.yaml` (separate step, separate path).
- [x] 2.4 Tests: `TestInit_BareProject` asserts the seeded `.jjay/config.yaml` contains launch+resume+`--resume`; `TestInit_ReRun`/`YesDoesNotClobber` cover non-destructive.

## 3. Launch ≠ resume intent in spawn

- [x] 3.1 Added `Intent` (IntentLaunch | IntentResume) to the window-open path; `OpenWindow` takes the intent.
- [x] 3.2 `OpenWindow` selects `resolveLaunch()` vs `resolveResume()` by intent; the "cannot diverge" comment inverted to "share setup, diverge on command".
- [x] 3.3 `Spawn` resolves launch via `internal/config` (`resolveLaunch`); `--agent` flag still overrides at highest priority. (Proposal explore/propose templates kept as consts — only `claude` apply is config-populated, per scope.)
- [x] 3.4 `resolveAgentCommand` continues to substitute `{change}`/`{prompt}`/`{wsdir}` for the selected template.

## 4. tmux-open command + session-open delegation

- [x] 4.1 `spawn.Reopen(name, wsDir, opts)` = `OpenWindow(..., IntentResume)` — the single reopen primitive (recreates window/panes, runs resume, does not touch the jj workspace).
- [x] 4.2 New `jjay tmux-open <workspace>` command (`spawn.TmuxOpen`: verifies workspace exists + not attached, resolves wsDir, calls Reopen); current session or `--session`; arg completion over existing workspaces.
- [x] 4.3 `session.reopenSpawns` now passes `spawn.Reopen` to `reopenDetached` (resume intent); best-effort/non-fatal/skip-attached behavior unchanged.
- [x] 4.4 Tests: existing `TestReopenDetached_*` cover loop/skip-attached/non-fatal; `TestLaunchResumeDiverge` asserts resume ≠ launch and resume omits `/opsx:apply`.

## 5. Docs & assets

- [x] 5.1 Added `internal/init/assets/commands/jjay/tmux-open.md` (+ mirrored to live `.claude/`; drift test green).
- [x] 5.2 Updated `session-open.md` (both copies) to state reopen resumes, not re-applies.
- [x] 5.3 Updated `SKILL.md` (both copies) reopen line: resumes via `resume`; mentions `tmux-open`.
- [x] 5.4 Documented `.jjay/config.yaml` (+ global path) in README as jjay's own config, distinct from `openspec/config.yaml`.

## 6. Manual verification (HUMAN — interactive, cannot be automated)

- [ ] 6.1 **Verify `--resume` scoping rides on cwd.** Spawn an apply workspace, let Claude start, kill the tmux session, run `jjay session-open`; confirm the reopened window's `claude --resume` picker is scoped to *that* workspace. (Requires a live interactive Claude session — must be run by a human.)
- [ ] 6.2 Confirm a workspace with no prior session opens cleanly (empty picker / Claude's own message), never silently re-running apply.
- [ ] 6.3 `session-open` with multiple detached spawns opens one window per spawn, each on its resume command.

## 7. Testing — full suite green

- [x] 7.1 `go test ./...` green (config resolver, init seeding, spawn intent/divergence, reopen path).
- [x] 7.2 `go test -tags integration ./...` green (spawn/session lifecycle unaffected; build verified the `--agent` override path).

## 7b. session-open bug fixes (surfaced during implementation, bundled here)

- [x] 7b.1 **Dot/colon in session names (jjay-e3bx).** `SessionName` sanitizes `.`/`:` → `_` via `sanitizeTmuxName`, so the create-name and every `-t` target agree (was: `switch-client -t jjay->mip.rs` → "can't find pane: rs"). Regression test in `session_test.go`.
- [x] 7b.2 **Cross-project reopen leak (jjay-02nr).** `session-open <path>` enumerated workspaces from cwd (the caller repo), reopening another project's spawns into the new session. Added `status.ListIn(repoRoot, …)` (queries `jj -R <repoRoot>`, anchors mainRoot there); `session.Open` threads `absPath` → `reopenSpawns` → `ListIn`. Verified `jj -R <target>` lists the target repo's workspaces, not jjay's.

## 8. Beans & changelog

- [x] 8.1 jjay-tzl3 → `in-progress`, linked (done at proposal creation).
- [x] 8.2 Noted on jjay-iex3 / jjay-euup: config foundation now exists; they extend it.
- [x] 8.3 CHANGELOG updated (headline ≤ 80 chars).
