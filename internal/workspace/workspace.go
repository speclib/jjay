package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// WindowName returns the tmux window name for a given change.
func WindowName(changeName string) string {
	return "ws-" + changeName
}

// WorkspaceDir returns the absolute path for the workspace directory:
// ../<project-name>-workspaces/<change-name>
func WorkspaceDir(changeName string) (string, error) {
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
