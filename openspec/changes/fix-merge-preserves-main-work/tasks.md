## 1. Reproduce & confirm root cause

- [ ] 1.1 Write a failing repro against the real lifecycle: spawn a workspace, commit a new dir in main `@` (ahead of the `main` bookmark), run `jjay merge`, observe the dir missing from main.
- [ ] 1.2 Confirm the mechanism with `jj log -r 'main..@'` and `jj op log` — verify the lost work was ahead of the bookmark and excluded from both merge parents.

## 2. Fix merge to integrate the full main line

- [ ] 2.1 In `internal/merge/merge.go`, before the merge commit, detect whether `@` (main working copy) is ahead of the `main` bookmark.
- [ ] 2.2 If ahead, include that work in the merge target (advance `main` to the committed main head, or use that head as the rebase/merge destination) so the merge tree is `full-main-line ∪ change@`.
- [ ] 2.3 If ahead-of-bookmark work cannot be safely included (e.g. dirty/uncommitted `@`), abort with a clear message and do NOT move the `main` bookmark.
- [ ] 2.4 Preserve existing rebase-before-merge conflict behavior (no regression of jjay-30gc fix).

## 3. Regression tests

- [ ] 3.1 Add `TestMerge_MainAddsNewFiles` to `internal/merge/merge_integration_test.go`: main commits a new file/dir AHEAD of the bookmark (do not `jj bookmark set main` to it); workspace commits unrelated work; after merge BOTH exist on main.
- [ ] 3.2 Add an abort test: unsafe-to-include main state → merge fails non-zero, bookmark unchanged, no loss.
- [ ] 3.3 Confirm the new test FAILS on pre-fix code and PASSES after the fix.

## 4. Docs & beans

- [ ] 4.1 README/CHANGELOG: note merge now preserves ahead-of-bookmark main work and aborts rather than dropping.
- [ ] 4.2 Confirm ADR-010 reflects the implemented approach; flip to Accepted on archive.
- [ ] 4.3 Set `jjay-ug7y` to `in-progress`; add `openspec-link` on archive.

## 5. Verify

- [ ] 5.1 Run `make test-integration`; confirm all merge scenarios pass including the new ones.
- [ ] 5.2 Manually: spawn a change, create a new proposal dir in main, merge the spawn, confirm the new proposal dir survives on main.
