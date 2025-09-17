package syntax

import (
	"os"
	"path/filepath"
	"testing"
)

func repoPathSyntax(rel string) string {
	return filepath.Join("..", "..", rel)
}

func TestParser_ParseSimple(t *testing.T) {
	path := repoPathSyntax(filepath.Join("test-data", "simple.gvt"))
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	file, diags := ParseFile(path, string(src))
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d: %#v", len(diags), diags)
	}
	if file == nil {
		t.Fatalf("expected AST file, got nil")
	}
	if file.Package != "test" {
		t.Fatalf("expected package 'test', got %q", file.Package)
	}
	if len(file.Decls) < 1 {
		t.Fatalf("expected at least 1 declaration, got %d", len(file.Decls))
	}
}

func TestParser_ParseHello(t *testing.T) {
	path := repoPathSyntax(filepath.Join("test-data", "hello.gvt"))
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	file, diags := ParseFile(path, string(src))
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d: %#v", len(diags), diags)
	}
	if want := "main"; file.Package != want {
		t.Fatalf("expected package %q, got %q", want, file.Package)
	}
	// Expect at least 2 functions: add and main
	funCount := 0
	for _, d := range file.Decls {
		if _, ok := d.(*FunDecl); ok {
			funCount++
		}
	}
	if funCount < 2 {
		t.Fatalf("expected at least 2 functions, got %d", funCount)
	}
}
