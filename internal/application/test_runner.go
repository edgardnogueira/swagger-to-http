package application

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// TestRunnerService implements the TestRunner interface
type TestRunnerService struct {
	httpExecutor    HTTPExecutor
	snapshotManager SnapshotManager
	fileWriter      FileWriter
}

// NewTestRunnerService creates a new TestRunnerService
func NewTestRunnerService(executor HTTPExecutor, snapshotManager SnapshotManager, fileWriter FileWriter) *TestRunnerService {
	return &TestRunnerService{
		httpExecutor:    executor,
		snapshotManager: snapshotManager,
		fileWriter:      fileWriter,
	}
}

// RunTests runs tests based on HTTP files and options
func (s *TestRunnerService) RunTests(ctx context.Context, patterns []string, options models.TestRunOptions) (*models.TestReport, error) {
	// Find all test files matching the patterns and filter
	files, err := s.FindTests(ctx, patterns, options.Filter)
	if err != nil {
		return nil, fmt.Errorf("error finding tests: %w", err)
	}

	// Create the test report
	report := &models.TestReport{
		Name:       "HTTP Tests",
		Summary:    models.TestSummary{},
		Results:    []models.TestResult{},
		CreatedAt:  time.Now(),
		Environment: make(map[string]string),
	}

	// Add environment variables to report
	for k, v := range options.EnvironmentVars {
		report.Environment[k] = v
	}

	// Set start time for the test run
	report.Summary.StartTime = time.Now()

	// Run the tests (parallel or sequential)
	if options.Parallel && options.MaxConcurrent > 0 {
		results, err := s.runTestsParallel(ctx, files, options)
		if err != nil {
			return nil, err
		}
		report.Results = results
	} else {
		results, err := s.runTestsSequential(ctx, files, options)
		if err != nil {
			return nil, err
		}
		report.Results = results
	}

	// Set end time for the test run
	report.Summary.EndTime = time.Now()
	report.Summary.DurationMs = report.Summary.EndTime.Sub(report.Summary.StartTime).Milliseconds()

	// Calculate summary statistics
	s.calculateSummary(&report.Summary, report.Results)

	return report, nil
}

// RunTest runs a single test
func (s *TestRunnerService) RunTest(ctx context.Context, request *models.HTTPRequest, options models.TestRunOptions) (*models.TestResult, error) {
	startTime := time.Now()

	// Create test result
	result := &models.TestResult{
		Name:     request.Name,
		Request:  request,
		FilePath: request.Path,
		Tags:     []string{request.Tag},
		Status:   models.TestStatusSkipped,
	}

	// Execute the request
	response, err := s.httpExecutor.Execute(ctx, request, options.EnvironmentVars)
	if err != nil {
		result.Status = models.TestStatusError
		result.Error = err.Error()
		return result, nil
	}

	result.Response = response
	result.Duration = time.Since(startTime)

	// Generate snapshot path
	snapshotPath := s.generateSnapshotPath(request, options)

	// Check if we should compare with snapshot
	if options.UpdateSnapshots == "all" {
		// Save snapshot directly without comparison
		err = s.snapshotManager.SaveSnapshot(ctx, response, snapshotPath)
		if err != nil {
			result.Status = models.TestStatusError
			result.Error = fmt.Sprintf("failed to save snapshot: %v", err)
		} else {
			result.Status = models.TestStatusPassed
			if result.SnapshotResult == nil {
				result.SnapshotResult = &models.SnapshotResult{
					Updated: true,
					Passed:  true,
				}
			}
		}
		return result, nil
	}

	// Try to load existing snapshot
	_, err = s.snapshotManager.LoadSnapshot(ctx, snapshotPath)
	if err != nil {
		// Snapshot doesn't exist
		if options.UpdateSnapshots == "missing" {
			// Create a new snapshot
			err = s.snapshotManager.SaveSnapshot(ctx, response, snapshotPath)
			if err != nil {
				result.Status = models.TestStatusError
				result.Error = fmt.Sprintf("failed to create snapshot: %v", err)
			} else {
				result.Status = models.TestStatusPassed
				if result.SnapshotResult == nil {
					result.SnapshotResult = &models.SnapshotResult{
						Created: true,
						Passed:  true,
					}
				}
			}
		} else if options.FailOnMissing {
			result.Status = models.TestStatusFailed
			result.Error = "snapshot missing"
		} else {
			result.Status = models.TestStatusPassed
			result.Error = "snapshot missing, not failing due to configuration"
		}
		return result, nil
	}

	// Compare with existing snapshot
	diff, err := s.snapshotManager.CompareSnapshots(ctx, response, snapshotPath)
	if err != nil {
		result.Status = models.TestStatusError
		result.Error = fmt.Sprintf("snapshot comparison error: %v", err)
		return result, nil
	}

	// Store snapshot result
	snapshotResult := &models.SnapshotResult{
		RequestPath:   request.Path,
		RequestMethod: request.Method,
		SnapshotPath:  snapshotPath,
		Diff:          diff,
		Passed:        !diff.HasDiff,
	}
	result.SnapshotResult = snapshotResult

	// Handle test result based on comparison
	if diff.HasDiff {
		// Test failed
		result.Status = models.TestStatusFailed

		// Check if we should update failed snapshots
		if options.UpdateSnapshots == "failed" {
			err = s.snapshotManager.SaveSnapshot(ctx, response, snapshotPath)
			if err != nil {
				result.Error = fmt.Sprintf("failed to update snapshot: %v", err)
			} else {
				snapshotResult.Updated = true
				result.Error = "snapshot updated"
			}
		} else {
			result.Error = "snapshot comparison failed"
		}
	} else {
		// Test passed
		result.Status = models.TestStatusPassed
	}

	return result, nil
}

