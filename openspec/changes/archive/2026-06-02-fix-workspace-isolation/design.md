## Context

`jjay spawn` currently uses `jj workspace add --revision @` to create a child workspace. The main workspace's `@` may have uncommitted changes. When the child workspace's agent creates jj operations, the main workspace becomes stale. Running `jj workspace update-stale` can lose the uncommitted work in `@`.

Discovered during manual testing on 2026-06-02 when a user lost uncommitted work and had to recover via `jj restore`.

## Goals / Non-Goals

**Goals:**
- Ensure the main workspace's uncommitted work is safe during concurrent agent work
- Child workspace still gets all files (including uncommitted openspec changes)

**Non-Goals:**
- Preventing staleness entirely (that's inherent to concurrent jj workspaces)
- Automatic recovery from staleness

## Decisions

### Snapshot via `jj new` before workspace creation

Run `jj new` in the main workspace before creating the child workspace. This:
1. Snapshots all uncommitted changes into `@-` (jj auto-snapshots on any command)
2. Creates a fresh empty `@` in the main workspace
3. The child workspace uses `--revision @-` to get the snapshot with all files

```
Before spawn:
  main @:  abc123 (uncommitted openspec change, code edits, etc.)

After jj new:
  main @:  def456 (empty, fresh — nothing to lose)
  main @-: abc123 (all files, committed via snapshot)

After workspace add --revision @-:
  main @:       def456 (empty, safe)
  child @:      ghi789 (new commit, parent = abc123, has all files)
```

_Alternative: run `jj new` in the child workspace after creation — rejected because the child's `@` would briefly share state with the main workspace's `@`, which is the original problem._

_Alternative: tell users to commit before spawning — rejected because it adds friction and is easy to forget._

### Preconditions before mutations

All precondition checks (tmux, openspec change, workspace exists, window exists) MUST run before `jj new`. If any check fails, the main workspace is untouched. This prevents the scenario where `jj new` runs but workspace creation fails, leaving the user on an unexpected empty `@`.

### Order of operations

1. Check TMUX env var
2. Check openspec change exists
3. Check jj workspace doesn't exist
4. Check tmux window name not taken
5. **`jj new`** ← new step, first mutation
6. Create parent directory
7. `jj workspace add --name <change> --revision @- <path>`
8. Create tmux window
9. Set up panes

## Risks / Trade-offs

- [User ends up on empty @] → Expected behavior. They can `jj edit @-` to get back to their work, or just continue working on the new `@`. Worth documenting.
- [jj new creates a commit in the log] → These are lightweight jj changes, not git commits. Normal jj workflow.
