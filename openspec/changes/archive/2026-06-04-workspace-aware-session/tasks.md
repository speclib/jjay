## 1. Shared workspace⋈window join

- [x] 1.1 Create `internal/status` package with a function that returns spawned workspaces with `{change, wsDir, attached}`, deriving from `jj workspace list` (reuse the parsing pattern in `spawn.go:checkWorkspaceNotExists`) intersected with `tmux list-windows` for the current session.
- [x] 1.2 Map change → window name via `workspace.WindowName`; exclude the `default` (main) workspace.
- [x] 1.3 Tolerate a missing tmux server (no error → treat all as detached), mirroring `session.go:checkSessionNotExists`.
- [x] 1.4 Unit tests for the join: attached, detached, no-tmux, default excluded.

## 2. `jjay status` command

- [x] 2.1 Add `statusCmd` (no args) to `cmd/jjay/main.go`, wired to the join helper.
- [x] 2.2 Render a table: change name, workspace dir, attached/detached status.
- [x] 2.3 Handle the empty case ("no running spawns", exit zero).
- [x] 2.4 Tests covering the spec scenarios (listed, detached, other-session, no-tmux, unexpected arg).

## 3. Reopen detached spawns on session-open

- [x] 3.1 Factor spawn's window-creation + `setupPanes` (agent launch) into a reusable helper so spawn and reopen share it.
- [x] 3.2 After `switchClient` in `session.Open`, enumerate detached spawns (via the join helper, scoped to the new session) and recreate a `ws-<change>` window + relaunch agent for each.
- [x] 3.3 Make reopen best-effort: a per-spawn failure logs and continues; session-open still succeeds and reports which spawns failed.
- [x] 3.4 Skip spawns whose window already exists (no duplicates).
- [x] 3.5 Tests covering the spec scenarios (reopen all, none, no-duplicate, one-fails-non-fatal).

## 4. ADR & docs

- [x] 4.1 Confirm ADR-006 (workspace = single source of truth) reflects the implemented behavior; flip its status to Accepted on archive.
- [x] 4.2 Update README with `jjay status` and the session-open reopen behavior.
- [x] 4.3 Update CHANGELOG.

## 5. Beans

- [x] 5.1 Set `jjay-72pw` status to `scrapped` (superseded by this change; principle captured in ADR-006).
- [x] 5.2 Set `jjay-7spt` and `jjay-cedd` status to `in-progress`; on archive, add `openspec-link` to each.
