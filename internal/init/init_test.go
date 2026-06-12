package init

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot walks up from the package dir to the repo root (the dir containing
// go.mod), so the drift test can find the live .claude/ content.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repo root (go.mod) above %s", dir)
		}
		dir = parent
	}
}

// TestEmbeddedAssetsMatchLiveClaude is the drift guard (task 3.3): the embedded
// assets jjay installs must be byte-identical to the live .claude/ content they
// are copied from, so a dogfooded change to .claude/ cannot silently diverge
// from what `jjay init` writes into a target.
func TestEmbeddedAssetsMatchLiveClaude(t *testing.T) {
	root := repoRoot(t)
	err := fs.WalkDir(assets, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel("assets", p)
		if err != nil {
			return err
		}
		embedded, err := assets.ReadFile(p)
		if err != nil {
			t.Fatalf("read embedded %s: %v", p, err)
		}
		live, err := os.ReadFile(filepath.Join(root, ".claude", rel))
		if err != nil {
			t.Errorf("live .claude/%s missing or unreadable: %v", rel, err)
			return nil
		}
		if !bytes.Equal(embedded, live) {
			t.Errorf("embedded assets/%s differs from live .claude/%s — re-copy to fix drift", rel, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}

// fakeOpenspec puts a fake `openspec` on PATH that creates a minimal openspec/
// tree in the target so the openspec step can run without the real binary.
func fakeOpenspec(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("openspec"); err == nil {
		// A real binary is present; tests that need openspec can use it. But to
		// keep tests hermetic we still shadow it with a fake.
	}
	bin := t.TempDir()
	script := `#!/usr/bin/env bash
# fake openspec init: positional path is $2 (init <path> --tools claude ...)
if [ "$1" = "init" ]; then
  target="$2"
  mkdir -p "$target/openspec"
  printf 'schema: spec-driven\n' > "$target/openspec/config.yaml"
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(bin, "openspec"), []byte(script), 0o755); err != nil {
		t.Fatalf("write fake openspec: %v", err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func mustExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to exist: %v", path, err)
	}
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("expected %s to NOT exist", path)
	}
}

// TestInit_BareProject: a bare project becomes fully initialized (task 6.1).
func TestInit_BareProject(t *testing.T) {
	fakeOpenspec(t)
	target := t.TempDir()
	var buf bytes.Buffer

	if err := Init(target, InitOptions{Yes: true, Out: &buf}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	mustExist(t, filepath.Join(target, "openspec", "config.yaml"))
	mustExist(t, filepath.Join(target, ".claude", "commands", "jjay", "spawn.md"))
	mustExist(t, filepath.Join(target, ".claude", "skills", "jjay", "SKILL.md"))
	mustExist(t, filepath.Join(target, "AGENTS.md"))
	// jjay config is seeded from the built-in (ADR-014), carrying launch+resume.
	cfg := filepath.Join(target, ".jjay", "config.yaml")
	mustExist(t, cfg)
	data, _ := os.ReadFile(cfg)
	for _, want := range []string{"agents:", "claude:", "launch:", "resume:", "--resume"} {
		if !strings.Contains(string(data), want) {
			t.Errorf(".jjay/config.yaml missing %q; got:\n%s", want, data)
		}
	}
	// jj and hooks are opt-in: not present without their flags.
	mustNotExist(t, filepath.Join(target, ".jjay", "hooks.example.sh"))
}

// TestInit_ReRunIsNoOp: re-running on a prepared project leaves files unchanged
// and reports skipped (task 6.1).
func TestInit_ReRunIsNoOp(t *testing.T) {
	fakeOpenspec(t)
	target := t.TempDir()

	if err := Init(target, InitOptions{Yes: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("first Init: %v", err)
	}

	agents := filepath.Join(target, "AGENTS.md")
	// Mutate AGENTS.md; a no-op re-run must preserve the user's content.
	if err := os.WriteFile(agents, []byte("USER EDIT"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	var buf bytes.Buffer
	if err := Init(target, InitOptions{Yes: true, Out: &buf}); err != nil {
		t.Fatalf("second Init: %v", err)
	}

	got, _ := os.ReadFile(agents)
	if string(got) != "USER EDIT" {
		t.Errorf("re-run clobbered AGENTS.md: got %q", string(got))
	}
	if !bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Errorf("expected re-run to report skipped artifacts, output:\n%s", buf.String())
	}
}

// TestInit_PartialCompletesMissing: a project with openspec but no jjay Claude
// integration gets the missing pieces, openspec left unchanged (task 6.1).
func TestInit_PartialCompletesMissing(t *testing.T) {
	fakeOpenspec(t)
	target := t.TempDir()

	// Pre-create openspec/ with a sentinel config the run must not touch.
	if err := os.MkdirAll(filepath.Join(target, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(target, "openspec", "config.yaml")
	if err := os.WriteFile(cfg, []byte("SENTINEL"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Init(target, InitOptions{Yes: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	mustExist(t, filepath.Join(target, ".claude", "commands", "jjay", "spawn.md"))
	got, _ := os.ReadFile(cfg)
	if string(got) != "SENTINEL" {
		t.Errorf("existing openspec config was modified: got %q", string(got))
	}
}

// TestInit_YesDoesNotClobber: --yes creates missing files but never overwrites
// an existing one (task 6.2).
func TestInit_YesDoesNotClobber(t *testing.T) {
	fakeOpenspec(t)
	target := t.TempDir()

	agents := filepath.Join(target, "AGENTS.md")
	if err := os.WriteFile(agents, []byte("EXISTING"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := Init(target, InitOptions{Yes: true, Out: &buf}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	got, _ := os.ReadFile(agents)
	if string(got) != "EXISTING" {
		t.Errorf("--yes clobbered existing AGENTS.md: got %q", string(got))
	}
	if !bytes.Contains(buf.Bytes(), []byte("--force")) {
		t.Errorf("expected a 'use --force' hint for the existing file, output:\n%s", buf.String())
	}
}

// TestInit_ForceOverwrites: --force overwrites an existing file (task 6.2).
func TestInit_ForceOverwrites(t *testing.T) {
	fakeOpenspec(t)
	target := t.TempDir()

	agents := filepath.Join(target, "AGENTS.md")
	if err := os.WriteFile(agents, []byte("OLD"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Init(target, InitOptions{Yes: true, Force: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("Init: %v", err)
	}

	got, _ := os.ReadFile(agents)
	if string(got) == "OLD" {
		t.Errorf("--force did not overwrite AGENTS.md")
	}
	if !bytes.Contains(got, []byte("AGENTS.md")) {
		t.Errorf("AGENTS.md not replaced with template, got: %q", string(got))
	}
}

// TestInit_OpenspecMissingBinary: a clear error when openspec is absent (task 2.3).
func TestInit_OpenspecMissingBinary(t *testing.T) {
	// Empty PATH so neither openspec nor jj resolve.
	t.Setenv("PATH", t.TempDir())
	target := t.TempDir()

	err := Init(target, InitOptions{Yes: true, Out: &bytes.Buffer{}})
	if err == nil {
		t.Fatal("expected an error when openspec is unavailable, got nil")
	}
	if !bytes.Contains([]byte(err.Error()), []byte("openspec")) {
		t.Errorf("expected error to mention openspec, got: %v", err)
	}
}

// TestInit_SkipOpenspec lets the other steps run without an openspec binary,
// so we can test claude/agents in isolation.
func TestInit_SkipOpenspec(t *testing.T) {
	target := t.TempDir()
	if err := Init(target, InitOptions{Yes: true, NoOpenspec: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("Init: %v", err)
	}
	mustExist(t, filepath.Join(target, ".claude", "skills", "jjay", "SKILL.md"))
	mustExist(t, filepath.Join(target, "AGENTS.md"))
}

// TestInit_NoClaudeSkipsIntegration verifies --no-claude skips the .claude install.
func TestInit_NoClaudeSkipsIntegration(t *testing.T) {
	target := t.TempDir()
	if err := Init(target, InitOptions{Yes: true, NoOpenspec: true, NoClaude: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("Init: %v", err)
	}
	mustNotExist(t, filepath.Join(target, ".claude"))
}

// TestInit_WithHooksScaffolds verifies --with-hooks writes the example file.
func TestInit_WithHooksScaffolds(t *testing.T) {
	target := t.TempDir()
	if err := Init(target, InitOptions{Yes: true, NoOpenspec: true, WithHooks: true, Out: &bytes.Buffer{}}); err != nil {
		t.Fatalf("Init: %v", err)
	}
	mustExist(t, filepath.Join(target, ".jjay", "hooks.example.sh"))
}
