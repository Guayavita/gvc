package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidateFile checks if the given path points to a readable file
func ValidateFile(path string) error {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("invalid path or broken symlink: %w", err)
	}
	fileInfo, err := os.Stat(resolvedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file: %s", path)
	}

	return nil
}
