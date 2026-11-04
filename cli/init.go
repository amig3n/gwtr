package cli

import (
	"github.com/spf13/cobra"
)

func WrapperScriptContent() string {
	// TODO: fix up this script
	return `gwtr() {
	result="$(command gwtr "$@")" || return $?
	if [ -n "$result" ] && [ -d "$result" ]; then
		cd "$result" || return 1
	else
		printf '%s\n' "$result"
	fi

	return 0
}
	`
}

func (cli *CLI) addInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize GWTR",
		Long:  "Initialize GWTR by creating a state file and populating it with the current list of worktrees.",
		RunE: func(cmd *cobra.Command, args []string) error {
			wrapperFlag, err := cmd.Flags().GetBool("wrapper")
			if err != nil {
				return err
			}

			if wrapperFlag {
				// only return wrapper script content
				cmd.Print(WrapperScriptContent())

			} else {
				// prepare service
				service, err := cli.PrepareService(cli.logger)
				if err != nil {
					return err
				}

				// init blank state file if not exists
				err = service.InitState()
				if err != nil {
					return err
				}
			}

			return nil
		},

	}
	
	cmd.Flags().BoolP("wrapper", "w", false, "Only return wrapper script content")

	return cmd
}
