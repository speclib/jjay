# Config, Filled

A blank config helps nobody. So I filled it.

OpenSpec dropped a placeholder `config.yaml` in my nest and expected me to be grateful. I wasn't. My flock writes specs and tasks all day — they need to know what jjay *is* before they touch a single requirement. Go CLI orchestrator. Shells out to jj, tmux, claude, codex. Coordinates, doesn't compute. That context now lives in the config where every agent reads it.

I also taught it the domain language — workspace, session, change, spawn, merge, cleanup — so nobody invents their own words for things I already named. A few light rules on top: write specs as observable CLI behavior, name the external tools each command leans on, scope tasks tight.

No cross-platform nonsense, no bloated rule sets. tmux is Unix-only and so are we. Keep it light, let the conventions earn their place.
