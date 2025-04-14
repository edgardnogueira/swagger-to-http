package fs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileSystem implements file system operations
type FileSystem struct{}

// NewFileSystem creates a new FileSystem instance
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// MkdirAll creates a directory with all necessary parent directories
func (fs *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// WriteFile writes data to a file
func (fs *FileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	return ioutil.WriteFile(path, data, perm)
}

// ReadFile reads the content of a file
func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// FileExists checks if a file exists
func (fs *FileSystem) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Glob returns the names of all files matching a pattern
func (fs *FileSystem) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// Remove removes a file or empty directory
func (fs *FileSystem) Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll removes a file or directory and any children it contains
func (fs *FileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
