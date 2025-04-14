package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// HTTPGenerator implements the HTTPGenerator interface
type HTTPGenerator struct {
	baseURL     string
	defaultTag   string
	indentJSON   bool
	includeAuth  bool
	authHeader   string
	authToken    string
}

// HTTPGeneratorOption represents an option for configuring the HTTP generator
type HTTPGeneratorOption func(*HTTPGenerator)

// WithBaseURL sets the base URL for requests
func WithBaseURL(baseURL string) HTTPGeneratorOption {
	return func(g *HTTPGenerator) {
		g.baseURL = baseURL
	}
}

// WithDefaultTag sets the default tag for requests without tags
func WithDefaultTag(tag string) HTTPGeneratorOption {
	return func(g *HTTPGenerator) {
		g.defaultTag = tag
	}
}

// WithIndentJSON enables or disables JSON indentation
func WithIndentJSON(indent bool) HTTPGeneratorOption {
	return func(g *HTTPGenerator) {
		g.indentJSON = indent
	}
}

// WithAuth enables or disables authentication header
func WithAuth(include bool, header, token string) HTTPGeneratorOption {
	return func(g *HTTPGenerator) {
		g.includeAuth = include
		g.authHeader = header
		g.authToken = token
	}
}

// NewHTTPGenerator creates a new HTTPGenerator with options
func NewHTTPGenerator(opts ...HTTPGeneratorOption) *HTTPGenerator {
	generator := &HTTPGenerator{
		defaultTag: "default",
		indentJSON: true,
		authHeader: "Authorization",
	}

	for _, opt := range opts {
		opt(generator)
	}

	return generator
}

// Generate generates HTTP requests from a Swagger/OpenAPI document
func (g *HTTPGenerator) Generate(ctx context.Context, doc *models.SwaggerDoc) (*models.HTTPCollection, error) {
	collection := &models.HTTPCollection{
		Directories: []models.HTTPDirectory{},
		RootFiles:   []models.HTTPFile{},
	}

	// Use servers from OpenAPI 3.0 or host+basePath from Swagger 2.0
	baseURL := g.baseURL
	if baseURL == "" {
		if len(doc.Servers) > 0 && doc.Servers[0].URL != "" {
			baseURL = doc.Servers[0].URL
		} else if doc.Host != "" {
			scheme := "https"
			if len(doc.Schemes) > 0 {
				scheme = doc.Schemes[0]
			}
			baseURL = fmt.Sprintf("%s://%s%s", scheme, doc.Host, doc.BasePath)
		}
	}

	// Create a map to organize requests by tag
	requestsByTag := make(map[string][]models.HTTPRequest)

	// Process each path and operation
	for path, pathItem := range doc.Paths {
		if pathItem.Get != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "GET", pathItem.Get)
			if err == nil {
				tag := g.getTag(pathItem.Get)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Post != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "POST", pathItem.Post)
			if err == nil {
				tag := g.getTag(pathItem.Post)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Put != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "PUT", pathItem.Put)
			if err == nil {
				tag := g.getTag(pathItem.Put)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Delete != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "DELETE", pathItem.Delete)
			if err == nil {
				tag := g.getTag(pathItem.Delete)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Patch != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "PATCH", pathItem.Patch)
			if err == nil {
				tag := g.getTag(pathItem.Patch)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Options != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "OPTIONS", pathItem.Options)
			if err == nil {
				tag := g.getTag(pathItem.Options)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
		if pathItem.Head != nil {
			req, err := g.GenerateRequest(ctx, path, &pathItem, "HEAD", pathItem.Head)
			if err == nil {
				tag := g.getTag(pathItem.Head)
				requestsByTag[tag] = append(requestsByTag[tag], *req)
			}
		}
	}

	// Create HTTP files for each tag
	for tag, requests := range requestsByTag {
		directory := models.HTTPDirectory{
			Name:  tag,
			Path:  tag,
			Files: []models.HTTPFile{},
		}

		file := models.HTTPFile{
			Filename: fmt.Sprintf("%s.http", sanitizeFilename(tag)),
			Requests: requests,
		}

		// If it's the default tag, add to root files
		if tag == g.defaultTag {
			collection.RootFiles = append(collection.RootFiles, file)
		} else {
			directory.Files = append(directory.Files, file)
			collection.Directories = append(collection.Directories, directory)
		}
	}

	return collection, nil
}

// GenerateRequest generates an HTTP request from a path and operation
func (g *HTTPGenerator) GenerateRequest(ctx context.Context, path string, pathItem *models.PathItem, method string, operation *models.Operation) (*models.HTTPRequest, error) {
	if operation == nil {
		return nil, fmt.Errorf("operation is nil for method %s and path %s", method, path)
	}

	name := operation.OperationID
	if name == "" {
		name = fmt.Sprintf("%s_%s", method, strings.ReplaceAll(path, "/", "_"))
	}

	url := g.buildURL(path)
	headersMap := g.buildHeadersMap(operation)
	body := g.buildRequestBody(operation)
	comments := g.buildComments(operation)
	tag := g.getTag(operation)

	request := &models.HTTPRequest{
		Name:     name,
		Method:   method,
		URL:      url,
		Headers:  headersMap,
		Body:     body,
		Comments: comments,
		Tag:      tag,
		Path:     path,
	}

	return request, nil
}

