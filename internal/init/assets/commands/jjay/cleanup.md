---
name: "JJAY: Cleanup"
description: Tear down a spawned workspace, its tmux window, and its directory
category: Workflow
tags: [jjay, workspace, cleanup]
---

Tear down a spawned workspace (jj workspace + tmux window + directory) by running the `jjay` binary.

This is a **thin wrapper** over `jjay cleanup`. Do not reimplement the cleanup logic in this prompt — the binary owns that behavior.

**Input**: Optionally specify a change name (e.g., `/jjay:cleanup add-auth`).

**Steps**

1. **Resolve the change name**

   - If a name is provided as an argument, use it.
   - If omitted, run `jjay status` (or `openspec list --json`) to list candidates and use the **AskUserQuestion tool** to let the user select. Do NOT guess.

2. **Run the binary**

   ```bash
   jjay cleanup <change>
   ```

3. **Relay the output**

   Show the binary's output verbatim. Surface any error directly.

**Guardrails**
- Thin wrapper only: invoke `jjay cleanup <change>` and relay output.
- Cleanup is destructive (removes the workspace directory) — confirm the change name with the user before running if there is any ambiguity.
- Requires the `jjay` binary on `PATH`.
