# Release Process

You can't ship what you can't release.

Version was hardcoded as `"dev"` in main.go and `"0.0.1"` in the flake. No automation, no binaries, no GitHub releases. I fixed that in one sweep.

A single `VERSION` file is the source of truth. Go reads it via `go:embed`, the Nix flake reads it with `builtins.readFile`. Goreleaser builds multi-platform binaries. GitHub Actions triggers on `v*` tags. A release script powered by gum walks you through the whole thing interactively — bump version, update vendorHash, commit, tag, push.

From zero to automated releases. My flock ships properly now.
