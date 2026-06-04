# claude-commands

## Purpose

The Claude Code integration layer over the `jjay` binary: the `/jjay:*` slash-command set and the `jjay` orchestrator skill. Commands are thin wrappers that invoke the binary; the skill carries the policy that makes Claude prefer spawning an isolated workspace over applying a change in the main session. Committed under `.claude/`, so they propagate into every spawned workspace.

## Requirements

### Requirement: jjay slash commands wrap the binary
The project SHALL provide Claude Code slash commands under `.claude/commands/jjay/` for the jjay verbs `spawn`, `status`, `merge`, `cleanup`, and `session-open`. Each command SHALL invoke the corresponding `jjay <verb>` binary call with the user-supplied arguments and SHALL NOT reimplement the verb's logic in the prompt.

#### Scenario: Spawn command runs the binary
- **WHEN** a user invokes `/jjay:spawn add-foo`
- **THEN** the command runs `jjay spawn add-foo`
- **THEN** the command relays the binary's output

#### Scenario: Verb with no argument
- **WHEN** a user invokes `/jjay:status`
- **THEN** the command runs `jjay status` with no positional arguments

#### Scenario: Missing required argument
- **WHEN** a user invokes `/jjay:merge` without a change name
- **THEN** the command prompts for the change name (e.g. via AskUserQuestion or by listing candidates) rather than guessing

### Requirement: Commands are checked into the repo
The `.claude/commands/jjay/` command files SHALL be committed to the repository so they propagate into spawned jj workspaces.

#### Scenario: Command available in a spawned workspace
- **WHEN** `jjay spawn add-foo` creates a workspace and an agent runs in it
- **THEN** the `/jjay:*` commands are present in that workspace's `.claude/commands/jjay/`

### Requirement: jjay orchestrator skill encodes the lifecycle policy
The project SHALL provide a `jjay` skill at `.claude/skills/jjay/SKILL.md` whose `description` causes it to auto-load when the conversation concerns implementing or managing changes in this repository. The skill SHALL state that changes are implemented by spawning an isolated agent workspace (`/jjay:spawn`), NOT by running `/opsx:apply` in the main session, and SHALL document the lifecycle: explore → propose → spawn → status → merge → cleanup.

#### Scenario: Skill steers toward spawn
- **WHEN** the conversation is about implementing an existing OpenSpec change in this repo and the skill is loaded
- **THEN** the guidance directs the agent to `/jjay:spawn <change>` rather than `/opsx:apply <change>` in the main session

### Requirement: Skill distinguishes orchestrator from worker
The skill SHALL describe the orchestrator session (where `/jjay:spawn` is run) versus the spawned worker session (where the agent runs `/opsx:apply` inside its workspace), and SHALL instruct that a worker session does not itself spawn.

#### Scenario: Worker does not recursively spawn
- **WHEN** an agent is running inside a spawned workspace (a worker session) and is asked to implement the change
- **THEN** the guidance has it apply within its workspace, not invoke `/jjay:spawn` again
