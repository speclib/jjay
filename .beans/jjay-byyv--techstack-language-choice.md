---
# jjay-byyv
title: techstack / language choice
status: todo
type: task
priority: normal
created_at: 2026-06-02T14:30:00Z
updated_at: 2026-06-02T13:18:04Z
---

Decide on the implementation language for jjay CLI.

Candidates: Rust, Go, Zig, Python

Considerations:
- jjay is primarily an orchestrator (shells out to jj, tmux, claude/codex)
- CLI ecosystem maturity matters (arg parsing, TUI)
- This is also a learning project, I want to learn the workflow not the language
- My team likes go and python
- Binary size and cross-compilation are nice-to-haves
