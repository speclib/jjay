## Context

Greenfield Go project for jjay, a CLI orchestrator that shells out to jj, tmux, and AI agents. No code exists yet. Based on teejay's proven setup pattern, adapted for a cobra CLI instead of a bubbletea TUI.

## Goals / Non-Goals

**Goals:**
- Working Go module with cobra CLI that compiles and runs
- Makefile for dev workflow
- Nix flake for reproducible builds and dev shell
- Initial test to establish pattern

**Non-Goals:**
- Implementing any jjay features (spawn, merge, cleanup)
- Adding bubbletea TUI (comes later)
- Release automation (separate bean: jjay-fder)
- CI/CD setup

## Decisions

### Directory structure

```
jjay/
├── cmd/jjay/
│   └── main.go          # cobra root command + version subcommand
├── internal/             # private packages (empty for now)
├── go.mod
├── go.sum
├── Makefile
├── flake.nix
└── flake.lock
```

**Rationale**: Standard Go layout. `cmd/` for executables, `internal/` for private packages. Same pattern as teejay.

### Cobra over bare os.Args

Use cobra for the CLI framework from day one.
_Alternative: plain flag package — rejected because jjay will have multiple subcommands (spawn, merge, cleanup, status) and cobra handles that naturally._

### Module path: `jjay`

Simple module name, can update later if published to GitHub.

### Nix flake: buildGoModule with vendorHash

Standard approach. Start with `vendorHash = null`, get real hash from build error, update.
Dev shell includes go and gopls.

### Makefile targets

- `make build` — `go build -o jjay ./cmd/jjay`
- `make test` — `go test ./...`
- `make lint` — `go vet ./...`

No `make help` — keep it minimal.

## Risks / Trade-offs

- [vendorHash maintenance] Must update when go.mod changes. Low friction.
- [Module path] May need `go mod edit -module` later for GitHub publishing.
