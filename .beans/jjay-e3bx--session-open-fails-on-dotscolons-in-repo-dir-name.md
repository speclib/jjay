---
# jjay-e3bx
title: session-open fails on dots/colons in repo dir name (tmux target parsing)
status: completed
type: bug
priority: normal
created_at: 2026-06-12T13:02:57Z
updated_at: 2026-06-12T13:02:57Z
parent: jjay-hjjg
openspec-link: openspec/changes/archive/2026-06-12-resume-spawns-on-reopen
---

`./jjay session-open ~/cLinden/mip.rs/` failed:
  Error: failed to switch to tmux session "jjay->mip.rs": exit status 1

Root cause: tmux treats `.` and `:` as the `session:window.pane` separators in TARGET names. `tmux new-session -s "jjay->mip.rs"` silently stores the session as `jjay->mip_rs` (tmux's own normalization), but `switch-client -t "jjay->mip.rs"` parses `.rs` as a pane → "can't find pane: rs" → exit 1. The create-name and target-name disagreed.

Fix: SessionName sanitizes `.` and `:` → `_` (sanitizeTmuxName) in one place, so creation and every -t target use the same tmux-safe name. Regression test in session_test.go (mip.rs → jjay->mip_rs, trailing slash, colon case).

Bundled into resume-spawns-on-reopen (same session-open/reopen surface).
