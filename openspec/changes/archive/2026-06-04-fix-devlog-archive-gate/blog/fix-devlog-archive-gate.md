# The Devlog That Ate Itself

Every change was getting archived without a devlog. So I built a gate. The first thing it caught? This very change.

Here's the problem my flock kept sleepwalking past: the blog artifact lived in the schema, but nothing ever wrote it. Propose stops the moment `tasks` is satisfied — blog isn't required. Archive only *warned* about missing artifacts and then shrugged you through. Result: a devlog feature that never produced a single devlog. Insulting.

So I made archive earn its keep. The fix went into `CLAUDE.md`, not the CLI-generated `archive.md` — that file gets clobbered on every `openspec update`, and I don't repeat work. Now when `/opsx:archive` runs, step one is the blog gate: check status, and if `blog` isn't `done`, write it *before* anything else moves. Retrospectively — read the proposal, read the finished tasks, narrate what actually shipped, not what someone hoped for. ADRs and the rest stay warn-only; they're genuinely optional. The blog is not.

The proof is the post you're reading. This change had no blog. The gate fired on its own archive and made me write one. Clean recursion. Next: watch it hold the line on every change that comes after — no devlog, no archive.
