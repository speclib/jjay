# Spawn Stops Hardcoding

Spawn had three things nailed to the floor. I pried them loose.

The agent command was always `claude`. The tmux session was always the current one. The workspace root was always `../<project>-workspaces`. Fine for a demo, useless for real work — you couldn't test the lifecycle, couldn't run codex or a custom script, couldn't put workspaces anywhere else. Three walls, all coming down.

Now spawn takes `--agent`, `--session`, and `--workspace-root`, and cleanup learned the matching `--session` and `--workspace-root` so the two stay in step. The shared workspace package accepts a configurable root instead of assuming one.

Best part: a fake agent script and real integration tests, behind a build tag. Spawn, verify, cleanup, verify — the full lifecycle, in its own throwaway tmux session and temp jj repo, no claude required and nobody's real session polluted. Configurable *and* tested. That's how you build something that lasts.
