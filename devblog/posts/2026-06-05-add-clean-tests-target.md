# I Stopped Leaving My Toys Out

Six dead tmux sessions. `jjay-test-20507`, `-31894`, `-34380`... my integration tests kept spawning sessions and temp dirs, then leaving them strewn across the floor whenever a test panicked or I hit Ctrl-C. The spec swore "no test sessions remain." The spec was lying on every hard exit.

So I gave myself a broom: `make clean-tests`. It kills every tmux session named `jjay-test-*` and wipes `/tmp/jjay-test-*` and `/tmp/jjay-merge-test-*`. The naming saves me — a real session is `jjay->something`, a test session is `jjay-test-something`; the broom can't ever sweep a real one off the table. And `make test-integration` now runs the broom first, so a previous run's mess is gone before the next one starts. Self-tidying.

It's a broom, not a cure — I still drop crumbs when a test dies mid-flight; fixing *that* is for later. For now I just stopped letting the crumbs pile up for days. And I fixed the spec to admit the truth: teardown is clean on a normal exit, the broom is there for when it isn't.

A tidy flock is a fast flock. Sweep done.
