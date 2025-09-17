package compiler

import (
	"fmt"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"jmpeax.com/guayavita/gvc/internal/codegen"
	"jmpeax.com/guayavita/gvc/internal/fs"
	"jmpeax.com/guayavita/gvc/internal/syntax"
	"tinygo.org/x/go-llvm"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "c",
	Long:  "Compile a guayavita file",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// Check if benchmark mode is enabled
		benchmark, _ := cmd.Flags().GetBool("benchmark")

		// Start benchmark measurement
		var startTime time.Time
		var memStart runtime.MemStats
		if benchmark {
			startTime = time.Now()
			runtime.GC() // Force garbage collection for accurate memory measurement
			runtime.ReadMemStats(&memStart)
		}

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
		// Parse the source
		parsedFile, diagnostics := syntax.ParseFile(file, content)

		// Print diagnostics if any
		if len(diagnostics) > 0 {
			log.Error("Parsing errors found:")
			for _, diag := range diagnostics {
				log.Error(diag.Render(content))
			}
		} else {
			log.Info("Parsing completed successfully")
		}

		// Pretty-print the AST
		log.Debugf("AST:\n%s", syntax.PrintFile(parsedFile))

		// Get compilation flags
		syntaxOnly, _ := cmd.Flags().GetBool("syntax-only")
		if !syntaxOnly {
			target, _ := cmd.Flags().GetString("target")
			jit, _ := cmd.Flags().GetBool("jit")
			emitLLVM, _ := cmd.Flags().GetBool("emit-llvm")
			outputDir, _ := cmd.Flags().GetString("output-dir")

			// Create and configure the code builder
			builder := codegen.NewCodeBuilder()
			builder.SetInputFile(file).SetOutputDir(outputDir).SetSource(content)

			if target != "" {
				builder.SetTarget(target)
			} else {
				builder.SetDefaultTarget()
			}

			// Determine compilation mode
			if jit {
				builder.SetMode(codegen.ModeJIT)
				log.Info("Using JIT execution mode")
			} else if emitLLVM {
				builder.SetMode(codegen.ModeEmitLLVM)
				log.Info("Emitting LLVM IR")
			} else {
				builder.SetMode(codegen.ModeBinary)
				log.Info("Compiling to binary")
			}

			// Build the code
			if err := builder.Build(parsedFile); err != nil {
				log.Errorf("Code generation failed: %s", err)
				return
			}

			log.Info("Code generation completed successfully")
		} else {
			log.Info("Syntax-only mode: skipping code generation")
		}

		// Output benchmark results if enabled
		if benchmark {
			elapsed := time.Since(startTime)
			var memEnd runtime.MemStats
			runtime.ReadMemStats(&memEnd)

			memUsed := memEnd.TotalAlloc - memStart.TotalAlloc

			// Define lipgloss styles for benchmark output
			titleStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00D7FF")).
				Border(lipgloss.DoubleBorder()).
				Padding(0, 1)

			labelStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFF00"))

			valueStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF87"))

			fileStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF79C6"))

			fmt.Println()
			fmt.Println(titleStyle.Render("Compilation Benchmark Results"))
			fmt.Printf("%s %s\n", labelStyle.Render("File:"), fileStyle.Render(file))
			fmt.Printf("%s %s\n", labelStyle.Render("Compilation Time:"), valueStyle.Render(elapsed.String()))
			fmt.Printf("%s %s\n", labelStyle.Render("Memory Used:"), valueStyle.Render(fmt.Sprintf("%d bytes (%.2f KB)", memUsed, float64(memUsed)/1024.0)))
			fmt.Printf("%s %s\n", labelStyle.Render("Total Allocations:"), valueStyle.Render(fmt.Sprintf("%d", memEnd.Mallocs-memStart.Mallocs)))
			fmt.Printf("%s %s\n", labelStyle.Render("Number of GC Cycles:"), valueStyle.Render(fmt.Sprintf("%d", memEnd.NumGC-memStart.NumGC)))
		}

	},
}
var targets = &cobra.Command{
	Use:  "targets",
	Long: "List supported targets",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Supported targets:")
		log.Infof("Default target: %s", llvm.DefaultTargetTriple())
		supportedTargets := codegen.SupportedTriples()
		for _, target := range supportedTargets {
			log.Info(target)
		}
	},
}

func init() {
	compileCmd.Flags().Bool("syntax-only", false, "Syntax only, do not output binaries")
	compileCmd.Flags().Bool("benchmark", false, "Benchmark mode, print benchmark results of compilation")
	compileCmd.Flags().StringP("target", "t", "", "LLVM target triple (e.g., x86_64-pc-linux-gnu)")
	compileCmd.Flags().Bool("jit", false, "Use JIT execution mode instead of compilation")
	compileCmd.Flags().Bool("emit-llvm", false, "Output LLVM IR (.ll) file instead of executable binary")
	compileCmd.Flags().StringP("output-dir", "o", "./bin", "Output directory for generated files")
}
func Commands() []*cobra.Command {
	return []*cobra.Command{compileCmd, targets}
}
