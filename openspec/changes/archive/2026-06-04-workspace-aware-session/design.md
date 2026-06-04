## Context

jjay already creates spawns (`internal/spawn`) and tears them down (`internal/cleanup`), and opens tmux sessions (`internal/session`). What is missing is *visibility into* and *recovery of* the spawns that exist. The tmux view (windows + agents) is volatile; the jj workspaces are durable. This change adds the read path (`jjay status`) and the repair path (reopen on `session-open`), both founded on ADR-006: the jj workspace is the single source of truth, no state file.

Existing building blocks to reuse:
- `workspace.WindowName(change)` → `"ws-" + change` (`internal/workspace/workspace.go`).
- `jj workspace list` parsing — lines of `"<name>: <commit>"`, already parsed in `spawn.go:checkWorkspaceNotExists` (matches `fields[0] == change+":"`).
- `tmux list-windows -F '#{window_name}'` (optionally `-t <session>`) — already used in `spawn.go:checkWindowNotExists`.
- Agent launch — `spawn.DefaultAgentCommand`, `resolveAgentCommand`, `setupPanes` (`spawn.go`).

## Goals / Non-Goals

**Goals:**
- `jjay status`: list spawned jj workspaces and report each as attached/detached relative to the current tmux session. Read-only, no state file.
- `session-open`: after switching to the session, recreate `ws-<change>` windows (and relaunch agents) for workspaces lacking one.
- Extract the workspace⋈window join into one place both features use, so they cannot disagree.

**Non-Goals:**
- Cross-session window discovery (deferred; no bean per user direction).
- A TUI for status ([jjay-bc1c](../../.beans/jjay-bc1c--jjay-tui-mode.md) is separate).
- Per-spawn metadata that has no jj/tmux home (e.g. exact original agent command) — reopen re-derives the agent command from the shared default/template.
- Reporting agent liveness inside the pane (process-level health) — "attached" means a window exists, not that the agent is still running.

## Decisions

- **Single source of truth = jj workspace (ADR-006).** Both `status` and reopen enumerate `jj workspace list` first, then check tmux. No persisted state.
- **One shared join.** A small helper (e.g. `internal/status`) returns, per spawned workspace, `{change, wsDir, attached bool}` by intersecting `jj workspace list` with `tmux list-windows` (current session). `status` renders it; reopen acts on the `!attached` subset. Reuse `workspace.WindowName` for the `ws-` mapping rather than re-deriving the prefix.
- **Identifying spawned vs main workspace.** The main working copy appears in `jj workspace list` too. Distinguish spawns by name convention — spawned workspaces are named after their change (the `--name <change>` passed in `spawn.go:createWorkspace`), and a corresponding workspace directory under the workspace root resolves via `workspace.WorkspaceDir`. The default workspace name is `default`; treat `default` as the main copy and exclude it (or mark it distinctly).
- **Reopen reuses spawn's window+agent setup.** Recreating a window for a detached spawn should produce the same window/pane/agent layout as the original spawn — factor the window-creation + `setupPanes` logic so both spawn and reopen call it, rather than duplicating `tmux new-window`/`send-keys`.
- **Reopen is best-effort and non-fatal.** A failure to reopen one spawn logs and continues (consistent with ADR-003's no-rollback stance); `session-open` still succeeds.
- **tmux-absent tolerance.** `status` treats "no tmux server" as "all detached" (mirrors `session.go:checkSessionNotExists`, which treats a tmux error as "no sessions").
- **Paths anchored on the main repo root.** All workspace-dir resolution and display anchor on the main working copy root, not the current directory, so `jjay status` is correct when run from inside a child workspace. The main root is resolved by walking up to the nearest `.jj` and following its `repo` entry: a directory means "this is the main copy", a file is a pointer to `<main>/.jj/repo` (`workspace.MainRepoRoot`). Workspace directories are rendered relative to that root (`workspace.RelativeToMain`) rather than as absolute paths.
- **Task progress column.** `status` reads each spawn's `openspec/changes/<change>/tasks.md` and counts `- [x]` vs total checkboxes, rendering `done/total (pct%)`. A missing/unreadable `tasks.md` renders `-` and never errors — consistent with the no-fail, derive-live stance.

## Risks / Trade-offs

- **Name-collision between a change named `default` and the main workspace** — unlikely (openspec change names are kebab descriptions), but reopen/status must not treat the main copy as a spawn. Mitigate by excluding the `default` workspace name.
- **Reopen relaunches agents unconditionally** — reopening a detached spawn starts a fresh agent even if the previous one had made progress in-pane. Acceptable: the work lives in the jj workspace, not the agent's scrollback; the agent (`/opsx:apply`) resumes from the workspace state.
- **Current-session-only scope** may surprise a user who spawned into another session — `status` shows those as detached. Documented as intentional (ADR-006 consequence); cross-session is a future concern.
