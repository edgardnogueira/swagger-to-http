package models

// This file provides converters between different model types to maintain
// compatibility after the model consolidation

// ToHTTPFileRequest converts an HTTPRequest to an HTTPFileRequest
func (r *HTTPRequest) ToHTTPFileRequest() *HTTPFileRequest {
	// Convert headers map to slice
	var headers []HTTPHeader
	for name, value := range r.Headers {
		headers = append(headers, HTTPHeader{
			Name:  name,
			Value: value,
		})
	}

	return &HTTPFileRequest{
		Name:     "", // Populated by caller if needed
		Method:   r.Method,
		URL:      r.URL,
		Headers:  headers,
		Body:     r.Body,
		Comments: []string{}, // Empty by default
		Tag:      "",         // Empty by default
		Path:     "",         // Empty by default
	}
}

// ToHTTPRequest converts an HTTPFileRequest to an HTTPRequest
func (r *HTTPFileRequest) ToHTTPRequest() *HTTPRequest {
	// Convert headers slice to map
	headers := make(map[string]string)
	for _, header := range r.Headers {
		headers[header.Name] = header.Value
	}

	return &HTTPRequest{
		Method:  r.Method,
		URL:     r.URL,
		Headers: headers,
		Body:    r.Body,
		Auth:    nil, // Auth details not available in HTTPFileRequest
	}
}

// GetName returns the Name field or empty string for HTTPRequest
func (r *HTTPRequest) GetName() string {
	return "" // HTTPRequest doesn't have a Name field
}

// GetTag returns the Tag field or empty string for HTTPRequest
func (r *HTTPRequest) GetTag() string {
	return "" // HTTPRequest doesn't have a Tag field
}

// GetPath returns the Path field or empty string for HTTPRequest
func (r *HTTPRequest) GetPath() string {
	return "" // HTTPRequest doesn't have a Path field
}

// GetComments returns an empty slice as HTTPRequest doesn't have Comments
func (r *HTTPRequest) GetComments() []string {
	return []string{} // HTTPRequest doesn't have Comments field
}

// IsPassed determines if the snapshot result passed
func (r *SnapshotResult) IsPassed() bool {
	return r.Matches && !r.WasUpdated && !r.WasCreated
}

// IsUpdated determines if the snapshot was updated
func (r *SnapshotResult) IsUpdated() bool {
	return r.WasUpdated
}

// IsCreated determines if the snapshot was created
func (r *SnapshotResult) IsCreated() bool {
	return r.WasCreated
}

// GetRequestPath gets the request path from the diff if available
func (r *SnapshotResult) GetRequestPath() string {
	if r.Diff != nil && r.Diff.RequestPath != "" {
		return r.Diff.RequestPath
	}
	return ""
}

// GetRequestMethod gets the request method from the diff if available
func (r *SnapshotResult) GetRequestMethod() string {
	if r.Diff != nil && r.Diff.RequestMethod != "" {
		return r.Diff.RequestMethod
	}
	return ""
}

// FormatHTTPHeaders converts a slice of HTTPHeader to a map of strings
func FormatHTTPHeaders(headers []HTTPHeader) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		result[header.Name] = header.Value
	}
	return result
}

// ParseHTTPHeaders converts a map of strings to a slice of HTTPHeader
func ParseHTTPHeaders(headers map[string]string) []HTTPHeader {
	var result []HTTPHeader
	for name, value := range headers {
		result = append(result, HTTPHeader{
			Name:  name,
			Value: value,
		})
	}
	return result
}

// StringToBytes converts a string to a byte slice
func StringToBytes(s string) []byte {
	return []byte(s)
}

// BytesToString converts a byte slice to a string
func BytesToString(b []byte) string {
	return string(b)
}
