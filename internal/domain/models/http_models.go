package models

import (
	"time"
)

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

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Auth    *AuthDetails      `json:"auth,omitempty"`
}

// AuthDetails represents authentication details for an HTTP request
type AuthDetails struct {
	Type  string `json:"type"` // Basic, Bearer, etc.
	Value string `json:"value,omitempty"`
}
