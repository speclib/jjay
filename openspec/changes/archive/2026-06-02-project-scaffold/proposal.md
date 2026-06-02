# Proposal: Project scaffold

**Change**: project-scaffold
**Status**: proposed
**Bean**: [jjay-x3k5 — project setup for Go project](../../../.beans/jjay-x3k5--project-setup-for-go-project.md)

## Why

jjay needs a Go project foundation before any features can be built. No code exists yet — we need a module, entry point, directory structure, test framework, and Nix flake in one shot.

## What Changes

- Initialize Go module with cobra dependency
- Create `cmd/jjay/main.go` with cobra root command and `version` subcommand
- Set up directory structure (`cmd/`, `internal/`)
- Add Makefile for dev workflow (build, test, lint)
- Add `flake.nix` with buildGoModule, dev shell, multi-platform support
- Add initial test file to establish testing pattern

## Capabilities

### New Capabilities

- `project-scaffold`: Go module, directory layout, cobra CLI entry point with `jjay version`
- `test-infrastructure`: Go testing setup with Makefile dev workflow
- `nix-flake`: Nix flake for building, running, and developing jjay

### Modified Capabilities

_(none — greenfield project)_

## Impact

- New files: `go.mod`, `go.sum`, `cmd/jjay/main.go`, `internal/`, `Makefile`, `flake.nix`, `flake.lock`
- Modified files: `README.md` (Nix installation section)
