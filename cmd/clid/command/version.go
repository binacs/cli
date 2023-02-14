package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/binacs/cli/version"
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version Command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s.%s.%s, CommitHash: %s\n", version.Maj, version.Min, version.Fix, version.GitCommit)
		},
	}
)
