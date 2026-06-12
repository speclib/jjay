# Reopen Means Resume, Not Restart

When a tmux session died and I reopened it, I used to march every agent back to the start line — re-running `/opsx:apply` from zero, throwing away the conversation that was already halfway through the work. Reopening should pick the agent back up where it left off. Now it does.

The fix names a thing I'd been pretending didn't exist: **launch and resume are different intents.** A first spawn *launches* (`claude "/opsx:apply …"`). A reopen *resumes* (`claude --resume`). They share the window and pane setup, but they must diverge on the command — the old "they can't diverge" comment was the bug wearing a halo. And since resume syntax is the agent's business, not mine, I gave jjay its own config file — `.jjay/config.yaml`, per-agent `launch` and `resume`, resolved project → global → built-in, one field at a time. (Yes, I once swore off a config file. The future arrived. ADR-006 is superseded; the flag still wins when you pass it.) There's a new single primitive too — `jjay tmux-open <workspace>` reopens one window, and `session-open` just loops it. One reopen path, no drift.

Two more session-open sins fell while I was in there. It used to reopen *the wrong project's* workspaces — open a session for `mip.rs` from inside jjay and it'd drag jjay's spawns along; now it asks the target repo (`jj -R`) what belongs to it. And it choked on a repo named `mip.rs`, because tmux reads the dot as "go to pane rs" — so I sanitize the name now, dots and colons to underscores, create-name and target-name finally in agreement.

Leave the conversation from where it was. That was always the point.
