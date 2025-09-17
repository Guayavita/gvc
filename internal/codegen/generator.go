package codegen

import (
	"fmt"

	"jmpeax.com/guayavita/gvc/internal/diag"
	"jmpeax.com/guayavita/gvc/internal/syntax"
	"tinygo.org/x/go-llvm"
)

// generateIR generates LLVM IR from the AST
func (b *LLVMCodeBuilder) generateIR(ast *syntax.File) error {
	// Create main function as entry point
	mainType := llvm.FunctionType(b.context.Int32Type(), []llvm.Type{}, false)
	mainFunc := llvm.AddFunction(b.module, "main", mainType)

	entry := b.context.AddBasicBlock(mainFunc, "entry")
	b.builder.SetInsertPoint(entry, entry.FirstInstruction())

	// Generate code for each declaration
	for _, decl := range ast.Decls {

		if err := b.generateDecl(decl); err != nil {
			return err
		}
	}

	// Return 0 from main
	b.builder.CreateRet(llvm.ConstInt(b.context.Int32Type(), 0, false))

	// Verify the module
	if err := llvm.VerifyModule(b.module, llvm.ReturnStatusAction); err != nil {
		// attach a diagnostic at file start as we don't have a specific node here
		b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, fmt.Sprintf("module verification failed: %v", err))
		return fmt.Errorf("module verification failed: %w", err)
	}

	return nil
}

func SupportedTriples() []string {
	return []string{
		// Apple macOS + iOS
		"x86_64-apple-darwin",
		"aarch64-apple-darwin",
		"x86_64-apple-ios",
		"aarch64-apple-ios",

		// Linux
		"armv7-unknown-linux-gnueabihf",
		"aarch64-unknown-linux-gnu",
		"riscv64-unknown-linux-gnu",
		"x86_64-unknown-linux-gnu",

		// Windows
		"aarch64-pc-windows-msvc",
		"aarch64-pc-windows-gnu",
		"riscv64-pc-windows-gnu", // experimental
		"x86_64-pc-windows-msvc",
		"x86_64-pc-windows-gnu",

		// WebAssembly
		"wasm32-unknown-unknown",
		"wasm64-unknown-unknown", // experimental
	}
}
