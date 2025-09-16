package codegen

import (
	"fmt"

	"jmpeax.com/guayavita/gvc/internal/syntax"
)

// generateBlock generates LLVM IR for a block statement
func (b *LLVMCodeBuilder) generateBlock(block *syntax.Block) error {
	for _, stmt := range block.Stmts {
		if err := b.generateStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

// generateStmt generates LLVM IR for a statement
func (b *LLVMCodeBuilder) generateStmt(stmt syntax.Stmt) error {
	switch s := stmt.(type) {
	case *syntax.ExprStmt:
		_, err := b.generateExpr(s.X)
		return err
	case *syntax.VarDecl:
		return b.generateVarDecl(s)
	case *syntax.ReturnStmt:
		if s.Result != nil {
			value, err := b.generateExpr(s.Result)
			if err != nil {
				return err
			}
			b.builder.CreateRet(value)
		} else {
			b.builder.CreateRetVoid()
		}
		return nil
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}
