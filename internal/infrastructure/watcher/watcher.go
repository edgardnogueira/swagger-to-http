package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// TestWatcherService implements the TestWatcher interface
type TestWatcherService struct {
	testRunner   application.TestRunner
	testReporter application.TestReporter
	stopChan     chan struct{}
	wg           sync.WaitGroup
	logger       *log.Logger
}

// NewTestWatcherService creates a new TestWatcherService
func NewTestWatcherService(
	testRunner application.TestRunner,
	testReporter application.TestReporter,
) *TestWatcherService {
	return &TestWatcherService{
		testRunner:   testRunner,
		testReporter: testReporter,
		stopChan:     make(chan struct{}),
		logger:       log.New(os.Stdout, "[Watcher] ", log.LstdFlags),
	}
}

// Watch starts watching for changes and running tests
func (s *TestWatcherService) Watch(ctx context.Context, patterns []string, options models.TestRunOptions) error {
	// Stop any existing watches
	s.Stop()

	// Create a new stop channel
	s.stopChan = make(chan struct{})

	// Set a reasonable default interval if not specified
	if options.WatchIntervalMs <= 0 {
		options.WatchIntervalMs = 1000 // Default to 1 second
	}

	// Find files to watch
	watchPaths := make(map[string]time.Time)
	if len(options.WatchPaths) > 0 {
		// Use provided watch paths
		for _, path := range options.WatchPaths {
			// Get initial modification time
			fileInfo, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("failed to stat watch path %s: %w", path, err)
			}
			watchPaths[path] = fileInfo.ModTime()
		}
	} else {
		// Derive watch paths from test file patterns
		files, err := s.findFilesToWatch(patterns)
		if err != nil {
			return fmt.Errorf("failed to find files to watch: %w", err)
		}

		// Get initial modification times for all files
		for _, file := range files {
			fileInfo, err := os.Stat(file)
			if err != nil {
				return fmt.Errorf("failed to stat file %s: %w", file, err)
			}
			watchPaths[file] = fileInfo.ModTime()
		}
	}

	// Start a goroutine to watch for changes
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ticker := time.NewTicker(time.Duration(options.WatchIntervalMs) * time.Millisecond)
		defer ticker.Stop()

		// Initial run of the tests
		s.runTests(ctx, patterns, options)

		for {
			select {
			case <-ticker.C:
				// Check for file changes
				if s.checkForChanges(watchPaths) {
					s.logger.Println("Changes detected, running tests...")
					s.runTests(ctx, patterns, options)
				}
			case <-s.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Stop stops watching for changes
func (s *TestWatcherService) Stop() error {
	select {
	case <-s.stopChan:
		// Already stopped
	default:
		close(s.stopChan)
	}

	// Wait for the watcher to stop
	s.wg.Wait()
	return nil
}

// Helper methods

// findFilesToWatch finds all files to watch based on the test patterns
func (s *TestWatcherService) findFilesToWatch(patterns []string) ([]string, error) {
	var files []string

	for _, pattern := range patterns {
		// Find file matches
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		// Add matches to the list
		for _, match := range matches {
			// Only add .http files
			if filepath.Ext(match) == ".http" {
				files = append(files, match)
			}
		}
	}

	// Add swagger/OpenAPI files that might be referenced
	swaggerFiles, err := s.findSwaggerFiles()
	if err != nil {
		return nil, err
	}
	files = append(files, swaggerFiles...)

	return files, nil
}

// findSwaggerFiles finds all swagger/OpenAPI files in the current directory
func (s *TestWatcherService) findSwaggerFiles() ([]string, error) {
	var files []string

	// Check for common swagger file patterns
	patterns := []string{
		"*.json",
		"*.yaml",
		"*.yml",
		"swagger/*.json",
		"swagger/*.yaml",
		"swagger/*.yml",
		"api/*.json",
		"api/*.yaml",
		"api/*.yml",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		for _, match := range matches {
			// Simple heuristic to check if this is likely a swagger file
			isSwagger, err := s.isLikelySwaggerFile(match)
			if err != nil {
				s.logger.Printf("Warning: Error checking file %s: %v", match, err)
				continue
			}

			if isSwagger {
				files = append(files, match)
			}
		}
	}

	return files, nil
}

// isLikelySwaggerFile checks if a file is likely to be a swagger/OpenAPI file
func (s *TestWatcherService) isLikelySwaggerFile(filePath string) (bool, error) {
	// Read the first 1KB of the file to check common swagger identifiers
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	content := string(buffer[:n])

	// Look for common swagger/OpenAPI identifiers
	swaggerIndicators := []string{
		"\"swagger\"", "\"openapi\"",
		"swagger:", "openapi:",
		"\"paths\"", "paths:",
		"\"info\"", "info:",
	}

	matchCount := 0
	for _, indicator := range swaggerIndicators {
		if contains(content, indicator) {
			matchCount++
		}
	}

	// If we find at least 3 indicators, it's likely a swagger file
	return matchCount >= 3, nil
}

// checkForChanges checks if any of the watched files have changed
func (s *TestWatcherService) checkForChanges(watchPaths map[string]time.Time) bool {
	changed := false

	for path, lastModTime := range watchPaths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			// File might have been deleted, ignore
			s.logger.Printf("Warning: Failed to stat file %s: %v", path, err)
			continue
		}

		if fileInfo.ModTime().After(lastModTime) {
			s.logger.Printf("File changed: %s", path)
			watchPaths[path] = fileInfo.ModTime()
			changed = true
		}
	}

	return changed
}

// runTests runs the tests and reports the results
func (s *TestWatcherService) runTests(ctx context.Context, patterns []string, options models.TestRunOptions) {
	// Run the tests
	report, err := s.testRunner.RunTests(ctx, patterns, options)
	if err != nil {
		s.logger.Printf("Error running tests: %v", err)
		return
	}

	// Create console report options
	reportOptions := models.TestReportOptions{
		Format:           "console",
		ColorOutput:      true,
		IncludeRequests:  false,
		IncludeResponses: false,
	}

	// Print the results to stdout
	err = s.testReporter.PrintReport(ctx, report, reportOptions, os.Stdout)
	if err != nil {
		s.logger.Printf("Error printing report: %v", err)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) >= len(substr) && s[0:len(substr)] == substr
}
