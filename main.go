package main

import (
	"os"
	"github.com/amig3n/gwtr/cli"
)

func main() {
	rootCmd := cli.NewCLI()
	err := rootCmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
