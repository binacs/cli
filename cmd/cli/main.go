package main

import (
	"os"

	cmd "github.com/binacs/cli/cmd/cli/command"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.CosCmd,
		cmd.CryptoCmd,
		cmd.PastebinCmd,
		cmd.TinyurlCmd,
		cmd.UserCmd,
		cmd.VersionCmd,
	)
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
