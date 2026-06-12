## Context

`internal/spawn/spawn.go` carries one command template per spawn flow (`DefaultAgentCommand` for apply, `proposalExplore/ProposeCommand` for proposals) and deliberately routes **both** first-launch and reopen through `openWindow`/`setupPanes`. `OpenWindow`'s doc comment states the invariant as a feature: *"Both Spawn and session-open route through openWindow so they cannot diverge."*

`internal/session/session.go` `reopenDetached` calls `spawn.OpenWindow(s.Name, s.WSDir, opts)` for each detached workspace, which resolves the apply template and `send-keys` `claude "/opsx:apply <change>" …` into the window. So reopening a workspace whose tmux died starts a fresh Claude session and re-runs `/opsx:apply` from scratch — discarding the conversation that was mid-flight.

Configuration today is CLI-flag-only (ADR-006): `--agent` overrides a hardcoded Go const; there is no jjay config file. The only `config.yaml` jjay touches is `openspec/config.yaml`, which belongs to openspec (`schema: spec-driven`). ADR-006 explicitly deferred a config file "until usage patterns emerge" and called repetition "a future config file problem". Per-agent **resume** knowledge is that pattern emerging: a flag cannot carry two commands (launch + resume) per agent.

## Goals / Non-Goals

**Goals:**
- Reopen resumes the agent (per the user's configured `resume` command), never re-runs `/opsx:apply`.
- Per-agent `launch` + `resume`, user-configurable, so jjay does not hardcode what "resume" means.
- One reopen primitive (`tmux-open`); `session-open` loops it — no second code path that can drift.
- A config foundation that jjay-iex3/jjay-euup extend rather than replace.

**Non-Goals:**
- Deciding resume *semantics* — `--resume` vs `--continue` vs bare shell is the user's `resume` string.
- Auto-detecting "is there a session to resume" — the configured command runs as-is.
- Renaming existing verbs (`session-open`, `spawn apply`, …) — deferred.
- Supporting non-`claude` agents beyond leaving the door open.

## Decisions

- **Agent = profile `{Launch, Resume}`, not a string (ADR-014).** `SpawnOptions.Agent string` becomes a resolved profile. Both fields are command templates over `{change}`/`{prompt}`/`{wsdir}`.
- **Intent on `openWindow`.** Add an `Intent` (launch | resume). `Spawn`/`SpawnProposal` pass launch; `tmux-open`/`session-open` pass resume. `setupPanes` picks `profile.Launch` or `profile.Resume` accordingly. The "cannot diverge" comment is inverted to "must diverge on the command, share the window/pane setup".
- **3-layer, per-field config resolution.** For each agent field:
  `resolve(agent, field) = project[agent][field] ?? global[agent][field] ?? builtin[agent][field]`
  with empty string treated as "unset" (so a present-but-blank `resume` still falls back). Resolution is **per field**, not per profile: a project overriding only `launch` keeps the global/built-in `resume`.
  - global = `~/.config/jjay/config.yaml`
  - project = `<repo>/.jjay/config.yaml` (located from the main repo root, not the spawned workspace)
  - builtin = a Go const map; the current `DefaultAgentCommand` becomes `builtin["claude"].Launch`, with `Resume` = `claude --resume --add-dir {wsdir}`.
  The resolver SHALL accept an **ordered list of sources** so jjay-iex3 can later insert/extend layers without rewriting it (`resolveFields([project, global, builtin])`).
- **Built-in is the single source of truth.** `jjay init` serializes the built-in into `<repo>/.jjay/config.yaml` (idempotent, non-destructive — ADR-008). Seeded file and runtime fallback come from the same const, so they cannot drift.
- **`tmux-open` is the reopen primitive.** `jjay tmux-open <workspace>` reopens one workspace's window (resume intent) in the current session. `reopenDetached` becomes `for s in detached { tmuxOpen(s, session) }`. The naming is deliberately informative (says *what* opens); a holistic verb rename is deferred.
- **`--resume` default scoping rides on cwd.** Windows already open with `tmux new-window -c wsDir` and the agent is launched from there, so `claude --resume`'s picker is scoped to that workspace's directory. No extra wiring — but it must be verified end-to-end (ADR-010's lesson: validate against the real lifecycle, not by reasoning).

## Risks / Trade-offs

- **`--resume` is interactive.** The default resume hands a picker to the human ("leave it from there"). `session-open` with N detached spawns opens N windows each awaiting a human pick. This is the intended behavior, but it is a visible departure from the hands-off `--dangerously-skip-permissions` launch flow. A user who wants hands-off resume sets `resume: 'claude --continue …'` in config.
- **`--resume` scoping depends on cwd.** If Claude keys sessions differently than expected, the picker could surface unrelated sessions. Mitigated by the `-c wsDir` window cwd; requires a manual-verify task.
- **New precedence surface.** Two config file paths + a built-in is more to reason about than one flag. Mitigated by per-field fallback, an ordered-source resolver, and documenting that this is jjay's file, distinct from `openspec/config.yaml`.
- **Supersedes ADR-006.** Reversing a prior "no config file" decision must be explicit (ADR-014 records it). Flags remain valid as the highest-priority override path so integration tests are unaffected.
- **Reopen-as-resume changes a long-standing behavior.** Anyone relying on reopen to re-drive `/opsx:apply` loses that; it was never the intent (the task calls it a bug).

## Open Questions (resolve during apply)

- Exact YAML library: prefer whatever `openspec` deps already vend (likely `gopkg.in/yaml.v3` via `go.sum`) before adding a new dependency.
- Whether `tmux-open` should error vs warn-and-skip when the workspace has no jj workspace / no window slot — align with `session-open`'s existing best-effort, non-fatal reopen model (ADR-003/006-session).
