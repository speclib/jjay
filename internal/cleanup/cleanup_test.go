package cleanup

import (
	"os"
	"path/filepath"
	"testing"
)

// tempWorkspaceRoot returns a temp directory to use as an explicit WorkspaceRoot
// so tests never touch the real ../<project>-workspaces tree.
func tempWorkspaceRoot(t *testing.T) string {
	t.Helper()
	root, err := os.MkdirTemp("", "cleanup-test-")
	if err != nil {
		t.Fatalf("os.MkdirTemp() failed: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(root) })
	return root
}

func TestRemoveDirectory_RemovesExisting(t *testing.T) {
	root := tempWorkspaceRoot(t)
	change := "feat-payments"

	// removeDirectory resolves <root>/<change>; create that so it exists.
	wsDir := filepath.Join(root, change)
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatalf("failed to create workspace dir: %v", err)
	}

	removeDirectory(change, root)

	if _, err := os.Stat(wsDir); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed, stat err = %v", wsDir, err)
	}
}

func TestRemoveDirectory_TolerantOfMissing(t *testing.T) {
	root := tempWorkspaceRoot(t)
	change := "never-created-xyz"
	wsDir := filepath.Join(root, change)

	// Must not panic; the directory simply does not exist.
	removeDirectory(change, root)

	if _, err := os.Stat(wsDir); !os.IsNotExist(err) {
		t.Errorf("expected %s to remain absent, stat err = %v", wsDir, err)
	}
}

func TestTmuxTarget(t *testing.T) {
	if got := tmuxTarget("", "ws-foo"); got != "ws-foo" {
		t.Errorf("tmuxTarget(\"\", \"ws-foo\") = %q, want %q", got, "ws-foo")
	}
	if got := tmuxTarget("mysess", "ws-foo"); got != "mysess:ws-foo" {
		t.Errorf("tmuxTarget(\"mysess\", \"ws-foo\") = %q, want %q", got, "mysess:ws-foo")
	}
}

func TestKillWindow_TolerantOfMissing(t *testing.T) {
	// For a change that does not exist, killWindow returns (no panic) whether
	// or not a tmux server is present: absent server hits the Output() error
	// branch, present server hits the "not found" branch.
	killWindow("nonexistent-change-xyz-12345", "")
}

func TestForgetWorkspace_TolerantOfMissing(t *testing.T) {
	// Same tolerance contract as killWindow for jj.
	forgetWorkspace("nonexistent-change-xyz-12345")
}

func TestCleanup_AllMissingReturnsNil(t *testing.T) {
	root := tempWorkspaceRoot(t)
	err := Cleanup("nonexistent-change-xyz-12345", CleanupOptions{WorkspaceRoot: root})
	if err != nil {
		t.Errorf("Cleanup() with all-missing resources = %v, want nil", err)
	}
}
