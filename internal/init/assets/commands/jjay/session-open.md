---
name: "JJAY: Session Open"
description: Create and switch to a tmux session for a jj repo, reopening windows for spawned workspaces
category: Workflow
tags: [jjay, session, tmux]
---

Create and switch to a dedicated tmux session for a jj repo (reopening a window + agent for each spawned workspace) by running the `jjay` binary.

This is a **thin wrapper** over `jjay session-open`. Do not reimplement the session/reopen logic in this prompt — the binary owns that behavior. Reopened workspaces **resume** their agent (the configured `resume` command, e.g. `claude --resume`) — they do NOT re-run `/opsx:apply` from scratch.

**Input**: Optionally specify a repo path (e.g., `/jjay:session-open ../myproject`).

**Steps**

1. **Resolve the path**

   - If a path is provided as an argument, use it.
   - If omitted, use the **AskUserQuestion tool** to ask for the jj repo path (default to the current repo root if that is clearly intended). Do NOT guess a path silently.

2. **Run the binary**

   ```bash
   jjay session-open <path>
   ```

3. **Relay the output**

   Show the binary's output verbatim, including which spawns were reopened and any per-spawn reopen failures.

**Guardrails**
- Thin wrapper only: invoke `jjay session-open <path>` and relay output.
- Requires the `jjay` binary on `PATH`.
