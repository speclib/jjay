## 1. coverage target

- [x] 1.1 `coverage` now runs `go test -tags integration -coverpkg=./... -coverprofile=coverage.out ./...` (plain `go test`), then HTML report + total-percentage print.
- [x] 1.2 `coverage: clean-tests` — sweeps integration-spawned `jjay-test-*` sessions/temp dirs first.

## 2. coverage-unit fallback

- [x] 2.1 Added `coverage-unit`: `go test -coverpkg=./... -coverprofile=coverage.out ./...` (no `-tags integration`), same HTML + percentage. No `clean-tests` dep.
- [x] 2.2 `coverage-unit` added to `.PHONY`.

## 3. badge clarity

- [x] 3.1 `badge: coverage` unchanged (the README-updating target). Verified `coverage`/`coverage-unit` do NOT patch the README badge line — only `badge` does.

## 4. Docs

- [x] 4.1 README "Testing & coverage" section under Contributing: `coverage` needs tmux+jj (runs integration); `coverage-unit` is the fallback; `badge` updates the README.
- [x] 4.2 Ran `make badge` — README badge refreshed to the true whole-repo number (53.8% phantom → 75.7%).

## 5. Verify

- [x] 5.1 `make coverage`: real numbers (merge ~73%, smokeTest 89.5%, cleanup ~78%, total 75.7%), not ~5%; HTML + percentage emitted; pre-run `clean-tests` ran. (A session spawned *during* the run survives — the known teardown gap `clean-tests` brooms, not a coverage bug.)
- [x] 5.2 `make coverage-unit`: completed without integration (50.0%, no tmux/jj required).
- [x] 5.3 `coverage`/`coverage-unit` leave the README badge untouched; `make badge` patches it.

## 6. Beans

- [x] 6.1 `jjay-gwpc` → in-progress, linked; → completed on archive.
- [x] 6.2 Noted on `jjay-bi7i`: it should invoke this `coverage` (integration+coverpkg) and commit the badge.
