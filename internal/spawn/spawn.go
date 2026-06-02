package spawn

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Spawn creates a jj workspace, tmux window, and launches an agent for the given change.
func Spawn(changeName string) error {
	if err := checkTmuxSession(); err != nil {
		return err
	}
	if err := checkOpenspecChange(changeName); err != nil {
		return err
	}
	if err := checkWorkspaceNotExists(changeName); err != nil {
		return err
	}
	if err := checkWindowNotExists(changeName); err != nil {
		return err
	}

	wsDir, err := workspaceDir(changeName)
	if err != nil {
		return err
	}

	if err := createWorkspace(changeName, wsDir); err != nil {
		return err
	}
	if err := createWindow(changeName); err != nil {
		return err
	}
	if err := setupPanes(changeName, wsDir); err != nil {
		return err
	}

	fmt.Printf("Spawned workspace for change %q in %s\n", changeName, wsDir)
	return nil
}

func windowName(changeName string) string {
	return "ws-" + changeName
}

// workspaceDir returns the absolute path for the workspace directory:
// ../<project-name>-workspaces/<change-name>
func workspaceDir(changeName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	projectName := filepath.Base(cwd)
	absPath, err := filepath.Abs(filepath.Join(cwd, "..", projectName+"-workspaces", changeName))
	if err != nil {
		return "", fmt.Errorf("failed to resolve workspace path: %w", err)
	}
	return absPath, nil
}

// CheckTmuxSession verifies we're running inside a tmux session.
func checkTmuxSession() error {
	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("jjay must be run inside a tmux session")
	}
	return nil
}

type openspecChange struct {
	Name string `json:"name"`
}

type openspecList struct {
	Changes []openspecChange `json:"changes"`
}

func checkOpenspecChange(changeName string) error {
	out, err := exec.Command("openspec", "list", "--json").Output()
	if err != nil {
		return fmt.Errorf("failed to list openspec changes: %w", err)
	}

	var list openspecList
	if err := json.Unmarshal(out, &list); err != nil {
		return fmt.Errorf("failed to parse openspec output: %w", err)
	}

	for _, c := range list.Changes {
		if c.Name == changeName {
			return nil
		}
	}
	return fmt.Errorf("openspec change %q does not exist", changeName)
}

func checkWorkspaceNotExists(changeName string) error {
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return fmt.Errorf("failed to list jj workspaces: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == changeName+":" {
			return fmt.Errorf("jj workspace %q already exists", changeName)
		}
	}
	return nil
}

func checkWindowNotExists(changeName string) error {
	wn := windowName(changeName)
	out, err := exec.Command("tmux", "list-windows", "-F", "#{window_name}").Output()
	if err != nil {
		return fmt.Errorf("failed to list tmux windows: %w", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == wn {
			return fmt.Errorf("tmux window %q already exists", wn)
		}
	}
	return nil
}

func createWorkspace(changeName, wsDir string) error {
	if err := os.MkdirAll(filepath.Dir(wsDir), 0o755); err != nil {
		return fmt.Errorf("failed to create workspace parent directory: %w", err)
	}
	// Base the new workspace on @ so it includes uncommitted files
	// (e.g., the active openspec change directory)
	cmd := exec.Command("jj", "workspace", "add", "--name", changeName, "--revision", "@", wsDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create jj workspace: %w", err)
	}
	return nil
}

func createWindow(changeName string) error {
	wn := windowName(changeName)
	cmd := exec.Command("tmux", "new-window", "-d", "-n", wn)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux window %q: %w", wn, err)
	}
	return nil
}

func setupPanes(changeName, wsDir string) error {
	wn := windowName(changeName)

	// Left pane: cd to workspace and launch claude agent
	// Use --add-dir to grant access to workspace dir so claude trusts it
	agentCmd := fmt.Sprintf(
		"cd %s && claude \"/opsx:apply %s\" --dangerously-skip-permissions --add-dir %s",
		wsDir, changeName, wsDir,
	)
	cmd := exec.Command("tmux", "send-keys", "-t", wn, agentCmd, "Enter")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch agent in left pane: %w", err)
	}

	// Split to create right pane
	cmd = exec.Command("tmux", "split-window", "-h", "-t", wn)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split tmux window: %w", err)
	}

	// Right pane: cd to workspace
	cmd = exec.Command("tmux", "send-keys", "-t", wn+".1", "cd "+wsDir, "Enter")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set up shell pane: %w", err)
	}

	return nil
}
