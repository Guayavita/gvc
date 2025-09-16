package codegen

import (
	"fmt"

	"jmpeax.com/guayavita/gvc/internal/syntax"
	"tinygo.org/x/go-llvm"
)

// generateExpr generates LLVM IR for an expression
func (b *LLVMCodeBuilder) generateExpr(expr syntax.Expr) (llvm.Value, error) {
	switch e := expr.(type) {
	case *syntax.BasicLit:
		return b.generateBasicLit(e)
	case *syntax.BinaryExpr:
		return b.generateBinaryExpr(e)
	case *syntax.CallExpr:
		return b.generateCallExpr(e)
	default:
		return llvm.Value{}, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// generateBasicLit generates LLVM IR for a basic literal
func (b *LLVMCodeBuilder) generateBasicLit(lit *syntax.BasicLit) (llvm.Value, error) {
	switch lit.Kind {
	case "INT":
		// Parse integer value and create constant
		return llvm.ConstInt(b.context.Int32Type(), 42, false), nil // Simplified for now
	case "STRING":
		// Create string constant
		str := b.builder.CreateGlobalStringPtr(lit.Value, "str")
		return str, nil
	default:
		return llvm.Value{}, fmt.Errorf("unsupported literal kind: %s", lit.Kind)
	}
}

// generateBinaryExpr generates LLVM IR for a binary expression
func (b *LLVMCodeBuilder) generateBinaryExpr(expr *syntax.BinaryExpr) (llvm.Value, error) {
	left, err := b.generateExpr(expr.Left)
	if err != nil {
		return llvm.Value{}, err
	}

	right, err := b.generateExpr(expr.Right)
	if err != nil {
		return llvm.Value{}, err
	}

	switch expr.Op {
	case "+":
		return b.builder.CreateAdd(left, right, "add"), nil
	case "-":
		return b.builder.CreateSub(left, right, "sub"), nil
	case "*":
		return b.builder.CreateMul(left, right, "mul"), nil
	case "/":
		return b.builder.CreateSDiv(left, right, "div"), nil
	default:
		return llvm.Value{}, fmt.Errorf("unsupported binary operator: %s", expr.Op)
	}
}

// generateCallExpr generates LLVM IR for a function call
func (b *LLVMCodeBuilder) generateCallExpr(expr *syntax.CallExpr) (llvm.Value, error) {
	// Get function name from identifier
	ident, ok := expr.Fun.(*syntax.Ident)
	if !ok {
		return llvm.Value{}, fmt.Errorf("unsupported function call expression: %T", expr.Fun)
	}

	funcName := ident.Name

	// Handle print function specially
	if funcName == "print" {
		return b.generatePrintCall(expr)
	}

	// For other functions, try to find them in the module
	function := b.module.NamedFunction(funcName)
	if function.IsNil() {
		return llvm.Value{}, fmt.Errorf("undefined function: %s", funcName)
	}

	// Generate arguments
	var args []llvm.Value
	for _, arg := range expr.Args {
		argValue, err := b.generateExpr(arg)
		if err != nil {
			return llvm.Value{}, fmt.Errorf("failed to generate argument: %w", err)
		}
		args = append(args, argValue)
	}

	// Call the function
	return b.builder.CreateCall(function.GlobalValueType(), function, args, "call"), nil
}

// generatePrintCall generates LLVM IR for a print function call
func (b *LLVMCodeBuilder) generatePrintCall(expr *syntax.CallExpr) (llvm.Value, error) {
	if len(expr.Args) != 1 {
		return llvm.Value{}, fmt.Errorf("print function expects exactly 1 argument, got %d", len(expr.Args))
	}

	// Generate the string argument
	arg, err := b.generateExpr(expr.Args[0])
	if err != nil {
		return llvm.Value{}, fmt.Errorf("failed to generate print argument: %w", err)
	}

	// Get the print function
	printFunc := b.module.NamedFunction("print")
	if printFunc.IsNil() {
		return llvm.Value{}, fmt.Errorf("print function not found - external functions not initialized")
	}

	// Call print function
	b.builder.CreateCall(printFunc.GlobalValueType(), printFunc, []llvm.Value{arg}, "")

	// Print returns void, but we need to return something for expressions
	return llvm.ConstInt(b.context.Int32Type(), 0, false), nil
}
