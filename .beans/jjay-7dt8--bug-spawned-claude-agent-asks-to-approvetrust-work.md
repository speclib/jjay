---
# jjay-7dt8
title: 'BUG: spawned claude agent asks to approve/trust workspace directory'
status: todo
type: bug
priority: high
created_at: 2026-06-03T19:24:05Z
updated_at: 2026-06-03T19:24:05Z
---

When jjay spawn launches claude in a workspace, claude still asks the user to trust the directory despite --add-dir being passed. The --add-dir flag grants tool access but does not bypass the workspace trust dialog. This requires manual intervention in every spawned session, defeating the purpose of autonomous agent operation.

Discovered: 2026-06-02 during spawn testing.
Current workaround: user manually approves in each spawned session.
Needs investigation: how to pre-trust a directory for Claude Code.
