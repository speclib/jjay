package spawn

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jjay/internal/workspace"
)

// SpawnOptions holds configurable parameters for Spawn.
type SpawnOptions struct {
	Agent         string // agent command template (with {change} and {wsdir} placeholders)
	Session       string // tmux session name (empty = current)
	WorkspaceRoot string // override workspace root (empty = default)
}

// Spawn creates a jj workspace, tmux window, and launches an agent for the given change.
func Spawn(changeName string, opts SpawnOptions) error {
	// Only check TMUX env when not targeting a specific session.
	// When --session is set, we target that session directly.
	if opts.Session == "" {
		if err := checkTmuxSession(); err != nil {
			return err
		}
	}
	if err := checkOpenspecChange(changeName); err != nil {
		return err
	}
	if err := checkWorkspaceNotExists(changeName); err != nil {
		return err
	}
	if err := checkWindowNotExists(changeName, opts.Session); err != nil {
		return err
	}

	wsDir, err := workspace.WorkspaceDir(changeName, opts.WorkspaceRoot)
	if err != nil {
		return err
	}

	// Snapshot uncommitted work by creating a new empty change.
	// This moves all current work to @- (safe, committed) and leaves
	// @ empty so nothing is lost if the main workspace becomes stale.
	if err := snapshotMainWorkspace(); err != nil {
		return err
	}

	if err := createWorkspace(changeName, wsDir); err != nil {
		return err
	}
	if err := createWindow(changeName, opts.Session, wsDir); err != nil {
		return err
	}
	if err := setupPanes(changeName, wsDir, opts); err != nil {
		return err
	}

	fmt.Printf("Spawned workspace for change %q in %s\n", changeName, wsDir)
	fmt.Println("Main workspace is now on a fresh change. Your previous work is in @-.")
	return nil
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

func snapshotMainWorkspace() error {
	cmd := exec.Command("jj", "new")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to snapshot main workspace (jj new): %w", err)
	}
	return nil
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

func checkWindowNotExists(changeName, session string) error {
	wn := workspace.WindowName(changeName)
	args := []string{"list-windows", "-F", "#{window_name}"}
	if session != "" {
		args = append(args, "-t", session)
	}
	out, err := exec.Command("tmux", args...).Output()
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
	// Use @- to get the snapshot created by jj new (contains all files).
	// Using @ would create a child of the empty new change.
	cmd := exec.Command("jj", "workspace", "add", "--name", changeName, "--revision", "@-", wsDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create jj workspace: %w", err)
	}
	return nil
}

func createWindow(changeName, session, wsDir string) error {
	wn := workspace.WindowName(changeName)
	args := []string{"new-window", "-d", "-n", wn, "-c", wsDir}
	if session != "" {
		args = append(args, "-t", session+":")
	}
	cmd := exec.Command("tmux", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux window %q: %w", wn, err)
	}
	return nil
}

// DefaultAgentCommand is the default agent command template.
const DefaultAgentCommand = `claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}`

// resolveAgentCommand substitutes {change} and {wsdir} placeholders in the agent command.
func resolveAgentCommand(template, changeName, wsDir string) string {
	r := strings.NewReplacer("{change}", changeName, "{wsdir}", wsDir)
	return r.Replace(template)
}

func tmuxTarget(session, window string) string {
	if session != "" {
		return session + ":" + window
	}
	return window
}

func setupPanes(changeName, wsDir string, opts SpawnOptions) error {
	wn := workspace.WindowName(changeName)
	target := tmuxTarget(opts.Session, wn)

	// Resolve agent command
	agentTemplate := opts.Agent
	if agentTemplate == "" {
		agentTemplate = DefaultAgentCommand
	}
	agentCmd := resolveAgentCommand(agentTemplate, changeName, wsDir)

	// Left pane: launch agent (window already starts in wsDir via -c flag)
	cmd := exec.Command("tmux", "send-keys", "-t", target, agentCmd, "Enter")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to launch agent in left pane: %w", err)
	}

	// Split to create right pane (starts in wsDir via -c flag)
	cmd = exec.Command("tmux", "split-window", "-h", "-t", target, "-c", wsDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split tmux window: %w", err)
	}

	return nil
}
