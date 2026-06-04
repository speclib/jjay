## Why

jjay orchestrates parallel agent sessions in a project, but every project must first be *prepared*: openspec initialized with the right schema and a `config.yaml`, an `AGENTS.md` describing the archive/beans/jj conventions, the `/jjay:*` commands + skill installed into `.claude/`, and (optionally) jj. Today that is a manual checklist — exactly the kind of repetitive bootstrap jjay should automate.

Bean [jjay-ofk7](../../.beans/jjay-ofk7--bootstrap-task.md) ("init task", parent [jjay-5y1a](../../.beans/jjay-5y1a--backlog.md) backlog) defines this. It also absorbs the install handoff from `add-claude-commands`: the canonical `/jjay:*` command + skill content authored there is **installed into a target project's `.claude/` by `jjay init`** (decided: per-target, not user-scope — ADR-007).

## What Changes

- **New `jjay init [path]` command** that prepares a target project (default: cwd) to be orchestrated by jjay. It is **idempotent** and **non-destructive**: existing files are not overwritten without confirmation. Steps:
  - **openspec** — run `openspec init` for the target, selecting the `claude` tool and the project's schema; ensure `openspec/config.yaml` exists (seed from template, prompt for project context).
  - **jjay Claude integration** — install the `/jjay:*` command templates into `<target>/.claude/commands/jjay/` and the `jjay` skill into `<target>/.claude/skills/jjay/` (content from `add-claude-commands`).
  - **AGENTS.md** — write/extend an `AGENTS.md` documenting the jjay conventions (openspec archive flow, beans tasks, jj usage).
  - **jj (optional)** — initialize a jj repo if requested and not already present.
  - **hooks (optional)** — scaffold example hooks (e.g. a commented beans hook) the user can enable.
- **Flags** for non-interactive use: `--yes` (accept defaults), `--with-jj`, `--no-claude`, and per-step skips, mirroring `openspec init`'s non-interactive style (`--tools`, `--force`).

## Capabilities

### New Capabilities
- `init`: the `jjay init` command — idempotent project bootstrap (openspec, Claude integration, AGENTS.md, optional jj + hooks).

### Modified Capabilities
<!-- None: init creates project scaffolding; it does not change other commands' behavior. -->

## Impact

- **Code**: new `internal/init/` package and `init` command in `cmd/jjay/main.go`. Shells out to `openspec init`, `jj git init`/`jj init`, and copies embedded command/skill templates.
- **Tools**: `openspec` (init), `jj` (optional init). No tmux (init does not spawn).
- **Templates**: the `/jjay:*` command + skill files from `add-claude-commands` become embedded assets `jjay init` writes into the target. **Depends on `add-claude-commands` landing first** (it is the source of that content).
- **Specs**: new `init` capability.
- **ADRs**: ADR-008 (init is idempotent, non-destructive, orchestrates external initializers).
- **Beans**: ofk7 → in-progress, linked here.
