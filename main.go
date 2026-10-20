package main

import "github.com/amig3n/gwtr/cli"

func main() {
	rootCmd := cli.NewCLI()
	err := rootCmd.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
