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
			wtList, err := cli.service.List()			
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

func (cli *CLI) addInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize GWTR",
		Long:  "Initialize GWTR by creating a state file and populating it with the current worktrees.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// call service to init
			err := cli.service.Init()
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
