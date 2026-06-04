---
name: "JJAY: Spawn"
description: Spawn an isolated jj workspace + tmux window and launch an agent on an OpenSpec change
category: Workflow
tags: [jjay, workspace, spawn]
---

Spawn an isolated agent workspace for an OpenSpec change by running the `jjay` binary.

This is a **thin wrapper** over `jjay spawn`. Do not reimplement the spawn logic in this prompt — the binary owns that behavior.

**Input**: Optionally specify a change name (e.g., `/jjay:spawn add-auth`).

**Steps**

1. **Resolve the change name**

   - If a name is provided as an argument, use it.
   - If omitted, run `openspec list --json` to get available changes and use the **AskUserQuestion tool** to let the user select. Do NOT guess.

2. **Run the binary**

   ```bash
   jjay spawn <change>
   ```

3. **Relay the output**

   Show the binary's output verbatim (workspace path, tmux window, agent launch). Do not summarize away errors — if `jjay` is not on `PATH` or the spawn fails, surface that directly.

**Guardrails**
- Thin wrapper only: invoke `jjay spawn <change>` and relay output.
- Requires the `jjay` binary on `PATH`.
- If you are running **inside a spawned worker workspace**, do NOT spawn again — implement the change in place with `/opsx:apply`. See the `jjay` skill.
