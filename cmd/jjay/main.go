package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"jjay/internal/spawn"
)

// version is set via ldflags at build time (from VERSION file or git tag).
// Default value is used for `go run` without ldflags.
var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "jjay",
	Short: "Manage parallel AI agent sessions with jj, tmux, and openspec",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of jjay",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var spawnCmd = &cobra.Command{
	Use:   "spawn <change-name>",
	Short: "Create workspace + tmux window + launch agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return spawn.Spawn(args[0])
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(spawnCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
