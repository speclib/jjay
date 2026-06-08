package merge

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// Merge merges a workspace's work into the main bookmark as a verification-gated
// pipeline (ADR-013). It folds ahead-of-bookmark main work (ADR-010), defines the
// workspace's work as the ancestor frontier of <change>@ (so work in @- is included,
// not just @), rebases that onto main, then PROVES the work landed via a post-merge
// smoke test before declaring success. On proof it forgets the workspace (no live
// pointer ⇒ no staleness); on failure it keeps the workspace intact and emits a loud
// recovery handle rather than reporting false success (jjay-q6ko, jjay-rse4).
func Merge(changeName string) error {
	// 1. Verify workspace exists
	if err := checkWorkspaceExists(changeName); err != nil {
		return err
	}

	// 2. Warn if workspace is empty
	checkWorkspaceEmpty(changeName)

	// 2a. Capture a pre-merge recovery handle before any rewrite (rse4 Tension-1).
	preMergeOp := currentOpID()

	// 2b. Fold main-line work that is ahead of the `main` bookmark into the
	// bookmark before merging. The orchestrator often commits work in the main
	// working copy (new proposals, bean edits) that lands in commits ahead of
	// the `main` bookmark. Merging onto the lagging bookmark would leave that
	// work in neither merge parent and silently drop it (see ADR-010, jjay-ug7y).
	if err := advanceMainToHead(); err != nil {
		return err
	}

	// 2c. Capture the added/modified file set across the work frontier BEFORE
	// merge (rse4 L2). The frontier is ancestors(<change>@) not on main, so work
	// in @- is included — not just <change>@ (jjay-q6ko instance 2).
	expectedFiles, err := workFrontierFiles(changeName)
	if err != nil {
		return err
	}

	// 3. Rebase the workspace branch onto main
	revset := changeName + "@"
	fmt.Printf("Rebasing %s onto main...\n", changeName)
	cmd := exec.Command("jj", "rebase", "-b", revset, "-d", "main")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to rebase workspace onto main: %s", strings.TrimSpace(string(out)))
	}

	// 4. Check for conflicts after rebase
	if hasConflicts, err := checkConflicts(changeName); err != nil {
		return err
	} else if hasConflicts {
		return fmt.Errorf("rebase produced conflicts in workspace %q. Resolve them manually, then retry 'jjay merge %s'", changeName, changeName)
	}

	// 5. Create merge commit
	mergeMsg := fmt.Sprintf("merge %s into main", changeName)
	cmd = exec.Command("jj", "new", "main", revset, "-m", mergeMsg)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create merge commit: %s", strings.TrimSpace(string(out)))
	}

	// 6. Move main bookmark
	cmd = exec.Command("jj", "bookmark", "set", "main", "-r", "@")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update main bookmark: %s", strings.TrimSpace(string(out)))
	}

	// 7. Smoke test: PROVE the work landed before declaring success (rse4 L1+L2).
	// The bookmark now points at the merge commit; verify against it.
	if err := smokeTest(expectedFiles, preMergeOp); err != nil {
		// Unproven: non-destructive. Keep the workspace intact for recovery,
		// do not advance further, do not forget. Loud, exit non-zero.
		return err
	}

	// 8. Proven → forget the workspace so no stale working-copy pointer remains
	// (jjay-q6ko instance 1 closed structurally). Directory cleanup stays in
	// the separate `cleanup` command.
	forgetWorkspace(changeName)

	// 9. Fresh change for the user
	cmd = exec.Command("jj", "new")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create fresh change: %s", strings.TrimSpace(string(out)))
	}

	fmt.Printf("Rebased and merged %s into main (verified).\n", changeName)
	fmt.Println("Work confirmed on main; workspace forgotten. You're on a fresh change.")
	return nil
}

// currentOpID returns the current jj operation id, used as a recovery handle in
// failure messages (`jj op restore <id>`). Best-effort: returns "" on error.
func currentOpID() string {
	out, err := exec.Command("jj", "op", "log", "--no-graph", "--limit", "1",
		"-T", "id.short()").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// workFrontierFiles returns the set of files added or modified across the
// workspace's ancestor work frontier: ancestors(<change>@) that are not yet on
// main and are non-empty. This includes @- (jjay-q6ko instance 2), not just
// <change>@. A truly divergent sibling is unreachable from <change>@ and is not
// included here — it is caught by the smoke test instead (ADR-013).
func workFrontierFiles(changeName string) (map[string]bool, error) {
	revset := fmt.Sprintf("ancestors(%s@) & main.. & ~empty()", changeName)
	// --types shows the change type; -s gives names. Use --name-only-ish via
	// `jj diff --from main --to <frontier-tip>`? Simpler: list files changed by
	// each frontier commit. Use `jj diff -r <revset> --name-only` aggregated.
	out, err := exec.Command("jj", "log", "-r", revset, "--no-graph",
		"-T", "commit_id ++ \"\\n\"").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate work frontier: %w", err)
	}
	files := map[string]bool{}
	for _, commit := range strings.Fields(string(out)) {
		d, err := exec.Command("jj", "diff", "-r", commit, "--name-only").Output()
		if err != nil {
			continue
		}
		for _, f := range strings.Split(strings.TrimSpace(string(d)), "\n") {
			f = strings.TrimSpace(f)
			if f != "" {
				files[f] = true
			}
		}
	}
	return files, nil
}

