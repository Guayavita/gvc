package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"jmpeax.com/guayavita/gvc/cmd/compiler"
	"jmpeax.com/guayavita/gvc/cmd/help"
)

var guayavitaCmd = &cobra.Command{
	Use:   "guayavita",
	Short: "Guayavita Compiler and Tooling",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		isVerbose, _ := cmd.Flags().GetBool("verbose")
		if isVerbose {
			log.SetLevel(log.DebugLevel)
		}
	},
}

func init() {
	guayavitaCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	guayavitaCmd.AddCommand(help.Commands()...)
	guayavitaCmd.AddCommand(compiler.Commands()...)
	guayavitaCmd.SilenceErrors = true

	log.SetDefault(log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false, // show timestamps
		ReportCaller:    false, // show file + line numbers
	}))
}

func Execute() error {
	return guayavitaCmd.Execute()
}
