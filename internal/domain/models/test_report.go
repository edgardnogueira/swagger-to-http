package models

import (
	"time"
)

// TestReport represents a complete test execution report
type TestReport struct {
	Name        string            `json:"name"`
	Summary     TestSummary       `json:"summary"`
	Results     []TestResult      `json:"results"`
	Environment map[string]string `json:"environment"`
	CreatedAt   time.Time         `json:"createdAt"`
	Sequences   []TestSequenceResult `json:"sequences,omitempty"`
}

// TestSummary contains the summary statistics for a test run
type TestSummary struct {
	TotalTests       int         `json:"totalTests"`
	PassedTests      int         `json:"passedTests"`
	FailedTests      int         `json:"failedTests"`
	SkippedTests     int         `json:"skippedTests"`
	ErrorTests       int         `json:"errorTests"`
	DurationMs       int64       `json:"durationMs"`
	StartTime        time.Time   `json:"startTime"`
	EndTime          time.Time   `json:"endTime"`
	SnapshotsTotal   int         `json:"snapshotsTotal"`
	SnapshotsUpdated int         `json:"snapshotsUpdated"`
	SnapshotsCreated int         `json:"snapshotsCreated"`
	SchemaValidated  int         `json:"schemaValidated,omitempty"`
	SchemaFailed     int         `json:"schemaFailed,omitempty"`
	SequencesTotal   int         `json:"sequencesTotal,omitempty"`
	SequencesPassed  int         `json:"sequencesPassed,omitempty"`
	SequencesFailed  int         `json:"sequencesFailed,omitempty"`
}

// TestResult represents the result of a single test (HTTP request)
type TestResult struct {
	Name            string            `json:"name"`
	FilePath        string            `json:"filePath"`
	Request         *HTTPRequest      `json:"request"`
	Response        *HTTPResponse     `json:"response"`
	SnapshotResult  *SnapshotResult   `json:"snapshotResult,omitempty"`
	SchemaResult    *SchemaValidationResult `json:"schemaResult,omitempty"`
	Duration        time.Duration     `json:"duration"`
	Status          TestStatus        `json:"status"`
	Error           string            `json:"error,omitempty"`
	Tags            []string          `json:"tags"`
	MetaData        map[string]string `json:"metaData,omitempty"`
	ExtractedVars   map[string]string `json:"extractedVars,omitempty"`
	AssertionResults []TestAssertionResult `json:"assertionResults,omitempty"`
}

// TestStatus represents the status of a test
type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusSkipped TestStatus = "skipped"
	TestStatusError   TestStatus = "error"
)

// TestFilter defines criteria for filtering tests
type TestFilter struct {
	Tags        []string          // Filter by tags
	Paths       []string          // Filter by request paths
	Methods     []string          // Filter by HTTP methods
	StatusCodes []int             // Filter by response status codes
	Names       []string          // Filter by test names
	Metadata    map[string]string // Filter by metadata
}

// TestReportOptions defines options for generating test reports
type TestReportOptions struct {
	IncludeRequests   bool   // Include full request details in report
	IncludeResponses  bool   // Include full response details in report
	Format            string // Report format (json, html, junit, console)
	OutputPath        string // Path to write report file
	ColorOutput       bool   // Use colors in console output
	Detailed          bool   // Include detailed information
	IncludeExtracted  bool   // Include extracted variables in report
	IncludeAssertions bool   // Include assertion results in report
}

// TestRunOptions defines options for running tests
type TestRunOptions struct {
	UpdateSnapshots     string            // Update mode for snapshots (none, all, failed, missing)
	FailOnMissing       bool              // Fail when snapshot is missing
	IgnoreHeaders       []string          // Headers to ignore in comparison
	Timeout             time.Duration     // HTTP request timeout
	Parallel            bool              // Run tests in parallel
	MaxConcurrent       int               // Maximum number of concurrent tests
	StopOnFailure       bool              // Stop testing after first failure
	Filter              TestFilter        // Test filter criteria
	ReportOptions       TestReportOptions // Report generation options
	EnvironmentVars     map[string]string // Environment variables for tests
	ContinuousMode      bool              // Run in continuous (watch) mode
	WatchPaths          []string          // Paths to watch for changes
	WatchIntervalMs     int               // Interval between watch checks in milliseconds
	
	// Advanced testing features (issue #13)
	ValidateSchema      bool              // Validate responses against OpenAPI schema
	ValidationOptions   ValidationOptions // Options for schema validation
	SequentialRun       bool              // Run tests in sequence with dependencies
	ExtractVariables    bool              // Extract variables from responses for use in subsequent tests
	VariableFormat      string            // Format for variable substitution (default: ${varname})
	SaveVariables       bool              // Save extracted variables to a file
	VariablesPath       string            // Path to save/load variables 
	EnableAssertions    bool              // Enable assertions on responses
	FailFast            bool              // Stop sequence on first failure
	DebugMode           bool              // Print detailed debug information
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode     int                 `json:"statusCode"`
	Status         string              `json:"status"`
	Headers        map[string][]string `json:"headers"`
	Body           []byte              `json:"body"`
	ContentType    string              `json:"contentType"`
	ContentLength  int64               `json:"contentLength"`
	Protocol       string              `json:"protocol"`
	ReceivedAt     time.Time           `json:"receivedAt"`
}

// SnapshotDiff represents the difference between a response and a snapshot
type SnapshotDiff struct {
	HasDiff    bool   `json:"hasDiff"`
	DiffString string `json:"diffString"`
	HeaderDiff map[string][]string `json:"headerDiff"`
	BodyDiff   string `json:"bodyDiff"`
	StatusDiff bool   `json:"statusDiff"`
}
