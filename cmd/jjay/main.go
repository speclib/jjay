package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"jjay/internal/cleanup"
	"jjay/internal/completion"
	jjinit "jjay/internal/init"
	"jjay/internal/merge"
	"jjay/internal/session"
	"jjay/internal/spawn"
	"jjay/internal/status"
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

var (
	spawnAgent         string
	spawnSession       string
	spawnWorkspaceRoot string
	proposalMode       string
)

// spawnCmd is the parent of the apply/proposal verbs. It has no bare-argument
// form (ADR-011): invoking `jjay spawn` without a verb prints usage and exits
// non-zero. RunE returns an error so the exit code is non-zero.
var spawnCmd = &cobra.Command{
	Use:   "spawn <verb>",
	Short: "Spawn an agent workspace (apply an existing change, or seed a new proposal)",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = cmd.Usage()
		return fmt.Errorf("spawn requires a verb: 'apply' or 'proposal'")
	},
}

var spawnApplyCmd = &cobra.Command{
	Use:   "apply <change-name>",
	Short: "Isolate an existing openspec change and run /opsx:apply (workspace app-<change>)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return spawn.Spawn(args[0], spawn.SpawnOptions{
			Agent:         spawnAgent,
			Session:       spawnSession,
			WorkspaceRoot: spawnWorkspaceRoot,
		})
	},
}

var spawnProposalCmd = &cobra.Command{
	Use:   "proposal <prompt>",
	Short: "Seed a new proposal spawn from a free-text prompt (workspace prop-<slug>)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := spawn.Mode(proposalMode)
		if mode != spawn.ModeExplore && mode != spawn.ModePropose {
			return fmt.Errorf("invalid --mode %q: must be 'explore' or 'propose'", proposalMode)
		}
		return spawn.SpawnProposal(args[0], mode, spawn.SpawnOptions{
			Agent:         spawnAgent,
			Session:       spawnSession,
			WorkspaceRoot: spawnWorkspaceRoot,
		})
	},
}

var (
	initYes        bool
	initForce      bool
	initWithJJ     bool
	initWithHooks  bool
	initNoClaude   bool
	initNoOpenspec bool
	initNoAgents   bool
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Prepare a project for orchestration by jjay (openspec, Claude integration, AGENTS.md)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) == 1 {
			path = args[0]
		}
		return jjinit.Init(path, jjinit.InitOptions{
			Yes:        initYes,
			Force:      initForce,
			WithJJ:     initWithJJ,
			WithHooks:  initWithHooks,
			NoClaude:   initNoClaude,
			NoOpenspec: initNoOpenspec,
			NoAgents:   initNoAgents,
			Out:        cmd.OutOrStdout(),
		})
	},
}

var mergeCmd = &cobra.Command{
	Use:   "merge <change-name>",
	Short: "Merge a workspace's work into main",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return merge.Merge(args[0])
	},
}

var sessionOpenCmd = &cobra.Command{
	Use:   "session-open <path>",
	Short: "Create and switch to a tmux session for a jj repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return session.Open(args[0])
	},
}

var (
	statusSession       string
	statusWorkspaceRoot string
)
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List spawned workspaces and whether each has a tmux window",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		spawns, mainRoot, err := status.List(statusSession, statusWorkspaceRoot)
		if err != nil {
			return err
		}
		status.Render(cmd.OutOrStdout(), mainRoot, spawns)
		return nil
	},
}

var (
	cleanupSession       string
	cleanupWorkspaceRoot string
)
var cleanupCmd = &cobra.Command{
	Use:   "cleanup <change-name>",
	Short: "Tear down workspace + tmux window + directory for a change",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cleanup.Cleanup(args[0], cleanup.CleanupOptions{
			Session:       cleanupSession,
			WorkspaceRoot: cleanupWorkspaceRoot,
		})
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(spawnCmd)
	spawnCmd.AddCommand(spawnApplyCmd)
	spawnCmd.AddCommand(spawnProposalCmd)
	rootCmd.AddCommand(mergeCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(sessionOpenCmd)
	rootCmd.AddCommand(statusCmd)

	// Per-verb change-name completion. This is the only place that knows the
	// verb↔function mapping — the redesign seam (ADR-009). The change-name
	// completion attaches to `spawn apply`'s argument; `spawn proposal` takes a
	// free-text prompt and offers no candidates (cobra default).
	spawnApplyCmd.ValidArgsFunction = completion.Spawnable
	// `spawn proposal` takes a free-text prompt: offer no candidates and suppress
	// file-name fallback (a prompt is not a path).
	spawnProposalCmd.ValidArgsFunction = func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	mergeCmd.ValidArgsFunction = completion.Mergeable
	cleanupCmd.ValidArgsFunction = completion.Cleanable

	initCmd.Flags().BoolVar(&initYes, "yes", false, "accept creation defaults without prompting (does not authorize overwriting existing files)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite existing files")
	initCmd.Flags().BoolVar(&initWithJJ, "with-jj", false, "initialize a jj repo if not already present")
	initCmd.Flags().BoolVar(&initWithHooks, "with-hooks", false, "scaffold example (commented) hooks")
	initCmd.Flags().BoolVar(&initNoClaude, "no-claude", false, "skip installing the jjay Claude integration")
	initCmd.Flags().BoolVar(&initNoOpenspec, "no-openspec", false, "skip the openspec step")
	initCmd.Flags().BoolVar(&initNoAgents, "no-agents", false, "skip writing AGENTS.md")

	// Shared spawn flags live on both verb subcommands.
	for _, c := range []*cobra.Command{spawnApplyCmd, spawnProposalCmd} {
		c.Flags().StringVar(&spawnAgent, "agent", "", "agent command template (placeholders: {change}/{prompt}, {wsdir})")
		c.Flags().StringVar(&spawnSession, "session", "", "tmux session to target (default: current)")
		c.Flags().StringVar(&spawnWorkspaceRoot, "workspace-root", "", "workspace root directory (default: ../<project>-workspaces)")
	}
	spawnProposalCmd.Flags().StringVar(&proposalMode, "mode", string(spawn.DefaultMode), "seed mode: 'explore' or 'propose'")

	cleanupCmd.Flags().StringVar(&cleanupSession, "session", "", "tmux session to target (default: current)")
	cleanupCmd.Flags().StringVar(&cleanupWorkspaceRoot, "workspace-root", "", "workspace root directory (default: ../<project>-workspaces)")

	statusCmd.Flags().StringVar(&statusSession, "session", "", "tmux session to inspect (default: current)")
	statusCmd.Flags().StringVar(&statusWorkspaceRoot, "workspace-root", "", "workspace root directory (default: ../<project>-workspaces)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
