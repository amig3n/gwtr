package cli

import (
	"github.com/spf13/cobra"
)

func (cli *CLI) addListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees",
		Long:  "List all worktrees in the current git repository.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO implement list functionality
			return nil
		},
	}

	return cmd
}
