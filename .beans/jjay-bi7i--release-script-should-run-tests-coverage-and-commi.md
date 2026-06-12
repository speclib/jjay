---
# jjay-bi7i
title: release script should run tests + coverage and commit the badge update last
status: todo
type: task
priority: normal
created_at: 2026-06-05T00:40:03Z
updated_at: 2026-06-12T14:04:09Z
parent: jjay-hjjg
blocked_by:
    - jjay-gwpc
---

`scripts/release.sh` does NOT run tests, coverage, or the README badge update — it relies on the manual pre-release checklist in RELEASING.md. Automate it so a release can't ship with failing tests or a stale badge.

## Wanted
1. Early in release.sh (before bumping/tagging), run the gate and ABORT on failure:
   - `make test` (unit) and `make test-integration` — or whatever `jjay-gwpc` settles as the combined coverage run
   - `make build`, `make lint`
2. Run `make badge` so the README coverage badge is refreshed.
3. Fold the refreshed README badge into the FINAL release commit — i.e. run `make badge` BEFORE the `jj describe -m "release: v${NEW_VERSION}"` step so the badge change rides in that same commit (currently that describe bundles VERSION + CHANGELOG + vendorHash). The badge update must be the last thing committed, not a separate dangling change.

## Where in the script
release.sh currently: safety checks → prompt version → update VERSION → update CHANGELOG → update vendorHash → `jj describe`/bookmark/tag/push.
- Insert the test/build/lint gate near the top (after safety checks, before version prompt) so we fail fast.
- Insert `make badge` just before the `jj describe` so the README change is part of the release commit.

## Caveats / dependencies
- Integration tests need tmux + jj on PATH (see jjay-gwpc). If the combined coverage isn't available everywhere, the script should at least run what it can and be explicit about what it skipped.
- Depends on / coordinates with jjay-gwpc (make coverage including integration) — decide the coverage command there, then call it here.

Found while prepping a release: tests/coverage/badge are all manual today; easy to forget and ship a stale badge or untested code.



Note (2026-06-12): coverage-includes-integration (jjay-gwpc) shipped `make coverage` = integration + whole-repo `-coverpkg`, and `make badge` patches the README. When implementing this bean, the release gate should call `make coverage` (not the old unit-only) and `make badge`, folding the refreshed badge into the release commit. `make coverage-unit` exists for tmux/jj-free CI.
