package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// FileWriter implements the FileWriter interface
type FileWriter struct{}

// NewFileWriter creates a new FileWriter
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// WriteCollection writes an HTTP collection to the file system
func (w *FileWriter) WriteCollection(ctx context.Context, collection *models.HTTPCollection) error {
	if collection == nil {
		return fmt.Errorf("collection is nil")
	}

	// Create root directory if it doesn't exist
	if collection.RootDir != "" {
		if err := os.MkdirAll(collection.RootDir, 0755); err != nil {
			return fmt.Errorf("failed to create root directory %s: %w", collection.RootDir, err)
		}
	}

	// Write root files
	for _, file := range collection.RootFiles {
		if err := w.WriteFile(ctx, &file, collection.RootDir); err != nil {
			return err
		}
	}

	// Write directories and their files
	for _, dir := range collection.Directories {
		dirPath := filepath.Join(collection.RootDir, dir.Path)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}

		for _, file := range dir.Files {
			if err := w.WriteFile(ctx, &file, dirPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteFile writes an HTTP file to the file system
func (w *FileWriter) WriteFile(ctx context.Context, file *models.HTTPFile, dirPath string) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	filePath := filepath.Join(dirPath, file.Filename)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer f.Close()

	for i, request := range file.Requests {
		if i > 0 {
			// Add a separator between requests
			if _, err := f.WriteString("\n###\n\n"); err != nil {
				return fmt.Errorf("failed to write separator to file %s: %w", filePath, err)
			}
		}

		if err := w.writeRequest(f, &request); err != nil {
			return fmt.Errorf("failed to write request to file %s: %w", filePath, err)
		}
	}

	return nil
}

// writeRequest writes an HTTP request to a file
func (w *FileWriter) writeRequest(f *os.File, request *models.HTTPRequest) error {
	// Write comments
	for _, comment := range request.Comments {
		lines := strings.Split(comment, "\n")
		for _, line := range lines {
			if _, err := f.WriteString(fmt.Sprintf("# %s\n", line)); err != nil {
				return err
			}
		}
	}

	// Write request name as a comment if available
	if request.Name != "" {
		if _, err := f.WriteString(fmt.Sprintf("# @name %s\n", request.Name)); err != nil {
			return err
		}
	}

	// Write request line
	if _, err := f.WriteString(fmt.Sprintf("%s %s\n", request.Method, request.URL)); err != nil {
		return err
	}

	// Write headers
	for name, value := range request.Headers {
		if _, err := f.WriteString(fmt.Sprintf("%s: %s\n", name, value)); err != nil {
			return err
		}
	}

	// Write body if available
	if request.Body != "" {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
		if _, err := f.WriteString(request.Body); err != nil {
			return err
		}
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}
