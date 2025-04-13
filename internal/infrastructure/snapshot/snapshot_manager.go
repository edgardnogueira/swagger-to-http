package snapshot

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"encoding/json"

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
func (m *SnapshotManager) SaveSnapshot(ctx context.Context, response *models.HTTPResponse, path string) error {
	// Format the snapshot as JSON
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Ensure the snapshot directory exists
	dir := filepath.Dir(path)
	if err := m.fileWriter.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Write the snapshot file
	if err := m.fileWriter.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// LoadSnapshot loads a snapshot from a file
func (m *SnapshotManager) LoadSnapshot(ctx context.Context, path string) (*models.HTTPResponse, error) {
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

	// Parse the snapshot
	var response models.HTTPResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &response, nil
}

// CompareWithSnapshot compares a current response with a snapshot
func (m *SnapshotManager) CompareWithSnapshot(ctx context.Context, current *models.HTTPResponse, snapshotPath string) (*models.SnapshotDiff, error) {
	// Load the snapshot
	expected, err := m.LoadSnapshot(ctx, snapshotPath)
	if err != nil {
		return nil, err
	}

	// Compare responses
	diff := &models.SnapshotDiff{
		RequestPath:   current.Request.Path,
		RequestMethod: current.Request.Method,
		Equal:         true,
	}

	// Check status codes
	if expected.StatusCode != current.StatusCode {
		diff.Equal = false
		diff.StatusDiff = true
		diff.StatusDiffExt = &models.StatusDiff{
			Expected: expected.StatusCode,
			Actual:   current.StatusCode,
			Equal:    false,
		}
	}

	// Compare bodies
	if expected.Body != current.Body {
		diff.Equal = false
		diff.BodyDiff = "Bodies differ"
		diff.BodyDiffExt = &models.BodyDiff{
			ExpectedContent: expected.Body,
			ActualContent:   current.Body,
			Equal:           false,
		}
	}

	// Compare headers (simplified)
	if !headersEqual(expected.Headers, current.Headers) {
		diff.Equal = false
		diff.HeaderDiff = make(map[string][]string)
		diff.HeaderDiffExt = &models.HeaderDiff{
			Equal: false,
		}
	}

	return diff, nil
}

// headersEqual compares two header maps for equality
func headersEqual(expected, actual map[string][]string) bool {
	if len(expected) != len(actual) {
		return false
	}

	for name, values := range expected {
		actualValues, ok := actual[name]
		if !ok {
			return false
		}

		if len(values) != len(actualValues) {
			return false
		}

		// Compare values (ignoring order)
		for _, value := range values {
			found := false
			for _, actualValue := range actualValues {
				if value == actualValue {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
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
func (m *SnapshotManager) ListSnapshots(ctx context.Context, snapshotsDir string) ([]string, error) {
	return m.fileWriter.Glob(filepath.Join(snapshotsDir, "**/*.snap"))
}

// CleanupSnapshots removes orphaned snapshots that don't have corresponding HTTP requests
func (m *SnapshotManager) CleanupSnapshots(ctx context.Context, snapshotsDir string, activeSnapshots map[string]bool) error {
	// Get all snapshot files
	allSnapshots, err := m.ListSnapshots(ctx, snapshotsDir)
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
