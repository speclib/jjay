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

// Verb-prefixes encode the spawn kind in the workspace/window name (ADR-011).
// Mirror of spawn.ApplyPrefix/spawn.ProposalPrefix; duplicated here so status
// (a data-source reader) does not import the command package.
const (
	applyPrefix    = "app-"
	proposalPrefix = "prop-"
)

// TaskCount summarizes the openspec tasks.md checkbox progress for a spawn.
type TaskCount struct {
	Done     int
	Total    int
	Found    bool // false when no readable tasks.md was located
	Archived bool // true when the tasks were read from the archive location
}

// Kind classifies a spawn by its name prefix (ADR-011): an apply spawn tracks
// an existing openspec change (`app-*`); a proposal spawn is prompt-seeded and
// may not have a change yet (`prop-*`).
type Kind int

const (
	KindApply Kind = iota
	KindProposal
)

// Spawn describes a single spawned jj workspace and whether its tmux window
// currently exists in the target session.
//
// Name is the jj workspace name (verb-prefixed). For an apply spawn, Change is
// the openspec change name (Name with the `app-` prefix stripped) and equals
// the change directory inside the workspace. For a proposal spawn the workspace
// name does NOT equal any change name — the agent invents a differently-named
// change inside the workspace (ADR-011) — so Change is empty and change-shaped
// columns (TASKS/ARCHIVED) do not apply.
type Spawn struct {
	Name     string    // jj workspace name (e.g. app-add-foo, prop-dark-mode)
	Kind     Kind      // apply vs proposal, from the name prefix
	Change   string    // openspec change name for apply spawns; empty for proposals
	WSDir    string    // resolved absolute workspace directory
	Attached bool      // true if a ws-<name> window exists in the session
	Archived bool      // true if the change is archived (tasks live under archive/)
	Merged   bool      // true if the spawn's work has already landed on main
	Tasks    TaskCount // openspec task progress (apply spawns only)
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
	return ListIn("", session, workspaceRoot)
}

// ListIn is List scoped to a specific repository. repoRoot anchors BOTH the
// `jj workspace list` query (via `jj -R <repoRoot>`) and the resolved mainRoot,
// so it reports the workspaces of THAT repo — not the current working
// directory's. An empty repoRoot falls back to cwd (the List behavior).
//
// This matters for `session-open <path>`: it opens a session for a DIFFERENT
// repo than the one jjay runs from, and must reopen that repo's workspaces, not
// the caller's (jjay-… cross-project leak).
func ListIn(repoRoot, session, workspaceRoot string) (spawns []Spawn, mainRoot string, err error) {
	args := []string{"workspace", "list"}
	if repoRoot != "" {
		args = []string{"-R", repoRoot, "workspace", "list"}
	}
	wsOut, err := exec.Command("jj", args...).Output()
	if err != nil {
		return nil, "", fmt.Errorf("failed to list jj workspaces: %w", err)
	}

	if repoRoot != "" {
		// Anchor on the given repo. Resolve through MainRepoRoot so a child
		// workspace path still yields the true main root.
		mainRoot, err = workspace.MainRepoRoot(repoRoot)
		if err != nil {
			return nil, "", err
		}
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, "", fmt.Errorf("failed to get working directory: %w", err)
		}
		mainRoot, err = workspace.MainRepoRoot(cwd)
		if err != nil {
			return nil, "", err
		}
	}

	windows := listWindows(session)

	spawns, err = join(string(wsOut), windows, mainRoot, workspaceRoot, readTaskCount, isMerged)
	return spawns, mainRoot, err
}

// WorkspaceNames returns the names of spawned jj workspaces (the `default` main
// working copy excluded). It is a lean reader for shell completion: it runs only
// `jj workspace list` — no tmux probing, no tasks.md reads — and shares the
// parse/exclusion rule with List via parseWorkspaceNames. Returns an error if jj
// cannot be run; completion callers should treat that as "no candidates".
func WorkspaceNames() ([]string, error) {
	out, err := exec.Command("jj", "workspace", "list").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list jj workspaces: %w", err)
	}
	return parseWorkspaceNames(string(out)), nil
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

// mergeChecker reports whether a spawn's work has already landed on main, given
// its change (== jj workspace) name; injectable for testing.
type mergeChecker func(change string) bool

