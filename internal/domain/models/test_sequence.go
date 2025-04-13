package models

import "time"

// TestSequence represents a sequence of HTTP tests with dependencies
type TestSequence struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Steps       []TestSequenceStep  `json:"steps"`
	Variables   map[string]string   `json:"variables,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	Tags        []string            `json:"tags,omitempty"`
	Metadata    map[string]string   `json:"metadata,omitempty"`
	FilePath    string              `json:"filePath"`
}

// TestSequenceStep represents a single step in a test sequence
type TestSequenceStep struct {
	Name            string                  `json:"name"`
	Description     string                  `json:"description,omitempty"`
	Request         *HTTPRequest            `json:"request"`
	ExpectedStatus  int                     `json:"expectedStatus,omitempty"`
	Variables       []VariableExtraction    `json:"variables,omitempty"`
	WaitBefore      time.Duration           `json:"waitBefore,omitempty"`
	WaitAfter       time.Duration           `json:"waitAfter,omitempty"`
	Skip            bool                    `json:"skip,omitempty"`
	SkipCondition   string                  `json:"skipCondition,omitempty"`
	StopOnFail      bool                    `json:"stopOnFail,omitempty"`
	SchemaValidate  bool                    `json:"schemaValidate,omitempty"`
	Assertions      []TestAssertion         `json:"assertions,omitempty"`
	ExpectedResult  *TestSequenceStepResult `json:"expectedResult,omitempty"`
}

// TestSequenceResult represents the result of running a test sequence
type TestSequenceResult struct {
	Name          string                   `json:"name"`
	Success       bool                     `json:"success"`
	StepResults   []TestSequenceStepResult `json:"stepResults"`
	ExecutionTime time.Duration            `json:"executionTime"`
	Variables     map[string]string        `json:"variables,omitempty"`
	StartTime     time.Time                `json:"startTime"`
	EndTime       time.Time                `json:"endTime"`
	Error         string                   `json:"error,omitempty"`
}

// TestSequenceStepResult represents the result of a single step in a test sequence
type TestSequenceStepResult struct {
	Name            string               `json:"name"`
	Status          TestStatus           `json:"status"`
	Response        *HTTPResponse        `json:"response,omitempty"`
	Variables       map[string]string    `json:"variables,omitempty"`
	ExecutionTime   time.Duration        `json:"executionTime"`
	Error           string               `json:"error,omitempty"`
	ValidationError string               `json:"validationError,omitempty"`
	SchemaResult    *SchemaValidationResult `json:"schemaResult,omitempty"`
	AssertionResults []TestAssertionResult  `json:"assertionResults,omitempty"`
}

// VariableExtraction defines how to extract a variable from an HTTP response
type VariableExtraction struct {
	Name     string `json:"name"`
	Source   string `json:"source"` // "body", "header", "status"
	Path     string `json:"path,omitempty"`
	Regexp   string `json:"regexp,omitempty"`
	Default  string `json:"default,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// TestAssertion defines an assertion to be made against the HTTP response
type TestAssertion struct {
	Type       string   `json:"type"` // "contains", "equals", "matches", "exists", "notExists"
	Source     string   `json:"source"` // "body", "header", "status"
	Path       string   `json:"path,omitempty"`
	Value      string   `json:"value,omitempty"`
	Values     []string `json:"values,omitempty"`
	Not        bool     `json:"not,omitempty"`
	IgnoreCase bool     `json:"ignoreCase,omitempty"`
}

// TestAssertionResult represents the result of a test assertion
type TestAssertionResult struct {
	Type      string `json:"type"`
	Source    string `json:"source"`
	Path      string `json:"path,omitempty"`
	Succeeded bool   `json:"succeeded"`
	Actual    string `json:"actual,omitempty"`
	Expected  string `json:"expected,omitempty"`
	Message   string `json:"message,omitempty"`
}
