package models

import (
	"time"
)

// SnapshotData represents a stored snapshot of an HTTP response
type SnapshotData struct {
	Metadata SnapshotMetadata `json:"metadata"`
	Content  string           `json:"content"`
}

// SnapshotMetadata contains metadata about a snapshot
type SnapshotMetadata struct {
	RequestPath   string              `json:"requestPath"`
	RequestMethod string              `json:"requestMethod"`
	ContentType   string              `json:"contentType"`
	StatusCode    int                 `json:"statusCode"`
	Headers       map[string][]string `json:"headers"`
	CreatedAt     time.Time           `json:"createdAt"`
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

// SnapshotResult contains the result of a snapshot comparison
type SnapshotResult struct {
	RequestPath   string
	RequestMethod string
	SnapshotPath  string
	Diff          *SnapshotDiff
	Passed        bool
	Updated       bool
	Error         error
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
