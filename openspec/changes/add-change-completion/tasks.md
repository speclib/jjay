## 0. Prerequisite gate

- [x] 0.1 Confirm `workspace-aware-session` is merged to main and `internal/status` exists there. **Do not start before this** (ADR-009 build-order gate).

## 1. openspec name reader

- [x] 1.1 Create `internal/openspec` package with `ChangeNames() ([]string, error)` — lift `openspecChange`/`openspecList` and the `openspec list --json` parse out of `internal/spawn/spawn.go`.
- [x] 1.2 Reimplement `spawn.go`'s `checkOpenspecChange` as a membership test over `ChangeNames()` — behavior unchanged; spawn's existing tests still pass.
- [x] 1.3 Unit test `ChangeNames()` (parses list; tolerates failure).

## 2. workspace name reader

- [x] 2.1 Add `WorkspaceNames() ([]string, error)` to `internal/status` — parse `jj workspace list`, exclude `default`, return names. No tmux, no tasks.md (share the parse/exclusion with `List()` where clean).
- [x] 2.2 Unit test `WorkspaceNames()` (excludes default; tolerates jj failure).

## 3. completion package

- [x] 3.1 Create `internal/completion` with `Spawnable`, `Mergeable`, `Cleanable`, each `([]string, cobra.ShellCompDirective)` using `ShellCompDirectiveNoFileComp`.
- [x] 3.2 `Spawnable = ChangeNames() \ WorkspaceNames()` (set minus); `Mergeable = Cleanable = WorkspaceNames()` (named apart; shared helper allowed).
- [x] 3.3 Graceful degradation: reader error → empty candidates, never a shell-visible error.
- [x] 3.4 Import only `internal/openspec` and `internal/status` — never `spawn`/`merge`/`cleanup` (ADR-009 one-way deps).
- [x] 3.5 Unit tests: set-minus correctness, default excluded, empty/error cases.

## 4. wire into CLI

- [x] 4.1 In `cmd/jjay/main.go`, set `ValidArgsFunction` on `spawnCmd`/`mergeCmd`/`cleanupCmd` to the matching completion function (the only place that knows verb↔function).
- [x] 4.2 Verify `jjay __complete spawn ""` / `merge ""` / `cleanup ""` return the expected filtered names.

## 5. docs & beans

- [x] 5.1 README: note that `spawn`/`merge`/`cleanup` change names tab-complete (and that `jjay completion <shell>` installs the script).
- [x] 5.2 Update CHANGELOG.
- [x] 5.3 Confirm ADR-009 reflects the implementation; flip to Accepted on archive.
- [ ] 5.4 Set `jjay-nd10` status to `in-progress`; add `openspec-link` on archive.

## 6. verify

- [x] 6.1 With a spawned workspace present: `jjay spawn <TAB>` omits the spawned change; `jjay merge <TAB>` / `cleanup <TAB>` list the spawned workspaces.
- [x] 6.2 With no tmux server running, completion still works (proves no tmux dependency).
