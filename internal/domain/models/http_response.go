package models

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Status is the HTTP status text
	Status string

	// Headers contains the HTTP headers
	Headers map[string][]string

	// Body contains the response body
	Body string
}