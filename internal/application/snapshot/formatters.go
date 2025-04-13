package snapshot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// ResponseFormatter defines the interface for formatting HTTP responses in snapshots
type ResponseFormatter interface {
	// Format converts an HTTP response to a string representation
	Format(response *models.HTTPResponse) (string, error)

	// Parse converts a string representation back to an HTTP response
	Parse(content string) (*models.HTTPResponse, error)

	// Compare compares two HTTP responses and returns a comparison result
	Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error)
}

// ComparisonResult represents the result of comparing two HTTP responses
type ComparisonResult struct {
	Matches   bool
	Diff      string
	StatusMatch bool
	HeadersMatch bool
	BodyMatch    bool
}

// formatters is a map of content types to formatters
var formatters = map[string]ResponseFormatter{
	"json":    &JSONFormatter{},
	"xml":     &XMLFormatter{},
	"text":    &TextFormatter{},
	"html":    &HTMLFormatter{},
	"binary":  &BinaryFormatter{},
	"default": &DefaultFormatter{},
}

// GetFormatter returns a formatter for the specified content type
func GetFormatter(contentType string) (ResponseFormatter, error) {
	// Clean up the content type
	contentType = strings.ToLower(contentType)
	contentType = strings.Split(contentType, ";")[0]
	contentType = strings.TrimSpace(contentType)

	// Map content type to formatter type
	formatterType := "default"
	if strings.Contains(contentType, "json") {
		formatterType = "json"
	} else if strings.Contains(contentType, "xml") {
		formatterType = "xml"
	} else if strings.Contains(contentType, "text") {
		formatterType = "text"
	} else if strings.Contains(contentType, "html") {
		formatterType = "html"
	} else if strings.Contains(contentType, "octet-stream") || 
		      strings.Contains(contentType, "application/pdf") || 
		      strings.Contains(contentType, "image/") {
		formatterType = "binary"
	}

	// Get the formatter
	formatter, ok := formatters[formatterType]
	if !ok {
		return nil, fmt.Errorf("no formatter found for content type: %s", contentType)
	}

	return formatter, nil
}

// BaseFormatter provides common functionality for all formatters
type BaseFormatter struct{}

// Format formats the response headers and metadata
func (f *BaseFormatter) formatHeaders(response *models.HTTPResponse) string {
	var sb strings.Builder

	// Add status line
	sb.WriteString(fmt.Sprintf("HTTP %d %s\n", response.StatusCode, response.Status))
	
	// Add headers
	for key, values := range response.Headers {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}

	// Add empty line to separate headers from body
	sb.WriteString("\n")

	return sb.String()
}

// parseHeaders parses the response headers and metadata
func (f *BaseFormatter) parseHeaders(content string) (*models.HTTPResponse, string, error) {
	response := &models.HTTPResponse{
		Headers: make(map[string][]string),
	}

	// Split the content into lines
	lines := strings.Split(content, "\n")
	lineIdx := 0

	// Parse the status line
	if lineIdx < len(lines) {
		statusLine := lines[lineIdx]
		lineIdx++

		if strings.HasPrefix(statusLine, "HTTP ") {
			// Parse status code and text
			parts := strings.SplitN(statusLine, " ", 3)
			if len(parts) >= 3 {
				status := parts[2]
				statusCode := 0
				fmt.Sscanf(parts[1], "%d", &statusCode)

				response.Status = status
				response.StatusCode = statusCode
			}
		}
	}

	// Parse headers
	for lineIdx < len(lines) {
		line := lines[lineIdx]
		lineIdx++

		// Empty line indicates end of headers
		if line == "" {
			break
		}

		// Parse header line
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.TrimSpace(parts[1])

			// Add to headers
			response.Headers[key] = append(response.Headers[key], value)
		}
	}

	// Return the remaining content as the body
	body := strings.Join(lines[lineIdx:], "\n")

	return response, body, nil
}

