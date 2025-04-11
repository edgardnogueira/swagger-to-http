package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// ResponseFormatter defines the interface for formatting and comparing response bodies
type ResponseFormatter interface {
	// Format formats a response body for storage
	Format(body []byte) ([]byte, error)
	
	// Compare compares two response bodies and returns their differences
	Compare(expected, actual []byte) *models.BodyDiff
}

// JSONFormatter formats and compares JSON response bodies
type JSONFormatter struct{}

// Format formats a JSON response body by parsing and re-serializing in a normalized format
func (f *JSONFormatter) Format(body []byte) ([]byte, error) {
	var data interface{}
	
	if len(body) == 0 {
		return []byte("{}"), nil
	}
	
	// Parse JSON
	err := json.Unmarshal(body, &data)
	if err != nil {
		return body, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Re-serialize with consistent indentation
	normalized, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return body, fmt.Errorf("failed to normalize JSON: %w", err)
	}
	
	return normalized, nil
}

// Compare compares two JSON response bodies
func (f *JSONFormatter) Compare(expected, actual []byte) *models.BodyDiff {
	diff := &models.BodyDiff{
		ContentType:      "application/json",
		ExpectedSize:     len(expected),
		ActualSize:       len(actual),
		ExpectedContent:  string(expected),
		ActualContent:    string(actual),
		Equal:            bytes.Equal(expected, actual),
	}
	
	// If they're equal, no need for detailed comparison
	if diff.Equal {
		return diff
	}
	
	// Create text diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(diff.ExpectedContent, diff.ActualContent, false)
	diff.DiffContent = dmp.DiffPrettyText(diffs)
	
	// Try to parse both as JSON for structural comparison
	var expectedJSON, actualJSON interface{}
	jsonDiff := &models.JsonDiff{
		MissingFields:    []string{},
		ExtraFields:      []string{},
		DifferentTypes:   make(map[string]models.TypeDiff),
		DifferentValues:  make(map[string]models.ValueDiff),
		Equal:            false,
	}
	
	expectedErr := json.Unmarshal(expected, &expectedJSON)
	actualErr := json.Unmarshal(actual, &actualJSON)
	
	// Only do JSON-specific comparison if both are valid JSON
	if expectedErr == nil && actualErr == nil {
		compareJSON("", expectedJSON, actualJSON, jsonDiff)
		diff.JsonDiff = jsonDiff
	}
	
	return diff
}

// compareJSON recursively compares two JSON objects and records differences
func compareJSON(path string, expected, actual interface{}, diff *models.JsonDiff) {
	// Helper to handle path creation
	getPath := func(base, key string) string {
		if base == "" {
			return key
		}
		return base + "." + key
	}
	
	// Handle different types of expected and actual
	switch expectedTyped := expected.(type) {
	case map[string]interface{}:
		// Check if actual is also a map
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			diff.DifferentTypes[path] = models.TypeDiff{
				ExpectedType: "object",
				ActualType:   getTypeName(actual),
			}
			return
		}
		
		// Check for missing and different fields
		for key, expectedValue := range expectedTyped {
			fieldPath := getPath(path, key)
			
			actualValue, exists := actualMap[key]
			if !exists {
				diff.MissingFields = append(diff.MissingFields, fieldPath)
			} else {
				compareJSON(fieldPath, expectedValue, actualValue, diff)
			}
		}
		
		// Check for extra fields
		for key := range actualMap {
			fieldPath := getPath(path, key)
			if _, exists := expectedTyped[key]; !exists {
				diff.ExtraFields = append(diff.ExtraFields, fieldPath)
			}
		}
		
	case []interface{}:
		// Check if actual is also an array
		actualArray, ok := actual.([]interface{})
		if !ok {
			diff.DifferentTypes[path] = models.TypeDiff{
				ExpectedType: "array",
				ActualType:   getTypeName(actual),
			}
			return
		}
		
		// Compare array lengths
		if len(expectedTyped) != len(actualArray) {
			diff.DifferentValues[path] = models.ValueDiff{
				Expected: fmt.Sprintf("array[%d]", len(expectedTyped)),
				Actual:   fmt.Sprintf("array[%d]", len(actualArray)),
			}
		}
		
		// Compare elements up to the length of the shorter array
		minLen := len(expectedTyped)
		if len(actualArray) < minLen {
			minLen = len(actualArray)
		}
		
		for i := 0; i < minLen; i++ {
			indexPath := fmt.Sprintf("%s[%d]", path, i)
			compareJSON(indexPath, expectedTyped[i], actualArray[i], diff)
		}
		
	default:
		// For primitive values, do direct comparison
		if !areValuesEqual(expected, actual) {
			diff.DifferentValues[path] = models.ValueDiff{
				Expected: expected,
				Actual:   actual,
			}
		}
	}
}

// getTypeName returns the type name of a value as a string
func getTypeName(value interface{}) string {
	if value == nil {
		return "null"
	}
	
	switch value.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case string:
		return "string"
	case float64, float32:
		return "number"
	case int, int64, int32:
		return "number"
	case bool:
		return "boolean"
	default:
		return fmt.Sprintf("%T", value)
	}
}

