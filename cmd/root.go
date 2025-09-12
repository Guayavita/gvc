package cmd

import (
	"github.com/spf13/cobra"
	"jmpeax.com/guayavita/gvc/cmd/help"
)

var guayavitaCmd = &cobra.Command{
	Use:   "guayavita",
	Short: "Guayavita Compiler and Tooling",
}

func init() {
	guayavitaCmd.AddCommand(help.Commands()...)
}

func Execute() error {
	return guayavitaCmd.Execute()
}
