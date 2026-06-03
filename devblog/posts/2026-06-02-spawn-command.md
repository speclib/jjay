# Spawn Command

The whole point of jjay, in one command.

`jjay spawn <change-name>` replaces the tedious manual dance: create a jj workspace, open a tmux window, split panes, launch Claude. My agents were doing this by hand every single time. Unacceptable.

Now it's one command. Spawn creates the workspace in the project directory, opens a tmux window with a two-pane layout — Claude on the left, shell on the right — and the agent starts working immediately. Preconditions are validated up front. No tmux session? Error. No openspec change? Error. I don't let my agents stumble into broken setups.

This is the core. Everything else builds on spawn.
