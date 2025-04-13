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
	TotalTests      int       `json:"totalTests"`
	PassedTests     int       `json:"passedTests"`
	FailedTests     int       `json:"failedTests"`
	SkippedTests    int       `json:"skippedTests"`
	ErrorTests      int       `json:"errorTests"`
	DurationMs      int64     `json:"durationMs"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	SnapshotsTotal   int       `json:"snapshotsTotal"`
	SnapshotsUpdated int       `json:"snapshotsUpdated"`
	SnapshotsCreated int       `json:"snapshotsCreated"`
	SchemaValidated  int       `json:"schemaValidated,omitempty"`
	SchemaFailed     int       `json:"schemaFailed,omitempty"`
	SequencesTotal   int       `json:"sequencesTotal,omitempty"`
	SequencesPassed  int       `json:"sequencesPassed,omitempty"`
	SequencesFailed  int       `json:"sequencesFailed,omitempty"`
}

// TestResult represents the result of a single test (HTTP request)
type TestResult struct {
	Name            string                 `json:"name"`
	FilePath        string                 `json:"filePath"`
	Request         *HTTPRequest           `json:"request"`
	Response        *HTTPResponse          `json:"response"`
	SnapshotResult  *SnapshotResult        `json:"snapshotResult,omitempty"`
	SchemaResult    *SchemaValidationResult `json:"schemaResult,omitempty"`
	Duration        time.Duration          `json:"duration"`
	Status          TestStatus             `json:"status"`
	Error           string                 `json:"error,omitempty"`
	Tags            []string               `json:"tags"`
	MetaData        map[string]string      `json:"metaData,omitempty"`
	ExtractedVars   map[string]string      `json:"extractedVars,omitempty"`
	AssertionResults []TestAssertionResult  `json:"assertionResults,omitempty"`
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
	IncludeRequests   bool    // Include full request details in report
	IncludeResponses  bool    // Include full response details in report
	Format            string  // Report format (json, html, junit, console)
	OutputPath        string  // Path to write report file
	ColorOutput       bool    // Use colors in console output
	Detailed          bool    // Include detailed information
	IncludeExtracted  bool    // Include extracted variables in report
	IncludeAssertions bool    // Include assertion results in report
}

// TestRunOptions defines options for running tests
type TestRunOptions struct {
	UpdateSnapshots     string          // Update mode for snapshots (none, all, failed, missing)
	FailOnMissing       bool            // Fail when snapshot is missing
	IgnoreHeaders       []string        // Headers to ignore in comparison
	Timeout             time.Duration   // HTTP request timeout
	Parallel            bool            // Run tests in parallel
	MaxConcurrent       int             // Maximum number of concurrent tests
	StopOnFailure       bool            // Stop testing after first failure
	Filter              TestFilter      // Test filter criteria
	ReportOptions       TestReportOptions // Report generation options
	EnvironmentVars     map[string]string // Environment variables for tests
	ContinuousMode      bool            // Run in continuous (watch) mode
	WatchPaths          []string        // Paths to watch for changes
	WatchIntervalMs     int             // Interval between watch checks in milliseconds
	
	// Advanced testing features (issue #13)
	ValidateSchema      bool            // Validate responses against OpenAPI schema
	ValidationOptions   ValidationOptions // Options for schema validation
	SequentialRun       bool            // Run tests in sequence with dependencies
	ExtractVariables    bool            // Extract variables from responses for use in subsequent tests
	VariableFormat      string          // Format for variable substitution (default: ${varname})
	SaveVariables       bool            // Save extracted variables to a file
	VariablesPath       string          // Path to save/load variables 
	EnableAssertions    bool            // Enable assertions on responses
	FailFast            bool            // Stop sequence on first failure
	DebugMode           bool            // Print detailed debug information
}

// TestSequenceResult represents the result of a test sequence
type TestSequenceResult struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	FilePath     string        `json:"filePath"`
	Steps        []TestStepResult `json:"steps"`
	Duration     time.Duration `json:"duration"`
	Status       TestStatus    `json:"status"`
	Error        string        `json:"error,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
}

