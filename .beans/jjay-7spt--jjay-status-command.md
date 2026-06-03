---
# jjay-7spt
title: jjay status command
status: todo
type: task
created_at: 2026-06-03T11:48:11Z
updated_at: 2026-06-03T11:48:11Z
---

Implement jjay status to show running agent workspaces. Derive state from tmux (list windows matching ws-*) and jj (workspace list). Show: change name, workspace dir, tmux window status, jj workspace state. No state file needed — everything derived from tmux + jj.
