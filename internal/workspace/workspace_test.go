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
	dir, err := WorkspaceDir("feat-payments", "")
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

func TestWorkspaceDirFrom_AnchoredOnMainRoot(t *testing.T) {
	// Anchored on the given main root, independent of cwd: parent of mainRoot,
	// then <mainName>-workspaces/<change>.
	got, err := WorkspaceDirFrom("/home/u/proj", "feat-x", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "/home/u/proj-workspaces/feat-x"
	if got != want {
		t.Errorf("WorkspaceDirFrom() = %q, want %q", got, want)
	}
}

func TestWorkspaceDirFrom_ExplicitRootWins(t *testing.T) {
	got, err := WorkspaceDirFrom("/home/u/proj", "feat-x", "/custom/root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom/root/feat-x" {
		t.Errorf("WorkspaceDirFrom() = %q, want /custom/root/feat-x", got)
	}
}

func TestRelativeToMain(t *testing.T) {
	got := RelativeToMain("/repo/main", "/repo/ws/add-foo")
	if got != "../ws/add-foo" {
		t.Errorf("RelativeToMain() = %q, want ../ws/add-foo", got)
	}
}

func TestMainRepoRoot_MainWorkingCopy(t *testing.T) {
	// .jj/repo is a real directory → this dir is the main repo root.
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jj", "repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := MainRepoRoot(root)
	if err != nil {
		t.Fatalf("MainRepoRoot: %v", err)
	}
	// t.TempDir may contain symlinks (e.g. /var → /private/var); compare cleaned.
	wantResolved, _ := filepath.EvalSymlinks(root)
	gotResolved, _ := filepath.EvalSymlinks(got)
	if gotResolved != wantResolved {
		t.Errorf("MainRepoRoot() = %q, want %q", got, root)
	}
}

func TestMainRepoRoot_ChildWorkspace(t *testing.T) {
	// Layout: <base>/main (real .jj/repo dir) and <base>/ws (child workspace
	// whose .jj/repo file points at ../main/.jj/repo).
	base := t.TempDir()
	main := filepath.Join(base, "main")
	child := filepath.Join(base, "ws")
	if err := os.MkdirAll(filepath.Join(main, ".jj", "repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(child, ".jj"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Pointer is relative to the child's .jj dir.
	if err := os.WriteFile(filepath.Join(child, ".jj", "repo"), []byte("../../main/.jj/repo\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := MainRepoRoot(child)
	if err != nil {
		t.Fatalf("MainRepoRoot: %v", err)
	}
	wantResolved, _ := filepath.EvalSymlinks(main)
	gotResolved, _ := filepath.EvalSymlinks(got)
	if gotResolved != wantResolved {
		t.Errorf("MainRepoRoot(child) = %q, want %q", got, main)
	}
}
