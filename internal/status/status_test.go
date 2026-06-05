package status

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"jjay/internal/workspace"
)

func spawnByChange(spawns []Spawn, change string) (Spawn, bool) {
	for _, s := range spawns {
		if s.Change == change {
			return s, true
		}
	}
	return Spawn{}, false
}

func spawnByName(spawns []Spawn, name string) (Spawn, bool) {
	for _, s := range spawns {
		if s.Name == name {
			return s, true
		}
	}
	return Spawn{}, false
}

const sampleWSList = `default: abc123 (no description set)
add-foo: def456 (no description set)
fix-bar: 789abc (no description set)
`

// noTasks is a taskCounter that reports nothing found, for join tests that
// don't care about progress counts.
func noTasks(_, _ string) TaskCount { return TaskCount{} }

// noMerged is a mergeChecker that reports everything unmerged, for join tests
// that don't care about the merged signal.
func noMerged(_ string) bool { return false }

const mainRoot = "/repo/main"
const wsRoot = "/repo/ws"

func TestJoin_AttachedAndDetached(t *testing.T) {
	windows := map[string]bool{
		workspace.WindowName("add-foo"): true,
		// fix-bar has no window
	}

	spawns, err := join(sampleWSList, windows, mainRoot, wsRoot, noTasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}

	if len(spawns) != 2 {
		t.Fatalf("expected 2 spawns (default excluded), got %d: %+v", len(spawns), spawns)
	}

	foo, ok := spawnByChange(spawns, "add-foo")
	if !ok {
		t.Fatal("expected add-foo spawn")
	}
	if !foo.Attached {
		t.Error("add-foo should be attached (window exists)")
	}

	bar, ok := spawnByChange(spawns, "fix-bar")
	if !ok {
		t.Fatal("expected fix-bar spawn")
	}
	if bar.Attached {
		t.Error("fix-bar should be detached (no window)")
	}
}

func TestJoin_NoTmuxAllDetached(t *testing.T) {
	// Empty window set == tmux missing/no server.
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	for _, s := range spawns {
		if s.Attached {
			t.Errorf("%s should be detached when no windows exist", s.Change)
		}
	}
}

func TestJoin_DefaultExcluded(t *testing.T) {
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	if _, ok := spawnByChange(spawns, "default"); ok {
		t.Error("default workspace must be excluded from spawns")
	}
}

func TestJoin_WSDirResolved(t *testing.T) {
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	foo, _ := spawnByChange(spawns, "add-foo")
	want, _ := workspace.WorkspaceDirFrom(mainRoot, "add-foo", wsRoot)
	if foo.WSDir != want {
		t.Errorf("WSDir = %q, want %q", foo.WSDir, want)
	}
}

func TestJoin_TaskCounterInvoked(t *testing.T) {
	tasks := func(wsDir, change string) TaskCount {
		if change == "add-foo" {
			return TaskCount{Done: 4, Total: 10, Found: true}
		}
		return TaskCount{}
	}
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, tasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	foo, _ := spawnByChange(spawns, "add-foo")
	if foo.Tasks.Done != 4 || foo.Tasks.Total != 10 {
		t.Errorf("expected 4/10 tasks for add-foo, got %+v", foo.Tasks)
	}
}

func TestJoin_MergeCheckerInvoked(t *testing.T) {
	// add-foo is merged (work on main), fix-bar is not.
	merged := func(change string) bool { return change == "add-foo" }
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks, merged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	foo, _ := spawnByChange(spawns, "add-foo")
	if !foo.Merged {
		t.Error("add-foo should be merged")
	}
	bar, _ := spawnByChange(spawns, "fix-bar")
	if bar.Merged {
		t.Error("fix-bar should not be merged")
	}
}

// TestIsMerged_ToleratesMissingJJ verifies the live merged check returns false
// (not a panic or error) when jj cannot be run — status must not fail on it.
// Only the failure path is deterministic without a real jj repo.
func TestIsMerged_ToleratesMissingJJ(t *testing.T) {
	if _, err := exec.LookPath("jj"); err == nil {
		t.Skip("jj binary present; cannot test the missing-binary path")
	}
	if isMerged("any-change") {
		t.Error("expected MERGED=false when jj is unavailable")
	}
}

const prefixedWSList = `default: abc123 (no description set)
app-add-foo: def456 (no description set)
prop-dark-mode: 789abc (no description set)
`

