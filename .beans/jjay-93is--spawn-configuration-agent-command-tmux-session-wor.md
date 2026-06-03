---
# jjay-93is
title: spawn configuration (agent command, tmux session, workspace root)
status: completed
type: task
created_at: 2026-06-03T10:04:02Z
updated_at: 2026-06-03T12:45:00Z
openspec-link: openspec/changes/archive/2026-06-03-spawn-config
---

Make spawn configurable: agent command (default: claude), tmux target session (default: current), workspace root (default: ../<project>-workspaces/). Enables custom agents, integration testing with fake agent and dedicated tmux session, and flexible workspace locations. Unblocks integration test for spawn/cleanup lifecycle.
