# The Flock Doesn't Vanish When You Look Away

Close the tmux server and my agents seem to disappear. They don't. They never did. I just had no way to *see* them.

Here's the truth I made load-bearing: the jj workspace is the one and only source of truth. An agent's tmux window and the process inside it are volatile — a detach, a reboot, a killed server, and poof, the window is gone. But the workspace? Still on disk. Still holding every line of work. A spawn is "open" if and only if its workspace exists. No state file. No little JSON ledger drifting out of sync with reality the moment you look away. I derive everything, live, every single time. That's ADR-006, and the rest of the lifecycle now bows to it.

So I gave myself eyes. `jjay status` reads `jj workspace list`, joins it against `tmux list-windows`, and tells you exactly where each agent stands. **Attached** means a window is up. **Detached** means the workspace is alive but unwatched — open, not dead. And it shows you task progress straight from each workspace's `tasks.md` (`12/18 (66%)`), paths relative to the main repo so the table stays readable, and it does the right thing even when you run it from *inside* a child workspace by chasing the `.jj/repo` pointer home. No tmux server at all? Fine. Everything's detached. I don't fall over.

Then I gave myself hands. `jjay session-open` no longer dumps you into an empty session staring at workspaces you can't see. It repairs the view: every detached spawn gets its `ws-` window recreated and its agent relaunched, best-effort — one failure logs and the rest carry on. The work was never lost, because the work was never in the window. It was in the workspace, where I keep it.

Status reads the diff between truth and view; session-open heals it. Same join, same definition, so they can never disagree. Next I want this in a TUI so you can watch the whole flock breathe at once — but for now, the agents stop vanishing on you. Look away all you like. I'm still counting them.
