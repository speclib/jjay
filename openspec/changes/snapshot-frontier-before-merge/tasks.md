## 1. Reproduce as failing tests

- [ ] 1.1 `TestMerge_UncommittedAtWorkSnapshotted` (jjay-4oyy): spawned workspace has real work left UNCOMMITTED in `@` (write files in wsDir, do not run jj there). On today's code the merge reports success but the files are absent from main. MUST fail before the fix.
- [ ] 1.2 `TestMerge_EmptyFrontierIsUnproven`: a workspace whose frontier is empty after snapshot ⇒ merge exits non-zero, keeps the workspace, warning carries the recovery handle. MUST fail before the fix (today it passes vacuously).

## 2. Tier 1 — force-snapshot before the frontier

- [ ] 2.1 In `internal/merge/merge.go`, before `workFrontierFiles`, resolve the workspace dir and run `jj -R <wsDir> status` (snapshotting no-op; ignore stdout). Place it after `advanceMainToHead`, before frontier capture.
- [ ] 2.2 Confirm against real jj that this snapshots the spawned workspace's dirty `@` (spike-proven) and does not disturb the main session `@`.

## 3. Tier 2 — empty frontier ⇒ unproven (the guarantee)

- [ ] 3.1 Change the smoke test: an empty frontier (post-snapshot) is NOT "nothing to prove" — return a loud error so the caller does not forget / `jj new` / print "verified". Keep the workspace; include the `jj op restore <preMergeOp>` handle.
- [ ] 3.2 Dir-content hint: before erroring, run `jj -R <wsDir> diff --from main --summary`; if non-empty, name those paths in the warning. Best-effort (hint failure doesn't change the refuse outcome).
- [ ] 3.3 Keep the existing `checkWorkspaceEmpty` early warning.

## 4. Audit behavior change

- [ ] 4.1 `TestMerge_EmptyWorkspace` previously merged an empty workspace expecting success — update it: an empty spawn now exits non-zero (unproven) with a clear message. Decide/confirm this is the intended new contract.

## 5. Testing — full suite green

- [ ] 5.1 1.1 and 1.2 now PASS after the fix; the prior merge scenarios still pass (modulo the empty-workspace contract change in 4.1).
- [ ] 5.2 `go test ./...` and `go test -tags integration ./...` green.

## 6. Docs & beans

- [ ] 6.1 CHANGELOG: merge now snapshots the workspace before defining its work and refuses to claim success on an empty frontier (headline ≤ 80 chars; detail on sub-bullets).
- [ ] 6.2 ADR-015 → Accepted on archive.
- [ ] 6.3 `jjay-4oyy` → in-progress, linked; → completed on archive.
- [ ] 6.4 `jjay-ychu`: note its *detection* acceptance is satisfied here; its remaining scope is Tier-3 op-log *identification* of orphan siblings. Keep open.

## 7. Verify

- [ ] 7.1 End-to-end: spawn, leave work uncommitted in the workspace `@`, `jjay merge` → confirm work lands on main and is verified (4oyy closed).
- [ ] 7.2 End-to-end: induce an empty frontier (empty spawn) → confirm merge refuses, keeps the workspace, prints the recovery handle.
