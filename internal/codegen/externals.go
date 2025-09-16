package codegen

import (
	"fmt"
	"tinygo.org/x/go-llvm"
)

// ExternalFunction represents an external function declaration
type ExternalFunction struct {
	Name       string
	ReturnType llvm.Type
	ParamTypes []llvm.Type
	Variadic   bool
}

// ExternalRegistry manages external function declarations
type ExternalRegistry struct {
	functions map[string]ExternalFunction
}

// NewExternalRegistry creates a new external function registry
func NewExternalRegistry() *ExternalRegistry {
	return &ExternalRegistry{
		functions: make(map[string]ExternalFunction),
	}
}

// RegisterFunction registers an external function
func (r *ExternalRegistry) RegisterFunction(name string, returnType llvm.Type, paramTypes []llvm.Type, variadic bool) {
	r.functions[name] = ExternalFunction{
		Name:       name,
		ReturnType: returnType,
		ParamTypes: paramTypes,
		Variadic:   variadic,
	}
}

// GetFunction returns an external function declaration
func (r *ExternalRegistry) GetFunction(name string) (ExternalFunction, bool) {
	fn, exists := r.functions[name]
	return fn, exists
}

// initializeExternalFunctions sets up external function declarations for the module
func (b *LLVMCodeBuilder) initializeExternalFunctions() {
	if b.externals == nil {
		b.externals = NewExternalRegistry()
	}

	// Register printf from libc for print functionality
	i8PtrType := llvm.PointerType(b.context.Int8Type(), 0)
	b.externals.RegisterFunction("printf", b.context.Int32Type(), []llvm.Type{i8PtrType}, true)

	// Register puts from libc as alternative
	b.externals.RegisterFunction("puts", b.context.Int32Type(), []llvm.Type{i8PtrType}, false)
}

// declareExternalFunction declares an external function in the LLVM module
func (b *LLVMCodeBuilder) declareExternalFunction(name string) (llvm.Value, error) {
	if b.externals == nil {
		return llvm.Value{}, fmt.Errorf("external registry not initialized")
	}

	extFunc, exists := b.externals.GetFunction(name)
	if !exists {
		return llvm.Value{}, fmt.Errorf("external function '%s' not registered", name)
	}

	// Check if function is already declared
	existing := b.module.NamedFunction(name)
	if !existing.IsNil() {
		return existing, nil
	}

	// Create function type and declare it
	funcType := llvm.FunctionType(extFunc.ReturnType, extFunc.ParamTypes, extFunc.Variadic)
	function := llvm.AddFunction(b.module, name, funcType)

	return function, nil
}

// createPrintFunction creates a print function that wraps printf/puts
func (b *LLVMCodeBuilder) createPrintFunction() error {
	// Check if print function already exists
	existing := b.module.NamedFunction("print")
	if !existing.IsNil() {
		return nil
	}

	// Create print function type (takes string, returns void)
	i8PtrType := llvm.PointerType(b.context.Int8Type(), 0)
	printType := llvm.FunctionType(b.context.VoidType(), []llvm.Type{i8PtrType}, false)
	printFunc := llvm.AddFunction(b.module, "print", printType)

	// Create entry block
	entry := b.context.AddBasicBlock(printFunc, "entry")

	// Save current insert point
	currentBlock := b.builder.GetInsertBlock()

	// Set insert point to print function
	b.builder.SetInsertPoint(entry, entry.FirstInstruction())

	// Declare puts function
	putsFunc, err := b.declareExternalFunction("puts")
	if err != nil {
		return fmt.Errorf("failed to declare puts function: %w", err)
	}

	// Get the string parameter
	param := printFunc.Param(0)

	// Call puts with the string parameter
	b.builder.CreateCall(putsFunc.GlobalValueType(), putsFunc, []llvm.Value{param}, "")

	// Return void
	b.builder.CreateRetVoid()

	// Restore previous insert point
	if !currentBlock.IsNil() {
		b.builder.SetInsertPoint(currentBlock, currentBlock.FirstInstruction())
	}

	return nil
}
