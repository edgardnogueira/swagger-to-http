package snapshot

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Service provides high-level snapshot testing functionality
type Service struct {
	manager      *Manager
	options      models.SnapshotOptions
	usedSnapshots map[string]bool
	mu           sync.Mutex
	stats        *models.SnapshotStats
}

// NewService creates a new snapshot service
func NewService(manager *Manager, options models.SnapshotOptions) *Service {
	return &Service{
		manager:      manager,
		options:      options,
		usedSnapshots: make(map[string]bool),
		stats:        &models.SnapshotStats{
			StartTime: time.Now(),
		},
	}
}

// RunTest runs a snapshot test for a response
func (s *Service) RunTest(ctx context.Context, response *models.HTTPResponse, path string) (*models.SnapshotResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if response == nil {
		return nil, fmt.Errorf("cannot test nil response")
	}
	
	// Mark the snapshot as used
	snapshotPath := filepath.Join(filepath.Dir(path), getResponseSnapshotName(response, path))
	s.usedSnapshots[snapshotPath] = true
	
	// Initialize result
	result := &models.SnapshotResult{
		RequestPath:   response.Request.Path,
		RequestMethod: response.Request.Method,
		SnapshotPath:  snapshotPath,
		Passed:        false,
	}
	
	// Compare with snapshot
	diff, err := s.manager.CompareWithSnapshot(ctx, response, snapshotPath)
	if err != nil {
		if s.options.UpdateMode == "all" || s.options.UpdateMode == "missing" {
			// Create new snapshot
			if createErr := s.manager.SaveSnapshot(ctx, response, snapshotPath); createErr != nil {
				result.Error = fmt.Errorf("failed to create snapshot: %w", createErr)
				s.stats.Errors++
				return result, result.Error
			}
			
			result.Passed = true
			result.Updated = true
			s.stats.Created++
			s.stats.Passed++
			s.stats.Total++
			return result, nil
		}
		
		result.Error = fmt.Errorf("snapshot comparison failed: %w", err)
		s.stats.Errors++
		return result, result.Error
	}
	
	// Set result properties
	result.Diff = diff
	result.Passed = diff.Equal
	
	// Update stats
	s.stats.Total++
	if result.Passed {
		s.stats.Passed++
	} else {
		s.stats.Failed++
		
		if s.options.UpdateMode == "all" || s.options.UpdateMode == "failed" {
			// Update snapshot
			if updateErr := s.manager.SaveSnapshot(ctx, response, snapshotPath); updateErr != nil {
				result.Error = fmt.Errorf("failed to update snapshot: %w", updateErr)
				s.stats.Errors++
				return result, result.Error
			}
			
			result.Updated = true
			s.stats.Updated++
		}
	}
	
	return result, nil
}

// CleanupUnusedSnapshots removes snapshots that weren't used in tests
func (s *Service) CleanupUnusedSnapshots(ctx context.Context, directory string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.manager.CleanupSnapshots(ctx, directory, s.usedSnapshots)
}

// GetStats returns the current test statistics
func (s *Service) GetStats() *models.SnapshotStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create a copy to avoid race conditions
	stats := *s.stats
	stats.EndTime = time.Now()
	return &stats
}

// ResetStats resets the test statistics
func (s *Service) ResetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.stats = &models.SnapshotStats{
		StartTime: time.Now(),
	}
}

// getResponseSnapshotName generates a snapshot name based on the response
func getResponseSnapshotName(response *models.HTTPResponse, path string) string {
	method := response.Request.Method
	requestPath := response.Request.Path
	
	// Create base filename from original path
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	baseName := base
	if ext != "" {
		baseName = base[:len(base)-len(ext)]
	}
	
	// If request path is empty, use baseName
	if requestPath == "" {
		return fmt.Sprintf("%s_%s.snap.json", baseName, method)
	}
	
	// Otherwise use a sanitized version of the request path
	sanitized := sanitizePathForFilename(requestPath)
	return fmt.Sprintf("%s_%s.snap.json", sanitized, method)
}

// sanitizePathForFilename converts a path to a safe filename
func sanitizePathForFilename(path string) string {
	// Replace any non-alphanumeric characters with underscores
	result := ""
	for _, c := range path {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-' || c == '_' {
			result += string(c)
		} else {
			result += "_"
		}
	}
	
	// Limit length
	if len(result) > 100 {
		result = result[:100]
	}
	
	return result
}
