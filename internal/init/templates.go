package init

// configTemplate is a minimal openspec config.yaml jjay seeds when `openspec
// init` ran non-interactively (which skips writing one). It picks the
// spec-driven schema by default and leaves a placeholder context block for the
// project to fill in. It is written non-destructively — never over an existing
// config without --force.
const configTemplate = `schema: spec-driven

context: |
  # Describe this project for the agents: tech stack, build, platform,
  # domain language. (Seeded by 'jjay init' — edit to fit your project.)
`

// agentsTemplate is the AGENTS.md jjay writes into a target project. It
// documents the conventions a jjay-orchestrated project follows: the openspec
// archive flow, beans tasks, and jj usage.
const agentsTemplate = `# AGENTS.md

This project is orchestrated with [jjay](https://github.com/mipmip/jjay):
parallel AI agent sessions on top of **jj**, **tmux**, and **openspec**.

## jj (Jujutsu)

- This project uses **jj**, not plain git.
- Do not let the agent appear as committer or co-author in the history.
- Implementation work happens in isolated **jj workspaces** (one per change),
  spawned by jjay — not in the main working copy.

## OpenSpec changes

- Work is tracked as **openspec changes** (` + "`openspec/changes/<name>/`" + `),
  each with a proposal, specs, design, tasks, and (on archive) a blog entry.
- Implement a change by **spawning an isolated workspace** (` + "`/jjay:spawn <change>`" + `),
  not by running ` + "`/opsx:apply`" + ` in the main session. A worker agent already
  inside a spawned workspace applies in place and must not recursively spawn.

### Archive flow

When archiving a change (` + "`/opsx:archive`" + `), prepare it for merge into ` + "`main`" + `:

1. **Blog gate** — ensure the ` + "`blog`" + ` artifact is ` + "`done`" + ` (auto-created if
   missing), written retrospectively from the proposal and completed tasks.
2. ` + "`jj describe`" + ` the change.
3. Update the CHANGELOG.

## Beans tasks

- High-level tasks live under ` + "`.beans/`" + ` and act as epics for openspec proposals.
- When a bean is used to create a proposal, set its status to ` + "`in-progress`" + `.
- When a proposal is archived, add an ` + "`openspec-link`" + ` to the bean's frontmatter.

## Lifecycle

` + "```" + `
explore → propose → spawn → status → merge → cleanup
` + "```" + `

See ` + "`.claude/skills/jjay/SKILL.md`" + ` for the full orchestrator policy.
`

// hooksExample is a commented, opt-in example hooks file. It does nothing until
// the user uncomments and wires it up; it just shows the shape.
const hooksExample = `#!/usr/bin/env bash
# jjay example hooks — opt-in. This file does nothing as shipped.
#
# Copy or source the parts you want into your own hook setup. These are
# illustrative scaffolds, all commented out.

# Example: refresh beans before spawning a worker.
# pre_spawn() {
#   beans prime
# }

# Example: write a devblog entry after a merge.
# post_merge() {
#   echo "merged $1" >> devblog/merges.log
# }
`
