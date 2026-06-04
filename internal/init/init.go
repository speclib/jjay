// Package init implements `jjay init`: an idempotent, non-destructive
// orchestrator that prepares a target project for jjay (ADR-008). It delegates
// openspec scaffolding to `openspec init` and jj scaffolding to jj's own
// commands, and owns only the jjay-specific assets (the `/jjay:*` commands, the
// `jjay` skill, and AGENTS.md). Each step detects whether its artifact already
// exists and skips it when present, never clobbering a user file without
// --force.
package init

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// assets holds the canonical `/jjay:*` command files and the `jjay` skill,
// copied from the repo's `.claude/` (the source authored in add-claude-commands).
// A drift test asserts these match the live `.claude/` content. go:embed cannot
// reach outside the package directory, so the canonical copies live here and the
// test guards them against divergence.
//
//go:embed all:assets
var assets embed.FS

// InitOptions controls the init pipeline.
type InitOptions struct {
	Yes        bool // accept creation defaults without prompting (does NOT authorize overwrite)
	Force      bool // overwrite existing user files
	WithJJ     bool // initialize a jj repo if absent
	WithHooks  bool // scaffold example (commented) hooks
	NoClaude   bool // skip installing the jjay Claude integration
	NoOpenspec bool // skip the openspec step
	NoAgents   bool // skip the AGENTS.md step

	// Out is where progress is written. Defaults to os.Stdout when nil.
	Out io.Writer
}

func (o InitOptions) out() io.Writer {
	if o.Out != nil {
		return o.Out
	}
	return os.Stdout
}

// stepResult is what each step reports for one artifact.
type stepResult int

const (
	resultCreated        stepResult = iota // the artifact was created
	resultSkipped                          // already present and valid; left unchanged
	resultOverwritten                      // existed and was overwritten under --force
	resultWouldOverwrite                   // existed; left in place because --force was not set
)

// Init prepares the project at path (default: cwd when empty) for orchestration
// by jjay. Steps run in order: openspec → claude integration → AGENTS.md →
// (jj) → (hooks). A failed step is reported and aborts the run; completed steps
// remain, and re-running resumes via idempotency (ADR-008, no rollback).
func Init(path string, opts InitOptions) error {
	target, err := resolveTarget(path)
	if err != nil {
		return err
	}

	w := opts.out()
	fmt.Fprintf(w, "Initializing jjay project at %s\n", target)

	steps := []struct {
		name string
		run  func(string, InitOptions) error
		skip bool
	}{
		{"openspec", stepOpenspec, opts.NoOpenspec},
		{"claude", stepClaude, opts.NoClaude},
		{"agents", stepAgents, opts.NoAgents},
		{"jj", stepJJ, !opts.WithJJ},
		{"hooks", stepHooks, !opts.WithHooks},
	}

	for _, s := range steps {
		if s.skip {
			continue
		}
		if err := s.run(target, opts); err != nil {
			return fmt.Errorf("step %s: %w", s.name, err)
		}
	}

	fmt.Fprintln(w, "Done.")
	return nil
}

// resolveTarget returns an absolute path for the target dir, defaulting to cwd
// when path is empty, and verifies it exists and is a directory.
func resolveTarget(path string) (string, error) {
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to resolve current directory: %w", err)
		}
		path = cwd
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve %q: %w", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("target %q: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("target %q is not a directory", abs)
	}
	return abs, nil
}

// report prints a per-artifact result line in the cleanup-step style.
func report(opts InitOptions, label string, r stepResult) {
	w := opts.out()
	switch r {
	case resultCreated:
		fmt.Fprintf(w, "  %s: created\n", label)
	case resultSkipped:
		fmt.Fprintf(w, "  %s: already present, skipped\n", label)
	case resultOverwritten:
		fmt.Fprintf(w, "  %s: overwritten (--force)\n", label)
	case resultWouldOverwrite:
		fmt.Fprintf(w, "  %s: exists, left in place (use --force to overwrite)\n", label)
	}
}

