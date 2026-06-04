---
# jjay-4ulx
title: 'more spawns: spawn-apply spawn-explore spawn-proposal'
status: draft
type: task
priority: normal
created_at: 2026-06-04T20:55:13Z
updated_at: 2026-06-04T22:42:00Z
parent: jjay-qltp
---

spawn now spawns opsx:apply, but we also need spawn-explore and spawn-proposal. These new spawns spawn a new special tmux + agent window which might not even need a seperate workspace as conflict are almost imposable. This feature allowes the user continue with a new proposal when claude is busy in the main window. The claude session in the main window should be used for solving conflict and not creating proposals

we still might spawn these in workspaces as I now can think of race conditions if we for instance run a merge cmd while a claude session is creating files.
