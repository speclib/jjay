# Tasks: release-process

## 1. Version embedding

- [x] 1.1 Create `VERSION` file with `0.1.0`
- [x] 1.2 Update `cmd/jjay/main.go` to embed VERSION via `go:embed` (keep ldflags override for goreleaser)
- [x] 1.3 Update `flake.nix` to read version from VERSION file via `builtins.readFile`
- [x] 1.4 Verify `go run ./cmd/jjay version` prints `0.1.0`

## 2. goreleaser configuration

- [x] 2.1 Create `.goreleaser.yaml` with builds for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
- [x] 2.2 Configure ldflags to inject version from git tag
- [x] 2.3 Configure archive format (tar.gz) and checksum generation
- [x] 2.4 Add goreleaser and gum to `flake.nix` devShell
- [x] 2.5 Verify `goreleaser check` passes
- [x] 2.6 Verify `goreleaser build --snapshot --clean` creates binaries

## 3. GitHub Actions workflow

- [x] 3.1 Create `.github/workflows/release.yml` triggered on `v*` tag push
- [x] 3.2 Configure checkout, Go setup, and goreleaser action with `GITHUB_TOKEN`

## 4. Release script

- [x] 4.1 Create `scripts/release.sh` with safety checks (clean tree, on main, CHANGELOG has Unreleased, tag doesn't exist)
- [x] 4.2 Add interactive version bump via gum (major/minor/patch)
- [x] 4.3 Add VERSION file update, CHANGELOG.md update (Unreleased → version + date)
- [x] 4.4 Add nix vendorHash auto-update step (skip if nix not installed)
- [x] 4.5 Add git commit, tag creation, and push
- [x] 4.6 Make script executable

## 5. Documentation

- [x] 5.1 Create `RELEASING.md` with pre-release checklist, release steps, verification, troubleshooting

## 6. Verification

- [x] 6.1 Verify `make test`, `make build`, `make lint` all pass
- [x] 6.2 Verify `nix build` succeeds with new vendorHash
