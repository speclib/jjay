---
name: jjay
description: Use when implementing or managing an OpenSpec change in this jjay repository — deciding how to apply a change, spawn/merge/cleanup agent workspaces, or check workspace status. Steers toward spawning an isolated agent workspace (`/jjay:spawn`) instead of running `/opsx:apply` in the main session.
license: MIT
metadata:
  author: jjay
  version: "1.0"
---

# jjay orchestrator skill

jjay manages parallel AI agent sessions with **jj**, **tmux**, and **openspec**. Its whole purpose is that an agent's work happens in an **isolated spawned workspace**, not in the main working copy. This skill encodes the policy that makes that the default.

## The rule

**Implement a change by spawning an isolated agent workspace, not by applying in the main session.**

When the conversation is about implementing an existing OpenSpec change in this repo, drive it with:

```
/jjay:spawn <change>
```

Do **NOT** run `/opsx:apply <change>` in the main session to implement a change. `jjay spawn` creates an isolated jj workspace + a `ws-<change>` tmux window and launches an agent that runs `/opsx:apply` inside that workspace. Applying in the main session is exactly the un-isolated, manual flow jjay was built to replace.

## Lifecycle

```
explore → propose → spawn → status → merge → cleanup
```

1. **explore** — think through the problem (`/opsx:explore`).
2. **propose** — create the change with its artifacts (`/opsx:propose`).
3. **spawn** — `/jjay:spawn <change>` creates an isolated workspace + tmux window and launches a worker agent on it.
4. **status** — `/jjay:status` lists spawned workspaces, task progress, and tmux window state (attached/detached). Read-only; derived live.
5. **merge** — `/jjay:merge <change>` rebases the workspace onto current `main` and merges it in (after the change is implemented, verified, and archived).
6. **cleanup** — `/jjay:cleanup <change>` tears down the workspace, its tmux window, and its directory.

`/jjay:session-open <path>` recreates the tmux view (windows + agents) for spawned workspaces — use it to reattach after a detach or reboot.

## Orchestrator vs worker

There are two kinds of session, and they behave differently:

- **Orchestrator session** — the main session where you run `/jjay:spawn`, `/jjay:status`, `/jjay:merge`, `/jjay:cleanup`. It dispatches and oversees work; it does not implement the change itself.
- **Worker session** — the agent running *inside* a spawned workspace. It implements the change by running `/opsx:apply` **in its own workspace** (this is correct here — the worker is already isolated).

**No recursive spawn.** Because `.claude/` is checked into the repo, the `/jjay:*` commands and this skill propagate into every spawned workspace. A worker therefore sees `/jjay:spawn` too — but a worker must **NOT** spawn again. Nesting workspaces inside workspaces is never what you want.

- If you are a **worker** (running inside a spawned workspace) and are asked to implement the change → run `/opsx:apply` in place. Do not call `/jjay:spawn`.
- If you are the **orchestrator** and are asked to implement a change → `/jjay:spawn <change>` and let the worker do the apply.

How to tell which you are: if your working directory is a spawned workspace (e.g. a sibling `*-workspaces/<change>` directory created by jjay, with a `ws-<change>` window), you are a worker — apply in place.

## Preconditions

- The `jjay` binary must be on `PATH` in the session (the commands shell out to it). This is true in spawned workspaces too, since they run in the same environment. If `jjay` is not found, surface that instead of falling back to a manual flow.
