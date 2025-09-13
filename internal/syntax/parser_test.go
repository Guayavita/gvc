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
