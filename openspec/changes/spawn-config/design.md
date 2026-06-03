## Context

Spawn and cleanup hardcode the agent command, tmux session, and workspace root. This blocks integration testing and multi-agent support. All three were marked "configurable later" in the original spawn bean.

## Goals / Non-Goals

**Goals:**
- CLI flags for agent command, tmux session, workspace root
- Full lifecycle integration test (spawn → verify → cleanup → verify)
- Fake agent script for testing

**Non-Goals:**
- Config file (flags are sufficient for now)
- Agent discovery or registry
- Per-change agent configuration

## Decisions

### Configuration via CLI flags, not config file

Flags on spawn/cleanup commands. No config file yet — that's premature until we see usage patterns.

| Flag | Default | Spawn | Cleanup |
|------|---------|-------|---------|
| `--agent` | claude command | yes | no |
| `--session` | current tmux session | yes | yes |
| `--workspace-root` | `../<project>-workspaces` | yes | yes |

_Alternative: YAML config file — rejected as premature. Flags are simpler and sufficient._

### Agent command templating

The `--agent` flag accepts a command string with `{change}` and `{wsdir}` placeholders:

```
--agent "./fake-agent.sh {change}"
--agent "codex --prompt '/opsx:apply {change}'"
```

Default: `claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}`

The command is passed to `tmux send-keys` as-is after substitution.

### SpawnOptions / CleanupOptions structs

Refactor `Spawn(changeName string)` to `Spawn(changeName string, opts SpawnOptions)`. Same for Cleanup. The options struct holds the configurable values:

```go
type SpawnOptions struct {
    Agent         string // agent command template
    Session       string // tmux session name (empty = current)
    WorkspaceRoot string // override workspace root (empty = default)
}
```

### Fix: tmux pane working directories (jjay-ps3d)

Current code uses `send-keys "cd <wsDir>"` to set the working directory for the right pane. This races with shell initialization — fish/bash may not be ready when the `cd` arrives.

Fix: use tmux's `-c` flag to set the starting directory at pane creation time:
- `tmux new-window -d -n ws-X -c <wsDir>` — window starts in workspace dir
- `tmux split-window -h -t ws-X -c <wsDir>` — right pane starts in workspace dir
- Remove `send-keys cd` for the right pane entirely

The left pane still needs `send-keys` to launch the agent command, but the `cd` prefix is no longer needed since the window already starts in the workspace dir.

### Integration test setup

```
Test setup:
  1. Create temp dir
  2. Init jj repo in temp dir
  3. Init openspec with a dummy change
  4. Create dedicated tmux session: jjay-test-<random>
  5. Copy fake-agent.sh to temp dir

Test body:
  1. Spawn(changeName, opts) with fake agent, test session, temp workspace root
  2. Assert: tmux window exists, jj workspace exists, dir exists, agent marker file exists
  3. Assert: both panes have correct working directory (tmux display-message #{pane_current_path})
  3. Cleanup(changeName, opts) with test session, temp workspace root
  4. Assert: all gone

Test teardown:
  - Kill test tmux session
  - Remove temp dirs
  - jj workspace forget (if still exists)
```

### WorkspaceDir accepts optional root

`workspace.WorkspaceDir(changeName, root string)` — if root is empty, use the default `../<project>-workspaces`. If set, use `<root>/<changeName>`.

### Tmux session targeting

When `--session` is set, all tmux commands use `-t session:windowName` syntax:
- `tmux new-window -d -t session: -n ws-X`
- `tmux send-keys -t session:ws-X ...`
- `tmux kill-window -t session:ws-X`
- `tmux list-windows -t session`

When empty, use current session (no `-t` prefix for session, just window name).

## Risks / Trade-offs

- [Flag API may change] → Acceptable for alpha. Flags are easy to rename.
- [Integration test requires tmux installed] → Skip test if tmux not available.
- [Integration test requires jj installed] → Skip test if jj not available.
- [Fake agent timing] → Agent runs in tmux pane asynchronously. Test needs a short sleep or poll to wait for the marker file.
