package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jjay/internal/spawn"
	"jjay/internal/status"
)

func TestSessionName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/home/user/projects/myapp", "jjay->myapp"},
		{"/tmp/foo", "jjay->foo"},
		{"relative-dir", "jjay->relative-dir"},
	}

	for _, tt := range tests {
		got := SessionName(tt.path)
		if got != tt.want {
			t.Errorf("SessionName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestCheckJJRepo_Valid(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".jj"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := checkJJRepo(dir); err != nil {
		t.Errorf("expected no error for valid jj repo, got: %v", err)
	}
}

func TestCheckJJRepo_Invalid(t *testing.T) {
	dir := t.TempDir()

	if err := checkJJRepo(dir); err == nil {
		t.Error("expected error for non-jj directory, got nil")
	}
}

func TestCheckJJRepo_FileNotDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".jj"), []byte("not a dir"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := checkJJRepo(dir); err == nil {
		t.Error("expected error when .jj is a file not a directory, got nil")
	}
}

// recordingOpener records which changes it was asked to open and can be told
// to fail for specific changes.
type recordingOpener struct {
	opened  []string
	failFor map[string]bool
}

func (r *recordingOpener) open(change, wsDir string, _ spawn.SpawnOptions) error {
	if r.failFor[change] {
		return fmt.Errorf("boom for %s", change)
	}
	r.opened = append(r.opened, change)
	return nil
}

func TestReopenDetached_ReopensAllDetached(t *testing.T) {
	spawns := []status.Spawn{
		{Name: "app-add-foo", Change: "add-foo", WSDir: "/ws/add-foo", Attached: false},
		{Name: "app-fix-bar", Change: "fix-bar", WSDir: "/ws/fix-bar", Attached: false},
	}
	rec := &recordingOpener{}
	var out strings.Builder
	reopenDetached(spawns, "sess", rec.open, &out)

	if len(rec.opened) != 2 {
		t.Fatalf("expected 2 reopened, got %v", rec.opened)
	}
}

func TestReopenDetached_None(t *testing.T) {
	rec := &recordingOpener{}
	var out strings.Builder
	reopenDetached(nil, "sess", rec.open, &out)
	if len(rec.opened) != 0 {
		t.Errorf("expected nothing reopened, got %v", rec.opened)
	}
}

func TestReopenDetached_SkipsAttached_NoDuplicate(t *testing.T) {
	spawns := []status.Spawn{
		{Name: "app-add-foo", Change: "add-foo", WSDir: "/ws/add-foo", Attached: true},  // already has window
		{Name: "app-fix-bar", Change: "fix-bar", WSDir: "/ws/fix-bar", Attached: false}, // needs reopen
	}
	rec := &recordingOpener{}
	var out strings.Builder
	reopenDetached(spawns, "sess", rec.open, &out)

	if len(rec.opened) != 1 || rec.opened[0] != "app-fix-bar" {
		t.Errorf("expected only app-fix-bar reopened, got %v", rec.opened)
	}
}

func TestReopenDetached_OneFailsNonFatal(t *testing.T) {
	spawns := []status.Spawn{
		{Name: "app-add-foo", Change: "add-foo", WSDir: "/ws/add-foo", Attached: false},
		{Name: "app-fix-bar", Change: "fix-bar", WSDir: "/ws/fix-bar", Attached: false},
	}
	rec := &recordingOpener{failFor: map[string]bool{"app-add-foo": true}}
	var out strings.Builder
	reopenDetached(spawns, "sess", rec.open, &out)

	// fix-bar still reopened despite add-foo failing.
	if len(rec.opened) != 1 || rec.opened[0] != "app-fix-bar" {
		t.Errorf("expected app-fix-bar reopened, got %v", rec.opened)
	}
	// Failure is reported.
	if !strings.Contains(out.String(), "add-foo") {
		t.Errorf("expected failure report to mention add-foo, got: %q", out.String())
	}
}