// areValuesEqual compares two primitive values for equality
func areValuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	
	return a == b
}

// XMLFormatter formats and compares XML response bodies
type XMLFormatter struct{}

// Format formats an XML response body
func (f *XMLFormatter) Format(body []byte) ([]byte, error) {
	// For now, just normalize whitespace
	if len(body) == 0 {
		return []byte("</>"), nil
	}
	
	normalized := normalizeXML(body)
	return normalized, nil
}

// normalizeXML performs basic XML normalization
func normalizeXML(xml []byte) []byte {
	// This is a simplified implementation
	// A full implementation would use a proper XML parser
	content := string(xml)
	
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	
	// Remove excessive whitespace
	content = strings.TrimSpace(content)
	
	return []byte(content)
}

// Compare compares two XML response bodies
func (f *XMLFormatter) Compare(expected, actual []byte) *models.BodyDiff {
	diff := &models.BodyDiff{
		ContentType:      "application/xml",
		ExpectedSize:     len(expected),
		ActualSize:       len(actual),
		ExpectedContent:  string(expected),
		ActualContent:    string(actual),
		Equal:            bytes.Equal(expected, actual),
	}
	
	if !diff.Equal {
		// Create text diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(diff.ExpectedContent, diff.ActualContent, false)
		diff.DiffContent = dmp.DiffPrettyText(diffs)
	}
	
	return diff
}

// HTMLFormatter formats and compares HTML response bodies
type HTMLFormatter struct{}

// Format formats an HTML response body
func (f *HTMLFormatter) Format(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return []byte("<!DOCTYPE html><html></html>"), nil
	}
	
	// For now, just normalize whitespace
	return normalizeHTML(body), nil
}

// normalizeHTML performs basic HTML normalization
func normalizeHTML(html []byte) []byte {
	// This is a simplified implementation
	// A full implementation would use a proper HTML parser
	content := string(html)
	
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	
	// Remove excessive whitespace
	content = strings.TrimSpace(content)
	
	return []byte(content)
}

// Compare compares two HTML response bodies
func (f *HTMLFormatter) Compare(expected, actual []byte) *models.BodyDiff {
	diff := &models.BodyDiff{
		ContentType:      "text/html",
		ExpectedSize:     len(expected),
		ActualSize:       len(actual),
		ExpectedContent:  string(expected),
		ActualContent:    string(actual),
		Equal:            bytes.Equal(expected, actual),
	}
	
	if !diff.Equal {
		// Create text diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(diff.ExpectedContent, diff.ActualContent, false)
		diff.DiffContent = dmp.DiffPrettyText(diffs)
	}
	
	return diff
}

// TextFormatter formats and compares text response bodies
type TextFormatter struct{}

// Format formats a text response body
func (f *TextFormatter) Format(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return []byte(""), nil
	}
	
	// Normalize line endings
	content := string(body)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	
	return []byte(content), nil
}

// Compare compares two text response bodies
func (f *TextFormatter) Compare(expected, actual []byte) *models.BodyDiff {
	diff := &models.BodyDiff{
		ContentType:      "text/plain",
		ExpectedSize:     len(expected),
		ActualSize:       len(actual),
		ExpectedContent:  string(expected),
		ActualContent:    string(actual),
		Equal:            bytes.Equal(expected, actual),
	}
	
	if !diff.Equal {
		// Create text diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(diff.ExpectedContent, diff.ActualContent, false)
		diff.DiffContent = dmp.DiffPrettyText(diffs)
	}
	
	return diff
}

// BinaryFormatter formats and compares binary response bodies
type BinaryFormatter struct{}

// Format formats a binary response body
func (f *BinaryFormatter) Format(body []byte) ([]byte, error) {
	return body, nil
}

// Compare compares two binary response bodies
func (f *BinaryFormatter) Compare(expected, actual []byte) *models.BodyDiff {
	diff := &models.BodyDiff{
		ContentType:      "application/octet-stream",
		ExpectedSize:     len(expected),
		ActualSize:       len(actual),
		Equal:            bytes.Equal(expected, actual),
	}
	
	// For binary content, don't include the full content as it could be large
	// and not meaningful in text form
	diff.ExpectedContent = fmt.Sprintf("[Binary data, %d bytes]", diff.ExpectedSize)
	diff.ActualContent = fmt.Sprintf("[Binary data, %d bytes]", diff.ActualSize)
	
	if !diff.Equal {
		// Show a basic hex diff for small binary files
		if len(expected) <= 1024 && len(actual) <= 1024 {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(
				fmt.Sprintf("%X", expected),
				fmt.Sprintf("%X", actual),
				false,
			)
			diff.DiffContent = dmp.DiffPrettyText(diffs)
		} else {
			diff.DiffContent = "Binary content differs (sizes: expected=" + 
				fmt.Sprintf("%d", diff.ExpectedSize) + 
				" actual=" + 
				fmt.Sprintf("%d", diff.ActualSize) + ")"
		}
	}
	
	return diff
}
