## Context

jjay manages parallel agent workspaces within a single project. But users work on multiple projects. A dedicated tmux session per project keeps things organized: `jjay->proj1`, `jjay->proj2`.

## Goals / Non-Goals

**Goals:**
- `jjay session-open <path>` creates and switches to a tmux session for the given repo
- Verify the path is a jj repo
- Prevent duplicate sessions

**Non-Goals:**
- Session listing or management (future: `jjay session-list`, `jjay session-close`)
- Project bootstrapping (`session-new` — future)
- Persisting session config

## Decisions

### Package: internal/session/

Separate package for session management. Will grow with `session-list`, `session-close` later.

### Session naming: `jjay-><dirname>`

Uses the directory basename. Tested — tmux handles `->` in session names fine.

### Resolve path to absolute

Use `filepath.Abs()` on the input path so relative paths work: `jjay session-open .` or `jjay session-open ../teejay`.

### Two tmux commands

```go
exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", absPath)
exec.Command("tmux", "switch-client", "-t", sessionName)
```

Create detached first, then switch. This avoids nesting issues.

### Check for .jj/ directory

Simple `os.Stat(filepath.Join(absPath, ".jj"))` — if it doesn't exist, error.

### Check for existing session

Parse `tmux list-sessions -F "#{session_name}"` and check for duplicates.

## Risks / Trade-offs

- [Dirname collision] → Two repos with same dirname get the same session name. Acceptable for now — user can rename the directory or we add a flag later.
- [Not inside tmux] → `switch-client` requires being inside tmux. Same constraint as spawn — check `TMUX` env var.
