## 1. init command scaffold

- [x] 1.1 Create `internal/init` package and an `Init(path string, opts InitOptions)` entry point.
- [x] 1.2 Add `initCmd` (`jjay init [path]`, path defaults to cwd) to `cmd/jjay/main.go`, with flags `--yes`, `--force`, `--with-jj`, `--no-claude`, and per-step skips.
- [x] 1.3 Define the step pipeline (detect → act), each step reporting created / skipped / would-overwrite.

## 2. openspec step

- [x] 2.1 Detect existing `openspec/`; if absent, shell out to `openspec init <path> --tools claude` (pass `--force` only under jjay `--force`).
- [x] 2.2 Ensure `openspec/config.yaml` exists; seed from the schema template and prompt for project context (skip prompts under `--yes`).
- [x] 2.3 Surface a clear error if the `openspec` binary is not available.

## 3. Claude integration step

- [x] 3.1 Embed the `/jjay:*` command files and `jjay/SKILL.md` from the repo's `.claude/` via `go:embed` (source = the canonical content from `add-claude-commands`).
- [x] 3.2 Write them into `<target>/.claude/commands/jjay/` and `<target>/.claude/skills/jjay/`, non-destructively (skip existing unless `--force`).
- [x] 3.3 Test asserting the embedded assets match the live `.claude/` content (drift guard).

## 4. AGENTS.md step

- [x] 4.1 Write/extend `<target>/AGENTS.md` with jjay conventions (openspec archive flow, beans tasks, jj usage); preserve an existing file unless `--force`.

## 5. Optional steps

- [x] 5.1 `--with-jj`: initialize a jj repo via jj's own command if not already present.
- [x] 5.2 hooks: scaffold example (commented) hooks the user can enable; opt-in.

## 6. Idempotency & non-destructiveness

- [x] 6.1 Tests: bare project → fully initialized; re-run → no-op; partially-initialized → completes missing steps only.
- [x] 6.2 Tests: `--yes` creates but does not clobber; `--force` overwrites; existing AGENTS.md/config.yaml preserved without `--force`.

## 7. Docs & beans

- [x] 7.1 README: document `jjay init`, its flags, and that it installs the jjay Claude integration per-target.
- [x] 7.2 Update CHANGELOG.
- [x] 7.3 Confirm ADR-008 reflects the implemented behavior; flip to Accepted on archive.
- [x] 7.4 Set `jjay-ofk7` status to `in-progress`; add `openspec-link` on archive.

## 8. Verify

- [x] 8.1 In a scratch directory, run `jjay init --with-jj` and confirm openspec, `.claude/commands/jjay/`, `.claude/skills/jjay/`, `AGENTS.md`, and a jj repo all exist.
- [x] 8.2 Re-run `jjay init` and confirm it is a no-op with a clear "already initialized" report.
- [x] 8.3 Confirm `/jjay:spawn` resolves in the freshly-initialized scratch project.
