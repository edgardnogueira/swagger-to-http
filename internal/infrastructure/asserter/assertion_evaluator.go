package asserter

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/extractor"
)

// AssertionEvaluatorService implements the AssertionEvaluator interface
type AssertionEvaluatorService struct {
	variableExtractor *extractor.VariableExtractorService
}

// NewAssertionEvaluatorService creates a new AssertionEvaluatorService
func NewAssertionEvaluatorService() *AssertionEvaluatorService {
	return &AssertionEvaluatorService{
		variableExtractor: extractor.NewVariableExtractorService(),
	}
}

// Evaluate evaluates a list of assertions against a response
func (s *AssertionEvaluatorService) Evaluate(
	ctx context.Context,
	response *models.HTTPResponse,
	assertions []models.TestAssertion,
) ([]models.TestAssertionResult, error) {
	results := make([]models.TestAssertionResult, 0, len(assertions))
	
	for _, assertion := range assertions {
		result, err := s.EvaluateAssertion(ctx, response, assertion)
		if err != nil {
			return results, err
		}
		
		results = append(results, *result)
	}
	
	return results, nil
}

// EvaluateAssertion evaluates a single assertion against a response
func (s *AssertionEvaluatorService) EvaluateAssertion(
	ctx context.Context,
	response *models.HTTPResponse,
	assertion models.TestAssertion,
) (*models.TestAssertionResult, error) {
	// Get the actual value to assert against
	actualValue, err := s.getValueFromResponse(response, assertion.Source, assertion.Path)
	if err != nil {
		return &models.TestAssertionResult{
			Type:        assertion.Type,
			Source:      assertion.Source,
			Path:        assertion.Path,
			Passed:      false,
			Description: fmt.Sprintf("Error extracting value: %s", err),
		}, nil
	}
	
	// Initialize the result
	result := &models.TestAssertionResult{
		Type:    assertion.Type,
		Source:  assertion.Source,
		Path:    assertion.Path,
		Actual:  actualValue,
	}
	
	// Evaluate the assertion based on its type
	switch strings.ToLower(assertion.Type) {
	case "equals":
		expected := assertion.Value
		result.Expected = expected
		equals := strings.EqualFold(actualValue, expected)
		if assertion.IgnoreCase {
			equals = strings.EqualFold(strings.ToLower(actualValue), strings.ToLower(expected))
		}
		result.Passed = (equals != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not equal '%s'", expected)
			} else {
				result.Description = fmt.Sprintf("Expected '%s', got '%s'", expected, actualValue)
			}
		}
		
	case "contains":
		expected := assertion.Value
		result.Expected = expected
		contains := strings.Contains(actualValue, expected)
		if assertion.IgnoreCase {
			contains = strings.Contains(strings.ToLower(actualValue), strings.ToLower(expected))
		}
		result.Passed = (contains != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not contain '%s'", expected)
			} else {
				result.Description = fmt.Sprintf("Expected to contain '%s', got '%s'", expected, actualValue)
			}
		}
		
	case "matches":
		pattern := assertion.Value
		result.Expected = pattern
		matched, err := regexp.MatchString(pattern, actualValue)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
		result.Passed = (matched != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not match pattern '%s'", pattern)
			} else {
				result.Description = fmt.Sprintf("Expected to match pattern '%s', got '%s'", pattern, actualValue)
			}
		}
		
	case "exists":
		// For "exists", we just check if the value was successfully extracted
		result.Passed = (actualValue != "" != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = "Expected value to not exist"
			} else {
				result.Description = "Expected value to exist"
			}
		}
		
	case "notexists":
		// "notexists" is the opposite of "exists"
		result.Passed = (actualValue == "")
		if !result.Passed {
			result.Description = fmt.Sprintf("Expected value to not exist, got '%s'", actualValue)
		}
		
	case "in":
		// Check if the actual value is in the list of expected values
		var found bool
		for _, val := range assertion.Values {
			result.Expected += val + ", "
			equals := actualValue == val
			if assertion.IgnoreCase {
				equals = strings.EqualFold(strings.ToLower(actualValue), strings.ToLower(val))
			}
			if equals {
				found = true
				break
			}
		}
		result.Expected = strings.TrimSuffix(result.Expected, ", ")
		result.Passed = (found != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not be one of [%s]", result.Expected)
			} else {
				result.Description = fmt.Sprintf("Expected value to be one of [%s], got '%s'", result.Expected, actualValue)
			}
		}
		
	case "lessthan", "lt":
		expected, err := strconv.ParseFloat(assertion.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number for comparison: %w", err)
		}
		actual, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return nil, fmt.Errorf("value is not a number: %w", err)
		}
		result.Expected = assertion.Value
		result.Passed = ((actual < expected) != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not be less than %v", expected)
			} else {
				result.Description = fmt.Sprintf("Expected value to be less than %v, got %v", expected, actual)
			}
		}
		
	case "greaterthan", "gt":
		expected, err := strconv.ParseFloat(assertion.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number for comparison: %w", err)
		}
		actual, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return nil, fmt.Errorf("value is not a number: %w", err)
		}
		result.Expected = assertion.Value
		result.Passed = ((actual > expected) != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = fmt.Sprintf("Expected value to not be greater than %v", expected)
			} else {
				result.Description = fmt.Sprintf("Expected value to be greater than %v, got %v", expected, actual)
			}
		}
		
	case "null", "nil":
		isNull := (actualValue == "null" || actualValue == "")
		result.Passed = (isNull != assertion.Not)
		if !result.Passed {
			if assertion.Not {
				result.Description = "Expected value to not be null"
			} else {
				result.Description = fmt.Sprintf("Expected null, got '%s'", actualValue)
			}
		}
		
	default:
		return nil, fmt.Errorf("unsupported assertion type: %s", assertion.Type)
	}
	
	return result, nil
}

// getValueFromResponse retrieves a value from the response based on source and path
func (s *AssertionEvaluatorService) getValueFromResponse(
	response *models.HTTPResponse,
	source string,
	path string,
) (string, error) {
	switch strings.ToLower(source) {
	case "body":
		// For body, use JSON path extraction if path is provided
		if path != "" {
			// Create a variable extraction for JSON path
			extraction := models.VariableExtraction{
				Source: "body",
				Path:   path,
			}
			return s.variableExtractor.extractFromBody(response, extraction)
		}
		// Return the entire body as string
		return string(response.Body), nil
		
	case "header":
		// Path is the header name
		if path == "" {
			return "", fmt.Errorf("header name (path) is required for header assertions")
		}
		
		headerValues, ok := response.Headers[path]
		if !ok || len(headerValues) == 0 {
			return "", nil // Header doesn't exist, will be handled by the assertion logic
		}
		
		return headerValues[0], nil
		
	case "status":
		return strconv.Itoa(response.StatusCode), nil
		
	case "contenttype":
		return response.ContentType, nil
		
	default:
		return "", fmt.Errorf("unsupported assertion source: %s", source)
	}
}

// parseBody tries to parse the response body as JSON
func (s *AssertionEvaluatorService) parseBody(body []byte) (interface{}, error) {
	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse JSON body: %w", err)
	}
	return parsed, nil
}
