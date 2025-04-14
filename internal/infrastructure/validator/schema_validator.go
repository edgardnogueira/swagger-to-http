package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// SchemaValidatorService implements the SchemaValidator interface
type SchemaValidatorService struct{}

// NewSchemaValidatorService creates a new SchemaValidatorService
func NewSchemaValidatorService() *SchemaValidatorService {
	return &SchemaValidatorService{}
}

// ValidateResponse validates a response against a schema
func (s *SchemaValidatorService) ValidateResponse(
	ctx context.Context, 
	response *models.HTTPResponse, 
	schemaPath string, 
	options models.ValidationOptions,
) (*models.SchemaValidationResult, error) {
	// Load the schema
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Parse the schema
	var schema map[string]interface{}
	if err := json.Unmarshal(schemaBytes, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Parse the response body
	var responseBody interface{}
	if err := json.Unmarshal([]byte(response.Body), &responseBody); err != nil {
		return &models.SchemaValidationResult{
			Valid:          false,
			Errors:         []models.ValidationError{{
				Path:    "",
				Message: fmt.Sprintf("invalid JSON response: %v", err),
			}},
			SchemaPath:     schemaPath,
			ResponseStatus: response.StatusCode,
			ContentType:    response.ContentType,
		}, nil
	}

	// Validate the response against the schema
	result := s.validateAgainstSchema(responseBody, schema, options, "")
	
	return &models.SchemaValidationResult{
		Valid:          len(result) == 0,
		Errors:         result,
		SchemaPath:     schemaPath,
		ResponseStatus: response.StatusCode,
		ContentType:    response.ContentType,
	}, nil
}

// ValidateResponseWithSwagger validates a response against a swagger document
func (s *SchemaValidatorService) ValidateResponseWithSwagger(
	ctx context.Context, 
	response *models.HTTPResponse, 
	swaggerDoc *models.SwaggerDoc, 
	path string, 
	method string, 
	options models.ValidationOptions,
) (*models.SchemaValidationResult, error) {
	// Get the schema for the response
	schemaJson, err := s.GetSchemaForOperation(ctx, swaggerDoc, path, method, response.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema for operation: %w", err)
	}

	// Parse the schema
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJson), &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Parse the response body
	var responseBody interface{}
	if err := json.Unmarshal([]byte(response.Body), &responseBody); err != nil {
		return &models.SchemaValidationResult{
			Valid:          false,
			Errors:         []models.ValidationError{{
				Path:    "",
				Message: fmt.Sprintf("invalid JSON response: %v", err),
			}},
			SchemaPath:     fmt.Sprintf("%s %s - %d", method, path, response.StatusCode),
			ResponseStatus: response.StatusCode,
			ContentType:    response.ContentType,
		}, nil
	}

	// Validate the response against the schema
	result := s.validateAgainstSchema(responseBody, schema, options, "")
	
	return &models.SchemaValidationResult{
		Valid:          len(result) == 0,
		Errors:         result,
		SchemaPath:     fmt.Sprintf("%s %s - %d", method, path, response.StatusCode),
		ResponseStatus: response.StatusCode,
		ContentType:    response.ContentType,
	}, nil
}

