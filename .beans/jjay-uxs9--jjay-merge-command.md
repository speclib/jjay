---
# jjay-uxs9
title: jjay merge command
status: completed
openspec-link: openspec/changes/archive/2026-06-03-merge-command
type: task
priority: normal
created_at: 2026-06-03T11:48:04Z
updated_at: 2026-06-03T12:13:01Z
parent: jjay-qmpg
---

Implement jjay merge <change> to merge a workspace change into main. Steps: jj new main <workspace-change> -m 'merge <change> into main', jj bookmark set main -r @, jj new. Should verify the workspace exists and the change has work to merge.
