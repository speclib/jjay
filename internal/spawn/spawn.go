package spawn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jjay/internal/openspec"
	"jjay/internal/workspace"
)

// Name prefixes encode the spawn kind in the workspace/window name so `status`
// (and a human reading `jj workspace list`) can tell apply spawns from proposal
// spawns at a glance. See ADR-011.
const (
	ApplyPrefix    = "app-"
	ProposalPrefix = "prop-"
)

// SpawnOptions holds configurable parameters for Spawn.
type SpawnOptions struct {
	Agent         string // agent command template (with {change}/{prompt} and {wsdir} placeholders)
	Session       string // tmux session name (empty = current)
	WorkspaceRoot string // override workspace root (empty = default)
}

// Spawn creates a jj workspace, tmux window, and launches an agent for an
// existing openspec change. The workspace/window is named `app-<change>`.
//
// This is the apply flow: it validates the change exists, then isolates it and
// launches `/opsx:apply`. The seed is the change name (substituted into
// {change} in the agent template).
func Spawn(changeName string, opts SpawnOptions) error {
	if err := checkOpenspecChange(changeName); err != nil {
		return err
	}
	name := ApplyPrefix + changeName
	agentTemplate := opts.Agent
	if agentTemplate == "" {
		agentTemplate = DefaultAgentCommand
	}
	if err := isolateAndLaunch(name, changeName, agentTemplate, opts); err != nil {
		return err
	}

	fmt.Printf("Spawned workspace for change %q in %s\n", changeName, name)
	fmt.Println("Main workspace is now on a fresh change. Your previous work is in @-.")
	return nil
}

// Mode selects the seed command a proposal spawn launches.
type Mode string

const (
	ModeExplore Mode = "explore"
	ModePropose Mode = "propose"
)

// DefaultMode is the configurable default for `spawn proposal`. Explore is the
// earliest mode of a proposal (ADR-011), so a bare `spawn proposal` starts there.
const DefaultMode = ModeExplore

// SpawnProposal creates a jj workspace, tmux window, and launches an agent
// seeded from a free-text prompt to create new work. There is no openspec
// change at spawn time — the agent invents one inside the isolated workspace.
//
// The identity is a code-derived slug from the prompt (no AI), prefixed
// `prop-`, made unique against existing workspaces/windows. The slug is the
// immutable handle and display name; it is never remapped after the agent names
// its change (ADR-011).
func SpawnProposal(prompt string, mode Mode, opts SpawnOptions) error {
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("proposal prompt must not be empty")
	}

	slug := workspace.Slug(prompt)
	taken, err := takenSlugs(opts.Session)
	if err != nil {
		return err
	}
	slug = workspace.UniqueSlug(slug, taken)
	name := ProposalPrefix + slug

	agentTemplate := opts.Agent
	if agentTemplate == "" {
		agentTemplate = proposalAgentCommand(mode)
	}
	// The seed substituted into the agent template is the prompt itself.
	if err := isolateAndLaunch(name, prompt, agentTemplate, opts); err != nil {
		return err
	}

	fmt.Printf("Spawned proposal %q (mode %s) in %s\n", slug, mode, name)
	fmt.Println("Main workspace is now on a fresh change. Your previous work is in @-.")
	return nil
}

// isolateAndLaunch is the shared tail of both spawn flows: it runs the tmux/jj
// precondition checks, snapshots the main workspace, creates the jj workspace,
// opens the tmux window, and launches the agent. `name` is the already-prefixed
// workspace/window name; `seed` is substituted into the agent template's
// {change}/{prompt} placeholder. The flows differ only in how `name` and the
// agent template are derived (validate-vs-slug + apply-vs-proposal template).
func isolateAndLaunch(name, seed, agentTemplate string, opts SpawnOptions) error {
	// Only check TMUX env when not targeting a specific session.
	// When --session is set, we target that session directly.
	if opts.Session == "" {
		if err := checkTmuxSession(); err != nil {
			return err
		}
	}
	if err := checkWorkspaceNotExists(name); err != nil {
		return err
	}
	if err := checkWindowNotExists(name, opts.Session); err != nil {
		return err
	}

	wsDir, err := workspace.WorkspaceDir(name, opts.WorkspaceRoot)
	if err != nil {
		return err
	}

	// Snapshot uncommitted work by creating a new empty change.
	// This moves all current work to @- (safe, committed) and leaves
	// @ empty so nothing is lost if the main workspace becomes stale.
	if err := snapshotMainWorkspace(); err != nil {
		return err
	}

	if err := createWorkspace(name, wsDir); err != nil {
		return err
	}
	return openWindow(name, seed, wsDir, agentTemplate, opts)
}

// CheckTmuxSession verifies we're running inside a tmux session.
func checkTmuxSession() error {
	if os.Getenv("TMUX") == "" {
		return fmt.Errorf("jjay must be run inside a tmux session")
	}
	return nil
}

