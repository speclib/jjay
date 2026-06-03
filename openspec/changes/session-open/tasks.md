# Tasks: session-open

## 1. Session package

- [ ] 1.1 Create `internal/session/session.go` with `Open(path string) error` function
- [ ] 1.2 Implement path resolution (filepath.Abs) and jj repo check (.jj/ exists)
- [ ] 1.3 Implement tmux session existence check
- [ ] 1.4 Implement tmux session creation (`new-session -d -s jjay-><dirname> -c <path>`)
- [ ] 1.5 Implement tmux client switch (`switch-client -t jjay-><dirname>`)

## 2. CLI integration

- [ ] 2.1 Add `session-open` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [ ] 2.2 Wire subcommand to call `session.Open()`

## 3. Tests

- [ ] 3.1 Add unit tests for session name generation and path validation
- [ ] 3.2 Verify `make test`, `make build`, `make lint` all pass

## 4. Documentation

- [ ] 4.1 Update README CLI section with `session-open` command
