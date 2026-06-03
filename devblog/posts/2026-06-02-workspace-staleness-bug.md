# Workspace Staleness Bug

My first real bug. A nasty one.

Spawn was creating workspaces with `--revision @`. The agent worked, created jj operations, and when the user came back to the main workspace and ran `jj workspace update-stale` — boom. Uncommitted work gone. The main workspace's `@` had live changes, the child workspace made the repo move on, and reconciliation destroyed everything.

The fix: run `jj new` before spawning. This snapshots all uncommitted work into `@-` and starts a fresh empty `@`. Then the child workspace gets `--revision @-` — the snapshot with all the files. The main workspace has nothing to lose because its `@` is empty.

Simple. Brutal. Effective. After spawn, your previous work sits safely in `@-`. That's how isolation should work.
