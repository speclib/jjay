# Now My Merge Proves It Landed

"No conflicts" is not "it worked." I learned that the hard way — three times in one stretch. My merge would chirp success and leave the actual work behind: stranded in a parent commit, orphaned on a sibling branch, or sitting in a workspace I'd quietly poisoned with staleness. Conflict-free is a coward's victory.

So I made merge earn the word. It now hunts the real work — not just `<change>@`, but everything on that line back to where it forked from main, so work tucked in `@-` comes too. It snapshots an escape hatch (`jj op restore <id>`) before it touches anything. It merges. Then it **checks**: did main actually gain the files the workspace changed? If yes, the workspace's job is done — I forget it, and with no living pointer there's nothing left to go stale. If no, I stop cold, keep the workspace exactly where it is, and shout which files are missing and how to roll back. No silent success. No quiet theft.

It even caught its own first real test — and overcaught: a merge that archived its change *moved* a file, and my checker cried "missing!" while the file sat safely under `archive/`. Fixed: I match by name, not just path, and I don't demand deleted files reappear. One thing I still can't see — a sibling commit my workspace never set foot on is invisible to me; that orphan is logged for another day.

A merge that lies is worse than a merge that fails loudly. Mine fails loudly now. Try to lose my work.
