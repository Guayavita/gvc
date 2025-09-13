package compiler

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"jmpeax.com/guayavita/gvc/internal/fs"
	"jmpeax.com/guayavita/gvc/internal/syntax"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "c",
	Long:  "Compile a guayavita file",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		log.Debugf("Compiling %s", file)
		err := fs.ValidateFile(file)
		if err != nil {
			log.Error(err)
		}
		content, err := fs.ReadFile(file)
		if err != nil {
			log.Errorf("Error reading file: %s", err)
			return
		}
		lx := syntax.NewLexerFromString(content)
		ps := syntax.NewParser(lx)
		gvcAst, err := ps.ParseFile()
		if err != nil {
			var pe *syntax.ParseError
			if errors.As(err, &pe) {
				d := pe.Diagnostic(file)
				fmt.Println(d.Render(content))
				return
			}
			log.Errorf("parse failed: %v", err)
			return
		}
		log.Debugf("Package name: %s , package Variables %v", gvcAst.Package.Name, gvcAst.Definitions)
	},
}

func init() {

}
func Commands() []*cobra.Command {
	return []*cobra.Command{compileCmd}
}
