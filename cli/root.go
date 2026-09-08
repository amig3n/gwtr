package cli

import (
	"github.com/spf13/cobra"
	"github.com/amig3n/gwtr/service"
)

type CLI struct {
	RootCmd *cobra.Command
	service *service.Service
}

func NewCLI(service *service.Service) *CLI {
	// init CLI
	rootCmd := &cobra.Command{
		Use:   "gwtr",
		Short: "GWTR",
		Long:  "Worktres manager for git. For anybody who works with multiple branches at one time.",
	}
	
	// create CLI object with passing Service from outside
	var cli CLI = CLI{
		RootCmd: rootCmd,
		service: service,
	}

	// add subcommands
	rootCmd.AddCommand(cli.addListCmd())
	rootCmd.AddCommand(cli.addInitCmd())
	rootCmd.AddCommand(cli.addAddCmd())
	// TODO add more subcommands

	return &cli
}
