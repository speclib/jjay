---
# jjay-euup
title: settings in project for spawned tmux layout and agent
status: todo
type: task
priority: normal
created_at: 2026-06-03T19:45:06Z
updated_at: 2026-06-12T13:06:20Z
parent: jjay-hjjg
blocked_by:
    - jjay-iex3
---

setup could be done at init
- stationary tmux layout (tmux without workspaces)
- apply workspace conf
  - tmux layout
  - agent to use
  - other panes configuration: e.g. nix develop

(body recovered 2026-06-05 from off-main commit 60f1a517 — it had been lost to an empty version on main)



Note (2026-06-12): the jjay config foundation now EXISTS (internal/config, 3-layer per-field resolution, .jjay/config.yaml seeded by init — ADR-014, resume-spawns-on-reopen). This bean now EXTENDS it (add fields: tmux layout, extra panes, agent-to-use) rather than building the config mechanism.
