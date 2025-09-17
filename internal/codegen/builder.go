package codegen

import (
	"fmt"
	"runtime"

	"jmpeax.com/guayavita/gvc/internal/diag"
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
	Source    string // original source content for diagnostics rendering
}

// CodeBuilder interface defines the builder pattern for code generation
type CodeBuilder interface {
	SetTarget(target string) CodeBuilder
	SetMode(mode CompilationMode) CodeBuilder
	SetOutputDir(dir string) CodeBuilder
	SetInputFile(file string) CodeBuilder
	SetSource(src string) CodeBuilder
	Build(ast *syntax.File) error
	Diagnostics() []diag.Diagnostic
	SetDefaultTarget()
}

// LLVMCodeBuilder implements CodeBuilder using LLVM
type LLVMCodeBuilder struct {
	config          BuilderConfig
	context         llvm.Context
	module          llvm.Module
	builder         llvm.Builder
	externals       *ExternalRegistry
	printedLiterals []string // Track string literals for wrapper script
	diagnostics     []diag.Diagnostic
}

// NewCodeBuilder creates a new LLVM-based code builder
func NewCodeBuilder() CodeBuilder {
	return &LLVMCodeBuilder{}
}

func (b *LLVMCodeBuilder) SetDefaultTarget() {
	b.config.Target = generateDefaultLLVMTriple(runtime.GOOS, runtime.GOARCH)
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

// SetSource sets the original source content for diagnostics rendering
func (b *LLVMCodeBuilder) SetSource(src string) CodeBuilder {
	b.config.Source = src
	return b
}

// Diagnostics returns collected diagnostics
func (b *LLVMCodeBuilder) Diagnostics() []diag.Diagnostic {
	return b.diagnostics
}

// Build compiles the AST according to the builder configuration
func (b *LLVMCodeBuilder) Build(ast *syntax.File) error {
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
	// Initialize LLVM components
	if err := b.initializeLLVM(); err != nil {
		// No position info here; record a generic diagnostic
		b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, fmt.Sprintf("failed to initialize LLVM: %v", err))
		return fmt.Errorf("failed to initialize LLVM: %w", err)
	}
	defer b.cleanup()

	// Generate LLVM IR from AST
	if err := b.generateIR(ast); err != nil {
		// Attach to file position if available
		pos := diag.Position{File: b.config.InputFile, Line: 1, Column: 1}
		if ast != nil {
			pos = ast.Pos()
		}
		b.addDiagnostic(diag.Error, pos, fmt.Sprintf("failed to generate LLVM IR: %v", err))
		return err
	}

	// Process according to mode
	switch b.config.Mode {
	case ModeBinary:
		if err := b.compileToBinary(); err != nil {
			b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, fmt.Sprintf("compilation to binary failed: %v", err))
			return err
		}
		return nil
	case ModeJIT:
		if err := b.executeJIT(); err != nil {
			b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, fmt.Sprintf("JIT execution failed: %v", err))
			return err
		}
		return nil
	case ModeEmitLLVM:
		if err := b.emitLLVM(); err != nil {
			b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, fmt.Sprintf("emit LLVM IR failed: %v", err))
			return err
		}
		return nil
	default:
		err := fmt.Errorf("unsupported compilation mode: %d", b.config.Mode)
		b.addDiagnostic(diag.Error, diag.Position{File: b.config.InputFile, Line: 1, Column: 1}, err.Error())
		return err
	}
}

// addDiagnostic appends a diagnostic to the builder
func (b *LLVMCodeBuilder) addDiagnostic(sev diag.Severity, pos diag.Position, msg string) {
	d := diag.Diagnostic{
		Severity: sev,
		Message:  msg,
		Span: diag.Span{
			Start: pos,
			End:   pos,
		},
	}
	b.diagnostics = append(b.diagnostics, d)
}

// errorAt records a diagnostic at the given node position and returns an error
func (b *LLVMCodeBuilder) errorAt(node syntax.Node, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	pos := diag.Position{File: b.config.InputFile, Line: 1, Column: 1}
	if node != nil {
		pos = node.Pos()
	}
	b.addDiagnostic(diag.Error, pos, msg)
	return fmt.Errorf("%s", msg)
}

// GenerateLLVMTriple builds the LLVM triple for a given OS/arch pair.
func generateDefaultLLVMTriple(goos, goarch string) string {
	switch goos {
	case "darwin":
		if goarch == "amd64" {
			return "x86_64-apple-darwin"
		} else if goarch == "arm64" {
			return "arm64-apple-darwin"
		}
	case "ios":
		if goarch == "arm64" {
			return "arm64-apple-ios"
		}
	case "linux":
		switch goarch {
		case "amd64":
			return "x86_64-unknown-linux-gnu"
		case "arm64":
			return "aarch64-unknown-linux-gnu"
		case "riscv64":
			return "riscv64-unknown-linux-gnu"
		}
	case "windows":
		switch goarch {
		case "amd64":
			return "x86_64-pc-windows-msvc"
		case "arm64":
			return "aarch64-pc-windows-msvc"
		case "riscv64":
			return "riscv64-pc-windows-msvc"
		}
	case "wasm":
		if goarch == "wasm32" {
			return "wasm32-wasi"
		}
	}
	return "unknown-unknown-unknown"
}
