package fs

import (
	"context"
	"os"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// FileSystemAdapter adapts the FileSystem type to implement the FileWriter interface
type FileSystemAdapter struct {
	fileWriter *FileWriter
	fileSystem *FileSystem
}

// NewFileSystemAdapter creates a new FileSystemAdapter
func NewFileSystemAdapter(fw *FileWriter, fs *FileSystem) *FileSystemAdapter {
	return &FileSystemAdapter{
		fileWriter: fw,
		fileSystem: fs,
	}
}

// WriteCollection delegates to the underlying FileWriter
func (a *FileSystemAdapter) WriteCollection(ctx context.Context, collection *models.HTTPCollection) error {
	return (*a.fileWriter).WriteCollection(ctx, collection)
}

// WriteFile delegates to the underlying FileWriter
func (a *FileSystemAdapter) WriteFile(ctx context.Context, file *models.HTTPFile, dirPath string) error {
	return (*a.fileWriter).WriteFile(ctx, file, dirPath)
}

// MkdirAll delegates to the FileSystem implementation
func (a *FileSystemAdapter) MkdirAll(path string, perm os.FileMode) error {
	return a.fileSystem.MkdirAll(path, perm)
}

// WriteSnapshotFile delegates to the FileSystem implementation
func (a *FileSystemAdapter) WriteSnapshotFile(path string, data []byte, perm os.FileMode) error {
	return a.fileSystem.WriteFile(path, data, perm)
}

// ReadFile delegates to the FileSystem implementation
func (a *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	return a.fileSystem.ReadFile(path)
}

// FileExists delegates to the FileSystem implementation
func (a *FileSystemAdapter) FileExists(path string) (bool, error) {
	return a.fileSystem.FileExists(path)
}

// Glob delegates to the FileSystem implementation
func (a *FileSystemAdapter) Glob(pattern string) ([]string, error) {
	return a.fileSystem.Glob(pattern)
}

// Remove delegates to the FileSystem implementation
func (a *FileSystemAdapter) Remove(path string) error {
	return a.fileSystem.Remove(path)
}
