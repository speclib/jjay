# One command to rule your project sessions

jjay now knows how to open a dedicated tmux session for any jj repo.

Until now, using jjay meant you were already sitting in the right tmux session, in the right directory. That's fine for one project, but when you're juggling three repos with parallel agents in each, the manual setup adds up: create a session, name it something sensible, cd into the repo, repeat. `jjay session-open <path>` collapses that into a single step. Point it at a jj repo and you land in a tmux session named `jjay-><dirname>`, working directory already set.

The implementation is deliberately minimal — resolve the path, check for `.jj/`, check for duplicate sessions, create detached, switch. Two tmux commands under the hood. We went with `filepath.Base()` for the session name because it reads well in `tmux list-sessions` and matches the mental model of "one session per project." The known trade-off is dirname collisions across different parent directories, but that's a problem for later.

This is the first piece of a session management layer (`internal/session/`) that will grow into `session-list` and `session-close`. The bigger picture: jjay should handle the full lifecycle from "I want to work on this project" to "spin up five agents" to "merge and clean up" — and now it handles the first step.
