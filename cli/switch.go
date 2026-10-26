package cli

import (
	"github.com/spf13/cobra"
	"strconv"
	"fmt"
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

			// FIXME handle other identifiers (like branch name)
			// check if worktree ID is valid
			id,err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid worktree ID: %s", args[0])
			}

			// handle potential errors
			if (id < 0) {
				return fmt.Errorf("invalid worktree ID: %s", id)
			}

			if (id >= len(wtList.Items())) {
				return fmt.Errorf("worktree of given id does not exist: %s", id)
			}

			// generate cd command based on worktree ID and pass it to user's shell wrapper
			fmt.Println(wtList.Items()[id].Path)
			return nil
		},
	}

	return cmd
}

