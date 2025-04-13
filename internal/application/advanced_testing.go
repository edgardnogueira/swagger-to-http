package application

import (
	"context"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// SchemaValidator defines the interface for validating responses against schemas
type SchemaValidator interface {
	// ValidateResponse validates a response against a schema
	ValidateResponse(ctx context.Context, response *models.HTTPResponse, schemaPath string, options models.ValidationOptions) (*models.SchemaValidationResult, error)

	// ValidateResponseWithSwagger validates a response against a swagger document
	ValidateResponseWithSwagger(ctx context.Context, response *models.HTTPResponse, swaggerDoc *models.SwaggerDoc, 
		path string, method string, options models.ValidationOptions) (*models.SchemaValidationResult, error)
	
	// GetSchemaForOperation retrieves the schema for a specific operation
	GetSchemaForOperation(ctx context.Context, swaggerDoc *models.SwaggerDoc, 
		path string, method string, statusCode int) (string, error)
}

// VariableExtractor defines the interface for extracting variables from responses
type VariableExtractor interface {
	// Extract extracts variables from a response based on extraction definitions
	Extract(ctx context.Context, response *models.HTTPResponse, extractions []models.VariableExtraction) (map[string]string, error)
	
	// ReplaceVariables replaces variable placeholders in a string
	ReplaceVariables(input string, variables map[string]string, format string) string
	
	// ReplaceVariablesInRequest replaces variable placeholders in a request
	ReplaceVariablesInRequest(request *models.HTTPRequest, variables map[string]string, format string) (*models.HTTPRequest, error)
	
	// SaveVariables saves variables to a file
	SaveVariables(ctx context.Context, variables map[string]string, path string) error
	
	// LoadVariables loads variables from a file
	LoadVariables(ctx context.Context, path string) (map[string]string, error)
}

// SequenceRunner defines the interface for running test sequences
type SequenceRunner interface {
	// RunSequence runs a test sequence
	RunSequence(ctx context.Context, sequence *models.TestSequence, options models.TestRunOptions) (*models.TestSequenceResult, error)
	
	// ParseSequenceFile parses a test sequence from a file
	ParseSequenceFile(ctx context.Context, filePath string) (*models.TestSequence, error)
	
	// FindSequences finds all test sequences matching a pattern
	FindSequences(ctx context.Context, patterns []string, filter models.TestFilter) ([]*models.TestSequence, error)
}

// AssertionEvaluator defines the interface for evaluating test assertions
type AssertionEvaluator interface {
	// Evaluate evaluates a list of assertions against a response
	Evaluate(ctx context.Context, response *models.HTTPResponse, assertions []models.TestAssertion) ([]models.TestAssertionResult, error)
	
	// EvaluateAssertion evaluates a single assertion against a response
	EvaluateAssertion(ctx context.Context, response *models.HTTPResponse, assertion models.TestAssertion) (*models.TestAssertionResult, error)
}

// AdvancedTestRunner extends the basic TestRunner interface with advanced testing features
type AdvancedTestRunner interface {
	TestRunner
	
	// RunWithSchemaValidation runs tests with schema validation
	RunWithSchemaValidation(ctx context.Context, patterns []string, options models.TestRunOptions, swaggerDoc *models.SwaggerDoc) (*models.TestReport, error)
	
	// RunSequences runs all test sequences matching patterns
	RunSequences(ctx context.Context, patterns []string, options models.TestRunOptions) (*models.TestReport, error)
	
	// RunSequence runs a single test sequence
	RunSequence(ctx context.Context, sequence *models.TestSequence, options models.TestRunOptions) (*models.TestSequenceResult, error)
}
