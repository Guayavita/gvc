package codegen

import (
	"jmpeax.com/guayavita/gvc/internal/syntax"
)

// generateDecl generates LLVM IR for a declaration
func (b *LLVMCodeBuilder) generateDecl(decl syntax.Decl) error {
	switch d := decl.(type) {
	case *syntax.FunDecl:
		return b.generateFunctionDecl(d)
	case *syntax.VarDecl:
		return b.generateVarDecl(d)
	default:
		return b.errorAt(d, "unsupported declaration type: %T", decl)
	}
}

// generateFunctionDecl generates LLVM IR for a function declaration
func (b *LLVMCodeBuilder) generateFunctionDecl(decl *syntax.FunDecl) error {
	// For now, skip non-main functions - this is a basic implementation
	if decl.Name != "main" {
		return nil
	}

	// Generate statements in function body
	if decl.Body != nil {
		return b.generateBlock(decl.Body)
	}

	return nil
}

// generateVarDecl generates LLVM IR for a variable declaration
func (b *LLVMCodeBuilder) generateVarDecl(decl *syntax.VarDecl) error {
	// Basic variable declaration - allocate space and store initial value
	varType := b.context.Int32Type() // Default to int32 for now
	alloca := b.builder.CreateAlloca(varType, decl.Name)

	if decl.Init != nil {
		value, err := b.generateExpr(decl.Init)
		if err != nil {
			return err
		}
		b.builder.CreateStore(value, alloca)
	}

	return nil
}
