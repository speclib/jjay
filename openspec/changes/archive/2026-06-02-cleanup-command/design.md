## Context

`jjay spawn` creates three resources: jj workspace, tmux window, workspace directory. Tearing these down manually requires three separate commands. Cleanup reverses spawn.

Spawn and cleanup share naming conventions (`windowName`, `workspaceDir`). These currently live in `internal/spawn/`. Both packages will grow, so shared helpers should be extracted.

## Goals / Non-Goals

**Goals:**
- Working `jjay cleanup <change>` that tears down all spawn artifacts
- Tolerant execution — skip missing pieces, no errors for partial state
- Extract shared helpers to `internal/workspace/`

**Non-Goals:**
- Interactive confirmation (just do it, like rm -rf)
- Cleaning up jj changes/commits (that's merge territory)
- Batch cleanup of all workspaces

## Decisions

### Package structure

```
internal/
├── workspace/workspace.go  ← WindowName(), WorkspaceDir()
├── spawn/spawn.go           ← uses workspace package
└── cleanup/cleanup.go       ← uses workspace package
```

Shared helpers are exported functions in `internal/workspace/`. Both spawn and cleanup import this package.

_Alternative: keep helpers in spawn, have cleanup import spawn — rejected because it creates a dependency from cleanup to spawn that doesn't make semantic sense._

### Tolerant execution

Each cleanup step checks if the resource exists before acting. Missing resources are reported as "skipped" not errors. This handles:
- Partial spawn failures (some resources created, others not)
- User already killed the tmux window manually
- User already ran `jj workspace forget`

### Execution order: tmux → jj → directory

1. Kill tmux window (stops the running agent immediately)
2. Forget jj workspace (removes workspace from jj tracking)
3. Remove directory (cleans up filesystem)

Rationale: kill the agent first to prevent it from making more changes. Then clean up jj state. Finally remove files.

_Alternative: directory first — rejected because the agent might still be writing to it._

### Output format

```
Cleaning up change "feat-payments"...
  tmux window ws-feat-payments: killed
  jj workspace feat-payments: forgotten
  workspace directory: removed
```

Or for partial state:
```
Cleaning up change "feat-payments"...
  tmux window ws-feat-payments: not found, skipped
  jj workspace feat-payments: forgotten
  workspace directory: removed
```

## Risks / Trade-offs

- [No confirmation] → Acceptable for a dev CLI. The data in the workspace is in jj history anyway.
- [Agent may be mid-operation when killed] → tmux kill-window sends SIGHUP. Claude handles this gracefully. jj changes are committed via auto-snapshot.
