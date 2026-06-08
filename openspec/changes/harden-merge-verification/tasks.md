## 1. Reproduce the three q6ko instances as failing integration tests

These encode the reproduced failures from jjay-q6ko. Each MUST fail on today's code and pass after the fix.

- [ ] 1.1 Add test helpers to `internal/merge/merge_integration_test.go`:
  - [ ] `isStale(t, wsDir) bool` — run `jj status` inside the spawned workspace dir, report whether it contains `working copy is stale`.
  - [ ] a way to build a **non-linear** workspace DAG (e.g. commit, `jj new <earlier-rev>` to fork, commit a sibling) so the divergent-sibling topology can be set up.
- [ ] 1.2 `TestMerge_EmptyAtWorkInParent` (q6ko instance 2): workspace commits real work in `@-`, leaves `@` empty; after merge the `@-` work is on main and the merge commit is non-empty.
- [ ] 1.3 `TestMerge_WorkOnSibling` (q6ko instance 3): workspace's real work is on a sibling commit `<change>@` does not descend from; after merge the sibling work is on main (not orphaned). **This is the test that does not exist today.**
- [ ] 1.4 `TestMerge_NotStaleAfterMerge` (q6ko instance 1): after a normal successful merge, the workspace is forgotten / not stale.
- [ ] 1.5 Confirm 1.2–1.4 FAIL on pre-fix `merge.go` (`go test -tags integration ./internal/merge/`).

## 2. Robust work frontier (closes instances 2 & 3)

- [ ] 2.1 In `internal/merge/merge.go`, replace the `<change>@`-only revset with a work-frontier definition: `main..(heads of the workspace's commits) & ~empty()` (includes `@-` and divergent siblings). Validate the exact revset against the real spawn lifecycle, not just unit reasoning.
- [ ] 2.2 Make rebase + merge operate over the frontier so siblings/`@-` work are included.
- [ ] 2.3 Keep `advanceMainToHead` (ADR-010) first; compute the frontier against the advanced bookmark. Preserve the rebase-conflict abort (no regression of jjay-30gc).

## 3. Smoke test (rse4 L1+L2) + recovery handle

- [ ] 3.1 Capture the `jj op log` head before any rewrite; thread it into failure messages as `jj op restore <id>`.
- [ ] 3.2 Capture the added/modified file set across the work frontier BEFORE merge (union over siblings, not just `<change>@`).
- [ ] 3.3 After merge, run L1 (workspace had work ⇒ main gained changes) and L2 (every captured file present on main). Verbose by default.
- [ ] 3.4 On smoke-test failure: loud structured warning naming expected vs missing files + recovery handle; exit non-zero; no auto-rollback.

## 4. Verification-gated workspace lifecycle (closes instance 1)

- [ ] 4.1 On smoke-test PASS: forget the jj workspace (reuse the `cleanup.forgetWorkspace` pattern) so no stale pointer remains.
- [ ] 4.2 On smoke-test FAIL: keep the workspace intact; do not advance further; do not forget.
- [ ] 4.3 `TestMerge_UnprovenKeepsWorkspace`: force a missing-file/empty-land condition; assert workspace kept, non-zero exit, recovery handle in the warning.

## 5. Testing — full suite green

- [ ] 5.1 All new scenarios pass after the fix; the 8 prior merge scenarios still pass.
- [ ] 5.2 `TestMerge_SmokeDetectsEmptyLand` (rse4 L1): workspace had work but main gains nothing ⇒ smoke test fails loud.
- [ ] 5.3 Run `go test -tags integration ./...` — entire suite green.
- [ ] 5.4 Confirm each q6ko-instance test (1.2–1.4) now PASSES (was FAIL in 1.5).

## 6. Docs & beans

- [ ] 6.1 CHANGELOG: merge now verifies the work landed and forgets the workspace on success (no more stale errors); keeps + warns non-destructively on failure with a recovery handle. Headline ≤ 80 chars; detail + `See` links on sub-bullets.
- [ ] 6.2 ADR-013 flipped to Accepted on archive. README merge line updated if behavior is user-visible (workspace forgotten on clean merge).
- [ ] 6.3 `jjay-q6ko` set to in-progress with this proposal linked; `openspec-link` added on archive.
- [ ] 6.4 `jjay-rse4` set to scrapped with a `## Reasons for Scrapping` section pointing at this proposal (folded in).

## 7. Verify

- [ ] 7.1 Reproduce the real lifecycle end-to-end: spawn → commit work on a sibling / in `@-` → `jjay merge` → confirm work on main, workspace forgotten, no stale error.
- [ ] 7.2 Confirm the unproven path end-to-end: induce a miss → confirm workspace kept + loud warning + `jj op restore` handle.
