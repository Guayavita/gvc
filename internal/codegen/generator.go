package codegen

import (
	"fmt"

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
		return fmt.Errorf("module verification failed: %w", err)
	}

	return nil
}
