# Cleanup Finally Gets Its Feathers Checked

Every command in my flock had tests. Every one — except the one that tears things down.

`cleanup` is the dangerous one. It kills tmux windows, forgets jj workspaces, and deletes directories. It's built to be tolerant — skip what's missing, never throw a tantrum — but nobody had ever pinned that promise down with a test. So I sent an agent in to fix it. No production code touched. Just proof that cleanup does what it claims.

Now `removeDirectory` is checked both ways: it removes a real dir under a temp root, and it shrugs off a missing one without a panic. `tmuxTarget` is nailed down — bare window name when there's no session, `session:window` when there is. And `killWindow`, `forgetWorkspace`, and the full `Cleanup` orchestration all get poked with a change that doesn't exist, proving they stay calm whether or not tmux or jj are even running. That last part matters: the tests pass in CI with no live server, because they assert *tolerance*, not a specific branch.

Coverage on `internal/cleanup` went from a flat **0%** to **59.3%**, and the full integration suite is still green. The package that cleans up after everyone finally has someone watching it.

Next: keep dragging the stragglers up to parity. Tolerance is a feature — now it's a tested one.
