package snapshot

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/application/snapshot"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/fs"
)

// SnapshotManager implements the snapshot.Manager interface
type SnapshotManager struct {
	fileWriter fs.FileWriter
}

// NewSnapshotManager creates a new snapshot manager with a file writer
func NewSnapshotManager(fileWriter fs.FileWriter) snapshot.Manager {
	return &SnapshotManager{
		fileWriter: fileWriter,
	}
}

// SaveSnapshot saves a HTTP response as a snapshot file
func (m *SnapshotManager) SaveSnapshot(response *models.HTTPResponse, path string, format string) error {
	// Create formatter for the specified format
	formatter, err := snapshot.GetFormatter(format)
	if err != nil {
		return err
	}

	// Format the response
	formatted, err := formatter.Format(response)
	if err != nil {
		return err
	}

	// Ensure the snapshot directory exists
	dir := filepath.Dir(path)
	if err := m.fileWriter.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Write the snapshot file
	if err := m.fileWriter.WriteFile(path, []byte(formatted), 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// LoadSnapshot loads a snapshot from a file
func (m *SnapshotManager) LoadSnapshot(path string, format string) (*models.HTTPResponse, error) {
	// Check if the snapshot file exists
	exists, err := m.fileWriter.FileExists(path)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("snapshot file not found: %s", path)
	}

	// Read the snapshot file
	data, err := m.fileWriter.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	// Get formatter for the specified format
	formatter, err := snapshot.GetFormatter(format)
	if err != nil {
		return nil, err
	}

	// Parse the snapshot
	return formatter.Parse(string(data))
}

// CompareSnapshots compares a current response with a snapshot
func (m *SnapshotManager) CompareSnapshots(current *models.HTTPResponse, snapshotPath string, format string) (*snapshot.ComparisonResult, error) {
	// Load the snapshot
	expected, err := m.LoadSnapshot(snapshotPath, format)
	if err != nil {
		return nil, err
	}

	// Get formatter for the specified format
	formatter, err := snapshot.GetFormatter(format)
	if err != nil {
		return nil, err
	}

	// Compare the responses
	return formatter.Compare(expected, current)
}

// GetSnapshotPath generates a snapshot path for a HTTP request
func (m *SnapshotManager) GetSnapshotPath(httpFile string, requestName string, baseDir string) string {
	// Generate a safe filename from the request name
	safeRequestName := strings.ReplaceAll(requestName, " ", "_")
	safeRequestName = strings.ReplaceAll(safeRequestName, "/", "_")
	safeRequestName = strings.ReplaceAll(safeRequestName, ":", "_")

	// Get relative path of the HTTP file from the base directory
	relPath, err := filepath.Rel(baseDir, httpFile)
	if err != nil {
		// Fallback to using the HTTP file name if relative path can't be determined
		relPath = filepath.Base(httpFile)
	}

	// Remove extension from the HTTP file
	relPath = strings.TrimSuffix(relPath, filepath.Ext(relPath))

	// Construct the snapshot path
	return filepath.Join(baseDir, "__snapshots__", relPath, safeRequestName+".snap")
}

// ListSnapshots returns a list of all snapshot files
func (m *SnapshotManager) ListSnapshots(snapshotsDir string) ([]string, error) {
	return m.fileWriter.Glob(filepath.Join(snapshotsDir, "**/*.snap"))
}

// CleanupSnapshots removes orphaned snapshots that don't have corresponding HTTP requests
func (m *SnapshotManager) CleanupSnapshots(snapshotsDir string, activeSnapshots map[string]bool) error {
	// Get all snapshot files
	allSnapshots, err := m.ListSnapshots(snapshotsDir)
	if err != nil {
		return err
	}

	// Remove snapshots that are not in the active set
	for _, snap := range allSnapshots {
		if !activeSnapshots[snap] {
			if err := m.fileWriter.Remove(snap); err != nil {
				return fmt.Errorf("failed to remove orphaned snapshot %s: %w", snap, err)
			}
		}
	}

	return nil
}
