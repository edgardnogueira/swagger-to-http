package fs

import (
	"context"
	"os"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// FileWriter defines the interface for writing files
type FileWriter interface {
	// WriteCollection writes a collection of HTTP files to disk
	WriteCollection(ctx context.Context, collection *models.HTTPCollection) error
	
	// WriteFile writes a single HTTP file to disk
	WriteFile(ctx context.Context, file *models.HTTPFile, dirPath string) error
	
	// MkdirAll creates a directory with all necessary parent directories
	MkdirAll(path string, perm os.FileMode) error
	
	// WriteSnapshotFile writes data to a file for snapshot functionality
	WriteSnapshotFile(path string, data []byte, perm os.FileMode) error
	
	// ReadFile reads the content of a file
	ReadFile(path string) ([]byte, error)
	
	// FileExists checks if a file exists
	FileExists(path string) (bool, error)
	
	// Glob returns the names of all files matching a pattern
	Glob(pattern string) ([]string, error)
	
	// Remove removes a file or empty directory
	Remove(path string) error
}
