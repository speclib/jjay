## 1. Merged detection

- [ ] 1.1 Validate the detection revset against the spike oracle: archived changes ⇒ merged, active spawns ⇒ not merged. Start from `main..<change>@` empty ⇒ merged.
- [ ] 1.2 Settle the empty-workspace edge: a spawn with no work should read `MERGED=no` (reference `merge.go:checkWorkspaceEmpty`).
- [ ] 1.3 Confirm just-rebased-not-merged reads `MERGED=no`.

## 2. Implement

- [ ] 2.1 Add `Merged bool` to the `Spawn` struct in `internal/status/status.go`.
- [ ] 2.2 Compute `Merged` in `List`/`join` via the validated jj revset; tolerate evaluation failure as `MERGED=no` (don't fail status).
- [ ] 2.3 In `Render`: rename header `STATUS` → `TMUX`; add `MERGED` column. Final order: `CHANGE WORKSPACE TASKS TMUX MERGED ARCHIVED`.

## 3. Unit tests

- [ ] 3.1 Unit-test the merged derivation (merged / not-merged / empty / tolerant-of-jj-failure), mirroring existing status tests.
- [ ] 3.2 Update any test asserting the old `STATUS` header to expect `TMUX`, and assert the `MERGED` column.

## 4. Integration test (folds in jjay-y5yn)

- [ ] 4.1 Add a `status` subtest to `test/integration/full_lifecycle_test.go`, between `spawn` and `cleanup`, calling status against the live spawned env (e.g. `status.List` + `Render`, or the binary).
- [ ] 4.2 Assert the spawned change appears, the column set includes TMUX + MERGED, and the fresh spawn reports `MERGED=no`.
- [ ] 4.3 Confirm `make test-integration` is green with the new subtest.

## 5. Docs & beans

- [ ] 5.1 README/CHANGELOG: note the `MERGED` column and the `STATUS`→`TMUX` rename (headline ≤ 80 chars).
- [ ] 5.2 Set `jjay-nnjz` and `jjay-y5yn` to `in-progress`; add `openspec-link` to both on archive. (Hard half tracked in `jjay-atq6`.)
