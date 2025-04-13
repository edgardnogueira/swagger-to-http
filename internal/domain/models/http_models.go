package models

import (
	"time"
)

// HTTPRequest represents an HTTP request for executing
type HTTPRequest struct {
	// Core fields
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
	Auth    *AuthDetails      `json:"auth,omitempty"`
	
	// Fields for file format compatibility
	Name     string      `json:"name,omitempty"`
	Path     string      `json:"path,omitempty"`
	Tag      string      `json:"tag,omitempty"`
	Comments []string    `json:"comments,omitempty"`
	
	// Additional fields for variable handling
	FormValues  map[string]string `json:"formValues,omitempty"`
	QueryParams map[string]string `json:"queryParams,omitempty"`
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

// HTTPHeader represents an HTTP header
type HTTPHeader struct {
	Name  string
	Value string
}

// ConvertHeadersToMap converts HTTP headers from slice to map format
func ConvertHeadersToMap(headers []HTTPHeader) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		result[header.Name] = header.Value
	}
	return result
}

// ConvertMapToHeaders converts HTTP headers from map to slice format
func ConvertMapToHeaders(headers map[string]string) []HTTPHeader {
	var result []HTTPHeader
	for name, value := range headers {
		result = append(result, HTTPHeader{Name: name, Value: value})
	}
	return result
}

// HTTPFileRequest represents an HTTP request in .http file format (for backward compatibility)
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

// ToHTTPRequest converts an HTTPFileRequest to an HTTPRequest
func (r *HTTPFileRequest) ToHTTPRequest() *HTTPRequest {
	return &HTTPRequest{
		Method:   r.Method,
		URL:      r.URL,
		Headers:  ConvertHeadersToMap(r.Headers),
		Body:     r.Body,
		Name:     r.Name,
		Path:     r.Path,
		Tag:      r.Tag,
		Comments: r.Comments,
	}
}

// FromHTTPRequest creates an HTTPFileRequest from an HTTPRequest
func FromHTTPRequest(r *HTTPRequest) *HTTPFileRequest {
	return &HTTPFileRequest{
		Method:   r.Method,
		URL:      r.URL,
		Headers:  ConvertMapToHeaders(r.Headers),
		Body:     r.Body,
		Name:     r.Name,
		Path:     r.Path,
		Tag:      r.Tag,
		Comments: r.Comments,
	}
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
