package cli

import (
	"github.com/spf13/cobra"
)

type CLI struct {
	RootCmd *cobra.Command
}

func NewCLI() *CLI {
	// init CLI
	rootCmd := &cobra.Command{
		Use:   "gwtr",
		Short: "GWTR",
		Long:  "Go-based worktree manager for git. For anybody who works with multiple branches at one time.",
	}
	
	var cli CLI = CLI{
		RootCmd: rootCmd,
	}

	// add subcommands
	rootCmd.AddCommand(cli.addListCmd())
	// TODO add more subcommands

	return &cli
}
