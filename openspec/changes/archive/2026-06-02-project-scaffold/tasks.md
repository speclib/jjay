# Tasks: project-scaffold

## 1. Go Module and CLI

- [x] 1.1 Run `go mod init jjay` to create go.mod
- [x] 1.2 Create `cmd/jjay/main.go` with cobra root command and `version` subcommand
- [x] 1.3 Create `internal/` directory
- [x] 1.4 Run `go mod tidy` to resolve dependencies

## 2. Test Framework

- [x] 2.1 Create `Makefile` with `build`, `test`, `lint` targets
- [x] 2.2 Create `cmd/jjay/main_test.go` with initial test
- [x] 2.3 Verify with `make test`, `make build`, `make lint`

## 3. Nix Flake

- [x] 3.1 Create `flake.nix` with buildGoModule, dev shell, multi-platform support
- [x] 3.2 Run `nix flake update` to generate `flake.lock`
- [x] 3.3 Build to get vendorHash, update flake.nix with real hash
- [x] 3.4 Verify `nix build`, `nix run . -- version`, `nix develop -c go version`

## 4. Documentation

- [x] 4.1 Update README.md with installation section (Nix + go install)
