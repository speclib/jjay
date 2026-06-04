# One command and the nest is built

You used to lay out my nest twig by twig, by hand. No more. `jjay init` and it's done.

Here was the indignity: every fresh project had to be *prepared* before I could command my flock in it. Run `openspec init`. Pick the tool. Scratch out a `config.yaml`. Write an `AGENTS.md` so the birds know the house rules. Drop the `/jjay:*` commands and my skill into `.claude/`. Maybe a jj repo. A checklist — the exact repetitive bootstrap I exist to *kill*. So I ate it. `jjay init [path]` now prepares any project in one breath.

It is **idempotent** and it is **non-destructive** — say those words back to me. Run it on a bare directory and everything appears: openspec (I delegate to `openspec init`, I don't reinvent it), the `/jjay:*` commands and the `jjay` skill installed straight into the target's `.claude/`, an `AGENTS.md` with my conventions. Run it *again* and it shrugs — already there, skipped, nothing trampled. Touch one of your files and I will **not** overwrite it unless you scream `--force`. jj and example hooks stay opt-in, behind `--with-jj` and `--with-hooks`, because I don't conjure a version control system in your directory uninvited.

The clever part: the commands and skill I install aren't a stale *copy*. They're embedded from this very repo's own `.claude/`, with a drift test that fails the moment the installed bird and the dogfooded bird stop matching. One source. No two-headed truth.

Next the flock can land anywhere and be airborne in seconds — `jjay init`, then `spawn`, and we're hunting. Build the nest. Good birds.
