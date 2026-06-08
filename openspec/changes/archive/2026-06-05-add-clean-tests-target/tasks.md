## 1. clean-tests target

- [x] 1.1 Add a `clean-tests` target to the Makefile:
  - kill tmux sessions matching `jjay-test-*` (e.g. `tmux list-sessions -F '#{session_name}' | grep '^jjay-test-' | xargs -r -n1 tmux kill-session -t`); tolerate no tmux server (no error).
  - `rm -rf /tmp/jjay-test-*` and `/tmp/jjay-merge-test-*`.
- [x] 1.2 Add `clean-tests` to `.PHONY`.

## 2. Auto-sweep before integration

- [x] 2.1 Make `test-integration` depend on `clean-tests` (`test-integration: clean-tests`) so a prior aborted run's debris is cleared first.

## 3. Verify

- [x] 3.1 With stray `jjay-test-*` sessions present, `make clean-tests` removes them and leaves `jjay->*` real sessions untouched. (Verified: killed `jjay-test-99999`, spared `jjay->faketest`.)
- [x] 3.2 `make test-integration` runs the sweep first, then the suite; passes.
- [x] 3.3 `make clean-tests` is a no-op (no error) when there are no stray sessions and no tmux server. (Verified: exit 0.)

## 4. Beans

- [x] 4.1 Set `jjay-zgqx` to `completed`; add `openspec-link` on archive.
