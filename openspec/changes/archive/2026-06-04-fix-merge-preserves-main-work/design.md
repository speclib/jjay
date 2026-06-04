## Context

`jjay merge` (`internal/merge/merge.go`) integrates a workspace into the `main` **bookmark**:
```
  jj rebase -b <change>@ -d main      # onto the bookmark
  jj new main <change>@ -m "merge"    # parents: main bookmark + change@
  jj bookmark set main -r @           # advance bookmark to merge commit
  jj new                              # fresh change
```
The `rebase-before-merge` fix (jjay-30gc) made this safe for files modified on both sides. It did not address the case where the orchestrator does work in the main **working copy** (`@`) — new proposals, bean edits — that sits in commits *ahead of* the `main` bookmark and was never folded into it. Such work is in neither merge parent, so the merge commit excludes it; advancing the bookmark to the merge commit then orphans it. Observed live in jjay-ug7y: `add-init-command` and `add-change-completion` disappeared from main (recovered from orphan commit `60f1a51`).

The existing e2e scenario "Workspace adds new files, main adds different files" did *not* catch this because its test (`TestMerge_WorkspaceAddsNewFiles`) folds main's new file onto the bookmark (`jj bookmark set main`) before merging — it never exercises main work that is *ahead of* the bookmark.

## Goals / Non-Goals

**Goals:**
- No silent loss of committed main-line work during `jjay merge`, including ahead-of-bookmark `@` commits.
- Loud, non-destructive abort when safe inclusion isn't possible (bookmark not moved).
- A regression test (`TestMerge_MainAddsNewFiles`) that fails on today's code and passes after the fix.

**Non-Goals:**
- Reverting or replacing `rebase-before-merge` — this extends it.
- Handling uncommitted dirty state by auto-committing it (out of scope; abort instead).
- Cross-repo or push behavior (unchanged; merge still does not push).

## Decisions

- **Target the true main tip, not just the bookmark (ADR-010).** Before the merge commit, detect whether `@` is ahead of the `main` bookmark (e.g. `jj log -r 'main..@'` non-empty, or compare bookmark target to the main-line head). If so, fold that head into the integration target — advance the `main` bookmark to the committed main head first, or use that head as the rebase/merge destination — so the merge tree is `full-main-line ∪ change@`.
- **No abort needed for "dirty" state.** Implementation revealed jj auto-snapshots the working copy: uncommitted main edits become part of `@` before any jj command runs, so they show up as ahead-of-bookmark work and are folded in by the rule above. There is no unreachable dirty state to abort on (unlike git). The rebase-conflict abort already in `merge.go` is preserved for genuine content conflicts.
- **Investigation task first.** Before coding, reproduce against the real lifecycle (spawn → commit new dir in main `@` → merge) and confirm the exact divergence with `jj op log` / `jj log -r 'main..@'`, so the fix targets the verified mechanism rather than the hypothesis.
- **Test mirrors the existing one.** Add `TestMerge_MainAddsNewFiles` next to `TestMerge_WorkspaceAddsNewFiles` in `internal/merge/merge_integration_test.go`, but create the main-side file in a commit ahead of the bookmark (do NOT `jj bookmark set main` to it) — that's the distinction the current test misses.

## Risks / Trade-offs

- **Behavior change on dirty main:** merge may now refuse where it previously "succeeded" by losing data. Correct, but a visible change — document it.
- **jj revset/bookmark semantics** are the crux; the fix must be validated end-to-end, not just by unit reasoning (ADR-010 consequence). The investigation task de-risks this.
- **Interaction with spawn's `jj new` snapshot:** spawn already moves main work to `@-`; the fix must compose with that, ensuring "ahead of bookmark" is computed against the right base.
- **Recovery guidance:** until merged, advise that lost work is usually recoverable from orphan commits (`jj op log`, `git checkout <orphan> -- <path>`), as done for ug7y.
