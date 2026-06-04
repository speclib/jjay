package completion

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

// stub installs the data-source readers for a test and restores them after.
func stub(t *testing.T, changes, workspaces []string, cerr, werr error) {
	t.Helper()
	origC, origW := changeNames, workspaceNames
	t.Cleanup(func() { changeNames, workspaceNames = origC, origW })
	changeNames = func() ([]string, error) { return changes, cerr }
	workspaceNames = func() ([]string, error) { return workspaces, werr }
}

func contains(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}

func assertNoFileComp(t *testing.T, d cobra.ShellCompDirective) {
	t.Helper()
	if d&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Errorf("expected ShellCompDirectiveNoFileComp set, got %v", d)
	}
}

func TestSpawnable_SetMinus(t *testing.T) {
	// add-foo is spawned, so only fix-bar should be offered for spawn.
	stub(t, []string{"add-foo", "fix-bar"}, []string{"add-foo"}, nil, nil)

	names, d := Spawnable(nil, nil, "")
	assertNoFileComp(t, d)
	if contains(names, "add-foo") {
		t.Error("add-foo is already spawned; must not be offered for spawn")
	}
	if !contains(names, "fix-bar") {
		t.Errorf("fix-bar should be offered for spawn, got %v", names)
	}
}

func TestSpawnable_AllSpawned(t *testing.T) {
	stub(t, []string{"add-foo"}, []string{"add-foo"}, nil, nil)
	names, d := Spawnable(nil, nil, "")
	assertNoFileComp(t, d)
	if len(names) != 0 {
		t.Errorf("expected no spawnable changes, got %v", names)
	}
}

func TestSpawnable_NoneSpawned(t *testing.T) {
	stub(t, []string{"add-foo", "fix-bar"}, nil, nil, nil)
	names, _ := Spawnable(nil, nil, "")
	if len(names) != 2 {
		t.Errorf("expected both changes spawnable, got %v", names)
	}
}

func TestSpawnable_ReaderError(t *testing.T) {
	// openspec read fails → no candidates, no file fallback, no error to shell.
	stub(t, nil, []string{"add-foo"}, errors.New("openspec down"), nil)
	names, d := Spawnable(nil, nil, "")
	assertNoFileComp(t, d)
	if len(names) != 0 {
		t.Errorf("expected no candidates on reader error, got %v", names)
	}

	// jj read fails → same.
	stub(t, []string{"add-foo"}, nil, nil, errors.New("jj down"))
	names, d = Spawnable(nil, nil, "")
	assertNoFileComp(t, d)
	if len(names) != 0 {
		t.Errorf("expected no candidates on workspace reader error, got %v", names)
	}
}

func TestMergeableAndCleanable_ListWorkspaces(t *testing.T) {
	stub(t, nil, []string{"add-foo", "fix-bar"}, nil, nil)

	for _, fn := range []struct {
		name string
		f    func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
	}{
		{"Mergeable", Mergeable},
		{"Cleanable", Cleanable},
	} {
		names, d := fn.f(nil, nil, "")
		assertNoFileComp(t, d)
		if !contains(names, "add-foo") || !contains(names, "fix-bar") {
			t.Errorf("%s should offer both workspaces, got %v", fn.name, names)
		}
	}
}

func TestMergeable_NoWorkspaces(t *testing.T) {
	stub(t, nil, nil, nil, nil)
	names, d := Mergeable(nil, nil, "")
	assertNoFileComp(t, d)
	if len(names) != 0 {
		t.Errorf("expected no candidates when no workspaces, got %v", names)
	}
}

func TestMergeableAndCleanable_ReaderError(t *testing.T) {
	stub(t, nil, nil, nil, errors.New("jj down"))
	for _, f := range []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective){Mergeable, Cleanable} {
		names, d := f(nil, nil, "")
		assertNoFileComp(t, d)
		if len(names) != 0 {
			t.Errorf("expected no candidates on reader error, got %v", names)
		}
	}
}

func TestSetMinus(t *testing.T) {
	got := setMinus([]string{"a", "b", "c"}, []string{"b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Errorf("setMinus = %v, want [a c]", got)
	}
	// Subtracting nothing returns all, in order.
	if got := setMinus([]string{"a", "b"}, nil); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("setMinus with empty b = %v, want [a b]", got)
	}
	// Subtracting everything returns empty.
	if got := setMinus([]string{"a"}, []string{"a"}); len(got) != 0 {
		t.Errorf("setMinus all-removed = %v, want []", got)
	}
}
