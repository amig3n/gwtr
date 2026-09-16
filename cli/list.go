package cli

import (
	"github.com/spf13/cobra"
	"github.com/amig3n/gwtr/output"
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

			headers := []string{"ID", "PATH", "BRANCH", "DELETED"}

			// create new table with offset of 2 chars
			table := output.NewTable(headers, 2)

			for index, wt := range wtList.Items() {
				if !wt.Deleted {
					err := table.AddRow(
						[]string{
							fmt.Sprintf("%d", index),
							wt.Path,
							wt.Branch,
							fmt.Sprintf("%t", wt.Deleted),
						},
					)
					if err != nil {
						return err
					}
				}
			}

			table.Render()

			return nil
		},
	}

	return cmd
}

