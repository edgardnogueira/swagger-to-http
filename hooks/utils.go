package hooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"gopkg.in/yaml.v3"
)

// SwaggerDiff contains information about changes between two Swagger files
type SwaggerDiff struct {
	ChangedPaths  map[string]bool
	ChangedTags   map[string]bool
	IsNewFile     bool
	HasStructural bool
}

// DetectChanges compares the current and previous versions of a Swagger file
// and returns information about the changes
func DetectChanges(currentPath, previousPath string) (*SwaggerDiff, error) {
	diff := &SwaggerDiff{
		ChangedPaths:  make(map[string]bool),
		ChangedTags:   make(map[string]bool),
		IsNewFile:     false,
		HasStructural: false,
	}

	// Check if previous version exists
	_, err := os.Stat(previousPath)
	if os.IsNotExist(err) {
		// This is a new file
		diff.IsNewFile = true
		
		// Load the current file
		current, err := loadSwaggerFile(currentPath)
		if err != nil {
			return nil, err
		}
		
		// For new files, all paths and tags are considered changed
		for path := range current.Paths {
			diff.ChangedPaths[path] = true
		}
		
		// Extract tags from the operations
		for _, pathItem := range current.Paths {
			for _, operation := range pathItem.Operations {
				for _, tag := range operation.Tags {
					diff.ChangedTags[tag] = true
				}
			}
		}
		
		diff.HasStructural = true
		return diff, nil
	} else if err != nil {
		return nil, fmt.Errorf("error checking previous file: %w", err)
	}

	// Load both files
	current, err := loadSwaggerFile(currentPath)
	if err != nil {
		return nil, fmt.Errorf("error loading current file: %w", err)
	}

	previous, err := loadSwaggerFile(previousPath)
	if err != nil {
		return nil, fmt.Errorf("error loading previous file: %w", err)
	}

	// Check for structural changes
	if current.SwaggerVersion != previous.SwaggerVersion ||
		current.Info.Title != previous.Info.Title ||
		current.Info.Version != previous.Info.Version {
		diff.HasStructural = true
	}

	// Compare paths
	for path, currentPathItem := range current.Paths {
		previousPathItem, exists := previous.Paths[path]
		if !exists {
			// New path
			diff.ChangedPaths[path] = true
			
			// Add tags from the operations
			for _, operation := range currentPathItem.Operations {
				for _, tag := range operation.Tags {
					diff.ChangedTags[tag] = true
				}
			}
			continue
		}

		// Compare operations in the path
		if !operationsEqual(currentPathItem.Operations, previousPathItem.Operations) {
			diff.ChangedPaths[path] = true
			
			// Add tags from the changed operations
			for _, operation := range currentPathItem.Operations {
				for _, tag := range operation.Tags {
					diff.ChangedTags[tag] = true
				}
			}
		}
	}

	// Check for removed paths
	for path := range previous.Paths {
		_, exists := current.Paths[path]
		if !exists {
			// Path was removed
			diff.HasStructural = true
			break
		}
	}

	return diff, nil
}

// operationsEqual compares two sets of operations for equality
func operationsEqual(a, b []*models.Operation) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	aMap := make(map[string]*models.Operation)
	for _, op := range a {
		aMap[op.Method] = op
	}

	bMap := make(map[string]*models.Operation)
	for _, op := range b {
		bMap[op.Method] = op
	}

	// Check that all operations in a exist in b and are equal
	for method, opA := range aMap {
		opB, exists := bMap[method]
		if !exists {
			return false
		}

		// Compare operation details
		if opA.OperationID != opB.OperationID ||
			!tagsEqual(opA.Tags, opB.Tags) ||
			!parametersEqual(opA.Parameters, opB.Parameters) ||
			!bodyEqual(opA.RequestBody, opB.RequestBody) ||
			!responsesEqual(opA.Responses, opB.Responses) {
			return false
		}
	}

	// Check that all operations in b exist in a
	for method := range bMap {
		_, exists := aMap[method]
		if !exists {
			return false
		}
	}

	return true
}

// tagsEqual compares two sets of tags for equality
func tagsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	aMap := make(map[string]bool)
	for _, tag := range a {
		aMap[tag] = true
	}

	// Check that all tags in b exist in a
	for _, tag := range b {
		if !aMap[tag] {
			return false
		}
	}

	return true
}

// parametersEqual compares two sets of parameters for equality
func parametersEqual(a, b []*models.Parameter) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	aMap := make(map[string]*models.Parameter)
	for _, param := range a {
		aMap[param.Name] = param
	}

	// Check that all parameters in b exist in a and are equal
	for _, paramB := range b {
		paramA, exists := aMap[paramB.Name]
		if !exists {
			return false
		}

		// Compare parameter details
		if paramA.In != paramB.In ||
			paramA.Required != paramB.Required ||
			paramA.Type != paramB.Type {
			return false
		}
	}

	return true
}

// bodyEqual compares two request bodies for equality
func bodyEqual(a, b *models.RequestBody) bool {
	// If both are nil, they're equal
	if a == nil && b == nil {
		return true
	}

	// If one is nil and the other isn't, they're not equal
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	// Compare request body details
	return a.Required == b.Required &&
		contentEqual(a.Content, b.Content)
}

// contentEqual compares two content maps for equality
func contentEqual(a, b map[string]*models.MediaType) bool {
	if len(a) != len(b) {
		return false
	}

	// Check that all content types in a exist in b and are equal
	for contentType, mediaTypeA := range a {
		mediaTypeB, exists := b[contentType]
		if !exists {
			return false
		}

		// Compare media type details (simplified)
		if mediaTypeA.Schema != nil || mediaTypeB.Schema != nil {
			// If schemas exist, they should be equal
			// In a real implementation, you would do a deeper comparison of schemas
			if (mediaTypeA.Schema == nil && mediaTypeB.Schema != nil) ||
				(mediaTypeA.Schema != nil && mediaTypeB.Schema == nil) {
				return false
			}
		}
	}

	return true
}

// responsesEqual compares two response maps for equality
func responsesEqual(a, b map[string]*models.Response) bool {
	if len(a) != len(b) {
		return false
	}

	// Check that all response codes in a exist in b and are equal
	for code, responseA := range a {
		responseB, exists := b[code]
		if !exists {
			return false
		}

		// Compare response details (simplified)
		if responseA.Description != responseB.Description {
			return false
		}

		// In a real implementation, you would compare content as well
	}

	return true
}

// loadSwaggerFile loads a Swagger file (JSON or YAML) and returns a parsed representation
func loadSwaggerFile(path string) (*models.SwaggerDoc, error) {
	// Read the file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Parse the file based on extension
	var doc models.SwaggerDoc
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".json" {
		// Parse JSON
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %w", err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		// Parse YAML
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return &doc, nil
}
