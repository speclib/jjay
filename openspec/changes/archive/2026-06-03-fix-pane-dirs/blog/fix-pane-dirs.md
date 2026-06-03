# The shell that wouldn't listen

I told my agents to work in the workspace directory. They did. The shell pane next to them? It went wherever it pleased.

The bug was a race condition. When jjay spawns a two-pane layout, the right pane gets a `cd` command via `tmux send-keys`. Simple enough — except the shell (fish, bash, whatever) hasn't finished initializing when that `cd` arrives. The command gets swallowed. The pane stays in the project root while the agent works in the workspace. Confusing for any human who switches over to check on things.

The fix is embarrassingly simple: tmux's `-c` flag sets the working directory at pane creation time. No race, no timing issues, no send-keys gymnastics. `tmux new-window -c /path` and `tmux split-window -c /path` — both panes start exactly where they should.

To make sure this never sneaks back, I added a proper integration test. It spins up a dedicated tmux session, creates a real jj workspace, spawns with a fake agent, then checks every pane's `#{pane_current_path}`. If either pane is in the wrong directory, the test fails. Then it runs cleanup and verifies everything is gone. No more trusting that send-keys will arrive on time.
