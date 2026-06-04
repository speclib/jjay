## 1. Reproduce & confirm root cause

- [x] 1.1 Write a failing repro against the real lifecycle: spawn a workspace, commit a new dir in main `@` (ahead of the `main` bookmark), run `jjay merge`, observe the dir missing from main.
- [x] 1.2 Confirm the mechanism with `jj log -r 'main..@'` and `jj op log` — verify the lost work was ahead of the bookmark and excluded from both merge parents.

## 2. Fix merge to integrate the full main line

- [x] 2.1 In `internal/merge/merge.go`, before the merge commit, detect whether `@` (main working copy) is ahead of the `main` bookmark.
- [x] 2.2 If ahead, include that work in the merge target (advance `main` to the committed main head, or use that head as the rebase/merge destination) so the merge tree is `full-main-line ∪ change@`.
- [x] 2.3 ~~If ahead-of-bookmark work cannot be safely included (e.g. dirty/uncommitted `@`), abort.~~ N/A in jj: the working copy is auto-snapshotted into `@`, so there is no unreachable "dirty" state — uncommitted main work is captured and preserved by 2.2. See design.md update; spec abort-scenario revised.
- [x] 2.4 Preserve existing rebase-before-merge conflict behavior (no regression of jjay-30gc fix).

## 3. Regression tests

- [x] 3.1 Add `TestMerge_MainAddsNewFiles` to `internal/merge/merge_integration_test.go`: main commits a new file/dir AHEAD of the bookmark (do not `jj bookmark set main` to it); workspace commits unrelated work; after merge BOTH exist on main.
- [x] 3.2 ~~Add an abort test: unsafe-to-include main state.~~ N/A — jj auto-snapshots the working copy, so there is no unreachable dirty state to abort on (see 2.3 / design.md). The "abort" scenario describes a state jj cannot produce.
- [x] 3.3 Confirm the new test FAILS on pre-fix code and PASSES after the fix. (Verified: FAIL without `advanceMainToHead`, PASS with it; full merge suite green.)

## 4. Docs & beans

- [x] 4.1 CHANGELOG: noted merge now folds ahead-of-bookmark main work into main before merging (jj auto-snapshot preserves uncommitted edits too). README merge line unchanged (behavior is transparent to the documented `jjay merge <change>` usage).
- [x] 4.2 ADR-010 updated to reflect the implemented approach (`latest(main..@ & ~empty())`, no dirty-state abort); flip to Accepted on archive.
- [x] 4.3 `jjay-ug7y` is `in-progress`; `openspec-link` added on archive.

## 5. Verify

- [x] 5.1 Ran `go test -tags integration ./...` — entire suite green, all 8 merge scenarios pass (7 prior + `TestMerge_MainAddsNewFiles`).
- [x] 5.2 The spawn→new-proposal-in-main→merge lifecycle is exercised end-to-end by `TestMerge_MainAddsNewFiles` (workspace + main-ahead-of-bookmark new dir → both survive), and was reproduced live this session against the buggy code before the fix. Equivalent to the manual check.
