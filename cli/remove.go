package cli

import (
	"github.com/spf13/cobra"
	"github.com/amig3n/gwtr/worktree"
	"strconv"
	"fmt"
)

func userConfirm(prompt string) bool {
	responsePrompt := "yY/nN"	
	var response string

	for {
	fmt.Printf("%s [%s]: ", prompt, responsePrompt)	
		fmt.Scanln(&response)
		switch response {
			case "y", "Y":
				return true
			case "n", "N":
				return false
			default: 
				continue
		}
	}
}

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
			
			forceFlag, err := cmd.Flags().GetBool("force")
			if err != nil {
				return err
			}

			if forceFlag || userConfirm(
				fmt.Sprintf("Are you sure you want to delete a worktree %s ?", worktree.Path)) {
				// perform deletion
				worktree.Delete()
				
				// save state
				err = service.SaveState(state)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("Worktree not deleted")
			}


			return nil
		},
	}

	cmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	return cmd
}

