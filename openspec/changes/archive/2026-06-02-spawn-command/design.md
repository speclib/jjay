## Context

jjay has a working Go scaffold with cobra CLI. This is the first real feature: `jjay spawn <change>` orchestrates jj, tmux, and claude to set up a parallel agent workspace.

All three external tools (jj, tmux, claude) are invoked via `os/exec`. The agent (claude) and layout are hardcoded for now — configurability comes later.

## Goals / Non-Goals

**Goals:**
- Working `jjay spawn <change>` that creates workspace + tmux window + agent
- Clear error messages for all precondition failures
- Clean separation of spawn logic in `internal/spawn/`

**Non-Goals:**
- Configurable agent (always claude for now)
- Configurable layout (always two-pane horizontal split)
- Configurable permissions flags (always `--dangerously-skip-permissions`)
- State tracking (separate bean: jjay-72pw)
- Nix develop shell in the right pane (future)

## Decisions

### Package structure: `internal/spawn/`

Spawn logic lives in `internal/spawn/spawn.go`. The cobra command in `cmd/jjay/` calls into this package. This keeps the CLI layer thin and the logic testable.

_Alternative: everything in main.go — rejected because it won't scale as we add merge, cleanup, status._

### Precondition check order

1. Check TMUX env var (cheapest, no subprocess)
2. Check openspec change exists (`openspec list --json`)
3. Check jj workspace doesn't exist (`jj workspace list`)
4. Check tmux window name not taken (`tmux list-windows`)

Fail fast on the first error. This avoids partial state (e.g., creating a workspace but failing on tmux).

_Alternative: check all and report all errors — rejected for simplicity. One clear error is better than a wall of text._

### Subprocess execution sequence

After preconditions pass:
1. Create parent dir `../<project>-workspaces/` if needed
2. `jj workspace add --name <change> --revision @ ../<project>-workspaces/<change>` — create workspace based on current working copy (so uncommitted openspec changes are included)
3. `tmux new-window -d -n "ws-<change>"` — create window without switching focus
4. `tmux send-keys -t "ws-<change>" "cd <wsDir> && claude \"/opsx:apply <change>\" --dangerously-skip-permissions --add-dir <wsDir>" Enter` — launch agent with workspace access
5. `tmux split-window -h -t "ws-<change>"` — create right pane
6. `tmux send-keys -t "ws-<change>.1" "cd <wsDir>" Enter` — cd shell pane to workspace

If any step fails after workspace creation, we don't roll back — the user can clean up manually (bean jjay-uypj for cleanup command). Rolling back adds complexity with little value at this stage.

### Naming conventions

| Thing | Convention | Example |
|-------|-----------|---------|
| jj workspace | `<change>` | `feat-payments` |
| workspace dir | `../<project>-workspaces/<change>` | `../jjay-workspaces/feat-payments` |
| tmux window | `ws-<change>` | `ws-feat-payments` |

### Hardcoded agent command

```
claude "/opsx:apply <change>" --dangerously-skip-permissions --add-dir <wsDir>
```

`--add-dir` grants claude access to the workspace directory, bypassing the folder trust dialog. This will become configurable later. For now, claude is the only supported agent.

## Risks / Trade-offs

- [No rollback on partial failure] → Acceptable for v1. User can `jj workspace forget` + `tmux kill-window` manually.
- [Hardcoded agent] → Will need config before supporting codex/mistral. Bean exists for this implicitly via the spawn bean's "configurable later" notes.
- [`--dangerously-skip-permissions`] → Required for autonomous agent operation but worth flagging. Will be configurable.
