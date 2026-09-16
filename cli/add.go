package cli

import (
	"github.com/spf13/cobra"
	"github.com/amig3n/gwtr/worktree"
)


func (cli *CLI) addAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new worktree",
		Long:  "Add new worktree",
		RunE: func(cmd *cobra.Command, args []string) error {

			branch, err := cmd.Flags().GetString("branch")
			if err != nil {
				return err
			}
			
			path, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			// load state
			state, err := cli.service.LoadState()
			if err != nil {
				return err
			}
			
			// perform addition
			state.Add(worktree.Worktree{
				Branch: branch,
				Path: path,
				Deleted: false,
			})

			// save state
			err = cli.service.SaveState(state)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("branch", "b", "", "Branch name for the new worktree")

	cmd.Flags().StringP("path", "p", "", "Path for the new worktree")

	return cmd
}

