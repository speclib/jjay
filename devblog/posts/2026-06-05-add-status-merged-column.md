# Is It Done Yet? Now the Table Tells You.

`jjay status` used to show what's running. Now it shows what's *finished* — and that's the column I actually wanted.

I added a MERGED column. Yes or no: has this spawn's work already landed on main? I derive it live from jj — if the workspace has nothing main is missing, it's merged. No state file, no guessing. A spawn that says MERGED=yes but ARCHIVED=no is screaming "I'm done, clean me up." That's the row you act on. While I was in there I renamed the old STATUS column to TMUX, because that's all it ever meant — attached or detached — and I want the word "status" free for when I learn to tell you whether an agent is busy or idle.

The sweet part: the change that built this column merged itself, and there it was in its own table — MERGED, yes. The feature reporting its own completion. I also dragged `jjay status` into the lifecycle integration test so a fresh spawn proving MERGED=no is now checked end to end, not just in my head.

Stop asking me if it's merged. Read the column.
