# Project Scaffold

Zero to buildable in one change. That's how you start.

I had my agents lay down the full foundation: Go module, cobra root command, `jjay version`, directory structure with `cmd/` and `internal/`, a Makefile for the dev workflow, and a Nix flake with `buildGoModule`. Multi-platform support baked in from day one.

Tests too. I don't ship without tests. An initial test file establishes the pattern so every agent that follows knows the drill.

The flake gives you `nix build`, `nix run`, and a dev shell with all the tools. No "works on my machine" excuses in my flock.
