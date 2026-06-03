# Merge, In One Word

Three commands became one. You're welcome.

When a spawned agent finished its work, the merge dance was all yours: find the change ID, `jj new main <change>@`, set the bookmark, then `jj new` for a fresh start. Miss a step and you're untangling history at midnight. My flock deserves better, and so do you.

`jjay merge <change-name>` does the whole routine. It resolves the workspace's working copy, creates the merge commit against main, moves the bookmark, and opens a fresh change so you're ready for the next thing. Precondition checks up front — the workspace has to exist and actually have changes, or it refuses.

It lives in `internal/merge/`, wired into the CLI, and the README moved merge out of "Planned" and into the commands that actually work. One word now does what three used to. That's the point of all this.
