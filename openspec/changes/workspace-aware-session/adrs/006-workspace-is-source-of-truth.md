# ADR-006: Workspace existence is the single source of truth

**Status**: Proposed

## Context

A jjay spawn produces three things: a jj workspace (on disk), a tmux window named `ws-<change>`, and an agent process running inside it. These have very different lifetimes. The jj workspace persists across reboots and tmux server restarts. The tmux window and the agent are volatile — closing the tmux server, detaching, or rebooting destroys them while the workspace survives.

To report what is "running" (`jjay status`) and to restore the view after `jjay session-open`, jjay needs a definition of when a spawn is "open." The obvious-but-wrong approach is a state file that records spawns, which immediately drifts from reality: a tmux window dies and the file still says "running"; a workspace is forgotten and the file still lists it. ADR-003 already treats "workspace exists but no tmux window" as a recoverable partial state rather than corruption — this ADR makes that assumption explicit and load-bearing.

## Options Considered

- **State file** — persist a list of spawns (change, wsdir, pid, window). Authoritative-looking but drifts: every external `jj workspace forget`, `tmux kill-window`, reboot, or crash desyncs it. Requires reconciliation logic that is itself just re-derivation from jj + tmux.
- **tmux as source of truth** — a spawn is "open" iff its `ws-*` window exists. Wrong: a detached/rebooted session has no windows, yet the workspaces (and the user's work) are very much still there.
- **jj workspace as source of truth** — a spawn is "open" iff its jj workspace exists. tmux windows and the agent are a recreatable *view* of that truth, not part of it. No persisted state; everything derived live.

## Decision

The **jj workspace is the single source of truth**. A spawn is "open" if and only if its jj workspace exists (`jj workspace list` contains `<change>:`). jjay persists **no** state file.

- tmux windows (`ws-<change>`) and the agent process are a recreatable view of the workspace, not authoritative state.
- `jjay status` derives its report live by joining `jj workspace list` with `tmux list-windows`; a workspace with no window is reported as **detached** (open, not attached), not as closed.
- `jjay session-open` repairs the view by recreating `ws-<change>` windows for workspaces that lack one.
- Removing a spawn means removing its workspace — which is `jjay cleanup`'s job, not a state-file edit.

This is already how `spawn` behaves (`internal/spawn/spawn.go` checks `jj workspace list` and `tmux list-windows` live, persisting nothing); this ADR elevates that behavior to a stated invariant the rest of the lifecycle (`status`, `session-open`, `cleanup`) must honor.

## Consequences

- **Positive**: No state file to drift, reconcile, lock, or corrupt. jjay is stateless between invocations.
- **Positive**: Survives reboots, tmux server restarts, and out-of-band `jj`/`tmux` commands — the truth is recomputed every time.
- **Positive**: `status` (read the diff) and `session-open` reopen (repair the diff) share one definition and one join, so they cannot disagree.
- **Negative**: "Open" is scoped to current tmux session for window detection (see the session-open delta); cross-session discovery is deferred and not modeled here.
- **Negative**: Cannot record metadata that has no jj/tmux home (e.g. which agent command launched a spawn) without reintroducing persisted state — accepted for now; reopen re-derives the agent command from the shared default/template.
