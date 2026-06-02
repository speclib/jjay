# Proposal: Go as implementation language

**Change**: techstack-go
**Status**: proposed
**Bean**: [jjay-byyv — techstack / language choice](../../../.beans/jjay-byyv--techstack-language-choice.md)

## Summary

Use Go as the implementation language for the jjay CLI.

## Why Go

- **Team familiarity** — the team already knows Go, so contributions and reviews are frictionless
- **jjay is an orchestrator** — it shells out to jj, tmux, and AI agents; Go's `os/exec` and goroutines make this natural
- **CLI ecosystem** — cobra for arg parsing, bubbletea for TUI, well-trodden path
- **Single binary** — zero runtime dependencies, easy Nix packaging, simple distribution
- **Learning the workflow, not the language** — Go stays out of the way so focus stays on the product

## What this means concretely

- Initialize a Go module (`github.com/pim/jjay` or similar)
- Use cobra for CLI structure
- Target Go 1.24+
- Add a Nix flake for builds

## Non-goals

- No TUI in the first iteration (bubbletea comes later)
- No cross-compilation setup yet — build for the host platform first
