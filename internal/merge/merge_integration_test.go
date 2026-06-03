//go:build integration

package merge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// testRepo sets up a temp jj repo with a main bookmark and returns the path.
// Caller must defer cleanup.
func testRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "jjay-merge-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Init jj repo (git-backed)
	run(t, dir, "jj", "git", "init")
	// Create initial file and describe
	writeFile(t, dir, "initial.txt", "initial content")
	run(t, dir, "jj", "describe", "-m", "initial commit")
	// Set main bookmark
	run(t, dir, "jj", "bookmark", "create", "main", "-r", "@")
	// Create fresh change so main is immutable-ish
	run(t, dir, "jj", "new")

	return dir
}

// createWorkspace creates a jj workspace with the given name in a subdir.
func createWorkspace(t *testing.T, repoDir, name string) string {
	t.Helper()
	wsDir := filepath.Join(repoDir, name+"-ws")
	run(t, repoDir, "jj", "workspace", "add", "--name", name, wsDir)
	return wsDir
}

func run(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s %s\n%s\n%v", name, strings.Join(args, " "), string(out), err)
	}
	return strings.TrimSpace(string(out))
}

func runMayFail(t *testing.T, dir string, name string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", name, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func fileExists(t *testing.T, dir, name string) bool {
	t.Helper()
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

func readFile(t *testing.T, dir, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatalf("failed to read %s: %v", name, err)
	}
	return string(data)
}

// mergeInRepo runs the merge logic against a workspace in the given repo dir.
func mergeInRepo(t *testing.T, repoDir, changeName string) error {
	t.Helper()
	// We need to run jj commands from the repo dir, so temporarily chdir
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", repoDir, err)
	}
	defer os.Chdir(origDir)

	return Merge(changeName)
}

func TestMerge_CleanNoMainChanges(t *testing.T) {
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace and add a file
	wsDir := createWorkspace(t, dir, "feat")
	writeFile(t, wsDir, "feature.txt", "new feature")
	run(t, wsDir, "jj", "describe", "-m", "add feature")

	// Merge
	if err := mergeInRepo(t, dir, "feat"); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// Verify workspace file is present in main
	if !fileExists(t, dir, "feature.txt") {
		t.Error("feature.txt should exist after merge")
	}
	if !fileExists(t, dir, "initial.txt") {
		t.Error("initial.txt should still exist after merge")
	}
}

func TestMerge_MainMovedNoOverlap(t *testing.T) {
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace
	wsDir := createWorkspace(t, dir, "feat")

	// Add file in workspace
	writeFile(t, wsDir, "ws-file.txt", "workspace content")
	run(t, wsDir, "jj", "describe", "-m", "add ws-file")

	// Move main forward with a different file
	writeFile(t, dir, "main-file.txt", "main content")
	run(t, dir, "jj", "describe", "-m", "add main-file")
	run(t, dir, "jj", "bookmark", "set", "main", "-r", "@")
	run(t, dir, "jj", "new")

	// Merge
	if err := mergeInRepo(t, dir, "feat"); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// Both files should exist
	if !fileExists(t, dir, "ws-file.txt") {
		t.Error("ws-file.txt should exist after merge")
	}
	if !fileExists(t, dir, "main-file.txt") {
		t.Error("main-file.txt should exist after merge")
	}
}

func TestMerge_SameFileModified(t *testing.T) {
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace
	wsDir := createWorkspace(t, dir, "feat")

	// Modify same file in workspace
	writeFile(t, wsDir, "initial.txt", "workspace version")
	run(t, wsDir, "jj", "describe", "-m", "modify initial.txt in workspace")

	// Modify same file on main
	writeFile(t, dir, "initial.txt", "main version")
	run(t, dir, "jj", "describe", "-m", "modify initial.txt on main")
	run(t, dir, "jj", "bookmark", "set", "main", "-r", "@")
	run(t, dir, "jj", "new")

	// Merge should fail with conflict
	err := mergeInRepo(t, dir, "feat")
	if err == nil {
		t.Fatal("merge should have failed due to conflict")
	}
	if !strings.Contains(err.Error(), "conflict") {
		t.Errorf("error should mention conflict, got: %v", err)
	}
}

func TestMerge_WorkspaceAddsNewFiles(t *testing.T) {
	// THIS IS THE CRITICAL BUG FIX TEST
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace
	wsDir := createWorkspace(t, dir, "feat")

	// Workspace adds new files
	writeFile(t, wsDir, "blog-post.md", "# New blog post")
	writeFile(t, wsDir, "tasks-checked.md", "- [x] done")
	run(t, wsDir, "jj", "describe", "-m", "add blog post and checked tasks")

	// Main also moves forward with different files
	writeFile(t, dir, "main-change.txt", "main moved forward")
	run(t, dir, "jj", "describe", "-m", "main change")
	run(t, dir, "jj", "bookmark", "set", "main", "-r", "@")
	run(t, dir, "jj", "new")

	// Merge
	if err := mergeInRepo(t, dir, "feat"); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// ALL files should exist — this is the bug we're fixing
	if !fileExists(t, dir, "blog-post.md") {
		t.Error("blog-post.md should exist after merge — THIS WAS THE BUG")
	}
	if !fileExists(t, dir, "tasks-checked.md") {
		t.Error("tasks-checked.md should exist after merge")
	}
	if !fileExists(t, dir, "main-change.txt") {
		t.Error("main-change.txt should exist after merge")
	}
	if !fileExists(t, dir, "initial.txt") {
		t.Error("initial.txt should still exist after merge")
	}
}

func TestMerge_EmptyWorkspace(t *testing.T) {
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace but don't add any files
	_ = createWorkspace(t, dir, "feat")

	// Merge should succeed with warning (empty workspace)
	if err := mergeInRepo(t, dir, "feat"); err != nil {
		t.Fatalf("merge should succeed even with empty workspace: %v", err)
	}
}

func TestMerge_MultipleWorkspaceCommits(t *testing.T) {
	dir := testRepo(t)
	defer os.RemoveAll(dir)

	// Create workspace
	wsDir := createWorkspace(t, dir, "feat")

	// Multiple commits in workspace
	writeFile(t, wsDir, "file1.txt", "first commit")
	run(t, wsDir, "jj", "describe", "-m", "first change")
	run(t, wsDir, "jj", "new")

	writeFile(t, wsDir, "file2.txt", "second commit")
	run(t, wsDir, "jj", "describe", "-m", "second change")
	run(t, wsDir, "jj", "new")

	writeFile(t, wsDir, "file3.txt", "third commit")
	run(t, wsDir, "jj", "describe", "-m", "third change")

	// Merge
	if err := mergeInRepo(t, dir, "feat"); err != nil {
		t.Fatalf("merge failed: %v", err)
	}

	// All files from all commits should exist
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("file%d.txt", i)
		if !fileExists(t, dir, name) {
			t.Errorf("%s should exist after merge", name)
		}
	}
}
