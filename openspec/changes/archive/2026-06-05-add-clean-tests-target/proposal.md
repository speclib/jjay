## Why

Integration tests create tmux sessions (`jjay-test-<random>`) and temp dirs (`/tmp/jjay-test-*`, `/tmp/jjay-merge-test-*`). They clean up on normal teardown — but a panicked or interrupted (Ctrl-C) test skips its `defer`/`t.Cleanup`, leaking the session and dir. They accumulate: a single session left **6 orphaned `jjay-test-*` tmux sessions**. The existing `test-infrastructure` spec *promises* "no test tmux sessions remain" — that guarantee is real on clean exits but silently broken on hard interrupts, with no recovery path. Bean [jjay-zgqx](../../.beans/jjay-zgqx--make-clean.md) ("make clean") asks for cleanup; the pain is specifically the orphaned tmux sessions.

## What Changes

- **New `make clean-tests` target** — a broom that sweeps test debris by prefix:
  - kill every tmux session matching `jjay-test-*`
  - `rm -rf /tmp/jjay-test-*` and `/tmp/jjay-merge-test-*`
- **`make test-integration` auto-sweeps first** — depends on `clean-tests`, so a prior crashed run's debris is cleared before the suite runs (self-healing environment without touching test code).
- **Prefix safety:** `jjay-test-` can never match a real spawn session (those are `jjay-><dirname>` — `session.go`), so the sweep cannot kill a real session. Verified: 5th char `-` vs `>`.

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `test-infrastructure`: add the `clean-tests` Makefile target; note that `test-integration` sweeps stale debris first and that `clean-tests` is the recovery path when teardown is bypassed (panic/interrupt).

## Impact

- **Files**: `Makefile` only — new `clean-tests` target, `test-integration` gains a `clean-tests` prereq, `.PHONY` updated. No test-code or production-code changes.
- **Non-goal (broom-only):** this does NOT fix *why* tests leak (panic-safe teardown). That prevention is a possible follow-up; here we only sweep. The spec's "clean up regardless of pass/fail" stays aspirational for the normal path; `clean-tests` covers the abnormal path.
- **Accepted trade-off:** because `test-integration` sweeps all `jjay-test-*` first, running it kills any `jjay-test-*` session from a *concurrent* integration run. Fine for single-dev use — don't run two integration suites at once.
- **ADRs**: none (dev-tooling Makefile target; no architectural decision).
- **Beans**: zgqx → in-progress, linked here.
