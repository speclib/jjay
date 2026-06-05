//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"jjay/internal/spawn"
	"jjay/internal/status"
	"jjay/internal/workspace"
)

// TestProposalSpawn_NameDiffersFromChange covers the load-bearing ADR-011
// invariant: a proposal spawn's workspace name (prop-<slug>) is NOT the
// openspec change name the agent eventually creates. status must key the spawn
// off the workspace name and classify it as a proposal — it must NOT mis-key it
// by the produced change name, and must not read change-shaped task counts.
func TestProposalSpawn_NameDiffersFromChange(t *testing.T) {
	env := setupTestEnv(t, false)

	prompt := "add dark mode to the settings page"
	slug := workspace.Slug(prompt) // dark-mode-settings-page
	wsName := spawn.ProposalPrefix + slug
	wsDir := filepath.Join(env.WsRoot, wsName)

	opts := spawn.SpawnOptions{
		Agent:         env.FakeAgent + " {prompt}",
		Session:       env.SessionName,
		WorkspaceRoot: env.WsRoot,
	}
	if err := spawn.SpawnProposal(prompt, spawn.ModeExplore, opts); err != nil {
		t.Fatalf("SpawnProposal() failed: %v", err)
	}
	t.Cleanup(func() {
		exec.Command("tmux", "kill-window", "-t", env.SessionName+":"+workspace.WindowName(wsName)).Run()
		exec.Command("jj", "workspace", "forget", wsName).Run()
	})

	// jj workspace is named by the slug, prefixed — not by any change name.
	assertJJWorkspace(t, wsName)
	assertDirExists(t, wsDir)

	// Simulate the agent inventing a DIFFERENTLY-named openspec change inside
	// the isolated workspace.
	producedChange := "add-dark-mode" // intentionally != slug
	changeDir := filepath.Join(wsDir, "openspec", "changes", producedChange)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create produced change dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(changeDir, "tasks.md"), []byte("- [ ] do it\n"), 0o644); err != nil {
		t.Fatalf("failed to write produced tasks.md: %v", err)
	}

	// status must surface the spawn keyed on the workspace name, classified as a
	// proposal, with no inferred change name and no change-shaped task counts.
	spawns, _, err := status.List(env.SessionName, env.WsRoot)
	if err != nil {
		t.Fatalf("status.List() failed: %v", err)
	}

	var found *status.Spawn
	for i := range spawns {
		if spawns[i].Name == wsName {
			found = &spawns[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("proposal spawn %q not found in status; got %+v", wsName, spawns)
	}
	if found.Kind != status.KindProposal {
		t.Errorf("spawn %q Kind = %v, want proposal", wsName, found.Kind)
	}
	if found.Change != "" {
		t.Errorf("proposal spawn must not infer a change name, got %q", found.Change)
	}
	if found.Tasks.Found {
		t.Errorf("proposal spawn must not read change-shaped tasks, got %+v", found.Tasks)
	}
	// And it must NOT be mis-keyed under the produced change name.
	for _, s := range spawns {
		if s.Name == producedChange || s.Change == producedChange {
			t.Errorf("spawn mis-keyed by produced change name %q: %+v", producedChange, s)
		}
	}
}
