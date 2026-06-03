package spawn

import (
	"os"
	"testing"

	"jjay/internal/workspace"
)

func TestCheckTmuxSession_InsideTmux(t *testing.T) {
	original := os.Getenv("TMUX")
	defer os.Setenv("TMUX", original)

	os.Setenv("TMUX", "/tmp/tmux-1000/default,12345,0")
	if err := checkTmuxSession(); err != nil {
		t.Errorf("expected no error inside tmux, got: %v", err)
	}
}

func TestCheckTmuxSession_OutsideTmux(t *testing.T) {
	original := os.Getenv("TMUX")
	defer os.Setenv("TMUX", original)

	os.Unsetenv("TMUX")
	if err := checkTmuxSession(); err == nil {
		t.Error("expected error outside tmux, got nil")
	}
}

func TestWorkspacePackageIntegration(t *testing.T) {
	// Verify spawn can use workspace package functions
	wn := workspace.WindowName("feat-payments")
	if wn != "ws-feat-payments" {
		t.Errorf("workspace.WindowName() = %q, want %q", wn, "ws-feat-payments")
	}

	_, err := workspace.WorkspaceDir("feat-payments", "")
	if err != nil {
		t.Fatalf("workspace.WorkspaceDir() unexpected error: %v", err)
	}
}