// writeFile writes content to dst non-destructively: it creates dst if missing,
// skips it if present (unless force), and reports which happened. Parent
// directories are created as needed.
func writeFile(opts InitOptions, dst string, content []byte) (stepResult, error) {
	exists := false
	if _, err := os.Stat(dst); err == nil {
		exists = true
	} else if !os.IsNotExist(err) {
		return 0, fmt.Errorf("stat %s: %w", dst, err)
	}

	if exists && !opts.Force {
		return resultWouldOverwrite, nil
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", filepath.Dir(dst), err)
	}
	if err := os.WriteFile(dst, content, 0o644); err != nil {
		return 0, fmt.Errorf("write %s: %w", dst, err)
	}
	if exists {
		return resultOverwritten, nil
	}
	return resultCreated, nil
}

// stepOpenspec ensures openspec is initialized for the target via `openspec
// init`, then ensures config.yaml exists. It does not reimplement openspec
// scaffolding (ADR-008).
func stepOpenspec(target string, opts InitOptions) error {
	w := opts.out()
	fmt.Fprintln(w, "openspec:")

	if _, err := exec.LookPath("openspec"); err != nil {
		return fmt.Errorf("the 'openspec' binary is not on PATH; install openspec and re-run (https://github.com/Fission-AI/OpenSpec)")
	}

	openspecDir := filepath.Join(target, "openspec")
	if dirExists(openspecDir) {
		report(opts, "openspec/", resultSkipped)
	} else {
		args := []string{"init", target, "--tools", "claude"}
		if opts.Force {
			args = append(args, "--force")
		}
		cmd := exec.Command("openspec", args...)
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run 'openspec init': %w", err)
		}
		report(opts, "openspec/", resultCreated)
	}

	// Ensure config.yaml exists. `openspec init --tools <x>` runs
	// non-interactively and skips writing config.yaml, so seed a minimal one
	// from a template (the project can fill in its context later). This is
	// non-destructive: an existing config is never overwritten without --force.
	configPath := filepath.Join(openspecDir, "config.yaml")
	r, err := writeFile(opts, configPath, []byte(configTemplate))
	if err != nil {
		return err
	}
	report(opts, "openspec/config.yaml", r)
	return nil
}

// stepClaude installs the embedded `/jjay:*` commands and the `jjay` skill into
// the target's .claude/, non-destructively.
func stepClaude(target string, opts InitOptions) error {
	w := opts.out()
	fmt.Fprintln(w, "claude integration:")

	return fs.WalkDir(assets, "assets", func(p string, d fs.DirEntry, err error) error {
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
		content, err := assets.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", p, err)
		}
		dst := filepath.Join(target, ".claude", rel)
		r, err := writeFile(opts, dst, content)
		if err != nil {
			return err
		}
		report(opts, filepath.Join(".claude", rel), r)
		return nil
	})
}

// stepAgents writes AGENTS.md documenting the jjay conventions, non-destructively.
func stepAgents(target string, opts InitOptions) error {
	w := opts.out()
	fmt.Fprintln(w, "AGENTS.md:")

	dst := filepath.Join(target, "AGENTS.md")
	r, err := writeFile(opts, dst, []byte(agentsTemplate))
	if err != nil {
		return err
	}
	report(opts, "AGENTS.md", r)
	return nil
}

// stepJJ initializes a jj repo via jj's own command if not already present.
func stepJJ(target string, opts InitOptions) error {
	w := opts.out()
	fmt.Fprintln(w, "jj:")

	if dirExists(filepath.Join(target, ".jj")) {
		report(opts, ".jj/", resultSkipped)
		return nil
	}
	if _, err := exec.LookPath("jj"); err != nil {
		return fmt.Errorf("the 'jj' binary is not on PATH; install jujutsu and re-run, or omit --with-jj")
	}
	cmd := exec.Command("jj", "git", "init")
	cmd.Dir = target
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'jj git init': %w", err)
	}
	report(opts, ".jj/", resultCreated)
	return nil
}

// stepHooks scaffolds an example (commented) hooks file the user can enable.
func stepHooks(target string, opts InitOptions) error {
	w := opts.out()
	fmt.Fprintln(w, "hooks:")

	dst := filepath.Join(target, ".jjay", "hooks.example.sh")
	r, err := writeFile(opts, dst, []byte(hooksExample))
	if err != nil {
		return err
	}
	report(opts, filepath.Join(".jjay", "hooks.example.sh"), r)
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
