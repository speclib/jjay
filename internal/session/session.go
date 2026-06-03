package session

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SessionName returns the tmux session name for a given directory path.
func SessionName(dirPath string) string {
	return "jjay->" + filepath.Base(dirPath)
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

	return nil
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