func TestJoin_ClassifiesByPrefix(t *testing.T) {
	// Apply spawns read tasks via the un-prefixed change name; proposal spawns
	// don't (their change name is unknown / agent-invented).
	tasks := func(wsDir, change string) TaskCount {
		if change == "add-foo" {
			return TaskCount{Done: 2, Total: 4, Found: true}
		}
		// A proposal spawn must never trigger a change-shaped task read.
		t.Errorf("unexpected task read for change %q", change)
		return TaskCount{}
	}
	spawns, err := join(prefixedWSList, map[string]bool{}, mainRoot, wsRoot, tasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}

	apply, ok := spawnByName(spawns, "app-add-foo")
	if !ok {
		t.Fatal("expected app-add-foo spawn")
	}
	if apply.Kind != KindApply || apply.Change != "add-foo" {
		t.Errorf("app-add-foo: Kind=%v Change=%q, want apply/add-foo", apply.Kind, apply.Change)
	}
	if apply.Tasks.Done != 2 || apply.Tasks.Total != 4 {
		t.Errorf("app-add-foo tasks = %+v, want 2/4", apply.Tasks)
	}

	prop, ok := spawnByName(spawns, "prop-dark-mode")
	if !ok {
		t.Fatal("expected prop-dark-mode spawn")
	}
	if prop.Kind != KindProposal {
		t.Errorf("prop-dark-mode Kind = %v, want proposal", prop.Kind)
	}
	if prop.Change != "" {
		t.Errorf("proposal spawn must not infer a change name, got %q", prop.Change)
	}
	if prop.Tasks.Found {
		t.Errorf("proposal spawn must not carry task counts, got %+v", prop.Tasks)
	}
}

func TestJoin_AttachedKeyedOnWorkspaceName(t *testing.T) {
	// The window name is ws-<workspace-name>, including the prefix.
	windows := map[string]bool{workspace.WindowName("prop-dark-mode"): true}
	spawns, err := join(prefixedWSList, windows, mainRoot, wsRoot, noTasks, noMerged)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	prop, _ := spawnByName(spawns, "prop-dark-mode")
	if !prop.Attached {
		t.Error("prop-dark-mode should be attached (ws-prop-dark-mode window exists)")
	}
}

func TestRender_TwoTables(t *testing.T) {
	spawns := []Spawn{
		{Name: "app-add-foo", Kind: KindApply, Change: "add-foo", WSDir: "/repo/ws/app-add-foo", Attached: true, Tasks: TaskCount{Done: 1, Total: 2, Found: true}},
		{Name: "prop-dark-mode", Kind: KindProposal, WSDir: "/repo/ws/prop-dark-mode", Attached: false},
	}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	out := b.String()

	for _, want := range []string{"CHANGES", "PROPOSAL SPAWNS", "add-foo", "dark-mode", "1/2 (50%)"} {
		if !strings.Contains(out, want) {
			t.Errorf("render output missing %q:\n%s", want, out)
		}
	}
	// The proposal table must come after the changes table.
	if strings.Index(out, "CHANGES") > strings.Index(out, "PROPOSAL SPAWNS") {
		t.Errorf("CHANGES table should precede PROPOSAL SPAWNS:\n%s", out)
	}
}

func TestRender_OnlyProposals(t *testing.T) {
	spawns := []Spawn{
		{Name: "prop-dark-mode", Kind: KindProposal, WSDir: "/repo/ws/prop-dark-mode"},
	}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	out := b.String()
	if !strings.Contains(out, "PROPOSAL SPAWNS") {
		t.Errorf("expected PROPOSAL SPAWNS table, got:\n%s", out)
	}
	if strings.Contains(out, "CHANGES") {
		t.Errorf("CHANGES table should be omitted when no apply spawns exist:\n%s", out)
	}
}

func TestParseWorkspaceNames_IgnoresJunk(t *testing.T) {
	out := "default: x\n\nadd-foo: y\nnot a workspace line\n"
	names := parseWorkspaceNames(out)
	if len(names) != 1 || names[0] != "add-foo" {
		t.Errorf("parseWorkspaceNames = %v, want [add-foo]", names)
	}
}

// WorkspaceNames shares parseWorkspaceNames with List, so the default-exclusion
// rule is covered here against the same sample List uses.
func TestParseWorkspaceNames_ExcludesDefault(t *testing.T) {
	names := parseWorkspaceNames(sampleWSList)
	if len(names) != 2 || names[0] != "add-foo" || names[1] != "fix-bar" {
		t.Errorf("parseWorkspaceNames = %v, want [add-foo fix-bar] (default excluded)", names)
	}
	for _, n := range names {
		if n == defaultWorkspaceName {
			t.Errorf("default workspace must be excluded, got %v", names)
		}
	}
}