func checkOpenspecChange(changeName string) error {
	names, err := openspec.ChangeNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		if name == changeName {
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

func checkWorkspaceNotExists(name string) error {
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return fmt.Errorf("failed to list jj workspaces: %w", err)
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == name+":" {
			return fmt.Errorf("jj workspace %q already exists", name)
		}
	}
	return nil
}

func checkWindowNotExists(name, session string) error {
	wn := workspace.WindowName(name)
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

// takenSlugs returns the set of bare slugs already in use by existing proposal
// spawns (workspaces or tmux windows carrying the `prop-` prefix), so a new
// slug can be made unique against them. A missing tmux server contributes no
// window names (every proposal is then known only via its workspace).
func takenSlugs(session string) (map[string]bool, error) {
	taken := map[string]bool{}

	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list jj workspaces: %w", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		ws := strings.TrimSuffix(fields[0], ":")
		if s, ok := strings.CutPrefix(ws, ProposalPrefix); ok {
			taken[s] = true
		}
	}

	args := []string{"list-windows", "-F", "#{window_name}"}
	if session != "" {
		args = append(args, "-t", session)
	}
	if wout, werr := exec.Command("tmux", args...).Output(); werr == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(wout)), "\n") {
			// Window names are ws-<name>; strip ws- then the prop- prefix.
			wn := strings.TrimPrefix(strings.TrimSpace(line), "ws-")
			if s, ok := strings.CutPrefix(wn, ProposalPrefix); ok {
				taken[s] = true
			}
		}
	}

	return taken, nil
}

func createWorkspace(name, wsDir string) error {
	if err := os.MkdirAll(filepath.Dir(wsDir), 0o755); err != nil {
		return fmt.Errorf("failed to create workspace parent directory: %w", err)
	}
	// Base the new workspace on @ so it includes uncommitted files
	// (e.g., the active openspec change directory)
	// Use @- to get the snapshot created by jj new (contains all files).
	// Using @ would create a child of the empty new change.
	cmd := exec.Command("jj", "workspace", "add", "--name", name, "--revision", "@-", wsDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create jj workspace: %w", err)
	}
	return nil
}

// OpenWindow creates the `ws-<name>` tmux window for an existing apply-spawn
// workspace and launches the apply agent inside it. It is the reopen entry
// point used by session-open, which only knows the workspace name (the seed for
// the apply template is that same name). Both Spawn and session-open route
// through openWindow so they cannot diverge.
//
// It does not create or check the jj workspace — the caller owns that.
func OpenWindow(name, wsDir string, opts SpawnOptions) error {
	agentTemplate := opts.Agent
	if agentTemplate == "" {
		agentTemplate = DefaultAgentCommand
	}
	// session-open reopens by workspace name; the apply template's {change}
	// placeholder is the change embedded in the `app-` name. For proposal
	// spawns the seed prompt is gone, but the window/pane layout still reopens;
	// the agent template falls back to the name as the seed.
	seed := strings.TrimPrefix(name, ApplyPrefix)
	return openWindow(name, seed, wsDir, agentTemplate, opts)
}

// openWindow creates the tmux window and sets up the panes/agent. `name` is the
// prefixed workspace/window name; `seed` is the value substituted into the
// agent template's {change}/{prompt} placeholder.
func openWindow(name, seed, wsDir, agentTemplate string, opts SpawnOptions) error {
	if err := createWindow(name, opts.Session, wsDir); err != nil {
		return err
	}
	return setupPanes(name, seed, wsDir, agentTemplate, opts)
}

func createWindow(name, session, wsDir string) error {
	wn := workspace.WindowName(name)
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

// DefaultAgentCommand is the default apply-flow agent command template.
const DefaultAgentCommand = `claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}`

// proposalExploreCommand / proposalProposeCommand are the proposal-flow agent
// templates. The {prompt} placeholder is the free-text seed; {wsdir} points the
// agent at the isolated workspace so it writes openspec/changes/<ai-name>/
// inside its own workspace, never racing the main working copy.
const (
	proposalExploreCommand = `claude "/opsx:explore {prompt}" --dangerously-skip-permissions --add-dir {wsdir}`
	proposalProposeCommand = `claude "/opsx:propose {prompt}" --dangerously-skip-permissions --add-dir {wsdir}`
)

// proposalAgentCommand returns the agent template for the given proposal mode.
func proposalAgentCommand(mode Mode) string {
	if mode == ModePropose {
		return proposalProposeCommand
	}
	return proposalExploreCommand
}

// resolveAgentCommand substitutes {change}/{prompt} and {wsdir} placeholders in
// the agent command. The seed fills both {change} and {prompt} so a single
// template family covers both flows.
func resolveAgentCommand(template, seed, wsDir string) string {
	r := strings.NewReplacer("{change}", seed, "{prompt}", seed, "{wsdir}", wsDir)
	return r.Replace(template)
}

func tmuxTarget(session, window string) string {
	if session != "" {
		return session + ":" + window
	}
	return window
}

func setupPanes(name, seed, wsDir, agentTemplate string, opts SpawnOptions) error {
	wn := workspace.WindowName(name)
	target := tmuxTarget(opts.Session, wn)

	agentCmd := resolveAgentCommand(agentTemplate, seed, wsDir)

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
