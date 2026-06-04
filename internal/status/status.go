// Package status derives the live state of spawned jj workspaces by joining
// `jj workspace list` with `tmux list-windows`. It persists nothing; the jj
// workspace is the single source of truth (see ADR-006). A spawned workspace
// is "attached" if a matching `ws-<change>` tmux window exists in the target
// session, otherwise "detached" — meaning the workspace still exists on disk
// (the spawn is still open), there is just no live tmux window/agent for it.
package status

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	"jjay/internal/workspace"
)

// defaultWorkspaceName is the jj name of the main working copy, which appears
// in `jj workspace list` but is not a spawn and is excluded from status.
const defaultWorkspaceName = "default"

// TaskCount summarizes the openspec tasks.md checkbox progress for a spawn.
type TaskCount struct {
	Done     int
	Total    int
	Found    bool // false when no readable tasks.md was located
	Archived bool // true when the tasks were read from the archive location
}

// Spawn describes a single spawned jj workspace and whether its tmux window
// currently exists in the target session.
type Spawn struct {
	Change   string    // jj workspace name == openspec change name
	WSDir    string    // resolved absolute workspace directory
	Attached bool      // true if a ws-<change> window exists in the session
	Archived bool      // true if the change is archived (tasks live under archive/)
	Tasks    TaskCount // openspec task progress for this workspace
}

// List returns the spawned workspaces for the repository, scoped to the given
// tmux session (empty = current session). It tolerates a missing tmux server:
// in that case every spawn is reported as detached. workspaceRoot overrides the
// workspace root used to resolve each workspace directory (empty = default).
//
// Paths are anchored on the main repo root (resolved from the current
// directory's .jj pointer) so resolution is correct even when jjay is run from
// inside a child workspace.
func List(session, workspaceRoot string) (spawns []Spawn, mainRoot string, err error) {
	wsOut, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to list jj workspaces: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get working directory: %w", err)
	}
	mainRoot, err = workspace.MainRepoRoot(cwd)
	if err != nil {
		return nil, "", err
	}

	windows := listWindows(session)

	spawns, err = join(string(wsOut), windows, mainRoot, workspaceRoot, readTaskCount)
	return spawns, mainRoot, err
}

// listWindows returns the set of tmux window names in the target session.
// A missing tmux server (or any tmux error) is treated as "no windows",
// mirroring session.checkSessionNotExists — every spawn is then detached.
func listWindows(session string) map[string]bool {
	args := []string{"list-windows", "-F", "#{window_name}"}
	if session != "" {
		args = append(args, "-t", session)
	}
	out, err := exec.Command("tmux", args...).Output()
	if err != nil {
		return map[string]bool{}
	}
	return parseWindows(string(out))
}

// parseWindows turns `tmux list-windows -F '#{window_name}'` output into a set.
func parseWindows(out string) map[string]bool {
	set := map[string]bool{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			set[line] = true
		}
	}
	return set
}

// parseWorkspaceNames parses `jj workspace list` output into workspace names,
// reusing the convention from spawn.checkWorkspaceNotExists: each line is
// "<name>: <commit>" so the name is fields[0] with the trailing colon stripped.
// The "default" main working copy is excluded.
func parseWorkspaceNames(out string) []string {
	var names []string
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		name := strings.TrimSuffix(fields[0], ":")
		if name == "" || name == defaultWorkspaceName {
			continue
		}
		// Only count lines that actually carried the "<name>:" colon form.
		if !strings.HasSuffix(fields[0], ":") {
			continue
		}
		names = append(names, name)
	}
	return names
}

var (
	taskDoneRe = regexp.MustCompile(`(?m)^\s*-\s*\[[xX]\]`)
	taskOpenRe = regexp.MustCompile(`(?m)^\s*-\s*\[\s\]`)
)

// taskCounter reads task progress for a spawn given its absolute workspace dir
// and change name; injectable for testing.
type taskCounter func(wsDir, change string) TaskCount

