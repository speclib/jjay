## Context

The CLI is cobra-based and already exposes `jjay completion {bash,zsh,fish,powershell}` plus the `__complete` dispatch — but `jjay spawn <TAB>` returns no candidates (`__complete spawn ""` → empty, ShellCompDirectiveDefault → file fallback). Nothing feeds cobra the live names. This change attaches `ValidArgsFunction`s that supply per-verb-filtered change names. Per ADR-009: completion composes lean data-source readers and never depends on command packages.

The two candidate sources already exist in spirit:
- openspec change names — parsed privately in `internal/spawn/spawn.go` (`openspecList`/`checkOpenspecChange`).
- spawned-workspace names — `internal/status.List()` (from `workspace-aware-session`) returns `[]Spawn{Change,...}`, excluding the `default` workspace; but `List()` also probes tmux and reads tasks.md.

## Goals / Non-Goals

**Goals:**
- Per-verb filtered completion: spawn = changes without a workspace; merge/cleanup = existing workspaces.
- A dedicated `internal/completion` package depending only on data sources, redesign-safe.
- Fast, side-effect-free TAB: name-only reads, no tmux/tasks I/O.
- Reuse, lightly refactored: lift openspec parsing out of `spawn.go` so both spawn and completion share one reader.

**Non-Goals:**
- Flag-value completion (`--session`, `--workspace-root`) — later, no bean (cobra does dir completion for `--workspace-root` natively anyway).
- Generating/installing completion scripts — `jjay completion <shell>` already exists; install is a user/init concern.
- Changing spawn/merge/cleanup runtime behavior — only completion candidates are added.
- Reusing the heavy `status.List()` for completion — rejected for latency (ADR-009).

## Decisions

- **`internal/completion` package** with exported `Spawnable`, `Mergeable`, `Cleanable`, each returning `([]string, cobra.ShellCompDirective)` and using `ShellCompDirectiveNoFileComp`. Named apart even though `Mergeable == Cleanable` today; shared body may delegate to an unexported helper.
- **`internal/openspec.ChangeNames()`** — lift `openspecChange`/`openspecList` and the `openspec list --json` parse out of `spawn.go` into a reader returning `[]string`. `spawn.go`'s `checkOpenspecChange` is reimplemented on top of it (membership test) — no behavior change.
- **`internal/status.WorkspaceNames()`** — a lean reader: parse `jj workspace list`, drop `default`, return names. Shares the `default`-exclusion rule with `List()` (factor the parse if convenient), but performs no tmux or tasks.md work. `List()` is unchanged.
- **`Spawnable = ChangeNames() \ WorkspaceNames()`** (set minus); `Mergeable = Cleanable = WorkspaceNames()`.
- **Binding in `cmd/jjay/main.go`:** `spawnCmd.ValidArgsFunction = wrap(completion.Spawnable)`, etc. Only these three lines know the verb↔function mapping — the redesign seam.
- **Graceful degradation:** a reader error yields an empty candidate list + `ShellCompDirectiveNoFileComp|...Error` as appropriate, never a shell-visible failure.

## Risks / Trade-offs

- **Build-order gate:** `internal/status.WorkspaceNames()` presupposes `internal/status` on main — so `workspace-aware-session` must merge first. If this change is implemented before that merge, the import won't resolve. Stated in the proposal and ADR-009; sequence accordingly.
- **`spawn.go` refactor risk:** lifting the openspec parse touches working code. Mitigate by keeping `checkOpenspecChange`'s behavior identical (now a membership check over `ChangeNames()`), covered by spawn's existing tests.
- **Package naming:** `internal/openspec` reads the *openspec binary* output — potential confusion with the external tool. Accepted (it is the openspec-domain reader); alternative names (`internal/changes`) noted but not chosen.
- **Two readers per spawn TAB:** `Spawnable` calls openspec + jj. Still fast (two subprocess reads, no tmux/file-walk); acceptable.
- **Completion staleness within a keystroke:** candidates are recomputed each TAB, so they track reality; no caching needed.