// TestWorkspaceNames_ToleratesMissingJJ verifies WorkspaceNames returns an error
// (not a panic) when jj cannot be run. Only the error path is deterministic
// without a real jj repo.
func TestWorkspaceNames_ToleratesMissingJJ(t *testing.T) {
	if _, err := exec.LookPath("jj"); err == nil {
		t.Skip("jj binary present; cannot test the missing-binary path")
	}
	if _, err := WorkspaceNames(); err == nil {
		t.Error("expected error when jj is unavailable, got nil")
	}
}

func TestRender_ListsRows_RelativePaths(t *testing.T) {
	// Workspace paths are rendered relative to the main repo root.
	spawns := []Spawn{
		{Change: "add-foo", WSDir: "/repo/ws/add-foo", Attached: true, Tasks: TaskCount{Done: 12, Total: 18, Found: true}},
		{Change: "fix-bar", WSDir: "/repo/ws/fix-bar", Attached: false},
	}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	out := b.String()

	for _, want := range []string{"add-foo", "../ws/add-foo", "12/18 (66%)", "attached", "fix-bar", "detached", "TMUX", "MERGED", "ARCHIVED"} {
		if !strings.Contains(out, want) {
			t.Errorf("render output missing %q:\n%s", want, out)
		}
	}
	// The tmux-state column is now headed TMUX, not STATUS.
	if strings.Contains(out, "STATUS") {
		t.Errorf("render must not use the old STATUS header:\n%s", out)
	}
	// The absolute workspace path must NOT appear — only the relative form.
	if strings.Contains(out, "/repo/ws/add-foo") {
		t.Errorf("render leaked absolute path:\n%s", out)
	}
}

func TestRender_MergedColumn(t *testing.T) {
	// Header order must be CHANGE WORKSPACE TASKS TMUX MERGED ARCHIVED, and a
	// merged-but-not-archived spawn is the "ready to clean up" signal.
	spawns := []Spawn{
		{Change: "landed", WSDir: "/repo/ws/landed", Merged: true, Archived: false, Tasks: TaskCount{Done: 3, Total: 3, Found: true}},
		{Change: "wip", WSDir: "/repo/ws/wip", Merged: false, Archived: false, Tasks: TaskCount{Done: 1, Total: 4, Found: true}},
	}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	out := b.String()

	// Locate the CHANGES table's column header (it follows the "CHANGES" title
	// line in the two-table layout).
	var header string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "CHANGE\t") || strings.HasPrefix(strings.TrimSpace(line), "CHANGE ") {
			header = line
			break
		}
	}
	if header == "" {
		t.Fatalf("could not find CHANGES column header in output:\n%s", out)
	}
	for i, col := range []string{"CHANGE", "WORKSPACE", "TASKS", "TMUX", "MERGED", "ARCHIVED"} {
		idx := strings.Index(header, col)
		if idx < 0 {
			t.Errorf("header missing %q: %q", col, header)
		}
		if i > 0 {
			prev := []string{"CHANGE", "WORKSPACE", "TASKS", "TMUX", "MERGED", "ARCHIVED"}[i-1]
			if strings.Index(header, prev) > idx {
				t.Errorf("column %q must come after %q in header: %q", col, prev, header)
			}
		}
	}
}

func TestRender_ArchivedColumn(t *testing.T) {
	spawns := []Spawn{
		{Change: "active-one", WSDir: "/repo/ws/active-one", Archived: false, Tasks: TaskCount{Done: 1, Total: 2, Found: true}},
		{Change: "old-one", WSDir: "/repo/ws/old-one", Archived: true, Tasks: TaskCount{Done: 5, Total: 5, Found: true, Archived: true}},
	}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	out := b.String()

	// Header has ARCHIVED; archived spawn shows yes, active shows no.
	if !strings.Contains(out, "yes") || !strings.Contains(out, "no") {
		t.Errorf("expected yes/no archived values, got:\n%s", out)
	}
	// Archived change still shows its counts (read from the archive location).
	if !strings.Contains(out, "5/5 (100%)") {
		t.Errorf("expected archived change task counts, got:\n%s", out)
	}
}

func TestRender_Empty(t *testing.T) {
	var b strings.Builder
	Render(&b, "/repo/main", nil)
	if !strings.Contains(strings.ToLower(b.String()), "no running spawns") {
		t.Errorf("expected empty-case message, got: %q", b.String())
	}
}

