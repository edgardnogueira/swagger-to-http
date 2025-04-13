package snapshot

import (
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Manager defines the interface for snapshot management
type Manager interface {
	// SaveSnapshot saves a HTTP response as a snapshot file
	SaveSnapshot(response *models.HTTPResponse, path string, format string) error

	// LoadSnapshot loads a snapshot from a file
	LoadSnapshot(path string, format string) (*models.HTTPResponse, error)

	// CompareSnapshots compares a current response with a snapshot
	CompareSnapshots(current *models.HTTPResponse, snapshotPath string, format string) (*ComparisonResult, error)

	// GetSnapshotPath generates a snapshot path for a HTTP request
	GetSnapshotPath(httpFile string, requestName string, baseDir string) string

	// ListSnapshots returns a list of all snapshot files
	ListSnapshots(snapshotsDir string) ([]string, error)

	// CleanupSnapshots removes orphaned snapshots that don't have corresponding HTTP requests
	CleanupSnapshots(snapshotsDir string, activeSnapshots map[string]bool) error
}