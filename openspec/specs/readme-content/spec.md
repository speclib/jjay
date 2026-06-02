# Spec: README Content

## Hero

- Centered at top of file
- Source: `artwork/hero.png`
- No caption needed — the image contains the tagline

## Description

- One-liner: what jjay is and what it does
- Mention jj, tmux, openspec as the three pillars
- State alpha status clearly

## The workflow jjay automates

Document the manual 6-step process:

1. **Spawn** — create jj workspace + tmux window, cd into workspace, launch coding agent with openspec task
2. **Work** — agent implements the change (repeat step 1 for parallel agents)
3. **Test** — manual testing of the result
4. **Archive** — openspec archive, jj describe, changelog update
5. **Merge** — `jj new` to merge, `jj bookmark set main`
6. **Cleanup** — `jj workspace forget`, rm workspace dir, kill tmux window

Include representative shell commands for each step.

## Prerequisites

List as requirements:
- jj (Jujutsu) — version control / workspace isolation
- tmux — session/window management
- openspec — change tracking and task specs

## Installation

Single line: "Coming soon."

## CLI preview

Show planned commands:
- `jjay spawn <change>` — create workspace + tmux window + launch agent
- `jjay status` — show running agents and their state
- `jjay merge <change>` — merge workspace into main
- `jjay cleanup <change>` — forget workspace, remove dir, kill window

Mark clearly as planned/not yet implemented.

## Roadmap

Include at minimum:
- Core lifecycle commands (spawn, merge, cleanup)
- Agent status monitoring
- Multiple agent support (Claude, Codex, Mistral)
- Nix develop integration for workspace environments
- Configurable tmux layouts

## Contributing

- State that contributions are welcome
- Keep it brief — link to issues, standard fork+PR workflow
- No CONTRIBUTING.md file yet

## License

- Reference the LICENSE file in the repo
