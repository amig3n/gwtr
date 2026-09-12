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
			// init blank state file if not exists
			err := cli.service.InitState()
			if err != nil {
				return err
			}
			// TODO transform worktrees into state format

			// save the state imidiately
			//err := cli.service.state.Save(worktrees)
			//if err != nil {
			//	return err
			//}

			return nil
		},

	}

	return cmd
}
