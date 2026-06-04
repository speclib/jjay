---
# jjay-oeay
title: manually test shell completion in the installed jjay (0.4)
status: todo
type: task
priority: normal
created_at: 2026-06-04T23:52:58Z
updated_at: 2026-06-04T23:52:58Z
parent: jjay-hjjg
---

Shell completion (jjay-nd10 / add-change-completion) shipped with unit tests, but the actual SHELL integration was never tested in an installed build. Unit tests cover the candidate logic; they cannot exercise the real cobra completion script in a live shell.

Manual test checklist (in an installed jjay, e.g. after `go install` / nix build, with completion sourced via `jjay completion <shell>`):
- `jjay spawn <TAB>` → offers the verbs (apply, proposal) once add-spawn-verbs lands; before that, the filtered change names.
- `jjay spawn apply <TAB>` → offers un-spawned change names only (not already-spawned ones).
- `jjay merge <TAB>` / `jjay cleanup <TAB>` → offer existing spawned workspaces only (default excluded).
- No file-path fallback on these arguments.
- Works in at least bash + zsh (fish if used).
- Degrades silently (no error spew) when run outside a repo or with jj/openspec missing.

Why a bean and not a test task: shell completion can only be verified by a human in a real shell with the script installed — out of scope for the Go test suite. Tracked under the 0.4 milestone (jjay-hjjg).

Related: add-change-completion (archived 2026-06-05), add-spawn-verbs (adds the verbs that change `spawn <TAB>` behavior).
