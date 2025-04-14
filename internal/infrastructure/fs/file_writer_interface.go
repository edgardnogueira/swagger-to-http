package fs

import (
	"os"
)

// SnapshotFileWriter defines the interface for file operations used by snapshot functionality
type SnapshotFileWriter interface {
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
