---
name: "JJAY: Status"
description: List spawned jj workspaces, task progress, and tmux window state
category: Workflow
tags: [jjay, workspace, status]
---

List every spawned jj workspace with its task progress and tmux window state by running the `jjay` binary.

This is a **thin wrapper** over `jjay status`. Do not reimplement the status logic in this prompt — the binary derives everything live.

**Steps**

1. **Run the binary** (takes no positional arguments)

   ```bash
   jjay status
   ```

2. **Relay the output**

   Show the binary's table verbatim (CHANGE / WORKSPACE / TASKS / ARCHIVED / STATUS). Surface any error directly — e.g. if `jjay` is not on `PATH`.

**Guardrails**
- Thin wrapper only: invoke `jjay status` and relay output.
- No positional arguments.
- Requires the `jjay` binary on `PATH`.
