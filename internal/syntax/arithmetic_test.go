package syntax

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_ArithmeticDefinition(t *testing.T) {
	file := filepath.Join("..", "..", "test-data", "simple.gvt")
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("fixture read: %v", err)
	}
	lx := NewLexerFromString(string(data))
	ps := NewParser(lx)
	f, err := ps.ParseFile()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if f.Package == nil || f.Package.Name != "test" {
		t.Fatalf("expected package 'test'")
	}
	if len(f.Definitions) != 2 {
		t.Fatalf("expected 2 definitions (PI and a), got %d", len(f.Definitions))
	}
}
