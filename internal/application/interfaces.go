package application

import (
	"context"
	
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// SwaggerParser defines the interface for parsing Swagger/OpenAPI documents
type SwaggerParser interface {
	// Parse parses a Swagger/OpenAPI document from a byte array
	Parse(ctx context.Context, data []byte) (*models.SwaggerDoc, error)
	
	// ParseFile parses a Swagger/OpenAPI document from a file
	ParseFile(ctx context.Context, filePath string) (*models.SwaggerDoc, error)
	
	// ParseURL parses a Swagger/OpenAPI document from a URL
	ParseURL(ctx context.Context, url string) (*models.SwaggerDoc, error)
	
	// Validate validates a Swagger/OpenAPI document
	Validate(ctx context.Context, doc *models.SwaggerDoc) error
}

// HTTPGenerator defines the interface for generating HTTP requests
type HTTPGenerator interface {
	// Generate generates HTTP requests from a Swagger/OpenAPI document
	Generate(ctx context.Context, doc *models.SwaggerDoc) (*models.HTTPCollection, error)
	
	// GenerateRequest generates an HTTP request from a path and operation
	GenerateRequest(ctx context.Context, path string, pathItem *models.PathItem, method string, operation *models.Operation) (*models.HTTPRequest, error)
}

// FileWriter defines the interface for writing HTTP files
type FileWriter interface {
	// WriteCollection writes an HTTP collection to the file system
	WriteCollection(ctx context.Context, collection *models.HTTPCollection) error
	
	// WriteFile writes an HTTP file to the file system
	WriteFile(ctx context.Context, file *models.HTTPFile, dirPath string) error
}

// HTTPExecutor defines the interface for executing HTTP requests
type HTTPExecutor interface {
	// Execute executes an HTTP request and returns the response
	Execute(ctx context.Context, request *models.HTTPRequest, variables map[string]string) (*models.HTTPResponse, error)
	
	// ExecuteFile executes all requests in an HTTP file
	ExecuteFile(ctx context.Context, file *models.HTTPFile, variables map[string]string) ([]*models.HTTPResponse, error)
}

// SnapshotManager defines the interface for managing response snapshots
type SnapshotManager interface {
	// SaveSnapshot saves a response as a snapshot
	SaveSnapshot(ctx context.Context, response *models.HTTPResponse, path string) error
	
	// LoadSnapshot loads a snapshot from the file system
	LoadSnapshot(ctx context.Context, path string) (*models.HTTPResponse, error)
	
	// CompareWithSnapshot compares a response with a snapshot
	CompareWithSnapshot(ctx context.Context, response *models.HTTPResponse, snapshotPath string) (*models.SnapshotDiff, error)
}

// ConfigProvider defines the interface for retrieving configuration
type ConfigProvider interface {
	// GetString retrieves a string configuration value
	GetString(key string) string
	
	// GetInt retrieves an integer configuration value
	GetInt(key string) int
	
	// GetBool retrieves a boolean configuration value
	GetBool(key string) bool
	
	// GetStringMap retrieves a string map configuration value
	GetStringMap(key string) map[string]interface{}
	
	// GetStringSlice retrieves a string slice configuration value
	GetStringSlice(key string) []string
}
