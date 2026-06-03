# Tasks: session-open

## 1. Session package

- [x] 1.1 Create `internal/session/session.go` with `Open(path string) error` function
- [x] 1.2 Implement path resolution (filepath.Abs) and jj repo check (.jj/ exists)
- [x] 1.3 Implement tmux session existence check
- [x] 1.4 Implement tmux session creation (`new-session -d -s jjay-><dirname> -c <path>`)
- [x] 1.5 Implement tmux client switch (`switch-client -t jjay-><dirname>`)

## 2. CLI integration

- [x] 2.1 Add `session-open` cobra subcommand to `cmd/jjay/main.go` requiring exactly one argument
- [x] 2.2 Wire subcommand to call `session.Open()`

## 3. Tests

- [x] 3.1 Add unit tests for session name generation and path validation
- [x] 3.2 Verify `make test`, `make build`, `make lint` all pass

## 4. Documentation

- [x] 4.1 Update README CLI section with `session-open` command
