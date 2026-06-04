package merge

import (
	"fmt"
	"os/exec"
	"strings"
)

// Merge merges a workspace's work into the main bookmark.
// It rebases the workspace onto main first to prevent silent file drops,
// then creates a merge commit, moves the main bookmark, and creates a fresh change.
func Merge(changeName string) error {
	// 1. Verify workspace exists
	if err := checkWorkspaceExists(changeName); err != nil {
		return err
	}

	// 2. Warn if workspace is empty
	checkWorkspaceEmpty(changeName)

	// 2b. Fold main-line work that is ahead of the `main` bookmark into the
	// bookmark before merging. The orchestrator often commits work in the main
	// working copy (new proposals, bean edits) that lands in commits ahead of
	// the `main` bookmark. Merging onto the lagging bookmark would leave that
	// work in neither merge parent and silently drop it (see ADR-010, jjay-ug7y).
	if err := advanceMainToHead(); err != nil {
		return err
	}

	// 3. Rebase workspace branch onto main
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

	// 7. Fresh change
	cmd = exec.Command("jj", "new")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create fresh change: %s", strings.TrimSpace(string(out)))
	}

	// 8. Success message
	fmt.Printf("Rebased and merged %s into main.\n", changeName)
	fmt.Println("Main bookmark updated. You're on a fresh change.")
	return nil
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
