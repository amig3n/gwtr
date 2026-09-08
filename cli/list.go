package cli

import (
	"github.com/spf13/cobra"
	"fmt"
)

func (cli *CLI) addListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees",
		Long:  "List all worktrees in the current git repository.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// call service to list worktrees
			wtList, err := cli.service.LoadState()
			if err != nil {
				return err
			}

			fmt.Printf("Worktrees:\n")
			fmt.Printf("%v\n", wtList)
			// pass the result to the output generator
			return nil
		},
	}

	return cmd
}