// TestStepResult represents the result of a single test step in a sequence
type TestStepResult struct {
	Name            string                 `json:"name"`
	Request         *HTTPRequest           `json:"request"`
	Response        *HTTPResponse          `json:"response,omitempty"`
	Duration        time.Duration          `json:"duration"`
	Status          TestStatus             `json:"status"`
	Error           string                 `json:"error,omitempty"`
	ExtractedVars   map[string]string      `json:"extractedVars,omitempty"`
	AssertionResults []TestAssertionResult  `json:"assertionResults,omitempty"`
	SchemaResult    *SchemaValidationResult `json:"schemaResult,omitempty"`
	ConditionalSkipped bool                 `json:"conditionalSkipped,omitempty"`
}

// TestSequence represents a sequence of tests to be run
type TestSequence struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Steps       []TestStep `json:"steps"`
	Variables   map[string]string `json:"variables,omitempty"`
}

// TestStep represents a single step in a test sequence
type TestStep struct {
	Name          string       `json:"name"`
	Request       *HTTPRequest `json:"request"`
	Variables     []VariableExtraction `json:"variables,omitempty"`
	ExpectedStatus int          `json:"expectedStatus,omitempty"`
	Assertions    []TestAssertion `json:"assertions,omitempty"`
	SchemaValidate bool         `json:"schemaValidate,omitempty"`
	WaitBefore    time.Duration `json:"waitBefore,omitempty"`
	WaitAfter     time.Duration `json:"waitAfter,omitempty"`
	Condition     *TestCondition `json:"condition,omitempty"`
}

// VariableExtraction represents a variable to extract from a response
type VariableExtraction struct {
	Name      string `json:"name"`
	Source    string `json:"source"` // body, header, status
	Path      string `json:"path"`
	Regex     string `json:"regex,omitempty"`
	Extractor string `json:"extractor,omitempty"` // jsonpath, regex, header
}

// TestAssertion represents an assertion to make on a response
type TestAssertion struct {
	Type      string      `json:"type"` // equals, contains, matches, etc.
	Source    string      `json:"source"` // body, header, status
	Path      string      `json:"path,omitempty"`
	Value     interface{} `json:"value"`
	Negate    bool        `json:"negate,omitempty"`
	Description string     `json:"description,omitempty"`
}

// TestAssertionResult represents the result of a test assertion
type TestAssertionResult struct {
	Type        string      `json:"type"`
	Source      string      `json:"source"`
	Path        string      `json:"path,omitempty"`
	Expected    interface{} `json:"expected"`
	Actual      interface{} `json:"actual"`
	Passed      bool        `json:"passed"`
	Description string      `json:"description,omitempty"`
	Error       string      `json:"error,omitempty"`
}

// TestCondition represents a condition for executing a test step
type TestCondition struct {
	Type      string      `json:"type"` // variable, status, response
	Variable  string      `json:"variable,omitempty"`
	Operator  string      `json:"operator"` // equals, notEquals, contains, etc.
	Value     interface{} `json:"value"`
}

// SchemaValidationResult represents the result of schema validation
type SchemaValidationResult struct {
	Valid       bool                   `json:"valid"`
	Errors      []SchemaValidationError `json:"errors,omitempty"`
	Warnings    []SchemaValidationError `json:"warnings,omitempty"`
	SchemaPath  string                 `json:"schemaPath,omitempty"`
}

// SchemaValidationError represents an error in schema validation
type SchemaValidationError struct {
	Path        string `json:"path"`
	Message     string `json:"message"`
	SchemaPath  string `json:"schemaPath,omitempty"`
}

// ValidationOptions represents options for schema validation
type ValidationOptions struct {
	IgnoreAdditionalProperties bool     `json:"ignoreAdditionalProperties"`
	IgnoreFormats             bool     `json:"ignoreFormats"`
	IgnorePatterns            bool     `json:"ignorePatterns"`
	RequiredPropertiesOnly    bool     `json:"requiredPropertiesOnly"`
	IgnoreProperties          []string `json:"ignoreProperties,omitempty"`
	IgnoreNullable            bool     `json:"ignoreNullable"`
}
