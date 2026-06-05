## 1. Tooling

- [ ] 1.1 Add `gotestsum` to the `flake.nix` devShell `buildInputs` list (alongside `go`, `gopls`, `goreleaser`, `gum`)
- [ ] 1.2 Verify `gotestsum --version` resolves inside `nix develop`

## 2. Makefile target

- [ ] 2.1 Ensure `test-integration` is declared in `.PHONY`
- [ ] 2.2 Update the `test-integration` recipe to run `gotestsum --format standard-verbose -- -tags integration ./...` when `gotestsum` is on `PATH`
- [ ] 2.3 Add a fallback branch (`command -v gotestsum`) so the recipe runs plain `go test -tags integration ./...` when `gotestsum` is absent

## 3. Banner routing in the test helper

- [ ] 3.1 Change `runIn()` in `test/integration/helpers_test.go` to capture the subprocess's combined stdout/stderr (`cmd.CombinedOutput()`) instead of wiring `cmd.Stdout`/`cmd.Stderr` to `os.Stdout`/`os.Stderr`
- [ ] 3.2 Emit the captured output via `t.Logf` so it nests under the invoking test/subtest, and keep the existing `t.Fatalf` on non-zero exit (include the captured output in the failure message)

## 4. Verification

- [ ] 4.1 Run `make test-integration` inside the dev shell; confirm colored per-scenario PASS/FAIL lines and a `DONE N tests in Xs` summary
- [ ] 4.2 Confirm subprocess (jj/OpenSpec) detail appears nested under its scenario, not as free-floating lines
- [ ] 4.3 Run `make test-integration` with `gotestsum` off `PATH` (or `PATH` stripped) and confirm the fallback runs the suite and reports results
- [ ] 4.4 Confirm a deliberately failing scenario is clearly marked and its captured subprocess detail is visible
