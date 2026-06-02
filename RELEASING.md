# Releasing jjay

## Prerequisites

- On the `main` branch with a clean working tree
- `gum` installed (available in nix devShell)
- `nix` installed (optional, for vendorHash update)
- Push access to the repository

## Pre-release checklist

- [ ] All tests pass: `make test`
- [ ] Build succeeds: `make build`
- [ ] Lint passes: `make lint`
- [ ] CHANGELOG.md has an `## Unreleased` section with entries
- [ ] All changes are committed and pushed to main

## Release steps

Run the release script:

```bash
./scripts/release.sh
```

The script will:
1. Verify safety checks (clean tree, on main, changelog ready)
2. Prompt for version bump type (major/minor/patch)
3. Update `VERSION` file
4. Update `CHANGELOG.md` (moves Unreleased entries under new version heading)
5. Update `vendorHash` in `flake.nix` (if nix is installed)
6. Create a git commit and tag
7. Push to origin

GitHub Actions will then build and publish the release automatically.

## Verification

After the release workflow completes:

1. Check the [GitHub Actions](../../actions) tab for the release workflow run
2. Verify the [GitHub release](../../releases) page has the new version
3. Download a binary and verify: `./jjay version`
4. Verify `nix build` works with the updated flake (if vendorHash was updated)

## Troubleshooting

### goreleaser fails in CI
- Check that the tag matches `v*` pattern (e.g., `v0.2.0`)
- Verify `goreleaser check` passes locally
- Check GitHub Actions logs for build errors

### vendorHash mismatch after release
- Run `nix build` locally, copy the correct hash from the error
- Update `vendorHash` in `flake.nix` and push a follow-up commit

### Wrong version released
- Delete the tag: `git tag -d vX.Y.Z && git push origin :refs/tags/vX.Y.Z`
- Delete the GitHub release manually
- Fix the issue and re-run the release script
