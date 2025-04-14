package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// VariableExtractorService implements the VariableExtractor interface
type VariableExtractorService struct{}

// NewVariableExtractorService creates a new VariableExtractorService
func NewVariableExtractorService() *VariableExtractorService {
	return &VariableExtractorService{}
}

// Extract extracts variables from an HTTP response based on extraction rules
func (s *VariableExtractorService) Extract(
	ctx context.Context,
	response *models.HTTPResponse,
	extractions []models.VariableExtraction,
) (map[string]string, error) {
	result := make(map[string]string)
	
	for _, extraction := range extractions {
		value, err := s.extractValue(response, extraction)
		if err != nil {
			if extraction.Required {
				return result, fmt.Errorf("failed to extract required variable %s: %w", extraction.Name, err)
			}
			// Use default value if provided
			if extraction.Default != "" {
				result[extraction.Name] = extraction.Default
			}
			continue
		}
		
		result[extraction.Name] = value
	}
	
	return result, nil
}

// ReplaceVariables replaces variable placeholders in a string
func (s *VariableExtractorService) ReplaceVariables(input string, variables map[string]string, format string) string {
	if input == "" || len(variables) == 0 {
		return input
	}
	
	if format == "" {
		format = "${%s}" // Default format
	}
	
	result := input
	
	for name, value := range variables {
		placeholder := fmt.Sprintf(format, name)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	return result
}

// ReplaceVariablesInRequest replaces variable placeholders in a request
func (s *VariableExtractorService) ReplaceVariablesInRequest(
	request *models.HTTPRequest,
	variables map[string]string,
	format string,
) (*models.HTTPRequest, error) {
	if request == nil || len(variables) == 0 {
		return request, nil
	}
	
	// Make a copy of the request to avoid modifying the original
	requestCopy := *request
	
	// Replace in URL
	requestCopy.URL = s.ReplaceVariables(requestCopy.URL, variables, format)
	
	// Replace in path
	requestCopy.Path = s.ReplaceVariables(requestCopy.Path, variables, format)
	
	// Replace in headers
	if requestCopy.Headers != nil {
		newHeaders := make(map[string]string)
		for name, value := range requestCopy.Headers {
			newHeaders[name] = s.ReplaceVariables(value, variables, format)
		}
		requestCopy.Headers = newHeaders
	}
	
	// Replace in body
	if len(requestCopy.Body) > 0 {
		requestCopy.Body = s.ReplaceVariables(requestCopy.Body, variables, format)
	}
	
	// Replace in form values
	if requestCopy.FormValues != nil {
		newFormValues := make(map[string]string)
		for name, value := range requestCopy.FormValues {
			newFormValues[name] = s.ReplaceVariables(value, variables, format)
		}
		requestCopy.FormValues = newFormValues
	}
	
	// Replace in query parameters
	if requestCopy.QueryParams != nil {
		newQueryParams := make(map[string]string)
		for name, value := range requestCopy.QueryParams {
			newQueryParams[name] = s.ReplaceVariables(value, variables, format)
		}
		requestCopy.QueryParams = newQueryParams
	}
	
	return &requestCopy, nil
}

// SaveVariables saves variables to a file
func (s *VariableExtractorService) SaveVariables(
	ctx context.Context,
	variables map[string]string,
	path string,
) error {
	data, err := json.MarshalIndent(variables, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}
	
	// Create parent directory if it doesn't exist
	dir := string(os.PathSeparator)
	if lastSlash := strings.LastIndex(path, string(os.PathSeparator)); lastSlash != -1 {
		dir = path[:lastSlash]
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write variables to file: %w", err)
	}
	
	return nil
}

// LoadVariables loads variables from a file
func (s *VariableExtractorService) LoadVariables(
	ctx context.Context,
	path string,
) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read variables file: %w", err)
	}
	
	var variables map[string]string
	if err := json.Unmarshal(data, &variables); err != nil {
		return nil, fmt.Errorf("failed to parse variables: %w", err)
	}
	
	return variables, nil
}

