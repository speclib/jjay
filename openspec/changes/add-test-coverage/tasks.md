## 1. Gitignore

- [ ] 1.1 Add `coverage.out` and `coverage.html` entries to `.gitignore`

## 2. Makefile Targets

- [ ] 2.1 Add `make coverage` target: run `go test -coverprofile=coverage.out ./...`, generate HTML report, print total percentage
- [ ] 2.2 Add `make badge` target: extract percentage from coverage profile, determine badge color, `sed` the README badge URL
- [ ] 2.3 Update `.PHONY` line to include `coverage` and `badge`

## 3. README Badge

- [ ] 3.1 Add placeholder shields.io coverage badge to README between the hero image `</p>` and `# jjay` heading

## 4. Verify

- [ ] 4.1 Run `make coverage` and confirm `coverage.out`, `coverage.html`, and percentage output
- [ ] 4.2 Run `make badge` and confirm README badge is updated with correct percentage and color
