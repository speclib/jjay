//go:build integration

package integration

import (
	"os/exec"
	"strings"
	"testing"

	"jjay/internal/cleanup"
	"jjay/internal/status"
)

// TestFullLifecycle tests spawn → verify → cleanup → verify as subtests
// sharing a single environment.
func TestFullLifecycle(t *testing.T) {
	env := setupTestEnv(t, false)

	t.Run("spawn", func(t *testing.T) {
		assertSpawn(t, env)
	})

	t.Run("status", func(t *testing.T) {
		assertStatus(t, env)
	})

	t.Run("cleanup", func(t *testing.T) {
		cleanupOpts := cleanup.CleanupOptions{
			Session:       env.SessionName,
			WorkspaceRoot: env.WsRoot,
		}
		// Cleanup keys off the workspace name (verb-prefixed), not the change
		// name — this is what `status` surfaces and `/jjay:cleanup` is given.
		if err := cleanup.Cleanup(env.WSName, cleanupOpts); err != nil {
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

// assertStatus runs status.List + Render against the live spawned environment
// and verifies the spawn shows up, the table carries the TMUX and MERGED
// columns, and the freshly-spawned (unmerged) workspace reports MERGED=no.
func assertStatus(t *testing.T, env *testEnv) {
	t.Helper()

	spawns, mainRoot, err := status.List(env.SessionName, env.WsRoot)
	if err != nil {
		t.Fatalf("status.List() failed: %v", err)
	}

	var found *status.Spawn
	for i := range spawns {
		if spawns[i].Change == env.ChangeName {
			found = &spawns[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("spawn %q not found in status output: %+v", env.ChangeName, spawns)
	}

	// A freshly spawned workspace has its own commit ahead of main: not merged.
	if found.Merged {
		t.Errorf("fresh spawn %q should report MERGED=no, got Merged=true", env.ChangeName)
	}

	var b strings.Builder
	status.Render(&b, mainRoot, spawns)
	out := b.String()

	// In the two-table layout the CHANGES table's column header follows the
	// "CHANGES" title line, so find the header row rather than assuming line 1.
	var header string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "CHANGE ") || strings.HasPrefix(line, "CHANGE\t") {
			header = line
			break
		}
	}
	if header == "" {
		t.Fatalf("could not find CHANGES column header in status output:\n%s", out)
	}
	for _, col := range []string{"TMUX", "MERGED"} {
		if !strings.Contains(header, col) {
			t.Errorf("status header missing %q column: %q", col, header)
		}
	}
	if !strings.Contains(out, env.ChangeName) {
		t.Errorf("status output missing spawn %q:\n%s", env.ChangeName, out)
	}
}

// assertCleanedUp verifies all resources were removed.
func assertCleanedUp(t *testing.T, env *testEnv) {
	t.Helper()

	wn := env.WindowName()

	out, _ := exec.Command("tmux", "list-windows", "-t", env.SessionName, "-F", "#{window_name}").Output()
	if containsLine(string(out), wn) {
		t.Errorf("tmux window %q still exists after cleanup", wn)
	}

	assertNoJJWorkspace(t, env.WSName)
	assertDirNotExists(t, env.WsDir)
}
