# Tasks: rebase-before-merge

## 1. Add rebase step to merge

- [x] 1.1 Add `jj rebase -b <workspace>@ -d main` before the merge commit creation in `internal/merge/merge.go`
- [x] 1.2 Add conflict check after rebase via `jj log -r "<ws>@" -T 'if(conflict, ...)'`
- [x] 1.3 If conflicts detected, print clear error message and exit without merging
- [x] 1.4 Update success message to mention rebase happened

## 2. E2E test scenarios

- [x] 2.1 Create `internal/merge/merge_integration_test.go` with `//go:build integration`
- [x] 2.2 Implement test helper: create temp jj repo with main bookmark
- [x] 2.3 Implement test helper: create workspace with file changes
- [x] 2.4 Scenario 1: clean merge — no main changes, workspace files present
- [x] 2.5 Scenario 2: main moved forward, no overlap — files from both sides present
- [x] 2.6 Scenario 3: main and workspace modify same file — conflict detected, merge aborted
- [x] 2.7 Scenario 4: workspace adds new files — new files present after merge (THE BUG FIX)
- [x] 2.8 Scenario 5: empty workspace — warning printed, merge proceeds
- [x] 2.9 Scenario 6: multiple workspace commits — all changes present

## 3. Verification

- [x] 3.1 Verify `make test` passes (unit tests)
- [x] 3.2 Verify `make build` and `make lint` pass
- [x] 3.3 Run `make test-integration` and verify all 6 scenarios pass
