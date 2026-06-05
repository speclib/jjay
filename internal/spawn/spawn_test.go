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

func TestResolveAgentCommand(t *testing.T) {
	// Apply template uses {change}; the seed fills it.
	got := resolveAgentCommand(DefaultAgentCommand, "add-foo", "/ws/app-add-foo")
	want := `claude "/opsx:apply add-foo" --dangerously-skip-permissions --add-dir /ws/app-add-foo`
	if got != want {
		t.Errorf("resolveAgentCommand(apply) = %q, want %q", got, want)
	}

	// Proposal explore template uses {prompt}; the seed fills it.
	got = resolveAgentCommand(proposalExploreCommand, "dark mode", "/ws/prop-dark-mode")
	want = `claude "/opsx:explore dark mode" --dangerously-skip-permissions --add-dir /ws/prop-dark-mode`
	if got != want {
		t.Errorf("resolveAgentCommand(explore) = %q, want %q", got, want)
	}
}

func TestProposalAgentCommand(t *testing.T) {
	if proposalAgentCommand(ModeExplore) != proposalExploreCommand {
		t.Error("explore mode should map to the explore template")
	}
	if proposalAgentCommand(ModePropose) != proposalProposeCommand {
		t.Error("propose mode should map to the propose template")
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
