# Proposal: Spawn configuration and integration test

**Change**: spawn-config
**Status**: proposed
**Bean**: [jjay-93is — spawn configuration](../../../.beans/jjay-93is--spawn-configuration-agent-command-tmux-session-wor.md)

## Why

Spawn currently hardcodes the agent command (claude), tmux session (current), and workspace root (../<project>-workspaces/). This blocks three things:

1. **Integration testing** — can't run spawn/cleanup lifecycle tests without a real claude and without polluting the user's tmux session
2. **Multiple agents** — codex, mistral, or custom scripts can't be used
3. **Flexible workspace locations** — some users may want workspaces elsewhere

## What Changes

- **Fix: shell pane working directory** (jjay-ps3d) — use `tmux split-window -c <wsDir>` and `tmux new-window -c <wsDir>` instead of `send-keys cd`. Eliminates race condition where shell hasn't initialized when `cd` arrives.
- Add `--agent` flag to spawn (default: `claude "/opsx:apply <change>" --dangerously-skip-permissions --add-dir <wsDir>`)
- Add `--session` flag to spawn and cleanup (default: current tmux session)
- Add `--workspace-root` flag to spawn and cleanup (default: `../<project>-workspaces`)
- Update the shared workspace package to accept configurable root
- Create a fake agent script for testing
- Add Go integration tests (build tag `integration`) covering spawn → verify → cleanup → verify lifecycle
- Integration test creates its own tmux session and temp jj repo
- Integration test asserts both panes are in the correct workspace directory

## Capabilities

### New Capabilities

- `spawn-config`: Configurable agent command, tmux session, workspace root via CLI flags
- `integration-test`: Full lifecycle test for spawn + cleanup using fake agent and isolated tmux session

### Modified Capabilities

- `spawn`: Accept configuration flags
- `cleanup`: Accept `--session` and `--workspace-root` flags to match spawn

## Impact

- Modified: `internal/spawn/spawn.go`, `internal/cleanup/cleanup.go`, `internal/workspace/workspace.go`
- Modified: `cmd/jjay/main.go` (add flags to spawn and cleanup subcommands)
- New: `testdata/fake-agent.sh`
- New: `internal/spawn/spawn_integration_test.go`
- No new Go dependencies
