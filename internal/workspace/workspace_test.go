package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWindowName(t *testing.T) {
	tests := []struct {
		change string
		want   string
	}{
		{"feat-payments", "ws-feat-payments"},
		{"add-auth", "ws-add-auth"},
		{"x", "ws-x"},
	}
	for _, tt := range tests {
		got := WindowName(tt.change)
		if got != tt.want {
			t.Errorf("WindowName(%q) = %q, want %q", tt.change, got, tt.want)
		}
	}
}

func TestWorkspaceDir(t *testing.T) {
	dir, err := WorkspaceDir("feat-payments")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be absolute
	if !filepath.IsAbs(dir) {
		t.Errorf("WorkspaceDir() returned relative path: %q", dir)
	}

	// Should end with <project>-workspaces/feat-payments
	cwd, _ := os.Getwd()
	projectName := filepath.Base(cwd)
	wantSuffix := filepath.Join(projectName+"-workspaces", "feat-payments")
	if !strings.HasSuffix(dir, wantSuffix) {
		t.Errorf("WorkspaceDir() = %q, want suffix %q", dir, wantSuffix)
	}
}
