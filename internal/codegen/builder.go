package codegen

import (
	"fmt"

	"jmpeax.com/guayavita/gvc/internal/syntax"
	"tinygo.org/x/go-llvm"
)

// CompilationMode defines the compilation target mode
type CompilationMode int

const (
	ModeBinary   CompilationMode = iota // Compile to binary executable
	ModeJIT                             // Execute using JIT
	ModeEmitLLVM                        // Output LLVM IR file
)

// BuilderConfig holds configuration for the code builder
type BuilderConfig struct {
	Mode      CompilationMode
	Target    string // LLVM target triple
	OutputDir string
	InputFile string
}

// CodeBuilder interface defines the builder pattern for code generation
type CodeBuilder interface {
	SetTarget(target string) CodeBuilder
	SetMode(mode CompilationMode) CodeBuilder
	SetOutputDir(dir string) CodeBuilder
	SetInputFile(file string) CodeBuilder
	Build(ast *syntax.File) error
}

// LLVMCodeBuilder implements CodeBuilder using LLVM
type LLVMCodeBuilder struct {
	config          BuilderConfig
	context         llvm.Context
	module          llvm.Module
	builder         llvm.Builder
	externals       *ExternalRegistry
	printedLiterals []string // Track string literals for wrapper script
}

// NewCodeBuilder creates a new LLVM-based code builder
func NewCodeBuilder() CodeBuilder {
	return &LLVMCodeBuilder{}
}

func (b *LLVMCodeBuilder) SetTarget(target string) CodeBuilder {
	b.config.Target = target
	return b
}

func (b *LLVMCodeBuilder) SetMode(mode CompilationMode) CodeBuilder {
	b.config.Mode = mode
	return b
}

func (b *LLVMCodeBuilder) SetOutputDir(dir string) CodeBuilder {
	b.config.OutputDir = dir
	return b
}

func (b *LLVMCodeBuilder) SetInputFile(file string) CodeBuilder {
	b.config.InputFile = file
	return b
}

// Build compiles the AST according to the builder configuration
func (b *LLVMCodeBuilder) Build(ast *syntax.File) error {
	// Initialize LLVM components
	if err := b.initializeLLVM(); err != nil {
		return fmt.Errorf("failed to initialize LLVM: %w", err)
	}
	defer b.cleanup()

	// Generate LLVM IR from AST
	if err := b.generateIR(ast); err != nil {
		return fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Process according to mode
	switch b.config.Mode {
	case ModeBinary:
		return b.compileToBinary()
	case ModeJIT:
		return b.executeJIT()
	case ModeEmitLLVM:
		return b.emitLLVM()
	default:
		return fmt.Errorf("unsupported compilation mode: %d", b.config.Mode)
	}
}
