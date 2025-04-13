package snapshot

import (
	"context"
	
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Manager defines the interface for snapshot management
type Manager interface {
	// SaveSnapshot saves a HTTP response as a snapshot file
	SaveSnapshot(ctx context.Context, response *models.HTTPResponse, path string) error

	// LoadSnapshot loads a snapshot from a file
	LoadSnapshot(ctx context.Context, path string) (*models.HTTPResponse, error)

	// CompareWithSnapshot compares a current response with a snapshot
	CompareWithSnapshot(ctx context.Context, current *models.HTTPResponse, snapshotPath string) (*models.SnapshotDiff, error)

	// GetSnapshotPath generates a snapshot path for a HTTP request
	GetSnapshotPath(httpFile string, requestName string, baseDir string) string

	// ListSnapshots returns a list of all snapshot files
	ListSnapshots(ctx context.Context, snapshotsDir string) ([]string, error)

	// CleanupSnapshots removes orphaned snapshots that don't have corresponding HTTP requests
	CleanupSnapshots(ctx context.Context, snapshotsDir string, activeSnapshots map[string]bool) error
}