// createDiff creates a human-readable diff between two strings
func (f *BaseFormatter) createDiff(expected, actual string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected, actual, false)
	return dmp.DiffPrettyText(diffs)
}

// compareHeaders compares two sets of headers, ignoring unimportant ones
func (f *BaseFormatter) compareHeaders(expected, actual map[string][]string) bool {
	// Headers to ignore in comparison
	ignoreHeaders := map[string]bool{
		"Date":           true,
		"Content-Length": true,
		"Server":         true,
		"Connection":     true,
	}

	// Check if all expected headers are in actual
	for key, expectedValues := range expected {
		// Skip ignored headers
		if ignoreHeaders[key] {
			continue
		}

		actualValues, exists := actual[key]
		if !exists {
			return false
		}

		// Compare header values
		if !equalStringSlices(expectedValues, actualValues) {
			return false
		}
	}

	// Check if actual has additional headers not in expected
	for key := range actual {
		// Skip ignored headers
		if ignoreHeaders[key] {
			continue
		}

		_, exists := expected[key]
		if !exists {
			return false
		}
	}

	return true
}

// equalStringSlices compares two string slices for equality
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// JSONFormatter formats JSON responses
type JSONFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *JSONFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// Format JSON body for readability
	if response.Body != "" {
		var parsedJSON interface{}

		if err := json.Unmarshal([]byte(response.Body), &parsedJSON); err == nil {
			// Pretty-print the JSON
			formattedJSON, err := json.MarshalIndent(parsedJSON, "", "  ")
			if err == nil {
				sb.Write(formattedJSON)
			} else {
				// Fallback to original body if formatting fails
				sb.WriteString(response.Body)
			}
		} else {
			// Fallback to original body if parsing fails
			sb.WriteString(response.Body)
		}
	}

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *JSONFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, bodyContent, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// Set the body
	response.Body = bodyContent

	// Try to normalize JSON for consistent comparison
	if bodyContent != "" {
		var parsedJSON interface{}

		if err := json.Unmarshal([]byte(bodyContent), &parsedJSON); err == nil {
			// Normalize the JSON by re-serializing it
			normalizedJSON, err := json.Marshal(parsedJSON)
			if err == nil {
				response.Body = string(normalizedJSON)
			}
		}
	}

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *JSONFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// Compare bodies
	if expected.Body != "" || actual.Body != "" {
		// Parse both bodies as JSON for comparison
		var expectedJSON, actualJSON interface{}

		expectedErr := json.Unmarshal([]byte(expected.Body), &expectedJSON)
		actualErr := json.Unmarshal([]byte(actual.Body), &actualJSON)

		if expectedErr == nil && actualErr == nil {
			// Both are valid JSON, compare as objects
			expectedBytes, _ := json.MarshalIndent(expectedJSON, "", "  ")
			actualBytes, _ := json.MarshalIndent(actualJSON, "", "  ")

			if string(expectedBytes) == string(actualBytes) {
				result.BodyMatch = true
			} else {
				result.BodyMatch = false
				result.Matches = false
				result.Diff += "JSON body mismatch:\n"
				result.Diff += f.createDiff(string(expectedBytes), string(actualBytes))
			}
		} else {
			// Fall back to string comparison
			if expected.Body == actual.Body {
				result.BodyMatch = true
			} else {
				result.BodyMatch = false
				result.Matches = false
				result.Diff += "Body mismatch:\n"
				result.Diff += f.createDiff(expected.Body, actual.Body)
			}
		}
	} else {
		// Both bodies are empty
		result.BodyMatch = true
	}

	return result, nil
}

// XMLFormatter formats XML responses
type XMLFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *XMLFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// Add the XML body (no special formatting for now)
	sb.WriteString(response.Body)

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *XMLFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, bodyContent, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// Set the body
	response.Body = bodyContent

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *XMLFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// Compare bodies
	if expected.Body == actual.Body {
		result.BodyMatch = true
	} else {
		result.BodyMatch = false
		result.Matches = false
		result.Diff += "XML body mismatch:\n"
		result.Diff += f.createDiff(expected.Body, actual.Body)
	}

	return result, nil
}

