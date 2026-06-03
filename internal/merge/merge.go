package merge

import (
	"fmt"
	"os/exec"
	"strings"
)

// Merge merges a workspace's work into the main bookmark.
// It creates a merge commit, moves the main bookmark, and creates a fresh change.
func Merge(changeName string) error {
	// 1. Verify workspace exists
	if err := checkWorkspaceExists(changeName); err != nil {
		return err
	}

	// 2. Warn if workspace is empty
	checkWorkspaceEmpty(changeName)

	// 3. Create merge commit
	mergeMsg := fmt.Sprintf("merge %s into main", changeName)
	revset := changeName + "@"
	cmd := exec.Command("jj", "new", "main", revset, "-m", mergeMsg)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create merge commit: %s", strings.TrimSpace(string(out)))
	}

	// 4. Move main bookmark
	cmd = exec.Command("jj", "bookmark", "set", "main", "-r", "@")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update main bookmark: %s", strings.TrimSpace(string(out)))
	}

	// 5. Fresh change
	cmd = exec.Command("jj", "new")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create fresh change: %s", strings.TrimSpace(string(out)))
	}

	// 6. Success message
	fmt.Printf("Merged %s into main.\n", changeName)
	fmt.Println("Main bookmark updated. You're on a fresh change.")
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
