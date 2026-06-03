# Tasks: merge-command

## 1. Merge package

- [ ] 1.1 Create `internal/merge/merge.go` with `Merge(changeName string) error` function
- [ ] 1.2 Implement workspace existence check (parse `jj workspace list`)
- [ ] 1.3 Implement empty workspace warning (via `jj log -r "<change>@"` template)
- [ ] 1.4 Implement merge: `jj new main <change>@ -m "merge <change> into main"`
- [ ] 1.5 Implement bookmark update: `jj bookmark set main -r @`
- [ ] 1.6 Implement fresh change: `jj new`
- [ ] 1.7 Print success message with merge details

## 2. CLI integration

- [ ] 2.1 Add `merge` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [ ] 2.2 Wire merge subcommand to call `merge.Merge()`

## 3. Tests

- [ ] 3.1 Add unit tests for workspace existence check
- [ ] 3.2 Verify `make test`, `make build`, `make lint` all pass

## 4. Documentation

- [ ] 4.1 Update README: move merge from Planned to CLI section
