package main

import (
	"os"
	"github.com/amig3n/gwtr/cli"
)

func main() {
	rootCmd := cli.NewCLI()

	// check if the first argument is valid command
	// if not, prepend switch 
	cmd, _, err := rootCmd.RootCmd.Find(os.Args[1:])
	if err != nil  || cmd == nil {
		args := append([]string{"switch"}, os.Args[1:]...)
		rootCmd.RootCmd.SetArgs(args)
	}

	err = rootCmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
