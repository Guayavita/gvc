package syntax

import (
	"os"
	"path/filepath"
	"testing"

	"jmpeax.com/guayavita/gvc/internal/diag"
)

func TestParser_ReportsErrors(t *testing.T) {
	path := repoPathSyntax(filepath.Join("test-data", "error-test.gvt"))
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fixture missing: %v", err)
	}
	_, diags := ParseFile(path, string(src))
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics for invalid file, got 0")
	}
	// ensure at least one diagnostic is an error severity
	var hasError bool
	for _, d := range diags {
		if d.Severity == diag.Error {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Fatalf("expected at least one error diagnostic, got: %#v", diags)
	}
}
