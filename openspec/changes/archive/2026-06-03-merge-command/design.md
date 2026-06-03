## Context

The jjay lifecycle is spawn → (agent works) → merge → cleanup. Spawn and cleanup are implemented. Merge is the missing piece that brings the workspace's work into main.

Currently done manually:
```bash
jj new main <workspace-change> -m "merge <change> into main"
jj bookmark set main -r @
jj new
```

## Goals / Non-Goals

**Goals:**
- Working `jjay merge <change>` that merges workspace work into main
- Resolve workspace change ID automatically
- Precondition checks
- Clean user experience (fresh change after merge)

**Non-Goals:**
- Auto-push (user decides when to push)
- Auto-cleanup (separate command)
- Conflict resolution (jj handles this — if there are conflicts, jj will report them)
- Merging multiple workspaces at once

## Decisions

### Resolve workspace change via revset

Use `jj log -r "<change>@" --no-graph -T 'change_id'` to get the workspace's current working copy change ID. The `<change>@` revset is jj's way to reference a workspace's working copy.

_Alternative: parse `jj workspace list` output — rejected because the revset is cleaner and more reliable._

### Merge is separate from cleanup

Merge only brings work into main. It doesn't kill the tmux window, forget the workspace, or remove the directory. The user may want to inspect the result before cleaning up.

```
jjay merge feat-payments    ← work is in main
# user inspects, maybe runs tests
jjay cleanup feat-payments  ← tear down
jj git push                 ← push when ready
```

### Warn but don't block on empty workspace

If the workspace `@` is empty, print a warning but proceed. The user might have a reason (e.g., the work is in `@-`). We don't know enough to block.

### Check emptiness via jj

Use `jj log -r "<change>@" --no-graph -T 'if(empty, "empty", "has-changes")'` to check if the workspace has changes.

### Package: internal/merge/

Same pattern as spawn and cleanup — separate package for the merge logic.

```
internal/
├── workspace/workspace.go
├── spawn/spawn.go
├── cleanup/cleanup.go
└── merge/merge.go          ← new
```

### Execution sequence

1. Verify workspace exists (parse `jj workspace list`)
2. Check if workspace `@` is empty (warn if so)
3. `jj new main <change>@ -m "merge <change> into main"`
4. `jj bookmark set main -r @`
5. `jj new` (fresh change for user)
6. Print success message

## Risks / Trade-offs

- [Merge conflicts] → jj handles conflict markers. If the merge has conflicts, `jj new` will create a conflicted commit. The user resolves in their working copy. Not our problem to solve.
- [Workspace still running] → merge doesn't check if the agent is still active in the tmux window. The user should wait for the agent to finish. Future: `jjay status` could show this.
- [Main bookmark doesn't exist] → jj will error on `jj new main ...`. Clear enough error message from jj itself.