// TextFormatter formats plain text responses
type TextFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *TextFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// Add the text body
	sb.WriteString(response.Body)

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *TextFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, bodyContent, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// Set the body
	response.Body = bodyContent

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *TextFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// Compare bodies
	if expected.Body == actual.Body {
		result.BodyMatch = true
	} else {
		result.BodyMatch = false
		result.Matches = false
		result.Diff += "Text body mismatch:\n"
		result.Diff += f.createDiff(expected.Body, actual.Body)
	}

	return result, nil
}

// HTMLFormatter formats HTML responses
type HTMLFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *HTMLFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// Add the HTML body (no special formatting for now)
	sb.WriteString(response.Body)

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *HTMLFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, bodyContent, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// Set the body
	response.Body = bodyContent

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *HTMLFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// Compare bodies
	if expected.Body == actual.Body {
		result.BodyMatch = true
	} else {
		result.BodyMatch = false
		result.Matches = false
		result.Diff += "HTML body mismatch:\n"
		result.Diff += f.createDiff(expected.Body, actual.Body)
	}

	return result, nil
}

// BinaryFormatter formats binary responses
type BinaryFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *BinaryFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// For binary data, just add a placeholder
	sb.WriteString(fmt.Sprintf("[Binary data, %d bytes]", len(response.Body)))

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *BinaryFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, _, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// For binary data, we can't restore the actual content from the snapshot
	// So we just set the body to an empty string
	response.Body = ""

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *BinaryFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// For binary data, we only compare lengths
	expectedLen := len(expected.Body)
	actualLen := len(actual.Body)

	if expectedLen == actualLen {
		result.BodyMatch = true
	} else {
		result.BodyMatch = false
		result.Matches = false
		result.Diff += fmt.Sprintf("Binary data size mismatch: expected %d bytes, got %d bytes\n", 
			expectedLen, actualLen)
	}

	return result, nil
}

// DefaultFormatter is a fallback formatter for unknown content types
type DefaultFormatter struct {
	BaseFormatter
}

// Format converts an HTTP response to a string representation
func (f *DefaultFormatter) Format(response *models.HTTPResponse) (string, error) {
	var sb strings.Builder

	// Add headers
	sb.WriteString(f.formatHeaders(response))

	// Add the body as is
	sb.WriteString(response.Body)

	return sb.String(), nil
}

// Parse converts a string representation back to an HTTP response
func (f *DefaultFormatter) Parse(content string) (*models.HTTPResponse, error) {
	response, bodyContent, err := f.parseHeaders(content)
	if err != nil {
		return nil, err
	}

	// Set the body
	response.Body = bodyContent

	return response, nil
}

// Compare compares two HTTP responses and returns a comparison result
func (f *DefaultFormatter) Compare(expected, actual *models.HTTPResponse) (*ComparisonResult, error) {
	result := &ComparisonResult{
		Matches: true,
	}

	// Compare status codes
	result.StatusMatch = expected.StatusCode == actual.StatusCode
	if !result.StatusMatch {
		result.Matches = false
		result.Diff += fmt.Sprintf("Status code mismatch: expected %d, got %d\n", 
			expected.StatusCode, actual.StatusCode)
	}

	// Compare headers
	result.HeadersMatch = f.compareHeaders(expected.Headers, actual.Headers)
	if !result.HeadersMatch {
		result.Matches = false
		result.Diff += "Headers mismatch\n"
	}

	// Compare bodies
	if expected.Body == actual.Body {
		result.BodyMatch = true
	} else {
		result.BodyMatch = false
		result.Matches = false
		result.Diff += "Body mismatch:\n"
		result.Diff += f.createDiff(expected.Body, actual.Body)
	}

	return result, nil
}