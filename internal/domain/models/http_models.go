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
	Name     string    `json:"name,omitempty"`
	Path     string    `json:"path,omitempty"`
	Tag      string    `json:"tag,omitempty"`
	Comments []string  `json:"comments,omitempty"`
	
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

// GetHeaderValue gets a header value by name
func (r *HTTPRequest) GetHeaderValue(name string) string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers[name]
}

// SetHeaderValue sets a header value
func (r *HTTPRequest) SetHeaderValue(name, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[name] = value
}

// GetHeaderValue gets a header value by name from HTTPFileRequest
func (r *HTTPFileRequest) GetHeaderValue(name string) string {
	for _, header := range r.Headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}

// SetHeaderValue sets a header value in HTTPFileRequest
func (r *HTTPFileRequest) SetHeaderValue(name, value string) {
	// Check if header already exists
	for i, header := range r.Headers {
		if header.Name == name {
			r.Headers[i].Value = value
			return
		}
	}
	// Add new header
	r.Headers = append(r.Headers, HTTPHeader{Name: name, Value: value})
}

// Clone creates a deep copy of an HTTPRequest
func (r *HTTPRequest) Clone() *HTTPRequest {
	clone := &HTTPRequest{
		Method: r.Method,
		URL:    r.URL,
		Body:   r.Body,
		Name:   r.Name,
		Path:   r.Path,
		Tag:    r.Tag,
	}
	
	// Copy headers
	if r.Headers != nil {
		clone.Headers = make(map[string]string)
		for k, v := range r.Headers {
			clone.Headers[k] = v
		}
	}
	
	// Copy auth
	if r.Auth != nil {
		clone.Auth = &AuthDetails{
			Type:  r.Auth.Type,
			Value: r.Auth.Value,
		}
	}
	
	// Copy comments
	if len(r.Comments) > 0 {
		clone.Comments = make([]string, len(r.Comments))
		copy(clone.Comments, r.Comments)
	}
	
	// Copy form values
	if r.FormValues != nil {
		clone.FormValues = make(map[string]string)
		for k, v := range r.FormValues {
			clone.FormValues[k] = v
		}
	}
	
	// Copy query params
	if r.QueryParams != nil {
		clone.QueryParams = make(map[string]string)
		for k, v := range r.QueryParams {
			clone.QueryParams[k] = v
		}
	}
	
	return clone
}

// Clone creates a deep copy of an HTTPFileRequest
func (r *HTTPFileRequest) Clone() *HTTPFileRequest {
	clone := &HTTPFileRequest{
		Name:   r.Name,
		Method: r.Method,
		URL:    r.URL,
		Body:   r.Body,
		Tag:    r.Tag,
		Path:   r.Path,
	}
	
	// Copy headers
	if len(r.Headers) > 0 {
		clone.Headers = make([]HTTPHeader, len(r.Headers))
		copy(clone.Headers, r.Headers)
	}
	
	// Copy comments
	if len(r.Comments) > 0 {
		clone.Comments = make([]string, len(r.Comments))
		copy(clone.Comments, r.Comments)
	}
	
	return clone
}
