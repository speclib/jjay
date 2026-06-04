package openspec

import (
	"os/exec"
	"testing"
)

func TestParseChangeNames_ParsesList(t *testing.T) {
	out := []byte(`{"changes":[{"name":"add-foo"},{"name":"fix-bar"}]}`)
	names, err := parseChangeNames(out)
	if err != nil {
		t.Fatalf("parseChangeNames: %v", err)
	}
	if len(names) != 2 || names[0] != "add-foo" || names[1] != "fix-bar" {
		t.Errorf("parseChangeNames = %v, want [add-foo fix-bar]", names)
	}
}

func TestParseChangeNames_Empty(t *testing.T) {
	names, err := parseChangeNames([]byte(`{"changes":[]}`))
	if err != nil {
		t.Fatalf("parseChangeNames: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected no names, got %v", names)
	}
}

func TestParseChangeNames_InvalidJSON(t *testing.T) {
	if _, err := parseChangeNames([]byte("not json")); err == nil {
		t.Error("expected error on invalid JSON, got nil")
	}
}

// TestChangeNames_ToleratesMissingBinary verifies ChangeNames surfaces an error
// (rather than panicking) when the openspec binary cannot be run. We can only
// assert the error path deterministically; the happy path needs a real binary.
func TestChangeNames_ToleratesMissingBinary(t *testing.T) {
	if _, err := exec.LookPath("openspec"); err == nil {
		t.Skip("openspec binary present; cannot test the missing-binary path")
	}
	if _, err := ChangeNames(); err == nil {
		t.Error("expected error when openspec is unavailable, got nil")
	}
}
