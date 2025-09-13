package fs

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to repo root path from this test file location
func repoPath(rel string) string {
	// tests run from package dir; use relative path to project root
	return filepath.Join("..", "..", rel)
}

func TestValidateFile_NonExistent(t *testing.T) {
	err := ValidateFile("/path/that/does/not/exist.gvc")
	if err == nil {
		t.Fatalf("expected error for non-existent file")
	}
}

func TestValidateFile_IsDirectory(t *testing.T) {
	// Using the repository root, which is a directory
	dir := repoPath("")
	// Clean the path
	dir, _ = filepath.Abs(dir)
	err := ValidateFile(dir)
	if err == nil {
		t.Fatalf("expected error for directory path, got nil")
	}
}

func TestValidateFile_OK(t *testing.T) {
	file := repoPath(filepath.Join("test-data", "simple.gvt"))
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("test fixture missing: %v", err)
	}
	if err := ValidateFile(file); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReadFile(t *testing.T) {
	file := repoPath(filepath.Join("test-data", "simple.gvt"))
	content, err := ReadFile(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("expected some content in simple.gvt")
	}
	// Should start with the word 'package'
	wantPrefix := "package "
	if len(content) < len(wantPrefix) || content[:len(wantPrefix)] != wantPrefix {
		t.Fatalf("expected content to start with %q, got: %q", wantPrefix, content)
	}
}
