package main

import (
	"os"

	cmd "github.com/binacs/cli/cmd/clid/command"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.StartCmd,
		cmd.VersionCmd,
	)
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
