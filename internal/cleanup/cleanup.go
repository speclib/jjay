package cleanup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"jjay/internal/workspace"
)

// Cleanup tears down a spawned workspace: kills tmux window, forgets jj workspace,
// removes workspace directory. Each step is tolerant — missing resources are skipped.
func Cleanup(changeName string) error {
	fmt.Printf("Cleaning up change %q...\n", changeName)

	// Order: tmux → jj → directory (kill agent first, then clean up state)
	killWindow(changeName)
	forgetWorkspace(changeName)
	removeDirectory(changeName)

	return nil
}

func killWindow(changeName string) {
	wn := workspace.WindowName(changeName)

	// Check if window exists
	out, err := exec.Command("tmux", "list-windows", "-F", "#{window_name}").Output()
	if err != nil {
		fmt.Printf("  tmux window %s: not found, skipped\n", wn)
		return
	}

	found := false
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == wn {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("  tmux window %s: not found, skipped\n", wn)
		return
	}

	cmd := exec.Command("tmux", "kill-window", "-t", wn)
	if err := cmd.Run(); err != nil {
		fmt.Printf("  tmux window %s: failed to kill (%v), skipped\n", wn, err)
		return
	}

	fmt.Printf("  tmux window %s: killed\n", wn)
}

func forgetWorkspace(changeName string) {
	// Check if workspace exists
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		fmt.Printf("  jj workspace %s: not found, skipped\n", changeName)
		return
	}

	found := false
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == changeName+":" {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("  jj workspace %s: not found, skipped\n", changeName)
		return
	}

	cmd := exec.Command("jj", "workspace", "forget", changeName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("  jj workspace %s: failed to forget (%v), skipped\n", changeName, err)
		return
	}

	fmt.Printf("  jj workspace %s: forgotten\n", changeName)
}

func removeDirectory(changeName string) {
	wsDir, err := workspace.WorkspaceDir(changeName)
	if err != nil {
		fmt.Printf("  workspace directory: failed to resolve path (%v), skipped\n", err)
		return
	}

	if _, err := os.Stat(wsDir); os.IsNotExist(err) {
		fmt.Printf("  workspace directory: not found, skipped\n")
		return
	}

	if err := os.RemoveAll(wsDir); err != nil {
		fmt.Printf("  workspace directory: failed to remove (%v), skipped\n", err)
		return
	}

	fmt.Printf("  workspace directory: removed\n")
}
