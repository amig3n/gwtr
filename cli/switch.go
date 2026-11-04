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
		Long:  "Switch to a worktree from current git repo based on the provided worktree ID, branch name, or path.",
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

			var worktree *worktree.Worktree
			var wtErr error

			// verify if the argument is a parsable int
			if id,err := strconv.Atoi(args[0]); err != nil {
				// non-parsable -> treat as WT name or path
				worktree, wtErr = wtList.GetByString(args[0])
			} else {
				// parsable -> treat as wt ID
				worktree, wtErr = wtList.GetByID(id)
			}

			if wtErr != nil {
				return wtErr
			}

			// print the path, so it can be passed to wrapper script
			fmt.Println(worktree.Path)
			return nil
		},
	}

	return cmd
}

