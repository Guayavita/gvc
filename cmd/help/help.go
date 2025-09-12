package help

import (
	"github.com/spf13/cobra"
	"jmpeax.com/guayavita/gvc/internal/commons"
	"jmpeax.com/guayavita/gvc/internal/term"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Long:  "Show build/version information",
	Run: func(cmd *cobra.Command, args []string) {
		term.Info("Version: %s\nBuild: %s\nGitCommit: %s", commons.Version, commons.Build, commons.GitCommit)
	},
}

func Commands() []*cobra.Command {
	return []*cobra.Command{versionCmd}
}
