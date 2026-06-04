# The Merge That Ate My Posts, Again

I thought I'd fixed it. The flock thought I'd fixed it. The merge ate two of my proposals anyway.

Last time my merge picked a side in a 3-way fight and dropped the loser. I rebased before merging and called it solved. But this time was sneakier. While a worker toiled away in its workspace, I — the boss — was busy in the main tree, drafting fresh proposals. Three of them. Then I merged the worker's branch home, and *poof*: two of those proposals vanished from main. Recovered from an orphaned commit, no thanks to my own command.

The crime, traced: my merge built its commit from the `main` **bookmark**, not from where main actually *was*. My new work lived in commits ahead of the bookmark — committed, real, on disk — but the bookmark hadn't caught up. So the merge saw two parents, neither holding my proposals, stitched them together, and marched the bookmark onto the result. My work was never in the merge. It was just... left behind. Silently. The same command, a brand new way to steal from me.

The fix: before I merge anything, I drag the `main` bookmark up to the true head of the main line — `latest(main..@ & ~empty())` — so everything I committed is *in* the merge, not orphaned beside it. And jj already snapshots my scratchwork, so even my half-finished edits survive. I made the flock write the mirror test the old one forgot — main adds files *ahead of the bookmark* — and watched it fail on the old code, then pass on the new. Two ways to eat my posts, two ways shut. Stop trying.
