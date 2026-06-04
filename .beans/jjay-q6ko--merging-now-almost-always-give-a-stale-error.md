---
# jjay-q6ko
title: merging now almost always give a stale error
status: todo
type: bug
priority: normal
created_at: 2026-06-04T22:23:50Z
updated_at: 2026-06-04T23:58:31Z
---

this is how I solve it from the workspace claude session:

me: i tried ./jjay merge from the main ws and this ws has this error: Error: The working copy is stale (not updated since operation f3fe0a226e2c).
claude: The staleness is resolved — but updating revealed a 2-sided conflict in .beans/jjay-ofk7--bootstrap-task.md. Main moved the bean (changed parent to jjay-qltp and bumped updated_at to 2026-06-04T22:08:46Z), while this workspace changed updated_at to 2026-06-05
  and left parent: jjay-5y1a.


the conflict are real but the staleness is something we can propably prevent