// smokeTest proves the workspace's work landed on main (rse4 L1+L2). L1: if the
// workspace had work (expectedFiles non-empty) main must contain those files. L2:
// every expected file must be present in main's tree. On any miss it returns a
// loud error naming missing files and the recovery handle. Verbose by default.
func smokeTest(expectedFiles map[string]bool, preMergeOp string) error {
	if len(expectedFiles) == 0 {
		// Empty workspace — nothing to prove (checkWorkspaceEmpty already warned).
		return nil
	}

	mainFiles, err := filesOnMain()
	if err != nil {
		return fmt.Errorf("smoke test could not read main's files: %w (recover with: jj op restore %s)", err, preMergeOp)
	}

	// L1: the workspace had work — main must have gained content.
	if len(mainFiles) == 0 {
		return fmt.Errorf("merge smoke test FAILED (L1): workspace had work (%d files) but main gained nothing — the merge landed empty.\nRecover with: jj op restore %s", len(expectedFiles), preMergeOp)
	}

	// L2: every expected (added/modified) file must be present on main.
	var missing []string
	for f := range expectedFiles {
		if !mainFiles[f] {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Printf("merge smoke test: expected %d files from the work frontier; %d missing on main.\n", len(expectedFiles), len(missing))
		return fmt.Errorf("merge smoke test FAILED (L2): files missing from main after merge:\n  %s\nThe workspace is kept for recovery. Recover with: jj op restore %s",
			strings.Join(missing, "\n  "), preMergeOp)
	}

	fmt.Printf("merge smoke test passed: all %d work-frontier files present on main.\n", len(expectedFiles))
	return nil
}

// filesOnMain returns the set of files in the `main` bookmark's tree.
func filesOnMain() (map[string]bool, error) {
	out, err := exec.Command("jj", "file", "list", "-r", "main").Output()
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for _, f := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		f = strings.TrimSpace(f)
		if f != "" {
			set[f] = true
		}
	}
	return set, nil
}

// forgetWorkspace forgets the jj workspace (reuses the cleanup pattern). On a
// proven merge there is no remaining live working-copy pointer, so the workspace
// cannot go stale (jjay-q6ko instance 1). Best-effort: a failure to forget is
// reported but does not fail the (already successful) merge.
func forgetWorkspace(changeName string) {
	cmd := exec.Command("jj", "workspace", "forget", changeName)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("  note: could not forget workspace %q: %s\n", changeName, strings.TrimSpace(string(out)))
		return
	}
	fmt.Printf("  workspace %q forgotten (work is on main).\n", changeName)
}

// advanceMainToHead moves the `main` bookmark forward to include any committed
// main-line work that is ahead of it (commits in the default workspace's @ line
// that the bookmark has not yet reached). This prevents `jjay merge` from
// silently dropping main-side work created after a spawn (ADR-010).
//
// If there is no such work, it is a no-op. If the ahead-of-bookmark work cannot
// be safely included — i.e. the only thing ahead is the current empty working
// copy with nothing committed — there is nothing to fold and it is also a no-op.
func advanceMainToHead() error {
	// Find the latest non-empty commit between the main bookmark and @ (the
	// default workspace working copy). This is the true head of the main line.
	out, err := exec.Command("jj", "log", "-r", "latest(main..@ & ~empty())",
		"--no-graph", "-T", "commit_id").Output()
	if err != nil {
		return fmt.Errorf("failed to inspect main-line work ahead of bookmark: %w", err)
	}
	head := strings.TrimSpace(string(out))
	if head == "" {
		// Nothing committed ahead of the bookmark — nothing to fold.
		return nil
	}

	fmt.Println("Folding main-line work ahead of the bookmark into main...")
	cmd := exec.Command("jj", "bookmark", "set", "main", "-r", head)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to advance main bookmark to include ahead-of-bookmark work: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func checkWorkspaceExists(changeName string) error {
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == changeName+":" {
			return nil
		}
	}

	return fmt.Errorf("workspace %q does not exist", changeName)
}

func checkWorkspaceEmpty(changeName string) {
	revset := fmt.Sprintf("%s@", changeName)
	out, err := exec.Command("jj", "log", "-r", revset, "--no-graph", "-T", `if(empty, "empty", "has-changes")`).Output()
	if err != nil {
		return
	}

	if strings.TrimSpace(string(out)) == "empty" {
		fmt.Printf("Warning: workspace %q has no changes in its working copy.\n", changeName)
	}
}

func checkConflicts(changeName string) (bool, error) {
	revset := fmt.Sprintf("%s@", changeName)
	out, err := exec.Command("jj", "log", "-r", revset, "--no-graph", "-T", `if(conflict, "conflict", "clean")`).Output()
	if err != nil {
		return false, fmt.Errorf("failed to check for conflicts: %w", err)
	}

	return strings.TrimSpace(string(out)) == "conflict", nil
}
