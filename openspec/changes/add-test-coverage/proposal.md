## Why

The project has unit and integration tests but no way to measure or report code coverage. Adding a `make coverage` target and a README badge makes coverage visible, encouraging the team to maintain and improve it over time.

Bean: [jjay-znb0](../../../.beans/jjay-znb0--add-test-coverage-and-integration-test-makefile-ta.md)

## What Changes

- Add `make coverage` Makefile target that runs unit tests with `-coverprofile`, generates an HTML report, and prints the coverage percentage
- Add `make badge` Makefile target that patches the README with a shields.io coverage badge using the extracted percentage
- Add `coverage.out` and `coverage.html` to `.gitignore`
- Add a coverage badge placeholder at the top of the README

## Capabilities

### New Capabilities
- `coverage-reporting`: Coverage profiling, HTML report generation, percentage extraction, and README badge updating

### Modified Capabilities
- `test-infrastructure`: Adding `coverage` and `badge` Makefile targets to the existing dev target requirements

## Impact

- **Makefile**: Two new targets (`coverage`, `badge`)
- **.gitignore**: Two new entries
- **README.md**: Badge added after the hero image
- **No new dependencies**: Uses `go tool cover` (stdlib) and `sed` for badge patching; shields.io for badge rendering
