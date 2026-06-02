---
# jjay-w88o
title: 'CRITICAL: fix workspace staleness — spawn must isolate working copies'
status: completed
openspec-link: openspec/changes/archive/2026-06-02-fix-workspace-isolation
type: bug
priority: critical
created_at: 2026-06-02T19:47:13Z
updated_at: 2026-06-02T20:53:11Z
---

jjay spawn uses --revision @ which causes both workspaces to edit the same jj change. When the spawned agent makes changes, the main workspace becomes stale. Running jj workspace update-stale can lose uncommitted work.

Fix: after creating the workspace with --revision @, immediately run jj new in the new workspace so it gets its own change (parent = @). This way the main workspace's @ is untouched and the new workspace is isolated.

Root cause: two workspaces editing the same change concurrently.
Impact: user can lose uncommitted work in the main workspace.
Discovered: 2026-06-02 during manual testing of jjay spawn.