// readTaskCount counts done/total checkboxes in the spawn's openspec tasks.md.
// It looks first at the active location (<wsDir>/openspec/changes/<change>/
// tasks.md); if that is absent, the change has been archived, so it looks under
// <wsDir>/openspec/changes/archive/<date>-<change>/tasks.md (matched by suffix,
// since the archive dir carries a date prefix) and marks the result Archived.
// A missing/unreadable file in both locations yields Found=false rather than an
// error — status must not fail on it.
func readTaskCount(wsDir, change string) TaskCount {
	active := filepath.Join(wsDir, "openspec", "changes", change, "tasks.md")
	if data, err := os.ReadFile(active); err == nil {
		return countTasks(string(data))
	}

	// Archived: openspec/changes/archive/<date>-<change>/tasks.md
	matches, _ := filepath.Glob(filepath.Join(wsDir, "openspec", "changes", "archive", "*-"+change, "tasks.md"))
	for _, m := range matches {
		// Guard against a suffix collision (e.g. querying "foo" must not match
		// an archived "do-foo"): the dir must be exactly "<YYYY-MM-DD>-<change>".
		parent := filepath.Base(filepath.Dir(m))
		if !isArchiveDirFor(parent, change) {
			continue
		}
		if data, err := os.ReadFile(m); err == nil {
			tc := countTasks(string(data))
			tc.Archived = true
			return tc
		}
	}
	return TaskCount{}
}

// archiveDirRe matches an archive directory name "YYYY-MM-DD-<change>",
// capturing the change portion.
var archiveDirRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-(.+)$`)

// isArchiveDirFor reports whether dir is the dated archive directory for the
// given change (i.e. "<YYYY-MM-DD>-<change>" with the change portion matching
// exactly, so "do-foo" is not accepted for "foo").
func isArchiveDirFor(dir, change string) bool {
	m := archiveDirRe.FindStringSubmatch(dir)
	return m != nil && m[1] == change
}

// countTasks counts markdown task checkboxes in tasks.md content.
func countTasks(content string) TaskCount {
	done := len(taskDoneRe.FindAllString(content, -1))
	open := len(taskOpenRe.FindAllString(content, -1))
	return TaskCount{Done: done, Total: done + open, Found: true}
}

// formatTasks renders a TaskCount as "done/total (pct%)", or "-" when no
// tasks.md was found.
func formatTasks(t TaskCount) string {
	if !t.Found || t.Total == 0 {
		return "-"
	}
	pct := t.Done * 100 / t.Total
	return fmt.Sprintf("%d/%d (%d%%)", t.Done, t.Total, pct)
}

// Render writes a human-readable table of spawns to w. With no spawns it prints
// a "no running spawns" line. The WORKSPACE column is shown relative to
// mainRoot. Exit-zero behavior is the caller's concern.
func Render(w io.Writer, mainRoot string, spawns []Spawn) {
	if len(spawns) == 0 {
		fmt.Fprintln(w, "No running spawns.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CHANGE\tWORKSPACE\tTASKS\tARCHIVED\tSTATUS")
	for _, s := range spawns {
		state := "detached"
		if s.Attached {
			state = "attached"
		}
		archived := "no"
		if s.Archived {
			archived = "yes"
		}
		rel := workspace.RelativeToMain(mainRoot, s.WSDir)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", s.Change, rel, formatTasks(s.Tasks), archived, state)
	}
	tw.Flush()
}

// join builds the spawn list from raw `jj workspace list` output and the set of
// current tmux window names. It is pure aside from the injected taskCounter, so
// it can be unit-tested. Workspace dirs are anchored on mainRoot.
func join(wsListOut string, windows map[string]bool, mainRoot, workspaceRoot string, tasks taskCounter) ([]Spawn, error) {
	var spawns []Spawn
	for _, name := range parseWorkspaceNames(wsListOut) {
		wsDir, err := workspace.WorkspaceDirFrom(mainRoot, name, workspaceRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve workspace dir for %q: %w", name, err)
		}
		tc := tasks(wsDir, name)
		spawns = append(spawns, Spawn{
			Change:   name,
			WSDir:    wsDir,
			Attached: windows[workspace.WindowName(name)],
			Archived: tc.Archived,
			Tasks:    tc,
		})
	}
	return spawns, nil
}
