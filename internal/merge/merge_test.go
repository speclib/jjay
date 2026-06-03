package merge

import (
	"strings"
	"testing"
)

func TestCheckWorkspaceExists_Parsing(t *testing.T) {
	// checkWorkspaceExists calls jj workspace list and parses its output.
	// We test the parsing logic by calling it with a name that definitely
	// doesn't exist, which exercises the full path including the "not found" branch.
	err := checkWorkspaceExists("nonexistent-workspace-xyz-12345")
	if err == nil {
		t.Fatal("expected error for nonexistent workspace, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent-workspace-xyz-12345") {
		t.Errorf("error should mention workspace name, got: %v", err)
	}
}
