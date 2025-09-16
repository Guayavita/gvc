package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"tinygo.org/x/go-llvm"
)

// getOutputName generates the output filename based on input file
func (b *LLVMCodeBuilder) getOutputName() string {
	if b.config.InputFile == "" {
		return filepath.Join(b.config.OutputDir, "output")
	}

	base := filepath.Base(b.config.InputFile)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(b.config.OutputDir, name)
}

// getTargetTriple returns the target triple to use
func (b *LLVMCodeBuilder) getTargetTriple() string {
	if b.config.Target != "" {
		return b.config.Target
	}
	return llvm.DefaultTargetTriple()
}

// initializeLLVM sets up the LLVM context, module, and builder
func (b *LLVMCodeBuilder) initializeLLVM() error {
	b.context = llvm.NewContext()

	moduleName := "guayavita_module"
	if b.config.InputFile != "" {
		base := filepath.Base(b.config.InputFile)
		moduleName = strings.TrimSuffix(base, filepath.Ext(base))
	}

	b.module = b.context.NewModule(moduleName)
	b.builder = b.context.NewBuilder()

	// Set target triple if specified
	if b.config.Target != "" {
		b.module.SetTarget(b.config.Target)
	}

	// Initialize external functions
	b.initializeExternalFunctions()

	// Create print function
	if err := b.createPrintFunction(); err != nil {
		return fmt.Errorf("failed to create print function: %w", err)
	}

	return nil
}

// cleanup releases LLVM resources
func (b *LLVMCodeBuilder) cleanup() {
	if !b.builder.IsNil() {
		b.builder.Dispose()
	}
	if !b.module.IsNil() {
		b.module.Dispose()
	}
	if !b.context.IsNil() {
		b.context.Dispose()
	}
}
