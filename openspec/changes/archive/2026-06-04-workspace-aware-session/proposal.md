## Why

`jjay spawn` creates jj workspaces and tmux windows, but nothing reports *what is currently running*. There is no `jjay status`. Worse, the tmux view is volatile: closing the tmux server, detaching, or rebooting destroys the `ws-*` windows and the agents inside them — but the jj workspaces survive on disk. After `jjay session-open`, the user lands in a fresh session with no windows, even though several spawns are still "open" as far as jj is concerned.

Two beans capture this gap (both under epic [jjay-qltp](../../.beans/jjay-qltp--version-03.md)):
- [jjay-7spt](../../.beans/jjay-7spt--jjay-status-command.md) — `jjay status` to list running agent workspaces.
- [jjay-cedd](../../.beans/jjay-cedd--reopen-existing-spawns-after-jjay-session-open.md) — reopen existing spawns after `session-open`.

They are two halves of one idea — *report the diff between workspaces and windows* (status) and *repair that diff* (reopen) — and they share the same join logic. An earlier, vaguer bean [jjay-72pw](../../.beans/jjay-72pw--state-tracking-for-running-workspaces.md) ("state tracking for running workspaces") is fully superseded by this change and is being scrapped; its real content was the *principle* that there is no state file, which this change records as ADR-006.

## What Changes

- **New `jjay status` command** — lists every jj workspace and, for each, whether a matching `ws-<change>` tmux window currently exists. A workspace with no window is reported as **detached** (still open, just not attached). Read-only; derives everything live from `jj workspace list` + `tmux list-windows`, persisting nothing.
- **Reopen on `session-open`** — after creating/switching to the session, `session-open` recreates a `ws-<change>` window (and relaunches the agent) for every jj workspace that lacks one, restoring the tmux view to match the workspace truth.
- **ADR-006** records the founding invariant: *the jj workspace is the single source of truth; a spawn is "open" iff its workspace exists; tmux windows and the agent are recreatable views.*

## Capabilities

### New Capabilities
- `status`: the `jjay status` command — derive and report running spawns from jj + tmux, with no state file.

### Modified Capabilities
- `session-open`: after opening the session, reattach detached spawns by recreating their tmux windows.

## Impact

- **Code**: new `internal/status/` package and `status` command in `cmd/jjay/main.go`; extend `internal/session/` (or a new reopen helper) to recreate windows on `session-open`. Reuses `workspace.WindowName` and the `jj workspace list` / `tmux list-windows` parsing patterns already in `internal/spawn/spawn.go`.
- **Tools**: `jj`, `tmux` (read + window creation), and the agent command template (shared with spawn).
- **Specs**: new `status` capability; delta to existing `session-open` capability.
- **ADRs**: ADR-006 (workspace = single source of truth).
- **Beans**: 72pw → scrapped; 7spt + cedd → in-progress, linked here.
