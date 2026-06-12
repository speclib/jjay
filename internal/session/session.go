package session

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jjay/internal/spawn"
	"jjay/internal/status"
)

// SessionName returns the tmux session name for a given directory path.
//
// tmux treats `.` and `:` specially in target names (`session:window.pane`), so
// a dir like `mip.rs` would create a session tmux stores as `jjay->mip_rs` but
// then fail to target as `jjay->mip.rs` ("can't find pane: rs"). Normalize those
// characters to `_` so creation and every later -t target agree.
func SessionName(dirPath string) string {
	return "jjay->" + sanitizeTmuxName(filepath.Base(dirPath))
}

// sanitizeTmuxName replaces characters tmux reserves in target names (`.` and
// `:`) with `_`, matching tmux's own normalization of `.` on session creation.
func sanitizeTmuxName(name string) string {
	return strings.NewReplacer(".", "_", ":", "_").Replace(name)
}

// Open creates a tmux session for the given jj repo path and switches to it.
func Open(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := checkJJRepo(absPath); err != nil {
		return err
	}

	sessionName := SessionName(absPath)

	if err := checkSessionNotExists(sessionName); err != nil {
		return err
	}

	if err := createSession(sessionName, absPath); err != nil {
		return err
	}

	if err := switchClient(sessionName); err != nil {
		return err
	}

	reopenSpawns(absPath, sessionName, os.Stdout)

	return nil
}

// openWindowFunc reopens a tmux window + RESUMES the agent for a detached spawn.
// Matches spawn.Reopen (resume intent — must not re-run /opsx:apply); injectable
// for testing.
type openWindowFunc func(name, wsDir string, opts spawn.SpawnOptions) error

// reopenSpawns recreates a tmux window for every detached spawn of the repo at
// repoRoot (the session's target repo — NOT the directory jjay runs from). It is
// best-effort and non-fatal (ADR-003 / ADR-006): a per-spawn failure is logged
// and the rest continue; session-open still succeeds. Scoping to repoRoot is
// what prevents reopening another project's workspaces into this session.
func reopenSpawns(repoRoot, sessionName string, out io.Writer) {
	spawns, _, err := status.ListIn(repoRoot, sessionName, "")
	if err != nil {
		// Can't enumerate workspaces — nothing to reopen, not fatal.
		fmt.Fprintf(out, "session-open: could not enumerate spawns to reopen: %v\n", err)
		return
	}
	reopenDetached(spawns, sessionName, spawn.Reopen, out)
}

// reopenDetached acts on the detached subset of spawns using open(). Spawns that
// already have a window (Attached) are skipped — no duplicates. Pure aside from
// the injected open func, so it is unit-testable.
func reopenDetached(spawns []status.Spawn, sessionName string, open openWindowFunc, out io.Writer) {
	var failed []string
	for _, s := range spawns {
		if s.Attached {
			continue // window already exists; do not duplicate
		}
		opts := spawn.SpawnOptions{Session: sessionName}
		// Reopen by the (prefixed) workspace name, which keys the jj workspace,
		// directory, and window — not the openspec change name, which a proposal
		// spawn does not have (ADR-011).
		if err := open(s.Name, s.WSDir, opts); err != nil {
			fmt.Fprintf(out, "session-open: could not reopen spawn %q: %v\n", s.Name, err)
			failed = append(failed, s.Name)
			continue
		}
		fmt.Fprintf(out, "session-open: reopened spawn %q\n", s.Name)
	}
	if len(failed) > 0 {
		fmt.Fprintf(out, "session-open: %d spawn(s) failed to reopen: %s\n", len(failed), strings.Join(failed, ", "))
	}
}

func checkJJRepo(absPath string) error {
	info, err := os.Stat(filepath.Join(absPath, ".jj"))
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s is not a jj repository (no .jj/ directory)", absPath)
	}
	return nil
}

func checkSessionNotExists(sessionName string) error {
	out, err := exec.Command("tmux", "list-sessions", "-F", "#{session_name}").Output()
	if err != nil {
		// tmux may fail if no server is running — that means no sessions exist
		return nil
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == sessionName {
			return fmt.Errorf("tmux session %q already exists", sessionName)
		}
	}
	return nil
}

func createSession(sessionName, absPath string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", absPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session %q: %w", sessionName, err)
	}
	return nil
}

func switchClient(sessionName string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to switch to tmux session %q: %w", sessionName, err)
	}
	return nil
}
