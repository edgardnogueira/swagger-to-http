package models

import (
	"time"
)

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode    int
	Status        string
	Headers       map[string][]string
	Body          []byte
	ContentType   string
	ContentLength int64
	Duration      time.Duration
	Request       *HTTPRequest
	RequestID     string
	Timestamp     time.Time
}

// SnapshotDiff represents the difference between a response and a snapshot
type SnapshotDiff struct {
	RequestPath   string
	RequestMethod string
	StatusDiff    *StatusDiff
	HeaderDiff    *HeaderDiff
	BodyDiff      *BodyDiff
	Equal         bool
}

// StatusDiff represents the difference between two status codes
type StatusDiff struct {
	Expected int
	Actual   int
	Equal    bool
}

// HeaderDiff represents the difference between two sets of headers
type HeaderDiff struct {
	MissingHeaders  map[string][]string
	ExtraHeaders    map[string][]string
	DifferentValues map[string]HeaderValueDiff
	Equal           bool
}

// HeaderValueDiff represents the difference between two header values
type HeaderValueDiff struct {
	Expected []string
	Actual   []string
}

// BodyDiff represents the difference between two response bodies
type BodyDiff struct {
	ContentType     string
	ExpectedSize    int
	ActualSize      int
	ExpectedContent string
	ActualContent   string
	DiffContent     string
	JsonDiff        *JsonDiff
	Equal           bool
}

// JsonDiff represents the difference between two JSON objects
type JsonDiff struct {
	MissingFields  []string
	ExtraFields    []string
	DifferentTypes map[string]TypeDiff
	DifferentValues map[string]ValueDiff
	Equal          bool
}

// TypeDiff represents a difference in types
type TypeDiff struct {
	ExpectedType string
	ActualType   string
}

// ValueDiff represents a difference in values
type ValueDiff struct {
	Expected interface{}
	Actual   interface{}
}
