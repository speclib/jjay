# Proposal: jjay session-open command

**Change**: session-open
**Status**: proposed
**Bean**: [jjay-4zwb — jjay command to create a new session in a certain jj repo](../../../.beans/jjay-4zwb--jjay-command-to-create-a-new-session-in-a-certain.md)

## Why

Currently you must manually create a tmux session and cd into a repo before using jjay. When managing multiple projects with parallel agents, you want dedicated tmux sessions per project: `jjay->proj1`, `jjay->proj2`. Each session is a home base for spawning, merging, and cleaning up workspaces.

## What Changes

- Add `jjay session-open <path>` cobra subcommand
- Creates tmux session named `jjay-><dirname>` with working directory set to the given path
- Switches the current tmux client to the new session
- Preconditions: path is a jj repo, session doesn't already exist

## Capabilities

### New Capabilities

- `session-open`: Create and switch to a dedicated tmux session for a jj repo

## Impact

- New: `internal/session/session.go`
- Modified: `cmd/jjay/main.go` (add session-open subcommand)
- Modified: `README.md` (add session-open to CLI section)
