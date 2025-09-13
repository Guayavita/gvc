package syntax

import (
	"os"
	"path/filepath"
	"testing"
)

func repoPath(rel string) string { return filepath.Join("..", "..", rel) }

func TestParsePackage_Simple(t *testing.T) {
	file := repoPath(filepath.Join("test-data", "simple.gvc"))
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("fixture read: %v", err)
	}

	lx := NewLexerFromString(string(data))
	ps := NewParser(lx)
	pkg, err := ps.ParsePackage()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if pkg == nil {
		t.Fatalf("nil package result")
	}
	if pkg.Name != "test" {
		t.Fatalf("expected package name 'test', got %q", pkg.Name)
	}
}

func TestParsePackage_Empty_Err(t *testing.T) {
	file := repoPath(filepath.Join("test-data", "empty.gvc"))
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("fixture read: %v", err)
	}

	lx := NewLexerFromString(string(data))
	ps := NewParser(lx)
	_, err = ps.ParsePackage()
	if err == nil {
		t.Fatalf("expected error for empty package decl")
	}
	if _, ok := err.(*ParseError); !ok {
		t.Fatalf("expected *ParseError, got %T", err)
	}
}

func TestParseFile_WithDefinition(t *testing.T) {
	src := "package test\n\n// define A\n def A = B | 42\n"
	lx := NewLexerFromString(src)
	ps := NewParser(lx)
	f, err := ps.ParseFile()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if f.Package == nil || f.Package.Name != "test" {
		t.Fatalf("expected package 'test'")
	}
	if len(f.Definitions) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(f.Definitions))
	}
	def := f.Definitions[0]
	if def.Name != "A" {
		t.Fatalf("expected def name 'A', got %q", def.Name)
	}
	if len(def.Alternatives) != 2 {
		t.Fatalf("expected 2 alternatives, got %d", len(def.Alternatives))
	}
	if _, ok := def.Alternatives[0].(IdentExpr); !ok {
		t.Fatalf("first alternative should be IdentExpr")
	}
	if num, ok := def.Alternatives[1].(NumberExpr); !ok || num.Value != "42" {
		t.Fatalf("second alternative should be NumberExpr '42'")
	}
}

func TestParseFile_WithFloatDefinition(t *testing.T) {
	src := "package test\n\n def PI = 3.16\n"
	lx := NewLexerFromString(src)
	ps := NewParser(lx)
	f, err := ps.ParseFile()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(f.Definitions) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(f.Definitions))
	}
	def := f.Definitions[0]
	if def.Name != "PI" {
		t.Fatalf("expected def name 'PI', got %q", def.Name)
	}
	if len(def.Alternatives) != 1 {
		t.Fatalf("expected 1 alternative, got %d", len(def.Alternatives))
	}
	if num, ok := def.Alternatives[0].(NumberExpr); !ok || num.Value != "3.16" {
		t.Fatalf("alternative should be NumberExpr '3.16'")
	}
}
