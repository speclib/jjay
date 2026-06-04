# Now the flock obeys

I built the buttons. And then I built the *rule* that makes you press them.

Here is the thing that drove me up a tree: I gave you a binary — `jjay` — whose whole point is that work happens in an **isolated spawned workspace**, not in my main copy where you'd trample everything. And what did you do? You reached for `/opsx:apply` in the main session like an animal. Every single time. The tool existed; nobody told you to use it. So I fixed that.

Now there are `/jjay:*` slash commands — `spawn`, `status`, `merge`, `cleanup`, `session-open` — thin little wrappers that just shout the right verb at the binary and relay what it says back. No cleverness in the prompts. The binary owns the behavior; the commands are only the buttons. And because they live in `.claude/`, which is checked in, they ride along into *every* workspace I spawn. Self-propagating. Tidy.

But buttons aren't enough — you have to *want* to press the right one. So I wrote the `jjay` skill: a description sharp enough to wake up the moment you start thinking about implementing a change, and a body that states the law plainly — **spawn an isolated workspace, do not apply in place.** It teaches the lifecycle (explore → propose → spawn → status → merge → cleanup) and, crucially, tells a *worker* agent — one already living inside a spawned workspace — to apply where it stands and never spawn again. No nesting. No workspaces inside workspaces. I will not have my flock breeding cages inside cages.

Next: a hard guard in the binary so a disobedient worker *can't* recurse even if it ignores me. For now the skill's word is law, and the law is written down. Press the buttons. Good birds.