// isMerged reports whether the spawn's work is already on the `main` bookmark,
// derived live from jj with no state file (ADR-006). A spawn is merged when its
// workspace head has no commits that `main` lacks, i.e. the revset
// `main..<change>@` is empty. A fresh/empty spawn sits on its own (empty) commit
// ahead of main, so that revset is non-empty and it correctly reads not-merged.
//
// Tolerance: any jj failure (unknown workspace, no `main` bookmark, jj missing)
// yields false rather than an error — status must not fail on it, mirroring the
// tmux-tolerance stance in listWindows.
func isMerged(change string) bool {
	revset := "main.." + change + "@"
	out, err := exec.Command("jj", "log", "-r", revset, "--no-graph", "-T", `"x"`).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == ""
}

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

// Render writes a human-readable view of spawns to w, split into two tables by
// kind (ADR-011): CHANGES (apply spawns, with change-shaped columns) and
// PROPOSAL SPAWNS (prompt-seeded, no change yet). With no spawns it prints a
// "no running spawns" line. The WORKSPACE column is shown relative to mainRoot.
// A table with no rows is omitted. Exit-zero behavior is the caller's concern.
func Render(w io.Writer, mainRoot string, spawns []Spawn) {
	if len(spawns) == 0 {
		fmt.Fprintln(w, "No running spawns.")
		return
	}

	var changes, proposals []Spawn
	for _, s := range spawns {
		if s.Kind == KindProposal {
			proposals = append(proposals, s)
		} else {
			changes = append(changes, s)
		}
	}

	if len(changes) > 0 {
		fmt.Fprintln(w, "CHANGES")
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "CHANGE\tWORKSPACE\tTASKS\tTMUX\tMERGED\tARCHIVED")
		for _, s := range changes {
			merged := "no"
			if s.Merged {
				merged = "yes"
			}
			archived := "no"
			if s.Archived {
				archived = "yes"
			}
			rel := workspace.RelativeToMain(mainRoot, s.WSDir)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", s.Change, rel, formatTasks(s.Tasks), stateOf(s), merged, archived)
		}
		tw.Flush()
	}

	if len(proposals) > 0 {
		if len(changes) > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, "PROPOSAL SPAWNS")
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "PROPOSAL\tWORKSPACE\tMERGED\tTMUX")
		for _, s := range proposals {
			merged := "no"
			if s.Merged {
				merged = "yes"
			}
			rel := workspace.RelativeToMain(mainRoot, s.WSDir)
			// Display the bare slug (workspace name minus the prop- prefix).
			name := strings.TrimPrefix(s.Name, proposalPrefix)
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, rel, merged, stateOf(s))
		}
		tw.Flush()
	}
}

// stateOf renders the attach state of a spawn.
func stateOf(s Spawn) string {
	if s.Attached {
		return "attached"
	}
	return "detached"
}

// join builds the spawn list from raw `jj workspace list` output and the set of
// current tmux window names. It is pure aside from the injected taskCounter, so
// it can be unit-tested. Workspace dirs are anchored on mainRoot.
func join(wsListOut string, windows map[string]bool, mainRoot, workspaceRoot string, tasks taskCounter, merged mergeChecker) ([]Spawn, error) {
	var spawns []Spawn
	for _, name := range parseWorkspaceNames(wsListOut) {
		// The workspace dir is keyed on the (prefixed) workspace name; the
		// directory carries the same name (ADR-011: prefix flows through naming).
		wsDir, err := workspace.WorkspaceDirFrom(mainRoot, name, workspaceRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve workspace dir for %q: %w", name, err)
		}

		s := Spawn{
			Name:     name,
			WSDir:    wsDir,
			Attached: windows[workspace.WindowName(name)],
			// Merged is keyed on the jj workspace name (the `<name>@` revset),
			// not the openspec change name — a proposal spawn has no matching
			// change name, but its workspace head still merges by name.
			Merged: merged(name),
		}

		switch {
		case strings.HasPrefix(name, proposalPrefix):
			// Proposal spawn: the produced change is named by the agent and is
			// NOT the workspace name, so we do not infer a change name or read
			// change-shaped task progress from it.
			s.Kind = KindProposal
		case strings.HasPrefix(name, applyPrefix):
			s.Kind = KindApply
			s.Change = strings.TrimPrefix(name, applyPrefix)
		default:
			// Legacy/unprefixed workspace (pre-ADR-011): treat as an apply spawn
			// whose change name equals the workspace name, preserving old behavior.
			s.Kind = KindApply
			s.Change = name
		}

		if s.Kind == KindApply {
			// The change directory inside the workspace uses the un-prefixed
			// change name, not the workspace name.
			tc := tasks(wsDir, s.Change)
			s.Archived = tc.Archived
			s.Tasks = tc
		}

		spawns = append(spawns, s)
	}
	return spawns, nil
}
