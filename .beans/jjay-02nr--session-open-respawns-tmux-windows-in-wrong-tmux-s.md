---
# jjay-02nr
title: session-open respawns tmux windows in wrong tmux sessions
status: completed
type: bug
priority: normal
created_at: 2026-06-08T11:56:21Z
updated_at: 2026-06-12T00:00:00Z
parent: jjay-hjjg
openspec-link: openspec/changes/archive/2026-06-12-resume-spawns-on-reopen
---

Fix (bundled into resume-spawns-on-reopen): session-open enumerated workspaces from cwd (the caller repo), so opening a session for repo B from inside repo A reopened A's spawns into B's session. Now `status.ListIn(repoRoot, …)` queries `jj -R <target>` and `session.Open` threads the target path through, scoping reopen to the target repo.

I started two jjay sessions and it totally mixed the respawned windows
