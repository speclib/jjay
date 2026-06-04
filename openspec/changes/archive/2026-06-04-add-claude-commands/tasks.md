## 1. jjay slash commands

- [x] 1.1 Create `.claude/commands/jjay/spawn.md` — thin wrapper running `jjay spawn <change>`; prompt for change name (list via `openspec list --json`) if missing.
- [x] 1.2 Create `.claude/commands/jjay/status.md` — runs `jjay status` (no args).
- [x] 1.3 Create `.claude/commands/jjay/merge.md` — runs `jjay merge <change>`; prompt if missing.
- [x] 1.4 Create `.claude/commands/jjay/cleanup.md` — runs `jjay cleanup <change>`; prompt if missing.
- [x] 1.5 Create `.claude/commands/jjay/session-open.md` — runs `jjay session-open <path>`; prompt if missing.
- [x] 1.6 Match the frontmatter/body shape of existing `.claude/commands/opsx/*.md`.

## 2. jjay orchestrator skill

- [x] 2.1 Create `.claude/skills/jjay/SKILL.md` with a `description` that auto-triggers on implementing/managing changes in this repo.
- [x] 2.2 Body: state the rule — implement changes via `/jjay:spawn`, not `/opsx:apply` in the main session.
- [x] 2.3 Body: document the lifecycle (explore → propose → spawn → status → merge → cleanup).
- [x] 2.4 Body: document orchestrator vs worker, with the no-recursive-spawn rule for workers.

## 3. Docs

- [x] 3.1 README section documenting the `/jjay:*` commands and the skill, including the `jjay`-on-PATH precondition.
- [x] 3.2 Update CHANGELOG.

## 4. ADR & beans

- [x] 4.1 Confirm ADR-007 reflects the implemented split; flip to Accepted on archive.
- [x] 4.2 Set `jjay-8xuj` status to `in-progress`; add `openspec-link` on archive.

## 5. Verify

- [x] 5.1 In a session with the `jjay` binary on PATH, confirm `/jjay:spawn <change>` creates a workspace + window and launches the agent.
- [x] 5.2 Confirm `/jjay:status`, `/jjay:merge`, `/jjay:cleanup`, `/jjay:session-open` invoke the right binary calls.
- [x] 5.3 Confirm the commands are present inside a spawned workspace's `.claude/commands/jjay/`.
