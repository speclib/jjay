# ADR-001: Use Go as implementation language

**Status**: Accepted

## Context

jjay is a CLI orchestrator — it shells out to jj, tmux, and AI agents (claude, codex). It doesn't do heavy computation. The team needed to pick an implementation language from candidates: Rust, Go, Zig, Python.

This is also a learning project where the focus should be on the workflow and product, not fighting a new language.

## Options Considered

- **Go** — Team knows it, excellent CLI ecosystem (cobra, bubbletea), single binary, great subprocess handling via `os/exec` and goroutines
- **Rust** — Strong CLI ecosystem (clap, ratatui), single binary, but steep learning curve if new
- **Python** — Team knows it, fast prototyping, but no single binary without extra packaging
- **Zig** — Single binary, minimal runtime, but immature CLI ecosystem

## Decision

Go. The team already knows it, so contributions and reviews are frictionless. The CLI ecosystem is mature (cobra for args, bubbletea for future TUI). Single binary output simplifies Nix packaging and distribution. Go stays out of the way so focus stays on the product.

## Consequences

- **Positive**: Fast onboarding for team members, proven CLI tooling, easy cross-compilation later
- **Positive**: `os/exec` and goroutines are a natural fit for orchestrating subprocesses
- **Negative**: No sum types (error handling is verbose compared to Rust)
- **Negative**: Generics are still maturing (unlikely to matter for an orchestrator)
