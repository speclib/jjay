---
name: "JJAY: Merge"
description: Merge a spawned workspace into main
category: Workflow
tags: [jjay, workspace, merge]
---

Merge a spawned workspace into `main` by running the `jjay` binary.

This is a **thin wrapper** over `jjay merge`. Do not reimplement the merge logic in this prompt — the binary owns rebase-before-merge and conflict handling.

**Input**: Optionally specify a change name (e.g., `/jjay:merge add-auth`).

**Steps**

1. **Resolve the change name**

   - If a name is provided as an argument, use it.
   - If omitted, run `jjay status` (or `openspec list --json`) to list candidates and use the **AskUserQuestion tool** to let the user select. Do NOT guess.

2. **Run the binary**

   ```bash
   jjay merge <change>
   ```

3. **Relay the output**

   Show the binary's output verbatim. If a merge conflict aborts the merge, surface that directly rather than papering over it.

**Guardrails**
- Thin wrapper only: invoke `jjay merge <change>` and relay output.
- Requires the `jjay` binary on `PATH`.
