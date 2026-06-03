//go:build integration

package integration

import (
	"os/exec"
	"testing"

	"jjay/internal/cleanup"
)

// TestFullLifecycle tests spawn → verify → cleanup → verify as subtests
// sharing a single environment.
func TestFullLifecycle(t *testing.T) {
	env := setupTestEnv(t, false)

	t.Run("spawn", func(t *testing.T) {
		assertSpawn(t, env)
	})

	t.Run("cleanup", func(t *testing.T) {
		cleanupOpts := cleanup.CleanupOptions{
			Session:       env.SessionName,
			WorkspaceRoot: env.WsRoot,
		}
		if err := cleanup.Cleanup(env.ChangeName, cleanupOpts); err != nil {
			t.Fatalf("Cleanup() failed: %v", err)
		}
		assertCleanedUp(t, env)
	})
}

// TestSpawn runs spawn and leaves resources alive for manual inspection.
// Use `make test-spawn` to run this test in isolation.
func TestSpawn(t *testing.T) {
	env := setupTestEnv(t, true)
	assertSpawn(t, env)

	t.Log("Spawn succeeded. Resources left alive for inspection:")
	t.Logf("  tmux session:  %s", env.SessionName)
	t.Logf("  workspace dir: %s", env.WsDir)
	t.Logf("  project dir:   %s", env.ProjectDir)
	t.Log("")
	t.Logf("To clean up manually:")
	t.Logf("  tmux kill-session -t %s", env.SessionName)
	t.Logf("  rm -rf %s", env.TmpDir)
}

// assertCleanedUp verifies all resources were removed.
func assertCleanedUp(t *testing.T, env *testEnv) {
	t.Helper()

	wn := env.WindowName()

	out, _ := exec.Command("tmux", "list-windows", "-t", env.SessionName, "-F", "#{window_name}").Output()
	if containsLine(string(out), wn) {
		t.Errorf("tmux window %q still exists after cleanup", wn)
	}

	assertNoJJWorkspace(t, env.ChangeName)
	assertDirNotExists(t, env.WsDir)
}