// GetSchemaForOperation retrieves the schema for a specific operation
func (s *SchemaValidatorService) GetSchemaForOperation(
	ctx context.Context, 
	swaggerDoc *models.SwaggerDoc, 
	path string, 
	method string, 
	statusCode int,
) (string, error) {
	// Find the path in the swagger document
	pathItem, ok := swaggerDoc.Paths[path]
	if !ok {
		// Try to match with path parameters
		for swaggerPath, item := range swaggerDoc.Paths {
			// Convert swagger path parameters to regex pattern
			pattern := swaggerPath
			pattern = regexp.MustCompile(`\{[^/]+\}`).ReplaceAllString(pattern, `[^/]+`)
			pattern = fmt.Sprintf("^%s$", pattern)
			
			if regexp.MustCompile(pattern).MatchString(path) {
				pathItem = item
				ok = true
				break
			}
		}
		
		if !ok {
			return "", fmt.Errorf("path not found in swagger document: %s", path)
		}
	}

	// Find the operation for the method
	var operation *models.Operation
	switch strings.ToLower(method) {
	case "get":
		operation = pathItem.Get
	case "post":
		operation = pathItem.Post
	case "put":
		operation = pathItem.Put
	case "delete":
		operation = pathItem.Delete
	case "options":
		operation = pathItem.Options
	case "head":
		operation = pathItem.Head
	case "patch":
		operation = pathItem.Patch
	default:
		return "", fmt.Errorf("unsupported HTTP method: %s", method)
	}
	
	if operation == nil {
		return "", fmt.Errorf("operation not found for method: %s", method)
	}

	// Find the response for the status code
	response, ok := operation.Responses[fmt.Sprintf("%d", statusCode)]
	if !ok {
		// Try to find default response
		response, ok = operation.Responses["default"]
		if !ok {
			return "", fmt.Errorf("response not found for status code: %d", statusCode)
		}
	}

	// Extract the schema
	if response.Schema == nil {
		return "", fmt.Errorf("schema not found for response")
	}

	// Convert the schema to JSON
	schemaJson, err := json.Marshal(response.Schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(schemaJson), nil
}

// validateAgainstSchema validates data against a schema and returns validation errors
func (s *SchemaValidatorService) validateAgainstSchema(
	data interface{}, 
	schema map[string]interface{}, 
	options models.ValidationOptions,
	path string,
) []models.ValidationError {
	var errors []models.ValidationError

	// This is a simplified schema validation implementation
	// In a real implementation, you would use a proper JSON Schema validator
	// like github.com/xeipuuv/gojsonschema
	
	// For now, we'll just do some basic type validation based on the schema type
	schemaType, _ := schema["type"].(string)
	
	switch schemaType {
	case "object":
		// Check if data is an object
		dataObj, ok := data.(map[string]interface{})
		if !ok {
			errors = append(errors, models.ValidationError{
				Path:    path,
				Message: "expected object but got different type",
				Value:   fmt.Sprintf("%v", data),
				Schema:  fmt.Sprintf("%v", schema),
			})
			return errors
		}
		
		// Check required properties
		if required, ok := schema["required"].([]interface{}); ok && !options.RequiredPropertiesOnly {
			for _, req := range required {
				propName, _ := req.(string)
				if _, ok := dataObj[propName]; !ok {
					errors = append(errors, models.ValidationError{
						Path:    joinPath(path, propName),
						Message: "required property missing",
						Schema:  fmt.Sprintf("required: %v", required),
					})
				}
			}
		}
		
		// Check properties
		if properties, ok := schema["properties"].(map[string]interface{}); ok {
			for propName, propSchema := range properties {
				// Skip if property is in ignored list
				if containsString(options.IgnoredProperties, propName) {
					continue
				}
				
				// If property exists in data, validate it
				if propValue, ok := dataObj[propName]; ok {
					propSchemaMap, _ := propSchema.(map[string]interface{})
					subErrors := s.validateAgainstSchema(propValue, propSchemaMap, options, joinPath(path, propName))
					errors = append(errors, subErrors...)
				} else if isRequired(schema, propName) && !options.RequiredPropertiesOnly {
					// Property is required but missing
					errors = append(errors, models.ValidationError{
						Path:    joinPath(path, propName),
						Message: "required property missing",
						Schema:  fmt.Sprintf("%v", propSchema),
					})
				}
			}
		}
		
		// Check for additional properties if not allowed
		if additionalProps, ok := schema["additionalProperties"]; !options.IgnoreAdditionalProperties && ok {
			// If additionalProperties is false, check for undefined properties
			if allow, ok := additionalProps.(bool); ok && !allow {
				properties, _ := schema["properties"].(map[string]interface{})
				for propName := range dataObj {
					if _, ok := properties[propName]; !ok {
						errors = append(errors, models.ValidationError{
							Path:    joinPath(path, propName),
							Message: "additional property not allowed",
							Value:   fmt.Sprintf("%v", dataObj[propName]),
						})
					}
				}
			}
		}
		
	case "array":
		// Check if data is an array
		dataArr, ok := data.([]interface{})
		if !ok {
			errors = append(errors, models.ValidationError{
				Path:    path,
				Message: "expected array but got different type",
				Value:   fmt.Sprintf("%v", data),
				Schema:  fmt.Sprintf("%v", schema),
			})
			return errors
		}
		
		// Validate array items
		if items, ok := schema["items"].(map[string]interface{}); ok {
			for i, item := range dataArr {
				itemPath := fmt.Sprintf("%s[%d]", path, i)
				subErrors := s.validateAgainstSchema(item, items, options, itemPath)
				errors = append(errors, subErrors...)
			}
		}
		
	case "string":
		// Check if data is a string
		dataStr, ok := data.(string)
		if !ok {
			errors = append(errors, models.ValidationError{
				Path:    path,
				Message: "expected string but got different type",
				Value:   fmt.Sprintf("%v", data),
				Schema:  fmt.Sprintf("%v", schema),
			})
			return errors
		}
		
		// Check pattern if defined and pattern validation is enabled
		if pattern, ok := schema["pattern"].(string); ok && !options.IgnorePatterns {
			matched, err := regexp.MatchString(pattern, dataStr)
			if err != nil || !matched {
				errors = append(errors, models.ValidationError{
					Path:    path,
					Message: "string does not match pattern",
					Value:   dataStr,
					Schema:  fmt.Sprintf("pattern: %s", pattern),
				})
			}
		}
		
	case "number", "integer":
		// Check if data is a number
		_, ok1 := data.(float64)
		_, ok2 := data.(int)
		if !ok1 && !ok2 {
			errors = append(errors, models.ValidationError{
				Path:    path,
				Message: fmt.Sprintf("expected %s but got different type", schemaType),
				Value:   fmt.Sprintf("%v", data),
				Schema:  fmt.Sprintf("%v", schema),
			})
		}
		
	case "boolean":
		// Check if data is a boolean
		_, ok := data.(bool)
		if !ok {
			errors = append(errors, models.ValidationError{
				Path:    path,
				Message: "expected boolean but got different type",
				Value:   fmt.Sprintf("%v", data),
				Schema:  fmt.Sprintf("%v", schema),
			})
		}
	}
	
	return errors
}

// Helper functions

// isRequired checks if a property is required
func isRequired(schema map[string]interface{}, propertyName string) bool {
	required, ok := schema["required"].([]interface{})
	if !ok {
		return false
	}
	
	for _, req := range required {
		if req.(string) == propertyName {
			return true
		}
	}
	
	return false
}

// joinPath joins path segments
func joinPath(base, property string) string {
	if base == "" {
		return property
	}
	return fmt.Sprintf("%s.%s", base, property)
}

// containsString checks if a string slice contains a string
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
