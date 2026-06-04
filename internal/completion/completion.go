// Package completion supplies per-verb shell-completion candidates for the
// jjay CLI's change-name arguments. Per ADR-009 it depends only on data-source
// readers (internal/openspec, internal/status) and never on command packages
// (spawn/merge/cleanup); the commands bind to these functions via
// ValidArgsFunction, not the reverse.
//
// Each function returns ([]string, cobra.ShellCompDirective) and uses
// ShellCompDirectiveNoFileComp so a change-name argument never falls back to
// file-name completion. Reads are name-only (no tmux, no tasks.md), keeping a
// TAB press fast. Completion is advisory — the commands' own precondition checks
// remain authoritative — so any reader failure degrades to "no candidates"
// rather than a shell-visible error.
package completion

import (
	"github.com/spf13/cobra"

	"jjay/internal/openspec"
	"jjay/internal/status"
)

// changeNames and workspaceNames are indirected so tests can substitute the
// data-source readers without invoking openspec/jj.
var (
	changeNames    = openspec.ChangeNames
	workspaceNames = status.WorkspaceNames
)

// Spawnable suggests openspec changes that do not yet have a spawned workspace
// (ChangeNames \ WorkspaceNames). You cannot spawn an already-spawned change.
func Spawnable(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	changes, cerr := changeNames()
	spawned, werr := workspaceNames()
	if cerr != nil || werr != nil {
		return none()
	}
	return setMinus(changes, spawned), cobra.ShellCompDirectiveNoFileComp
}

// Mergeable suggests existing spawned workspaces — only those can be merged.
func Mergeable(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return spawnedWorkspaces()
}

// Cleanable suggests existing spawned workspaces — only those can be torn down.
// Named apart from Mergeable (though identical today) so the sets can diverge
// later without re-plumbing the CLI bindings (ADR-009).
func Cleanable(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return spawnedWorkspaces()
}

// spawnedWorkspaces is the shared body for Mergeable/Cleanable.
func spawnedWorkspaces() ([]string, cobra.ShellCompDirective) {
	names, err := workspaceNames()
	if err != nil {
		return none()
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// none is the graceful-degradation result: no candidates, no file fallback, no
// shell-visible error.
func none() ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// setMinus returns the elements of a that are not in b, preserving a's order.
func setMinus(a, b []string) []string {
	exclude := make(map[string]struct{}, len(b))
	for _, x := range b {
		exclude[x] = struct{}{}
	}
	var out []string
	for _, x := range a {
		if _, skip := exclude[x]; !skip {
			out = append(out, x)
		}
	}
	return out
}
