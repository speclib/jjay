<p align="center">
  <img src="artwork/hero.png" alt="jjay — Control the flock. Manage parallel agent sessions with jj, tmux and openspec." />
</p>

[![Coverage](https://img.shields.io/badge/coverage-75.7%25-yellow)](coverage.html)

# jjay

Manage parallel AI agent sessions with **jj**, **tmux**, and **openspec**.

> **Alpha** — jjay is under active development. Usage could be risky.

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

## Configuration

jjay has its own config file — `<repo>/.jjay/config.yaml` (project) and
`~/.config/jjay/config.yaml` (global) — **distinct from `openspec/config.yaml`**.
It holds per-agent command templates, resolved per field as
project → global → built-in:

```yaml
agents:
  claude:
    launch: 'claude "/opsx:apply {change}" --dangerously-skip-permissions --add-dir {wsdir}'
    resume: 'claude --resume --add-dir {wsdir}'
```

`launch` starts the work on a first spawn; `resume` is run when a workspace is
reopened (`jjay session-open` / `jjay tmux-open`) so the agent **resumes** its
conversation instead of re-running `/opsx:apply` from scratch. `jjay init` seeds
this file from the built-in defaults.

## Tech stack

- **Go** — single-binary CLI, cobra for commands, bubbletea for future TUI
- [jj (Jujutsu)](https://martinvonz.github.io/jj/) — version control and workspace isolation
- [tmux](https://github.com/tmux/tmux) — terminal session and window management
- [openspec](https://github.com/speclib/openspec) — change tracking and task specs

## Installation

### Nix

```bash
# Run directly
nix run github:mipmip/jjay -- version

# Or add to your flake inputs
```

### From source

```bash
go install ./cmd/jjay
```

## CLI

```
jjay init [path]              Prepare a project for orchestration by jjay
jjay session-open <path>      Create and switch to a tmux session for a jj repo
jjay spawn apply <change>     Isolate an existing change + launch /opsx:apply (app-<change>)
jjay spawn proposal <prompt>  Seed a new proposal spawn from a prompt (prop-<slug>)
jjay status                   List spawned workspaces, task progress, and window state
jjay merge <name>             Merge a spawned workspace into main
jjay cleanup <name>           Tear down workspace + tmux window + directory
jjay version                  Print version
```

### `jjay spawn`

`spawn` takes a **verb**; there is no bare `jjay spawn <change>` form (running
`spawn` with no verb prints usage and exits non-zero):

- **`jjay spawn apply <change>`** — isolate an existing openspec change and run
  `/opsx:apply` inside it. Workspace/window named **`app-<change>`**.
- **`jjay spawn proposal <prompt> [--mode explore|propose]`** — seed a *new*
  proposal from free text, so the orchestrator can keep working while
  exploration/proposal happens in its own window. No openspec change is required
  at spawn time; the agent creates one. `--mode` selects the seed command
  (`/opsx:explore` or `/opsx:propose`, default `explore`).
  - The identity is a **code-derived slug** of the prompt (no AI): lowercase,
    drop punctuation/stopwords, keep salient tokens, cap length, add a uniqueness
    suffix on collision. Workspace/window named **`prop-<slug>`**. The slug is the
    immutable handle — it is **not** renamed after the agent names its change, so
    a `prop-<slug>` workspace may contain a differently-named change directory.

### `jjay init`

Prepares a project (default: the current directory) so jjay can orchestrate it.
It is **idempotent** and **non-destructive**: safe to re-run, and it never
overwrites an existing file without `--force`. Steps run in order:

1. **openspec** — runs `openspec init <path> --tools claude` if `openspec/` is
   absent (delegating to openspec, not reimplementing it), then checks that
   `openspec/config.yaml` exists.
2. **Claude integration** — installs the `/jjay:*` slash commands into
   `<target>/.claude/commands/jjay/` and the `jjay` skill into
   `<target>/.claude/skills/jjay/`. These are embedded from this repo's own
   `.claude/`, so an installed copy is byte-identical to the dogfooded one.
3. **AGENTS.md** — writes an `AGENTS.md` documenting the jjay conventions
   (openspec archive flow, beans tasks, jj usage).
4. **jj** *(opt-in, `--with-jj`)* — initializes a jj repo via `jj git init` if
   none is present.
5. **hooks** *(opt-in, `--with-hooks`)* — scaffolds a commented example hooks
   file you can enable.

Each step reports whether its artifact was **created**, **skipped** (already
present), or left in place (with a hint to pass `--force`).

| Flag | Effect |
| --- | --- |
| `--yes` | Accept creation defaults without prompting. Does **not** authorize overwriting existing files. |
| `--force` | Overwrite existing files (`AGENTS.md`, commands, …). |
| `--with-jj` | Initialize a jj repo if absent. |
| `--with-hooks` | Scaffold the example hooks file. |
| `--no-claude` | Skip installing the jjay Claude integration. |
| `--no-openspec` | Skip the openspec step. |
| `--no-agents` | Skip writing `AGENTS.md`. |

`jjay init` installs the jjay Claude integration **per target project** (into
that project's `.claude/`), not user-wide. Requires the `openspec` binary on
`PATH`; `--with-jj` additionally requires `jj`.

### Shell completion

The arguments of `spawn`, `merge`, and `cleanup` tab-complete, each filtered to
the candidates that verb can actually act on:

- `jjay spawn <TAB>` → the verbs `apply` and `proposal`.
- `jjay spawn apply <TAB>` → openspec changes that do **not** yet have a
  workspace (you can't spawn an already-spawned change).
- `jjay spawn proposal <TAB>` → nothing (it takes a free-text prompt, not a
  candidate name; file-name completion is suppressed).
- `jjay merge <TAB>` / `jjay cleanup <TAB>` → existing spawned workspaces (the
  `default` main working copy is never offered).

Completion is fast and side-effect free — it reads only `openspec list` and
`jj workspace list` (no tmux, no task files) — and degrades silently to no
candidates if a source can't be read.

Install the completion script for your shell with `jjay completion <shell>`
(`bash`, `zsh`, `fish`, or `powershell`); follow the script's own header for
where to source it.

### `jjay status`

Lists every spawned jj workspace with its task progress and whether a matching
`ws-<name>` tmux window exists in the current session. Spawns are split by kind
(from the name prefix): **CHANGES** (`app-*`, tracking an openspec change) and
**PROPOSAL SPAWNS** (`prop-*`, prompt-seeded with no change yet). A table with
no rows is omitted.

```
CHANGES
CHANGE     WORKSPACE                            TASKS         TMUX      MERGED  ARCHIVED
add-foo    ../myproject-workspaces/app-add-foo  12/18 (66%)   attached  no      no
old-feat   ../myproject-workspaces/app-old-feat 5/5 (100%)    detached  yes     yes

PROPOSAL SPAWNS
PROPOSAL    WORKSPACE                               MERGED  TMUX
dark-mode   ../myproject-workspaces/prop-dark-mode  no      attached
```

Proposal spawns omit the change-shaped columns (TASKS/ARCHIVED), which are
meaningless before the agent creates a change.

- **WORKSPACE** is shown **relative to the main repo root**, and is resolved
  correctly even when `jjay status` is run from inside a child workspace.
- **TASKS** is `done/total (percent)`, read from the change's `tasks.md`; `-`
  means no tasks file was found.
- **TMUX** is **attached** when a `ws-<change>` window exists in the current
  session, otherwise **detached** — the workspace is still open on disk, there
  is just no live window/agent for it (e.g. after a detach or reboot). (This
  column was previously named **STATUS**.)
- **MERGED** is **yes** when the spawn's work has already landed on the `main`
  bookmark (derived live from jj: the workspace has no commits `main` lacks). A
  spawn that is **merged but not archived** is the "ready to clean up" signal.
- **ARCHIVED** is **yes** when the change has been archived. Task counts are
  then read from `openspec/changes/archive/<date>-<change>/tasks.md` instead of
  the active `openspec/changes/<change>/tasks.md`, so archived spawns still
  report their progress.

Status is read-only and derives everything live from `jj workspace list` +
`tmux list-windows`; it persists no state (see
[ADR-006](openspec/adrs/006-workspace-is-source-of-truth.md)). With no tmux
server running, every spawn is reported as detached.

### Reopen on `session-open`

The tmux view (windows + agents) is volatile, but jj workspaces are durable.
After creating and switching to the session, `jjay session-open` recreates a
`ws-<change>` window and relaunches the agent for every spawned workspace that
lacks one — restoring the view to match the workspaces on disk. Reopen is
best-effort: if one spawn fails to reopen, the rest still open and session-open
still succeeds, reporting which spawns failed.

## Claude Code integration

jjay ships a Claude Code integration layer under `.claude/` so agents (and you) drive the tool the way it was designed — spawning isolated workspaces rather than applying changes in place. Because `.claude/` is committed and spawned workspaces are jj copies of the repo, these propagate into every spawned workspace automatically.

### `/jjay:*` slash commands

Thin wrappers over the `jjay` binary — one per CLI verb:

| Command | Runs |
| --- | --- |
| `/jjay:spawn <change>` | `jjay spawn <change>` — workspace + tmux window + agent |
| `/jjay:status` | `jjay status` |
| `/jjay:merge <change>` | `jjay merge <change>` |
| `/jjay:cleanup <change>` | `jjay cleanup <change>` |
| `/jjay:session-open <path>` | `jjay session-open <path>` |

They reimplement no logic — the binary is the only source of behavior. Commands for change-name verbs prompt for the change (listing candidates via `openspec list --json` / `jjay status`) if you omit it.

### `jjay` orchestrator skill

`.claude/skills/jjay/SKILL.md` auto-loads when the conversation is about implementing or managing a change in this repo, and encodes the policy: **implement a change by spawning an isolated agent workspace (`/jjay:spawn`), not by running `/opsx:apply` in the main session.** It documents the lifecycle (explore → propose → spawn → status → merge → cleanup) and the **orchestrator-vs-worker** distinction — including the rule that a worker (an agent already running inside a spawned workspace) applies in place and must not recursively spawn.

### Precondition

The `/jjay:*` commands shell out to the `jjay` binary, so **`jjay` must be on `PATH`** in the session (true in spawned workspaces too, since they run in the same environment). See [Installation](#installation).

## Roadmap

- Core lifecycle commands (spawn, merge, cleanup)
- Agent status monitoring
- Multiple agent support (Claude, Codex, Mistral)
- Configurable tmux layouts
- Nix develop integration for workspace environments

## Contributing

Contributions are welcome. Fork the repo, create a branch, and open a pull request.

Found a bug or have an idea? [Open an issue](../../issues).

### Testing & coverage

- `make test` — fast unit tests.
- `make test-integration` — full lifecycle tests (spawn/merge/session); **requires `tmux` and `jj` on `PATH`**.
- `make coverage` — whole-repo coverage including the integration suite (`-tags integration -coverpkg=./...`), so spawn/merge/cleanup report their real numbers, not ~5%. Also requires `tmux` + `jj`, and sweeps test debris first.
- `make coverage-unit` — coverage without the integration tag, for environments lacking `tmux`/`jj` (e.g. bare CI).
- `make badge` — runs `coverage` and patches the README coverage badge. (`coverage` alone only prints the number; **`badge` is what updates the README.**)

## License

[MIT](LICENSE)
