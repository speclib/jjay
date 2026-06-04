package main

import "testing"

func TestVersionCmd(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
}

func TestRootCmdHasVersionSubcommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("root command should have a version subcommand")
	}
}

func TestRootCmdHasStatusSubcommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "status" {
			found = true
			break
		}
	}
	if !found {
		t.Error("root command should have a status subcommand")
	}
}

func TestStatusCmd_RejectsArgs(t *testing.T) {
	// status takes no positional arguments; an extra arg is a usage error.
	if err := statusCmd.Args(statusCmd, []string{"extra-arg"}); err == nil {
		t.Error("expected status to reject positional arguments, got nil")
	}
	// And with no args, the Args validator accepts.
	if err := statusCmd.Args(statusCmd, []string{}); err != nil {
		t.Errorf("expected status to accept zero args, got: %v", err)
	}
}
