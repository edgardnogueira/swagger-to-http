package models

import (
	"time"
)

// HTTPRequest represents an HTTP request for executing
type HTTPRequest struct {
	Method  string         `json:"method"`
	URL     string         `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string         `json:"body,omitempty"`
	Auth    *AuthDetails   `json:"auth,omitempty"`
}

// AuthDetails represents authentication details for an HTTP request
type AuthDetails struct {
	Type  string `json:"type"` // Basic, Bearer, etc.
	Value string `json:"value,omitempty"`
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	// StatusCode is the HTTP status code
	StatusCode int `json:"statusCode"`

	// Status is the HTTP status text
	Status string `json:"status"`

	// Headers contains the HTTP headers
	Headers map[string][]string `json:"headers"`

	// Body contains the response body
	Body string `json:"body,omitempty"`

	// For extended response information
	ContentType    string        `json:"contentType,omitempty"`
	ContentLength  int64         `json:"contentLength,omitempty"`
	Duration       time.Duration `json:"duration,omitempty"`
	Request        *HTTPRequest  `json:"request,omitempty"`
	RequestID      string        `json:"requestId,omitempty"`
	Timestamp      time.Time     `json:"timestamp,omitempty"`
	ReceivedAt     time.Time     `json:"receivedAt,omitempty"`
	Protocol       string        `json:"protocol,omitempty"`
}

// HTTPFileRequest represents an HTTP request in .http file format
type HTTPFileRequest struct {
	Name     string
	Method   string
	URL      string
	Headers  []HTTPHeader
	Body     string
	Comments []string
	Tag      string
	Path     string
}

// HTTPHeader represents an HTTP header
type HTTPHeader struct {
	Name  string
	Value string
}

// HTTPFile represents a collection of HTTP requests to be written to a .http file
type HTTPFile struct {
	Filename string
	Requests []HTTPFileRequest
}

// HTTPDirectory represents a directory containing HTTP files
type HTTPDirectory struct {
	Name  string
	Path  string
	Files []HTTPFile
}

// HTTPCollection represents a collection of directories and files
type HTTPCollection struct {
	RootDir      string
	Directories  []HTTPDirectory
	RootFiles    []HTTPFile
}
