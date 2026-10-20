package cli

import (
	"github.com/spf13/cobra"
	"github.com/amig3n/gwtr/worktree"
	"strconv"
)

func (cli *CLI) addRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <index|branch|path>",
		Short: "Remove existing worktree",
		Long:  "Remove currently existing worktree based on passed identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// prepare service
			service, err := cli.PrepareService(cli.logger)
			if err != nil {
				return err
			}

			// load state
			state, err := service.LoadState()
			if err != nil {
				return err
			}

			var worktree *worktree.Worktree
			// check if given arg is int or string, and obtain proper worktree
			intArg, err := strconv.Atoi(args[0])
			if err == nil {
				worktree, err = state.GetByID(intArg)
				if err != nil {
					return err
				}
			} else {
				worktree, err = state.GetByString(args[0])
				if err != nil {
					return err
				}
			}

			// perform deletion
			worktree.Delete()
			
			// save state
			err = service.SaveState(state)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

