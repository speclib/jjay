---
name: "JJAY: Tmux Open"
description: Reopen one spawned workspace's tmux window and resume its agent
category: Workflow
tags: [jjay, tmux, resume]
---

Reopen a single existing spawned workspace's tmux window — resuming its agent (not re-running `/opsx:apply`) — by running the `jjay` binary.

This is a **thin wrapper** over `jjay tmux-open`. Do not reimplement the reopen/resume logic in this prompt — the binary owns that behavior (it runs the agent's configured `resume` command, e.g. `claude --resume`).

**Input**: Optionally specify a workspace name (e.g., `/jjay:tmux-open app-add-foo`).

**Steps**

1. **Resolve the workspace name**

   - If a name is provided as an argument, use it.
   - If omitted, run `jjay status` to list spawned workspaces and use the **AskUserQuestion tool** to let the user pick one that has no live window. Do NOT guess.

2. **Run the binary**

   ```bash
   jjay tmux-open <workspace>
   ```

3. **Relay the output**

   Show the binary's output verbatim. The window opens in the workspace directory and resumes the agent; it does NOT re-run `/opsx:apply` from scratch.

**Guardrails**
- Thin wrapper only: invoke `jjay tmux-open <workspace>` and relay output.
- Requires the `jjay` binary on `PATH`.
- Reopens an already-existing workspace; it does not create or modify the jj workspace.
