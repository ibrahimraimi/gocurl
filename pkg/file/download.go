package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Download downloads content from a reader to a file
func Download(content io.Reader, filePath string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
