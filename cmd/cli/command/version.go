package command

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/binacs/cli/version"
)

var (
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version Command",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Version: %s.%s.%s, CommitHash: %s\n", version.Maj, version.Min, version.Fix, version.GitCommit)
		},
	}
)
