---
# jjay-l80s
title: archive should be done after merge
status: draft
type: task
priority: normal
created_at: 2026-06-04T23:35:19Z
updated_at: 2026-06-05T10:10:56Z
parent: jjay-5y1a
blocked_by:
    - jjay-rse4
---

after merge we should run smoke tests to confirm the feature landen correctly in the main code. We should make this part of our main workflow. Maybe in cleanup.



Note (2026-06-05): the smoke-test half of this bean is now detailed in jjay-rse4 (post-merge smoke test). This bean retains the WORKFLOW question — when/where archive happens relative to merge (the body says 'maybe in cleanup'), and how archive-after-merge ties into the lifecycle. Blocked by rse4 (the smoke-test mechanism should exist first). Related: jjay-fd0z (ask to archive BEFORE merge).