// Helper function to extract a value based on extraction rule
func (s *VariableExtractorService) extractValue(response *models.HTTPResponse, extraction models.VariableExtraction) (string, error) {
	switch strings.ToLower(extraction.Source) {
	case "body":
		return s.extractFromBody(response, extraction)
	case "header":
		return s.extractFromHeader(response, extraction)
	case "status":
		return s.extractFromStatus(response, extraction)
	default:
		return "", fmt.Errorf("unsupported extraction source: %s", extraction.Source)
	}
}

// Extract value from response body
func (s *VariableExtractorService) extractFromBody(response *models.HTTPResponse, extraction models.VariableExtraction) (string, error) {
	// Check if a JSON path is specified
	if extraction.Path != "" {
		bodyBytes := []byte(response.Body)
		return s.extractFromJsonPath(bodyBytes, extraction.Path)
	}
	
	// Check if a regular expression is specified
	if extraction.Regexp != "" {
		return s.extractWithRegexp(response.Body, extraction.Regexp)
	}
	
	// Return the entire body as a string
	return response.Body, nil
}

// Extract value from response header
func (s *VariableExtractorService) extractFromHeader(response *models.HTTPResponse, extraction models.VariableExtraction) (string, error) {
	// Path is used as the header name
	if extraction.Path == "" {
		return "", fmt.Errorf("header name (path) is required for header extraction")
	}
	
	headerValues, ok := response.Headers[extraction.Path]
	if !ok {
		return "", fmt.Errorf("header not found: %s", extraction.Path)
	}
	
	if len(headerValues) == 0 {
		return "", fmt.Errorf("header has no values: %s", extraction.Path)
	}
	
	// Use the first value
	headerValue := headerValues[0]
	
	// Extract with regexp if specified
	if extraction.Regexp != "" {
		return s.extractWithRegexp(headerValue, extraction.Regexp)
	}
	
	return headerValue, nil
}

// Extract value from response status
func (s *VariableExtractorService) extractFromStatus(response *models.HTTPResponse, extraction models.VariableExtraction) (string, error) {
	// Use the status code as a string
	return strconv.Itoa(response.StatusCode), nil
}

// Extract a value from JSON using a path
func (s *VariableExtractorService) extractFromJsonPath(data []byte, path string) (string, error) {
	// Parse the JSON
	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Navigate through the path
	segments := strings.Split(strings.Trim(path, "."), ".")
	current := parsed
	
	for _, segment := range segments {
		// Handle array indexing
		if strings.HasSuffix(segment, "]") {
			openBracket := strings.Index(segment, "[")
			if openBracket == -1 {
				return "", fmt.Errorf("invalid array index syntax: %s", segment)
			}
			
			// Get the property name and array index
			propName := segment[:openBracket]
			indexStr := segment[openBracket+1 : len(segment)-1]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return "", fmt.Errorf("invalid array index: %s", indexStr)
			}
			
			// Get the array
			obj, ok := current.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("expected object but got: %T", current)
			}
			
			arr, ok := obj[propName].([]interface{})
			if !ok {
				return "", fmt.Errorf("expected array but got: %T", obj[propName])
			}
			
			// Check array bounds
			if index < 0 || index >= len(arr) {
				return "", fmt.Errorf("array index out of bounds: %d", index)
			}
			
			current = arr[index]
		} else {
			// Handle simple object property
			obj, ok := current.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("expected object but got: %T", current)
			}
			
			value, ok := obj[segment]
			if !ok {
				return "", fmt.Errorf("property not found: %s", segment)
			}
			
			current = value
		}
	}
	
	// Convert the value to a string
	switch v := current.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	case nil:
		return "null", nil
	default:
		// For objects and arrays, return the JSON representation
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to convert value to string: %w", err)
		}
		return string(bytes), nil
	}
}

// Extract a value using a regular expression
func (s *VariableExtractorService) extractWithRegexp(input string, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression: %w", err)
	}
	
	matches := re.FindStringSubmatch(input)
	if len(matches) == 0 {
		return "", fmt.Errorf("no match found for pattern: %s", pattern)
	}
	
	// If the pattern has capture groups, use the first capture group
	// Otherwise, use the entire match
	if len(matches) > 1 {
		return matches[1], nil
	}
	
	return matches[0], nil
}
