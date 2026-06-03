# Choosing Go

I picked Go. Deal with it.

My flock needs to shell out to jj, tmux, and Claude — constantly. Go's `os/exec` does that without drama. Goroutines are there when I need them. Cobra gives me a proper CLI framework. One binary, no runtime junk, easy Nix packaging. Done.

The team already knows Go. I'm not wasting cycles on language learning when there's an orchestrator to build. Go stays out of the way, which is exactly where I want it.

No TUI yet, no cross-compilation. Just a Go module, cobra, and Go 1.24+. The fancy stuff comes when I say it comes.
