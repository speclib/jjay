## Why

`jjay merge` silently deletes work that was created in the **main working copy after a spawn** — a critical, data-losing regression class. Bean [jjay-ug7y](../../.beans/jjay-ug7y--critical-merge-deletes-active-sibling-change-dirs.md) captured it after two active change-dirs (`add-init-command`, `add-change-completion`) vanished from main when `workspace-aware-session` was merged. They were recoverable from an orphan jj commit, but the loss was silent.

This is a sibling of [jjay-30gc](../../.beans/jjay-30gc--critical-openspec-artifacts-diverge-on-merge-files.md) (fixed by `rebase-before-merge`). That fix solved files **modified on both sides**. It does **not** cover new work added to main *after* the spawn snapshot — because the bug is not in the rebase, it's in what `merge` treats as "main".

**Root cause.** `internal/merge/merge.go` operates entirely on the `main` **bookmark**:
```
  jj rebase -b <change>@ -d main      # rebase onto the bookmark
  jj new main <change>@ -m "merge"    # merge commit's parents: main bookmark + change@
  jj bookmark set main -r @           # advance bookmark to the merge commit
```
The merge commit's tree is `main-bookmark-tree ∪ change@-tree`. But the lost work lived in the main **working copy** (`@`) in commits *ahead of* the `main` bookmark — created during the session, never folded into the bookmark. That work is in **neither** merge parent, so it is excluded; then `main` advances to the merge commit, and the un-bookmarked `@` work is orphaned. jj does not warn — it just isn't there anymore.

```
   main bookmark ──▶ (merge parent)        @ (main working copy, AHEAD of bookmark)
        │                                   │  add-init-command/      ← created post-spawn
        │                                   │  add-change-completion/  ← never on bookmark
        ▼                                   ▼
   merge commit = main ∪ change@   ✗ excludes @'s ahead-of-bookmark work
        │
        ▼  bookmark set main -r @(merge)
   main now MISSING the post-spawn work — silently
```

## What Changes

- **`jjay merge` preserves main-working-copy work created after spawn.** Before merging, it SHALL account for commits between the `main` bookmark and the current main working copy, so no committed-but-un-bookmarked main work is dropped. Concretely: detect that `@` (main working copy) is ahead of the `main` bookmark, and fold that work into the merge (e.g. rebase the workspace onto the actual head of main's working-copy line, or advance the merge base to include the ahead-of-bookmark commits) rather than onto the lagging bookmark alone.
- **Abort-with-explanation instead of silent loss.** If `merge` cannot safely include ahead-of-bookmark main work, it SHALL fail with a clear message (and not move the bookmark), never silently discard.
- **Regression test** `TestMerge_MainAddsNewFiles` (the mirror of the existing `TestMerge_WorkspaceAddsNewFiles`): main adds a new dir/file after the workspace base; after merge, both the main-only addition and the workspace work exist.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `merge`: merge SHALL preserve main-working-copy work that is ahead of the `main` bookmark, not just the bookmark's tree.

## Impact

- **Code**: `internal/merge/merge.go` (the rebase/merge-base/bookmark sequence); new scenario in `internal/merge/merge_integration_test.go`.
- **Severity**: CRITICAL — silent data loss of in-flight proposals; should land before further merges of the staged changes.
- **Relation**: extends the `rebase-before-merge` fix (archived 2026-06-04); does not revert it.
- **ADRs**: ADR-010 (merge integrates the full main line, not only the bookmark).
- **Beans**: ug7y → in-progress, linked here.
