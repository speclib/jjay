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

const sampleWSList = `default: abc123 (no description set)
add-foo: def456 (no description set)
fix-bar: 789abc (no description set)
`

// noTasks is a taskCounter that reports nothing found, for join tests that
// don't care about progress counts.
func noTasks(_, _ string) TaskCount { return TaskCount{} }

const mainRoot = "/repo/main"
const wsRoot = "/repo/ws"

func TestJoin_AttachedAndDetached(t *testing.T) {
	windows := map[string]bool{
		workspace.WindowName("add-foo"): true,
		// fix-bar has no window
	}

	spawns, err := join(sampleWSList, windows, mainRoot, wsRoot, noTasks)
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
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks)
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
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	if _, ok := spawnByChange(spawns, "default"); ok {
		t.Error("default workspace must be excluded from spawns")
	}
}

func TestJoin_WSDirResolved(t *testing.T) {
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, noTasks)
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
	spawns, err := join(sampleWSList, map[string]bool{}, mainRoot, wsRoot, tasks)
	if err != nil {
		t.Fatalf("join: %v", err)
	}
	foo, _ := spawnByChange(spawns, "add-foo")
	if foo.Tasks.Done != 4 || foo.Tasks.Total != 10 {
		t.Errorf("expected 4/10 tasks for add-foo, got %+v", foo.Tasks)
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

	for _, want := range []string{"add-foo", "../ws/add-foo", "12/18 (66%)", "attached", "fix-bar", "detached", "ARCHIVED"} {
		if !strings.Contains(out, want) {
			t.Errorf("render output missing %q:\n%s", want, out)
		}
	}
	// The absolute workspace path must NOT appear — only the relative form.
	if strings.Contains(out, "/repo/ws/add-foo") {
		t.Errorf("render leaked absolute path:\n%s", out)
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
