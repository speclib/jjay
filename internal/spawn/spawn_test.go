package spawn

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestWorkspaceDir(t *testing.T) {
	dir, err := workspaceDir("feat-payments")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)
	wantSuffix := filepath.Join(projectName+"-workspaces", "feat-payments")
	if !strings.HasSuffix(dir, wantSuffix) {
		t.Errorf("workspaceDir() = %q, want suffix %q", dir, wantSuffix)
	}
}

func TestWindowName(t *testing.T) {
	tests := []struct {
		change string
		want   string
	}{
		{"feat-payments", "ws-feat-payments"},
		{"add-auth", "ws-add-auth"},
	}
	for _, tt := range tests {
		got := windowName(tt.change)
		if got != tt.want {
			t.Errorf("windowName(%q) = %q, want %q", tt.change, got, tt.want)
		}
	}
}
