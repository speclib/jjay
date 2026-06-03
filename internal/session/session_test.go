package session

import (
	"os"
	"path/filepath"
	"testing"
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
