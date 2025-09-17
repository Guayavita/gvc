package syntax

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrinter_BasicOutput(t *testing.T) {
	path := repoPathSyntax(filepath.Join("test-data", "hello.gvt"))
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	file, diags := ParseFile(path, string(src))
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %d: %#v", len(diags), diags)
	}
	out := PrintFile(file)
	if out == "" {
		t.Fatalf("expected printed output, got empty string")
	}
	// We expect to see some of the function names in the printed AST output
	if !strings.Contains(out, "Name: add") {
		t.Fatalf("printed output does not include function 'add':\n%s", out)
	}
	if !strings.Contains(out, "Name: main") {
		t.Fatalf("printed output does not include function 'main':\n%s", out)
	}
}
