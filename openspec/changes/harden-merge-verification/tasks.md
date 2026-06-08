## 1. Reproduce the q6ko instances as tests

These encode the reproduced failures from jjay-q6ko.

- [x] 1.1 Add test helpers to `internal/merge/merge_integration_test.go`:
  - [x] `workspaceExists(t, repoDir, name)` — whether the jj workspace is still registered (used to assert forget-on-success / keep-on-failure). (`isStale` was considered but forget-on-success makes "gone" the cleaner assertion than "stale".)
  - [x] non-linear DAG setup via `jj new <earlier-rev>` to fork a sibling.
- [x] 1.2 `TestMerge_EmptyAtWorkInParent` (q6ko instance 2): work in `@-`, `@` empty; after merge the `@-` work is on main, non-empty. PASSES (frontier reaches `@-`).
- [x] 1.3 Instance 3 (true orphan sibling): spike + test proved it is UNDETECTABLE by merge (frontier empty ⇒ smoke test has nothing to expect). jj links only one `@` per workspace. Reframed: covered by `TestMerge_SmokeDetectsMissingFile` for the detectable case; true-orphan op-log detection deferred to jjay-ychu.
- [x] 1.4 `TestMerge_NotStaleAfterMerge` (q6ko instance 1): after a proven merge the workspace is forgotten (no stale pointer possible). PASSES.

## 2. Robust work frontier (closes instances 2; bounds 3)

- [x] 2.1 Replaced `<change>@`-only with the ancestor frontier `ancestors(<change>@) & main.. & ~empty()` (includes `@-`). Spike-validated against empty-`@`. Divergent siblings are out of frontier scope (detected, not auto-included; see jjay-ychu).
- [x] 2.2 Frontier file capture (`workFrontierFiles`) aggregates added/modified files across the frontier commits.
- [x] 2.3 `advanceMainToHead` (ADR-010) kept first; rebase-conflict abort preserved (no regression — `TestMerge_SameFileModified` still passes).

## 3. Smoke test (rse4 L1+L2) + recovery handle

- [x] 3.1 `currentOpID()` captures the `jj op log` head before any rewrite; threaded into failure messages as `jj op restore <id>`.
- [x] 3.2 `workFrontierFiles` captures the added/modified set across the frontier before merge.
- [x] 3.3 `smokeTest` runs L1 (had work ⇒ main gained content) and L2 (every captured file present on main). Verbose by default.
- [x] 3.4 On failure: loud structured error naming missing files + recovery handle; returns error (non-zero exit); no auto-rollback.

## 4. Verification-gated workspace lifecycle (closes instance 1)

- [x] 4.1 On smoke-test PASS: `forgetWorkspace` forgets the jj workspace (no stale pointer).
- [x] 4.2 On smoke-test FAIL: merge returns the error before forgetting / before `jj new` — workspace kept, not advanced.
- [x] 4.3 `TestSmokeTest_L1L2` asserts a missing expected file fails loudly with the recovery handle.

## 5. Testing — full suite green

- [x] 5.1 New scenarios pass; the prior merge scenarios still pass (12 merge tests green).
- [x] 5.2 `TestMerge_SmokeDetectsMissingFile` covers the reachable-work L2 path; `TestSmokeTest_L1L2` covers L1/L2 directly.
- [x] 5.3 `go test -tags integration ./...` — entire suite green (incl. `test/integration` lifecycle).
- [x] 5.4 Instance tests (1.2, 1.4) pass; instance 3 reframed per the spike finding.

## 6. Docs & beans

- [x] 6.1 CHANGELOG updated (headline ≤ 80 chars).
- [ ] 6.2 ADR-013 → Accepted on archive. README merge line updated if user-visible.
- [ ] 6.3 `jjay-q6ko` → in-progress, linked; `openspec-link` on archive.
- [x] 6.4 `jjay-rse4` already scrapped with reasons pointing here (done at proposal time).
- [x] 6.5 Filed jjay-ychu for the undetectable true-orphan (op-log forensics).

## 7. Verify

- [x] 7.1 End-to-end via integration tests (spawn-equivalent temp repo → work in `@-` → merge → on main, workspace forgotten).
- [x] 7.2 Unproven path covered by `TestSmokeTest_L1L2` (miss → kept + loud + `jj op restore`).