// RunTestFile runs all tests in a file
func (s *TestRunnerService) RunTestFile(ctx context.Context, file *models.HTTPFile, options models.TestRunOptions) ([]*models.TestResult, error) {
	var results []*models.TestResult

	for _, request := range file.Requests {
		// Check if the test meets the filter criteria
		if !s.matchesFilter(&request, options.Filter) {
			continue
		}

		// Set the file path in the request
		request.Path = file.Filename

		// Run the test
		result, err := s.RunTest(ctx, &request, options)
		if err != nil {
			return nil, err
		}

		results = append(results, result)

		// Stop on failure if configured
		if options.StopOnFailure && (result.Status == models.TestStatusFailed || result.Status == models.TestStatusError) {
			break
		}
	}

	return results, nil
}

// FindTests finds all tests matching the provided patterns
func (s *TestRunnerService) FindTests(ctx context.Context, patterns []string, filter models.TestFilter) ([]*models.HTTPFile, error) {
	var files []*models.HTTPFile

	// Process each pattern
	for _, pattern := range patterns {
		// Find matches using filepath.Glob
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error matching pattern %s: %w", pattern, err)
		}

		for _, match := range matches {
			// Skip if not a .http file
			if !strings.HasSuffix(match, ".http") {
				continue
			}

			// Parse the HTTP file
			httpFile, err := s.parseHTTPFile(ctx, match)
			if err != nil {
				return nil, fmt.Errorf("error parsing HTTP file %s: %w", match, err)
			}

			// Apply filter to the HTTP file
			if !s.fileMatchesFilter(httpFile, filter) {
				continue
			}

			files = append(files, httpFile)
		}
	}

	return files, nil
}

// Helper functions

// runTestsSequential runs tests sequentially
func (s *TestRunnerService) runTestsSequential(ctx context.Context, files []*models.HTTPFile, options models.TestRunOptions) ([]models.TestResult, error) {
	var results []models.TestResult

	for _, file := range files {
		fileResults, err := s.RunTestFile(ctx, file, options)
		if err != nil {
			return nil, err
		}

		// Convert pointers to values and add to results
		for _, result := range fileResults {
			results = append(results, *result)
		}

		// Stop on failure if configured
		if options.StopOnFailure {
			for _, result := range fileResults {
				if result.Status == models.TestStatusFailed || result.Status == models.TestStatusError {
					return results, nil
				}
			}
		}
	}

	return results, nil
}

// runTestsParallel runs tests in parallel
func (s *TestRunnerService) runTestsParallel(ctx context.Context, files []*models.HTTPFile, options models.TestRunOptions) ([]models.TestResult, error) {
	var results []models.TestResult
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	resultChan := make(chan *models.TestResult)
	done := make(chan struct{})

	// Calculate how many workers to use
	numWorkers := options.MaxConcurrent
	if numWorkers <= 0 {
		numWorkers = 5 // Default to 5 concurrent tests
	}

	// Get total number of requests
	var totalRequests int
	for _, file := range files {
		totalRequests += len(file.Requests)
	}

	// Don't create more workers than requests
	if numWorkers > totalRequests {
		numWorkers = totalRequests
	}

	// Create work queue
	workQueue := make(chan workItem, totalRequests)

	// Fill work queue
	for _, file := range files {
		for i := range file.Requests {
			workQueue <- workItem{
				file:       file,
				requestIdx: i,
			}
		}
	}
	close(workQueue)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for work := range workQueue {
				// Check for context cancellation
				select {
				case <-ctx.Done():
					return
				default:
					// Continue processing
				}

				file := work.file
				req := file.Requests[work.requestIdx]

				// Check if the test meets the filter criteria
				if !s.matchesFilter(&req, options.Filter) {
					continue
				}

				// Set the file path in the request
				req.Path = file.Filename

				// Clone the options to avoid race conditions
				localOpts := options

				// Run the test
				result, err := s.RunTest(ctx, &req, localOpts)
				if err != nil {
					select {
					case errChan <- err:
						// Error sent
						return
					default:
						// Another error was already sent, just return
						return
					}
				}

				// Send result
				select {
				case resultChan <- result:
					// Result sent
				case <-ctx.Done():
					return
				}

				// Stop on failure if configured
				if options.StopOnFailure && (result.Status == models.TestStatusFailed || result.Status == models.TestStatusError) {
					return
				}
			}
		}()
	}

	// Collect results
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all tests to complete or an error to occur
	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case err := <-errChan:
			return results, err
		case result := <-resultChan:
			mu.Lock()
			results = append(results, *result)
			mu.Unlock()
		case <-done:
			return results, nil
		}
	}
}

