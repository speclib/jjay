## 1. Slug helper

- [x] 1.1 Add a deterministic slug function (prompt → short kebab handle): lowercase, strip punctuation, drop stopwords, keep N salient tokens, cap length.
- [x] 1.2 Uniqueness: if `prop-<slug>` collides with an existing jj workspace or tmux window, append `-2`/`-3`/…
- [x] 1.3 Unit tests: typical prompt, punctuation/stopwords, length cap, collision suffix.

## 2. Spawn verbs

- [x] 2.1 Add cobra subcommands `spawn apply <change>` and `spawn proposal <prompt>` in `cmd/jjay/main.go`; `spawn` with no verb prints usage and exits non-zero (no bare-arg form).
- [x] 2.2 `--mode explore|propose` on `proposal` (default from config); maps to `/opsx:explore`/`/opsx:propose` agent templates with a `{prompt}` placeholder.
- [x] 2.3 Apply flow: existing `Spawn`, name prefixed `app-<change>`.
- [x] 2.4 Proposal flow: skip `checkOpenspecChange`; derive slug; name `prop-<slug>`; isolate; launch the mode's seed command.
- [x] 2.5 Share the isolate → window → panes tail between both flows; branch only on validate-vs-slug + agent template.

## 3. Decouple workspace-name from change-name

- [x] 3.1 Audit `merge`, `status`, `cleanup` for the `workspace == change` assumption; generalize where a proposal spawn would break it.
- [x] 3.2 Where a command needs the produced change name, read it from the workspace, not infer from the workspace name.

## 4. status two-table view

- [x] 4.1 Classify each spawn by `app-`/`prop-` prefix in `internal/status`.
- [x] 4.2 Render two tables: CHANGES (with MERGED/ARCHIVED/TASKS) and PROPOSAL SPAWNS (no change-shaped columns).

## 4b. completion for the verbs

- [x] 4b.1 `spawn <TAB>` completes the verbs (`apply`, `proposal`) — cobra does this for subcommands automatically; verify.
- [x] 4b.2 Move the existing `completion.Spawnable` `ValidArgsFunction` from `spawnCmd` to the new `apply` subcommand; `proposal` gets none (free-text prompt).
- [x] 4b.3 Verify `jjay __complete spawn apply ""` returns filtered changes and `__complete spawn proposal ""` returns nothing.

## 5. Caller migration & naming

- [x] 5.1 Update callers from bare `spawn <change>` to `spawn apply <change>`: the `/jjay:spawn` command file (`.claude/commands/jjay/spawn.md`) and the spawn integration tests.
- [x] 5.2 Note the naming change: spawns are now `app-<change>` (was `<change>`). Confirm cleanup/status/merge lookups use the prefixed name consistently; decide whether any in-flight unprefixed workspaces need handling (likely none — document).

## 6. Tests & docs

- [x] 6.1 Integration: a proposal spawn whose produced change name differs from its slug; assert no mis-keying in status/merge.
- [x] 6.2 README/CHANGELOG: document `spawn apply`/`spawn proposal`, `--mode`, prefixes, two-table status (changelog headline ≤ 80 chars).
- [x] 6.3 Confirm ADR-011 reflects the implementation; flip to Accepted on archive.

## 7. Beans

- [x] 7.1 Set `jjay-4ulx` to `in-progress`; add `openspec-link` on archive. (Lifecycle enum tracked in `jjay-mg00`.)
