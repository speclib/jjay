## ADDED Requirements

### Requirement: Coverage profiling target
The Makefile SHALL have a `coverage` target that runs unit tests with coverage profiling and produces both a machine-readable profile and an HTML report.

#### Scenario: Generate coverage profile
- **WHEN** `make coverage` is executed
- **THEN** `go test -coverprofile=coverage.out ./...` runs
- **THEN** a `coverage.out` file is produced in the project root

#### Scenario: Generate HTML report
- **WHEN** `make coverage` is executed
- **THEN** `go tool cover -html=coverage.out -o coverage.html` runs
- **THEN** a `coverage.html` file is produced in the project root

#### Scenario: Print coverage percentage
- **WHEN** `make coverage` is executed
- **THEN** the total coverage percentage is printed to stdout (e.g. `Coverage: 74.2%`)

### Requirement: Coverage badge target
The Makefile SHALL have a `badge` target that updates the README with a shields.io coverage badge reflecting the current coverage percentage.

#### Scenario: Badge updates README
- **WHEN** `make badge` is executed
- **THEN** the coverage percentage is extracted from the coverage profile
- **THEN** the README.md coverage badge URL is updated with the current percentage
- **THEN** the badge color reflects the coverage level (green ≥80%, yellow ≥60%, red <60%)

#### Scenario: Badge depends on coverage
- **WHEN** `make badge` is executed without a prior `make coverage` run
- **THEN** coverage is generated first before updating the badge

### Requirement: Coverage artifacts excluded from VCS
The `.gitignore` SHALL include `coverage.out` and `coverage.html` so generated coverage files are not committed.

#### Scenario: Gitignore entries present
- **WHEN** the `.gitignore` file is inspected
- **THEN** it contains entries for `coverage.out` and `coverage.html`

### Requirement: README coverage badge
The README SHALL display a coverage badge after the hero image, showing the latest locally-measured coverage percentage.

#### Scenario: Badge visible in README
- **WHEN** the README.md is viewed
- **THEN** a shields.io coverage badge image is present between the hero image and the `# jjay` heading
