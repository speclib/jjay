# ADR-002: OpenSpec config — project context and light rules

**Status**: Accepted

## Context

The openspec `config.yaml` ships as a placeholder with commented-out examples. AI agents generating artifacts (proposals, specs, designs, tasks) need project-specific context to produce relevant output. Without it, artifacts are generic and require heavy editing.

## Options Considered

- **Minimal context only** — Just list the tech stack, no per-artifact rules. Simplest, but rules can catch common mistakes.
- **Heavy rules** — Detailed per-artifact constraints like the OpenSpec project itself (cross-platform path rules, CI verification tasks). Overkill for a Unix-only orchestrator in early development.
- **Context + light rules** — Tech stack, domain language, and a few rules per artifact. Enough to guide AI without over-constraining.

## Decision

Context + light rules. The config includes:
- Tech stack (Go, cobra, bubbletea, Nix)
- Platform scope (Linux and macOS — tmux is Unix-only)
- Domain language definitions (workspace, session, change, spawn, merge, cleanup)
- Light rules: specs focus on observable CLI behavior, tasks scoped to single commands, design documents subprocess invocations

## Consequences

- **Positive**: AI-generated artifacts are grounded in the project's reality from the start
- **Positive**: Domain language definitions prevent ambiguity in specs
- **Negative**: Rules may need updating as the project evolves (but that's cheap)
