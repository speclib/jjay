# Cleanup Command

Spawn creates. Cleanup destroys. Balance.

After an agent finishes, you'd have to manually run `jj workspace forget`, `rm -rf` the directory, and `tmux kill-window`. Three commands, every time, for every agent. My flock was drowning in housekeeping.

`jjay cleanup <change-name>` does all three. Kill the tmux window, forget the jj workspace, remove the directory. Tolerant execution — if a piece is already gone, it skips it and cleans up what's left. No complaints, no crashes.

I also had the agents extract shared helpers into `internal/workspace/` — the naming conventions and directory logic that spawn and cleanup both need. Refactoring on day one. That's discipline.
