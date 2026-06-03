## ADDED Requirements

### Requirement: Merge creates merge commit
The `jjay merge <change>` command SHALL create a merge commit with two parents: the current `main` bookmark and the workspace's working copy (`<change>@`). The commit message SHALL be `merge <change> into main`.

#### Scenario: Successful merge
- **WHEN** `jjay merge feat-payments` is executed and workspace `feat-payments` exists with changes
- **THEN** `jj new main feat-payments@ -m "merge feat-payments into main"` is run
- **THEN** a merge commit exists with main and the workspace's change as parents

### Requirement: Merge updates main bookmark
After creating the merge commit, the command SHALL move the `main` bookmark to the new merge commit and create a fresh empty change for the user.

#### Scenario: Bookmark moved
- **WHEN** `jjay merge feat-payments` completes successfully
- **THEN** `main` bookmark points to the merge commit
- **THEN** the user's working copy is a fresh empty change on top of main

### Requirement: Merge requires workspace to exist
The command SHALL verify the jj workspace exists before proceeding.

#### Scenario: Workspace exists
- **WHEN** `jjay merge feat-payments` is executed and workspace `feat-payments` exists
- **THEN** merge proceeds

#### Scenario: Workspace does not exist
- **WHEN** `jjay merge feat-payments` is executed and workspace `feat-payments` does not exist
- **THEN** jjay exits with non-zero exit code
- **THEN** an error message indicates the workspace does not exist

### Requirement: Merge warns on empty workspace
The command SHALL warn if the workspace's working copy is empty (no changes). It SHALL still proceed after the warning.

#### Scenario: Empty workspace
- **WHEN** `jjay merge feat-payments` is executed and the workspace's `@` is empty
- **THEN** a warning is printed indicating the workspace has no changes
- **THEN** merge proceeds anyway

### Requirement: Merge requires change name argument
The command SHALL require exactly one argument: the change name.

#### Scenario: No argument
- **WHEN** `jjay merge` is executed without arguments
- **THEN** cobra prints usage help and exits with non-zero exit code

### Requirement: Merge does not push
The command SHALL NOT push to any remote. Pushing is a separate user action.

#### Scenario: No push
- **WHEN** `jjay merge feat-payments` completes successfully
- **THEN** no `jj git push` is executed
- **THEN** the user can push manually when ready