// buildURL builds the URL for a request
func (g *HTTPGenerator) buildURL(path string) string {
	if g.baseURL == "" {
		return path
	}
	
	baseURL := g.baseURL
	if !strings.HasSuffix(baseURL, "/") && !strings.HasPrefix(path, "/") {
		baseURL = baseURL + "/"
	}
	if strings.HasSuffix(baseURL, "/") && strings.HasPrefix(path, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	
	return baseURL + path
}

// buildHeadersMap builds the headers map for a request
func (g *HTTPGenerator) buildHeadersMap(operation *models.Operation) map[string]string {
	headers := make(map[string]string)
	
	// Add default headers
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	// Add authentication header if enabled
	if g.includeAuth && g.authToken != "" {
		headers[g.authHeader] = g.authToken
	}

	return headers
}

// buildHeaders builds the headers for a request (legacy format)
func (g *HTTPGenerator) buildHeaders(operation *models.Operation) []models.HTTPHeader {
	headers := []models.HTTPHeader{
		{Name: "Content-Type", Value: "application/json"},
		{Name: "Accept", Value: "application/json"},
	}

	// Add authentication header if enabled
	if g.includeAuth && g.authToken != "" {
		headers = append(headers, models.HTTPHeader{
			Name:  g.authHeader,
			Value: g.authToken,
		})
	}

	return headers
}

// buildRequestBody builds the body for a request
func (g *HTTPGenerator) buildRequestBody(operation *models.Operation) string {
	// For OpenAPI 3.0
	if operation.RequestBody != nil && operation.RequestBody.Content != nil {
		if mediaType, ok := operation.RequestBody.Content["application/json"]; ok && mediaType.Schema != nil {
			return g.generateExampleFromSchema(mediaType.Schema)
		}
	}

	// For Swagger 2.0
	for _, param := range operation.Parameters {
		if param.In == "body" && param.Schema != nil {
			return g.generateExampleFromSchema(param.Schema)
		}
	}

	return ""
}

// buildComments builds comments for a request
func (g *HTTPGenerator) buildComments(operation *models.Operation) []string {
	comments := []string{}

	if operation.Summary != "" {
		comments = append(comments, operation.Summary)
	}

	if operation.Description != "" {
		comments = append(comments, operation.Description)
	}

	return comments
}

// getTag gets the tag for an operation
func (g *HTTPGenerator) getTag(operation *models.Operation) string {
	if len(operation.Tags) > 0 {
		return operation.Tags[0]
	}
	return g.defaultTag
}

// generateExampleFromSchema generates an example from a schema
func (g *HTTPGenerator) generateExampleFromSchema(schema *models.Schema) string {
	if schema == nil {
		return ""
	}

	// If there's an example, use it
	if schema.Example != nil {
		jsonBytes, err := json.Marshal(schema.Example)
		if err == nil {
			return string(jsonBytes)
		}
	}

	// Generate example based on schema type
	example := g.generateExample(schema)
	
	if g.indentJSON {
		jsonBytes, err := json.MarshalIndent(example, "", "  ")
		if err == nil {
			return string(jsonBytes)
		}
	}
	
	jsonBytes, err := json.Marshal(example)
	if err == nil {
		return string(jsonBytes)
	}
	
	return ""
}

// generateExample generates an example for a schema
func (g *HTTPGenerator) generateExample(schema *models.Schema) interface{} {
	if schema == nil {
		return nil
	}

	// Handle $ref
	if schema.Ref != "" {
		// In a real implementation, we would resolve the reference
		// For simplicity, we'll return a placeholder
		return map[string]interface{}{"__ref": schema.Ref}
	}

	switch schema.Type {
	case "object":
		return g.generateObjectExample(schema)
	case "array":
		return g.generateArrayExample(schema)
	case "string":
		return g.generateStringExample(schema)
	case "integer", "number":
		return 0
	case "boolean":
		return false
	default:
		return nil
	}
}

// generateObjectExample generates an example for an object schema
func (g *HTTPGenerator) generateObjectExample(schema *models.Schema) interface{} {
	if schema.Properties == nil {
		return map[string]interface{}{}
	}

	example := map[string]interface{}{}
	for name, propSchema := range schema.Properties {
		example[name] = g.generateExample(propSchema)
	}
	return example
}

// generateArrayExample generates an example for an array schema
func (g *HTTPGenerator) generateArrayExample(schema *models.Schema) interface{} {
	if schema.Items == nil {
		return []interface{}{}
	}

	// Convert Items to Schema
	itemSchema := &models.Schema{
		Type:       schema.Items.Type,
		Format:     schema.Items.Format,
		Properties: nil, // Items doesn't have Properties directly
		Items:      schema.Items.Items,
	}

	// Just return a single example item
	return []interface{}{g.generateExample(itemSchema)}
}

// generateStringExample generates an example for a string schema
func (g *HTTPGenerator) generateStringExample(schema *models.Schema) interface{} {
	switch schema.Format {
	case "date":
		return "2025-01-01"
	case "date-time":
		return "2025-01-01T12:00:00Z"
	case "email":
		return "user@example.com"
	case "uuid":
		return "00000000-0000-0000-0000-000000000000"
	default:
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		return "string"
	}
}

// sanitizeFilename sanitizes a filename
func sanitizeFilename(name string) string {
	// Replace invalid characters with underscore
	name = strings.Map(func(r rune) rune {
		if r == ' ' || r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			return '_'
		}
		return r
	}, name)
	
	// Convert to lowercase for better cross-platform compatibility
	return strings.ToLower(name)
}

// cleanPath ensures the path has a leading slash
func cleanPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

// joinPath joins directory and filename
func joinPath(dir, file string) string {
	return filepath.Join(dir, file)
}
