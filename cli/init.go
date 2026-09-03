package cli

import (
	"github.com/spf13/cobra"
)

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
