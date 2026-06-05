//go:build integration

package integration

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"jjay/internal/spawn"
	"jjay/internal/workspace"
)

// testEnv holds all shared state for an integration test run.
type testEnv struct {
	ChangeName  string // openspec change name (no prefix)
	WSName      string // jj workspace name (verb-prefixed, e.g. app-test-change)
	SessionName string
	ProjectDir  string
	WsRoot      string
	WsDir       string
	FakeAgent   string
	TmpDir      string
}

func (e *testEnv) WindowName() string {
	return workspace.WindowName(e.WSName)
}

// setupTestEnv creates a temp jj repo, openspec change, and tmux session.
// When keepAlive is true, resources are left running for manual inspection
// and the cleanup command is logged. When false, everything is torn down
// automatically when the test finishes.
func setupTestEnv(t *testing.T, keepAlive bool) *testEnv {
	t.Helper()

	requireCmd(t, "tmux")
	requireCmd(t, "jj")
	requireCmd(t, "openspec")

	changeName := "test-change"
	wsName := spawn.ApplyPrefix + changeName
	sessionName := fmt.Sprintf("jjay-test-%d", rand.Intn(100000))

	var tmpDir string
	if keepAlive {
		var err error
		tmpDir, err = os.MkdirTemp("", "jjay-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
	} else {
		tmpDir = t.TempDir()
	}

	projectDir := filepath.Join(tmpDir, "testproject")
	wsRoot := filepath.Join(tmpDir, "workspaces")
	wsDir := filepath.Join(wsRoot, wsName)

	// Find fake-agent.sh
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	fakeAgent := filepath.Join(cwd, "..", "..", "testdata", "fake-agent.sh")
	if _, err := os.Stat(fakeAgent); err != nil {
		t.Fatalf("fake-agent.sh not found at %s: %v", fakeAgent, err)
	}

	// Init jj repo
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}
	runIn(t, projectDir, "jj", "git", "init")

	if err := os.WriteFile(filepath.Join(projectDir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to write README: %v", err)
	}

	// Init openspec
	runIn(t, projectDir, "openspec", "init", "--tools", "none")
	runIn(t, projectDir, "openspec", "new", "change", changeName)

	// Create tmux session
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create tmux session %q: %v", sessionName, err)
	}

	// Chdir to project
	origDir, _ := os.Getwd()
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("failed to chdir to project: %v", err)
	}

	if keepAlive {
		// Only restore cwd; leave tmux session and temp dir alive
		t.Cleanup(func() {
			os.Chdir(origDir)
		})
	} else {
		// Full teardown
		t.Cleanup(func() {
			os.Chdir(origDir)
			exec.Command("tmux", "kill-session", "-t", sessionName).Run()
			cmd := exec.Command("jj", "workspace", "forget", wsName)
			cmd.Dir = projectDir
			cmd.Run()
		})
	}

	return &testEnv{
		ChangeName:  changeName,
		WSName:      wsName,
		SessionName: sessionName,
		ProjectDir:  projectDir,
		WsRoot:      wsRoot,
		WsDir:       wsDir,
		FakeAgent:   fakeAgent,
		TmpDir:      tmpDir,
	}
}

// assertSpawn runs Spawn and verifies all resources were created.
func assertSpawn(t *testing.T, env *testEnv) {
	t.Helper()

	opts := spawn.SpawnOptions{
		Agent:         env.FakeAgent + " {change}",
		Session:       env.SessionName,
		WorkspaceRoot: env.WsRoot,
	}

	if err := spawn.Spawn(env.ChangeName, opts); err != nil {
		t.Fatalf("Spawn() failed: %v", err)
	}

	wn := env.WindowName()

	// tmux window exists
	out, err := exec.Command("tmux", "list-windows", "-t", env.SessionName, "-F", "#{window_name}").Output()
	if err != nil {
		t.Fatalf("failed to list tmux windows: %v", err)
	}
	if !containsLine(string(out), wn) {
		t.Errorf("tmux window %q not found in session %q, got: %s", wn, env.SessionName, out)
	}

	// jj workspace exists (named with the app- prefix)
	assertJJWorkspace(t, env.WSName)

	// workspace directory exists
	assertDirExists(t, env.WsDir)

	// both panes have correct working directory
	assertPaneDir(t, env.SessionName, wn+".0", env.WsDir)
	assertPaneDir(t, env.SessionName, wn+".1", env.WsDir)

	// agent marker file
	markerFile := filepath.Join(env.WsDir, "agent-was-here.txt")
	if !waitForFile(markerFile, 5*time.Second) {
		t.Errorf("agent marker file %q not found after waiting", markerFile)
	}
}

func requireCmd(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s not available, skipping integration test", name)
	}
}

func runIn(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v failed in %s: %v", name, args, dir, err)
	}
}

func containsLine(output, line string) bool {
	for _, l := range strings.Split(strings.TrimSpace(output), "\n") {
		if l == line {
			return true
		}
	}
	return false
}

func waitForFile(path string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

func assertJJWorkspace(t *testing.T, changeName string) {
	t.Helper()
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		t.Fatalf("failed to list jj workspaces: %v", err)
	}
	if !strings.Contains(string(out), changeName+":") {
		t.Errorf("jj workspace %q not found, got: %s", changeName, out)
	}
}

func assertNoJJWorkspace(t *testing.T, changeName string) {
	t.Helper()
	out, _ := exec.Command("jj", "workspace", "list").Output()
	if strings.Contains(string(out), changeName+":") {
		t.Errorf("jj workspace %q still exists", changeName)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("directory %q does not exist", path)
	}
}

func assertPaneDir(t *testing.T, session, pane, expectedDir string) {
	t.Helper()
	target := session + ":" + pane
	out, err := exec.Command("tmux", "display-message", "-p", "-t", target, "#{pane_current_path}").Output()
	if err != nil {
		t.Fatalf("failed to get pane dir for %q: %v", target, err)
	}
	got := strings.TrimSpace(string(out))
	if got != expectedDir {
		t.Errorf("pane %q working dir = %q, want %q", pane, got, expectedDir)
	}
	t.Logf("pane %q working dir OK: %s", pane, got)
}

func assertDirNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("directory %q still exists", path)
	}
}
