# Proposal: Fill in openspec/config.yaml

**Change**: config-yaml
**Status**: proposed
**Bean**: [jjay-qzin — config.yaml](../../../.beans/jjay-qzin--configyaml.md)

## Summary

Replace the placeholder openspec config with project-specific context and light rules tailored to jjay's scope as a Go CLI orchestrator.

## Config content

```yaml
schema: spec-driven

context: |
  Tech stack: Go (1.24+), cobra (CLI), bubbletea (future TUI)
  Build: single binary, Nix flake
  Platform: Linux and macOS (tmux is Unix-only)

  jjay is a CLI orchestrator for parallel AI agent sessions.
  It shells out to jj, tmux, and AI agents (claude, codex).
  It does not do heavy computation — it coordinates.

  Domain language:
  - "workspace" = jj workspace (isolated working copy)
  - "session" = tmux window running an agent
  - "change" = openspec change being worked on
  - "spawn" = create workspace + session + launch agent
  - "merge" = integrate workspace back into main
  - "cleanup" = forget workspace, remove dir, kill session

rules:
  specs:
    - Write requirements as observable CLI behavior and outcomes
    - Specify which external tools (jj, tmux, agent) each command depends on
  tasks:
    - Scope tasks to a single command or lifecycle phase
  design:
    - Document subprocess invocations and their expected outputs
```

## Non-goals

- No cross-platform rules (tmux is Unix-only, no Windows support planned)
- No heavy per-artifact rule sets — keep it light and let conventions emerge
