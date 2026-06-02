<p align="center">
  <img src="artwork/hero.png" alt="jjay — Control the flock. Manage parallel agent sessions with jj, tmux and openspec." />
</p>

# jjay

Manage parallel AI agent sessions with **jj**, **tmux**, and **openspec**.

> **Alpha** — jjay is under active development. Nothing works yet.

## What jjay automates

Running multiple coding agents in parallel (Claude, Codex, Mistral) requires a repetitive manual workflow. This is the process jjay will replace:

### 1. Spawn a workspace

```bash
# Create a new tmux window
tmux new-window -n "feat/payments"

# Create an isolated jj workspace
jj workspace add ../myproject-workspaces/feat-payments
cd ../myproject-workspaces/feat-payments

# Launch a coding agent on the task
claude "/opsx:apply feat-payments" --dangerously-skip-permissions
```

### 2. Repeat for parallel agents

Spin up as many workspaces as you need — each agent works in isolation.

### 3. Test

Manually verify the results in each workspace.

### 4. Archive the change

```bash
openspec archive --change feat-payments
jj describe -m "feat: add payment processing"
```

### 5. Merge into main

```bash
jj new main feat-payments -m "merge feat-payments into main"
jj bookmark set main -r @
```

### 6. Cleanup

```bash
jj workspace forget feat-payments
rm -rf ../myproject-workspaces/feat-payments
tmux kill-window -t "feat/payments"
```

jjay will handle all of this with a single command.

## Prerequisites

- [jj (Jujutsu)](https://martinvonz.github.io/jj/) — version control and workspace isolation
- [tmux](https://github.com/tmux/tmux) — terminal session and window management
- [openspec](https://github.com/speclib/openspec) — change tracking and task specs

## Installation

Coming soon.

## CLI preview

Planned commands (not yet implemented):

```
jjay spawn <change>     Create workspace + tmux window + launch agent
jjay status             Show running agents and their state
jjay merge <change>     Merge workspace into main
jjay cleanup <change>   Forget workspace, remove dir, kill window
```

## Roadmap

- Core lifecycle commands (spawn, merge, cleanup)
- Agent status monitoring
- Multiple agent support (Claude, Codex, Mistral)
- Configurable tmux layouts
- Nix develop integration for workspace environments

## Contributing

Contributions are welcome. Fork the repo, create a branch, and open a pull request.

Found a bug or have an idea? [Open an issue](../../issues).

## License

[MIT](LICENSE)
