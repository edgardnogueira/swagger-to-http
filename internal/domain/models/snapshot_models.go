package models

import (
	"errors"
	"time"
)

// SnapshotDiff represents the difference between a response and a snapshot
type SnapshotDiff struct {
	// Simple version for test reports
	HasDiff    bool              `json:"hasDiff"`
	DiffString string            `json:"diffString"`
	HeaderDiff map[string][]string `json:"headerDiff"`
	BodyDiff   string            `json:"bodyDiff"`
	StatusDiff bool              `json:"statusDiff"`
	
	// Extended version for detailed analysis
	RequestPath    string      `json:"requestPath,omitempty"`
	RequestMethod  string      `json:"requestMethod,omitempty"`
	StatusDiffExt  *StatusDiff  `json:"statusDiffExt,omitempty"`
	HeaderDiffExt  *HeaderDiff  `json:"headerDiffExt,omitempty"`
	BodyDiffExt    *BodyDiff    `json:"bodyDiffExt,omitempty"`
	Equal          bool        `json:"equal,omitempty"`
}

// StatusDiff represents the difference between two status codes
type StatusDiff struct {
	Expected int  `json:"expected"`
	Actual   int  `json:"actual"`
	Equal    bool `json:"equal"`
}

// HeaderDiff represents the difference between two sets of headers
type HeaderDiff struct {
	MissingHeaders  map[string][]string `json:"missingHeaders"`
	ExtraHeaders    map[string][]string `json:"extraHeaders"`
	DifferentValues map[string]HeaderValueDiff `json:"differentValues"`
	Equal           bool `json:"equal"`
}

// HeaderValueDiff represents the difference between two header values
type HeaderValueDiff struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
}

// BodyDiff represents the difference between two response bodies
type BodyDiff struct {
	ContentType     string `json:"contentType"`
	ExpectedSize    int    `json:"expectedSize"`
	ActualSize      int    `json:"actualSize"`
	ExpectedContent string `json:"expectedContent"`
	ActualContent   string `json:"actualContent"`
	DiffContent     string `json:"diffContent"`
	JsonDiff        *JsonDiff `json:"jsonDiff,omitempty"`
	Equal           bool   `json:"equal"`
}

// JsonDiff represents the difference between two JSON objects
type JsonDiff struct {
	MissingFields   []string `json:"missingFields"`
	ExtraFields     []string `json:"extraFields"`
	DifferentTypes  map[string]TypeDiff `json:"differentTypes"`
	DifferentValues map[string]ValueDiff `json:"differentValues"`
	Equal           bool `json:"equal"`
}

// TypeDiff represents a difference in types
type TypeDiff struct {
	ExpectedType string `json:"expectedType"`
	ActualType   string `json:"actualType"`
}

// ValueDiff represents a difference in values
type ValueDiff struct {
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}

// SnapshotData represents a stored snapshot of an HTTP response
type SnapshotData struct {
	Metadata SnapshotMetadata `json:"metadata"`
	Content  string           `json:"content"`
}

// SnapshotMetadata contains metadata about a snapshot
type SnapshotMetadata struct {
	RequestPath    string                `json:"requestPath"`
	RequestMethod  string                `json:"requestMethod"`
	ContentType    string                `json:"contentType"`
	StatusCode     int                   `json:"statusCode"`
	Headers        map[string][]string   `json:"headers"`
	CreatedAt      time.Time             `json:"createdAt"`
}

// SnapshotOptions defines options for saving and comparing snapshots
type SnapshotOptions struct {
	// UpdateExisting determines whether to update existing snapshots
	UpdateExisting bool
	
	// IgnoreHeaders is a list of headers to ignore during comparison
	IgnoreHeaders []string
	
	// IgnoreFields is a list of JSON fields to ignore during comparison
	IgnoreFields []string
	
	// UpdateMode specifies the update mode (all, failed, none)
	UpdateMode string
	
	// BasePath is the base path for storing snapshots
	BasePath string
}

// SnapshotResult represents the result of a snapshot comparison
type SnapshotResult struct {
	SnapshotPath  string        `json:"snapshotPath"`
	Exists        bool          `json:"exists"`
	Matches       bool          `json:"matches"`
	Diff          *SnapshotDiff `json:"diff,omitempty"`
	WasUpdated    bool          `json:"wasUpdated"`
	WasCreated    bool          `json:"wasCreated"`
	UpdateMode    string        `json:"updateMode"`
	Error         string        `json:"error,omitempty"`
	
	// Fields for compatibility
	Passed        bool          `json:"passed,omitempty"`
	Updated       bool          `json:"updated,omitempty"`
	Created       bool          `json:"created,omitempty"`
	RequestPath   string        `json:"requestPath,omitempty"`
	RequestMethod string        `json:"requestMethod,omitempty"`
}

// SnapshotStats tracks statistics for snapshot testing
type SnapshotStats struct {
	Total     int
	Passed    int
	Failed    int
	Updated   int
	Created   int
	Errors    int
	StartTime time.Time
	EndTime   time.Time
}

// GetError returns the error as an error interface
func (sr *SnapshotResult) GetError() error {
	if sr.Error != "" {
		return errors.New(sr.Error)
	}
	return nil
}

// SetError sets the error from an error interface
func (sr *SnapshotResult) SetError(err error) {
	if err != nil {
		sr.Error = err.Error()
	} else {
		sr.Error = ""
	}
}

// SyncCompatibilityFields ensures all compatibility fields are updated
func (sr *SnapshotResult) SyncCompatibilityFields() {
	sr.Passed = sr.Matches
	sr.Updated = sr.WasUpdated
	sr.Created = sr.WasCreated
	
	if sr.Diff != nil {
		sr.RequestPath = sr.Diff.RequestPath
		sr.RequestMethod = sr.Diff.RequestMethod
	}
}
