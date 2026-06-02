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
