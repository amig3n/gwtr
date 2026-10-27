package cli

import (
	"github.com/spf13/cobra"
	"strconv"
	"fmt"
	"github.com/amig3n/gwtr/worktree"
)

func (cli *CLI) addSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch to a worktree",
		Long:  "Switch to a worktree from currrent git repo based on provided worktree ID.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// prepare service
			service, err := cli.PrepareService(cli.logger)
			if err != nil {
				return err
			}

			// call service to list worktrees
			wtList, err := service.LoadState()
			if err != nil {
				return err
			}

			// check if identifier is a number or string
			var worktree *worktree.Worktree
			var wtErr error

			if id,err := strconv.Atoi(args[0]); err != nil {
				// treat as string
				worktree, wtErr = wtList.GetByString(args[0])

			} else {
				worktree, wtErr = wtList.GetByID(id)
			}

			if wtErr != nil {
				return wtErr
			}

			fmt.Println(worktree.Path)
			return nil
		},
	}

	return cmd
}

