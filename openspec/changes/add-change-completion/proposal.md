## Why

Typing change names by hand for `jjay spawn`/`merge`/`cleanup` is error-prone — the names are long kebab strings (`workspace-aware-session`), and a typo fails the command's precondition check. The CLI is built on cobra, which already ships the completion *scaffolding* (`jjay completion {bash,zsh,fish,powershell}` and the `__complete` dispatch), but no candidates are wired in: `jjay spawn <TAB>` currently falls back to filename completion. Nothing feeds cobra the live openspec/workspace names.

Bean [jjay-nd10](../../.beans/jjay-nd10--read-openspec-changes-to-autocomplete-after-spawn.md) ("read openspec changes to autocomplete after spawn etc", parent [jjay-qltp](../../.beans/jjay-qltp--version-03.md)) asks for exactly this. "etc" = the other verbs for which a change-name argument is relevant: `merge` and `cleanup`.

## What Changes

- **Change-name completion** for the positional argument of `jjay spawn`, `jjay merge`, and `jjay cleanup`, with **per-verb filtering** so each suggests only valid candidates:
  - `spawn <TAB>` → openspec changes that do **not** yet have a workspace (you can't spawn an already-spawned change).
  - `merge <TAB>` → existing spawned workspaces (only those are mergeable).
  - `cleanup <TAB>` → existing spawned workspaces (only those can be torn down).
- **A dedicated `internal/completion` package** holding the candidate logic, depending only on data-source packages (never on the command packages). Three named functions — `Spawnable`, `Mergeable`, `Cleanable` — kept distinct even where two coincide today, so the CLI can evolve without spaghetti.
- **Lean name-only readers** so a TAB press stays fast: `internal/openspec.ChangeNames()` (lifted from `spawn.go`'s existing parsing, which then reuses it) and a `WorkspaceNames()` reader (jj workspace list only — no tmux/tasks probing).

## Capabilities

### New Capabilities
- `completion`: change-name shell completion for spawn/merge/cleanup, filtered per verb by workspace state.

### Modified Capabilities
<!-- None at the requirement level. spawn/merge/cleanup behavior is unchanged; only completion candidates are added. -->

## Impact

- **GATED ON `workspace-aware-session` MERGING FIRST.** Completion's `Mergeable`/`Cleanable` (and the `spawn` set-minus) need the spawned-workspace list. The lean `WorkspaceNames()` reader lives in `internal/status`, which currently exists only in the unmerged `workspace-aware-session` workspace. Until that merges to main, `internal/status` is not present and this change cannot build. See `design.md` and ADR-009.
- **Code**: new `internal/completion/`; new `internal/openspec/` (lifted from `spawn.go` — `spawn.go` reuses it, no behavior change); new `WorkspaceNames()` in `internal/status`; `ValidArgsFunction` bindings in `cmd/jjay/main.go`.
- **No new dependencies, no shell scripts** — uses cobra's built-in dynamic completion; the existing `jjay completion <shell>` script is the install path (out of scope here).
- **Out of scope** (deliberately, no bean): flag-value completion (`--session`, `--workspace-root`); generating/installing completion scripts (cobra + a future init concern).
- **ADRs**: ADR-009 (completion package depends on data sources, not commands; lean readers; named-apart verb funcs).
- **Beans**: nd10 → in-progress, linked here.
