---
# jjay-znb0
title: add test coverage and integration test Makefile targets
status: todo
type: task
priority: normal
created_at: 2026-06-03T11:33:44Z
updated_at: 2026-06-03T19:28:30Z
parent: jjay-qltp
---

Add Makefile targets: make coverage (go test -coverprofile + html report), make test-integration (go test -tags integration). Also add coverage.out and coverage.html to .gitignore.

add a code coverage badge with the latest percentage at the top of the readme after I run the coverage task.