// calculateSummary calculates the summary statistics for a test run
func (s *TestRunnerService) calculateSummary(summary *models.TestSummary, results []models.TestResult) {
	summary.TotalTests = len(results)
	
	for _, result := range results {
		switch result.Status {
		case models.TestStatusPassed:
			summary.PassedTests++
		case models.TestStatusFailed:
			summary.FailedTests++
		case models.TestStatusSkipped:
			summary.SkippedTests++
		case models.TestStatusError:
			summary.ErrorTests++
		}

		if result.SnapshotResult != nil {
			summary.SnapshotsTotal++
			if result.SnapshotResult.Updated {
				summary.SnapshotsUpdated++
			}
			if result.SnapshotResult.Created {
				summary.SnapshotsCreated++
			}
		}
	}
}

// generateSnapshotPath generates a path for storing a snapshot
func (s *TestRunnerService) generateSnapshotPath(request *models.HTTPRequest, options models.TestRunOptions) string {
	// Use the configured snapshot directory or default
	snapshotDir := options.Filter.Paths[0]
	if snapshotDir == "" {
		snapshotDir = ".snapshots"
	}

	// Generate a filename based on the request
	filename := fmt.Sprintf("%s_%s", request.Method, strings.Replace(request.URL, "/", "_", -1))
	filename = strings.Replace(filename, ":", "", -1)
	filename = strings.Replace(filename, "?", "_", -1)
	filename = strings.Replace(filename, "&", "_", -1)
	filename = strings.Replace(filename, "=", "_", -1)
	filename = strings.Replace(filename, ".", "_", -1)
	filename = strings.Replace(filename, " ", "_", -1)
	filename = strings.TrimSuffix(filename, "_")

	// Include the tag in the path if available
	if request.Tag != "" {
		return filepath.Join(snapshotDir, request.Tag, filename+".json")
	}

	return filepath.Join(snapshotDir, filename+".json")
}

// matchesFilter checks if a request matches the filter criteria
func (s *TestRunnerService) matchesFilter(request *models.HTTPRequest, filter models.TestFilter) bool {
	// Filter by tag
	if len(filter.Tags) > 0 {
		tagMatch := false
		for _, tag := range filter.Tags {
			if request.Tag == tag {
				tagMatch = true
				break
			}
		}
		if !tagMatch {
			return false
		}
	}

	// Filter by path
	if len(filter.Paths) > 0 {
		pathMatch := false
		for _, path := range filter.Paths {
			if strings.Contains(request.URL, path) {
				pathMatch = true
				break
			}
		}
		if !pathMatch {
			return false
		}
	}

	// Filter by method
	if len(filter.Methods) > 0 {
		methodMatch := false
		for _, method := range filter.Methods {
			if strings.EqualFold(request.Method, method) {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			return false
		}
	}

	// Filter by name
	if len(filter.Names) > 0 {
		nameMatch := false
		for _, name := range filter.Names {
			if strings.Contains(request.Name, name) {
				nameMatch = true
				break
			}
		}
		if !nameMatch {
			return false
		}
	}

	return true
}

// fileMatchesFilter checks if an HTTP file matches the filter criteria
func (s *TestRunnerService) fileMatchesFilter(file *models.HTTPFile, filter models.TestFilter) bool {
	// If no requests, the file doesn't match
	if len(file.Requests) == 0 {
		return false
	}

	// Check if any request in the file matches the filter
	for _, request := range file.Requests {
		if s.matchesFilter(&request, filter) {
			return true
		}
	}

	return false
}

// parseHTTPFile parses an HTTP file from the file system
func (s *TestRunnerService) parseHTTPFile(ctx context.Context, filePath string) (*models.HTTPFile, error) {
	// TODO: Implement HTTP file parser (stub for now)
	return &models.HTTPFile{
		Filename: filePath,
		Requests: []models.HTTPRequest{},
	}, nil
}

// workItem represents a unit of work for parallel processing
type workItem struct {
	file       *models.HTTPFile
	requestIdx int
}