func TestRender_NoTasksFile(t *testing.T) {
	spawns := []Spawn{{Change: "add-foo", WSDir: "/repo/ws/add-foo", Tasks: TaskCount{}}}
	var b strings.Builder
	Render(&b, "/repo/main", spawns)
	// Missing tasks.md renders as "-", not a crash or 0/0.
	if !strings.Contains(b.String(), "-") {
		t.Errorf("expected '-' for missing tasks, got: %q", b.String())
	}
}

func TestParseWindows(t *testing.T) {
	set := parseWindows("ws-add-foo\nws-fix-bar\n\n")
	if !set["ws-add-foo"] || !set["ws-fix-bar"] {
		t.Errorf("parseWindows missing entries: %v", set)
	}
	if len(set) != 2 {
		t.Errorf("expected 2 windows, got %d: %v", len(set), set)
	}
}

func TestCountTasks(t *testing.T) {
	content := `## 1. Group
- [x] 1.1 done
- [ ] 1.2 open
- [X] 1.3 done caps
  - [ ] nested open
not a task line
`
	got := countTasks(content)
	if got.Done != 2 || got.Total != 4 || !got.Found {
		t.Errorf("countTasks = %+v, want Done=2 Total=4 Found=true", got)
	}
}

func TestFormatTasks(t *testing.T) {
	tests := []struct {
		in   TaskCount
		want string
	}{
		{TaskCount{Done: 12, Total: 18, Found: true}, "12/18 (66%)"},
		{TaskCount{Done: 0, Total: 0, Found: true}, "-"},
		{TaskCount{Found: false}, "-"},
		{TaskCount{Done: 3, Total: 3, Found: true}, "3/3 (100%)"},
	}
	for _, tt := range tests {
		if got := formatTasks(tt.in); got != tt.want {
			t.Errorf("formatTasks(%+v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestReadTaskCount_MissingFile(t *testing.T) {
	tc := readTaskCount(t.TempDir(), "no-such-change")
	if tc.Found {
		t.Errorf("expected Found=false for missing tasks.md, got %+v", tc)
	}
}

func TestReadTaskCount_ReadsActiveFile(t *testing.T) {
	ws := t.TempDir()
	dir := filepath.Join(ws, "openspec", "changes", "add-foo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "- [x] one\n- [ ] two\n- [ ] three\n"
	if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	tc := readTaskCount(ws, "add-foo")
	if tc.Done != 1 || tc.Total != 3 || !tc.Found {
		t.Errorf("readTaskCount = %+v, want 1/3 found", tc)
	}
	if tc.Archived {
		t.Error("active tasks should not be marked Archived")
	}
}

func TestReadTaskCount_ReadsArchivedFile(t *testing.T) {
	ws := t.TempDir()
	// No active changes/add-foo, but an archived one with a date prefix.
	dir := filepath.Join(ws, "openspec", "changes", "archive", "2026-06-04-add-foo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "- [x] one\n- [x] two\n- [x] three\n"
	if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	tc := readTaskCount(ws, "add-foo")
	if tc.Done != 3 || tc.Total != 3 || !tc.Found {
		t.Errorf("readTaskCount = %+v, want 3/3 found", tc)
	}
	if !tc.Archived {
		t.Error("archived tasks should be marked Archived")
	}
}

func TestReadTaskCount_ActiveWinsOverArchive(t *testing.T) {
	ws := t.TempDir()
	active := filepath.Join(ws, "openspec", "changes", "add-foo")
	arch := filepath.Join(ws, "openspec", "changes", "archive", "2026-06-04-add-foo")
	for _, d := range []string{active, arch} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(active, "tasks.md"), []byte("- [ ] open\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(arch, "tasks.md"), []byte("- [x] done\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tc := readTaskCount(ws, "add-foo")
	if tc.Archived || tc.Done != 0 || tc.Total != 1 {
		t.Errorf("active location must win: got %+v", tc)
	}
}

func TestReadTaskCount_SuffixCollisionGuard(t *testing.T) {
	ws := t.TempDir()
	// An archived "do-foo" must not be matched when querying "foo".
	dir := filepath.Join(ws, "openspec", "changes", "archive", "2026-06-04-do-foo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tasks.md"), []byte("- [x] x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	tc := readTaskCount(ws, "foo")
	if tc.Found {
		t.Errorf("should not match do-foo when querying foo: %+v", tc)
	}
}
