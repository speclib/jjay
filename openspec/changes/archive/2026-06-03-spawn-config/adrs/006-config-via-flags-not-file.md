# ADR-006: Configuration via CLI flags, not config file

**Status**: Proposed

## Context

Spawn hardcodes the agent command, tmux session, and workspace root. These need to be configurable for integration testing (fake agent, dedicated session) and multi-agent support (codex, mistral). The question is where configuration lives.

## Options Considered

- **CLI flags** — `--agent`, `--session`, `--workspace-root`. Simple, no parsing, no file to manage. Sufficient for current needs.
- **YAML config file** (e.g., `.jjay.yaml`) — more powerful, supports defaults per project. But premature — we don't know what the config surface will look like yet.
- **Environment variables** — `JJAY_AGENT`, `JJAY_SESSION`. Middle ground but less discoverable than flags.

## Decision

CLI flags on spawn and cleanup commands. No config file until usage patterns emerge. Flags are easy to add, easy to discover (via `--help`), and don't require a file format decision.

## Consequences

- **Positive**: Zero config files to manage, flags are self-documenting
- **Positive**: Integration tests can pass flags directly
- **Negative**: Repetitive for users who always use the same agent — but that's a future config file problem
