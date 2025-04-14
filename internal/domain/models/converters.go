package models

// This file provides converters between different model types to maintain
// compatibility after the model consolidation

// ToHTTPFileRequest converts an HTTPRequest to an HTTPFileRequest
func (r *HTTPRequest) ToHTTPFileRequest() *HTTPFileRequest {
	// Create a new HTTPFileRequest
	fileReq := &HTTPFileRequest{
		Name:     r.Name,
		Method:   r.Method,
		URL:      r.URL,
		Headers:  make(map[string]string),
		Body:     r.Body,
		Comments: make([]string, len(r.Comments)),
		Tag:      r.Tag,
		Path:     r.Path,
	}
	
	// Copy headers
	if r.Headers != nil {
		for k, v := range r.Headers {
			fileReq.Headers[k] = v
		}
	}
	
	// Copy comments
	copy(fileReq.Comments, r.Comments)
	
	return fileReq
}

// Helper function to split an Authorization header into type and value
func splitAuthHeader(header string) []string {
	for i := 0; i < len(header); i++ {
		if header[i] == ' ' {
			return []string{header[:i], header[i+1:]}
		}
	}
	return []string{header}
}

// GetName returns the Name field or empty string for HTTPRequest
func (r *HTTPRequest) GetName() string {
	return r.Name
}

// GetTag returns the Tag field or empty string for HTTPRequest
func (r *HTTPRequest) GetTag() string {
	return r.Tag
}

// GetPath returns the Path field or empty string for HTTPRequest
func (r *HTTPRequest) GetPath() string {
	return r.Path
}

// GetComments returns the Comments slice from HTTPRequest
func (r *HTTPRequest) GetComments() []string {
	return r.Comments
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
	return r.RequestPath
}

// GetRequestMethod gets the request method from the diff if available
func (r *SnapshotResult) GetRequestMethod() string {
	if r.Diff != nil && r.Diff.RequestMethod != "" {
		return r.Diff.RequestMethod
	}
	return r.RequestMethod
}

// ConvertHeadersToMap converts a slice of HTTPHeader to a map
func ConvertHeadersToMap(headers []HTTPHeader) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		result[header.Name] = header.Value
	}
	return result
}

// ConvertMapToHeaders converts a map to a slice of HTTPHeader
func ConvertMapToHeaders(headers map[string]string) []HTTPHeader {
	var result []HTTPHeader
	for name, value := range headers {
		result = append(result, HTTPHeader{
			Name:  name,
			Value: value,
		})
	}
	return result
}
