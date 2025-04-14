package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/asserter"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/extractor"
)

// SequenceRunnerService implements the SequenceRunner interface
type SequenceRunnerService struct {
	httpExecutor       application.HTTPExecutor
	variableExtractor  *extractor.VariableExtractorService
	assertionEvaluator *asserter.AssertionEvaluatorService
	schemaValidator    application.SchemaValidator
}

// NewSequenceRunnerService creates a new SequenceRunnerService
func NewSequenceRunnerService(
	httpExecutor application.HTTPExecutor,
	schemaValidator application.SchemaValidator,
) *SequenceRunnerService {
	return &SequenceRunnerService{
		httpExecutor:       httpExecutor,
		variableExtractor:  extractor.NewVariableExtractorService(),
		assertionEvaluator: asserter.NewAssertionEvaluatorService(),
		schemaValidator:    schemaValidator,
	}
}

// RunSequence runs a test sequence
func (s *SequenceRunnerService) RunSequence(
	ctx context.Context,
	sequence *models.TestSequence,
	options models.TestRunOptions,
) (*models.TestSequenceResult, error) {
	// Initialize test sequence result
	result := &models.TestSequenceResult{
		Name:        sequence.Name,
		Success:     true,
		StepResults: make([]models.TestSequenceStepResult, 0, len(sequence.Steps)),
		StartTime:   time.Now(),
		Variables:   make(map[string]string),
	}
	
	// Copy initial variables from sequence and options
	for k, v := range sequence.Variables {
		result.Variables[k] = v
	}
	for k, v := range options.EnvironmentVars {
		result.Variables[k] = v
	}
	
	// Run each step in the sequence
	for _, step := range sequence.Steps {
		// Check if step should be skipped
		if step.Skip {
			result.StepResults = append(result.StepResults, models.TestSequenceStepResult{
				Name:   step.Name,
				Status: models.TestStatusSkipped,
			})
			continue
		}
		
		// Check skip condition if provided
		if step.SkipCondition != "" {
			// Evaluate the skip condition (basic implementation - supports ${var} == "value" syntax)
			skipCondition := s.variableExtractor.ReplaceVariables(step.SkipCondition, result.Variables, "${%s}")
			if s.evaluateSkipCondition(skipCondition) {
				result.StepResults = append(result.StepResults, models.TestSequenceStepResult{
					Name:   step.Name,
					Status: models.TestStatusSkipped,
					Error:  fmt.Sprintf("Skipped due to condition: %s", step.SkipCondition),
				})
				continue
			}
		}
		
		// Wait before step if specified
		if step.WaitBefore > 0 {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(step.WaitBefore):
				// Continue execution after waiting
			}
		}
		
		// Create a copy of the request with variables replaced
		requestWithVars, err := s.variableExtractor.ReplaceVariablesInRequest(
			step.Request,
			result.Variables,
			"${%s}",
		)
		if err != nil {
			result.StepResults = append(result.StepResults, models.TestSequenceStepResult{
				Name:   step.Name,
				Status: models.TestStatusError,
				Error:  fmt.Sprintf("Error replacing variables in request: %v", err),
			})
			result.Success = false
			if options.FailFast || step.StopOnFail {
				break
			}
			continue
		}
		
		// Execute the request
		startTime := time.Now()
		response, err := s.httpExecutor.Execute(ctx, requestWithVars, result.Variables)
		executionTime := time.Since(startTime)
		
		// Initialize step result
		stepResult := models.TestSequenceStepResult{
			Name:          step.Name,
			ExecutionTime: executionTime,
			Variables:     make(map[string]string),
		}
		
		// Handle request execution error
		if err != nil {
			stepResult.Status = models.TestStatusError
			stepResult.Error = fmt.Sprintf("Error executing request: %v", err)
			result.StepResults = append(result.StepResults, stepResult)
			result.Success = false
			if options.FailFast || step.StopOnFail {
				break
			}
			continue
		}
		
		// Store the response
		stepResult.Response = response
		
		// Check expected status code if specified
		if step.ExpectedStatus != 0 && response.StatusCode != step.ExpectedStatus {
			stepResult.Status = models.TestStatusFailed
			stepResult.Error = fmt.Sprintf(
				"Expected status code %d but got %d",
				step.ExpectedStatus,
				response.StatusCode,
			)
			result.StepResults = append(result.StepResults, stepResult)
			result.Success = false
			if options.FailFast || step.StopOnFail {
				break
			}
			continue
		}
		
		// Validate schema if enabled
		if (options.ValidateSchema || step.SchemaValidate) && s.schemaValidator != nil {
			// Check if a SwaggerDoc is available from validation options
			var swaggerDoc *models.SwaggerDoc
			
			// Initialize validation options if needed
			validationOptions := options.ValidationOptions
			
			if swaggerDoc != nil {
				// Validate response against swagger schema
				validationResult, err := s.schemaValidator.ValidateResponseWithSwagger(
					ctx,
					response,
					swaggerDoc,
					requestWithVars.Path,
					requestWithVars.Method,
					validationOptions,
				)
				
				if err != nil {
					stepResult.ValidationError = fmt.Sprintf("Schema validation error: %v", err)
				} else {
					stepResult.SchemaResult = validationResult
					if !validationResult.Valid {
						stepResult.ValidationError = fmt.Sprintf(
							"Schema validation failed with %d errors",
							len(validationResult.Errors),
						)
					}
				}
				
				// If schema validation failed and we need to stop on failure
				if stepResult.ValidationError != "" && (options.FailFast || step.StopOnFail) {
					stepResult.Status = models.TestStatusFailed
					result.StepResults = append(result.StepResults, stepResult)
					result.Success = false
					break
				}
			}
		}
		
		// Evaluate assertions if provided
		if len(step.Assertions) > 0 {
			assertionResults, err := s.assertionEvaluator.Evaluate(ctx, response, step.Assertions)
			if err != nil {
				stepResult.Status = models.TestStatusError
				stepResult.Error = fmt.Sprintf("Error evaluating assertions: %v", err)
				result.StepResults = append(result.StepResults, stepResult)
				result.Success = false
				if options.FailFast || step.StopOnFail {
					break
				}
				continue
			}
			
			stepResult.AssertionResults = assertionResults
			
			// Check if any assertions failed
			for _, assertionResult := range assertionResults {
				if !assertionResult.Passed {
					stepResult.Status = models.TestStatusFailed
					stepResult.Error = fmt.Sprintf(
						"Assertion failed: %s - %s",
						assertionResult.Type,
						assertionResult.Description,
					)
					result.Success = false
					if options.FailFast || step.StopOnFail {
						result.StepResults = append(result.StepResults, stepResult)
						break
					}
				}
			}
			
			// If we broke out of the loop above due to a failed assertion, break the outer loop
			if stepResult.Status == models.TestStatusFailed && (options.FailFast || step.StopOnFail) {
				break
			}
		}
		
		// Extract variables if provided
		if len(step.Variables) > 0 {
			extractedVars, err := s.variableExtractor.Extract(ctx, response, step.Variables)
			if err != nil {
				stepResult.Status = models.TestStatusError
				stepResult.Error = fmt.Sprintf("Error extracting variables: %v", err)
				result.StepResults = append(result.StepResults, stepResult)
				result.Success = false
				if options.FailFast || step.StopOnFail {
					break
				}
				continue
			}
			
			// Store extracted variables
			for k, v := range extractedVars {
				result.Variables[k] = v
				stepResult.Variables[k] = v
			}
		}
		
		// If step status wasn't set by assertions or errors, it passed
		if stepResult.Status == "" {
			stepResult.Status = models.TestStatusPassed
		}
		
		// Add step result to the sequence result
		result.StepResults = append(result.StepResults, stepResult)
		
		// Wait after step if specified
		if step.WaitAfter > 0 {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(step.WaitAfter):
				// Continue execution after waiting
			}
		}
	}
	
	// Set end time and calculate execution time
	result.EndTime = time.Now()
	result.ExecutionTime = result.EndTime.Sub(result.StartTime)
	
	// Save variables to file if configured
	if options.SaveVariables && options.VariablesPath != "" {
		if err := s.variableExtractor.SaveVariables(ctx, result.Variables, options.VariablesPath); err != nil {
			// Log error but continue
			fmt.Printf("Error saving variables: %v\n", err)
		}
	}
	
	return result, nil
}

// ParseSequenceFile parses a test sequence from a file
func (s *SequenceRunnerService) ParseSequenceFile(
	ctx context.Context,
	filePath string,
) (*models.TestSequence, error) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sequence file: %w", err)
	}
	
	// Parse the JSON
	var sequence models.TestSequence
	if err := json.Unmarshal(data, &sequence); err != nil {
		return nil, fmt.Errorf("failed to parse sequence file: %w", err)
	}
	
	// Set the file path
	sequence.FilePath = filePath
	
	// Set default values for steps if needed
	for i := range sequence.Steps {
		if sequence.Steps[i].Variables == nil {
			sequence.Steps[i].Variables = make([]models.VariableExtraction, 0)
		}
		if sequence.Steps[i].Assertions == nil {
			sequence.Steps[i].Assertions = make([]models.TestAssertion, 0)
		}
	}
	
	return &sequence, nil
}

// FindSequences finds all test sequences matching the provided patterns
func (s *SequenceRunnerService) FindSequences(
	ctx context.Context,
	patterns []string,
	filter models.TestFilter,
) ([]*models.TestSequence, error) {
	var sequences []*models.TestSequence
	
	for _, pattern := range patterns {
		// Find files matching the pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %s: %w", pattern, err)
		}
		
		for _, match := range matches {
			// Only process JSON files
			if !strings.HasSuffix(match, ".json") {
				continue
			}
			
			// Parse the sequence file
			sequence, err := s.ParseSequenceFile(ctx, match)
			if err != nil {
				// Log error but continue with other files
				fmt.Printf("Error parsing sequence file %s: %v\n", match, err)
				continue
			}
			
			// Apply filter
			if s.matchesFilter(sequence, filter) {
				sequences = append(sequences, sequence)
			}
		}
	}
	
	return sequences, nil
}

// matchesFilter checks if a sequence matches the filter criteria
func (s *SequenceRunnerService) matchesFilter(sequence *models.TestSequence, filter models.TestFilter) bool {
	// Filter by tags
	if len(filter.Tags) > 0 {
		tagMatch := false
		for _, tag := range filter.Tags {
			for _, sequenceTag := range sequence.Tags {
				if tag == sequenceTag {
					tagMatch = true
					break
				}
			}
			if tagMatch {
				break
			}
		}
		if !tagMatch {
			return false
		}
	}
	
	// Filter by name
	if len(filter.Names) > 0 {
		nameMatch := false
		for _, name := range filter.Names {
			if strings.Contains(sequence.Name, name) {
				nameMatch = true
				break
			}
		}
		if !nameMatch {
			return false
		}
	}
	
	// Filter by metadata
	if len(filter.Metadata) > 0 {
		for key, value := range filter.Metadata {
			if sequenceValue, ok := sequence.Metadata[key]; !ok || sequenceValue != value {
				return false
			}
		}
	}
	
	return true
}

// evaluateSkipCondition evaluates a skip condition expression
// Supports basic syntax: "value1 == value2", "value1 != value2"
func (s *SequenceRunnerService) evaluateSkipCondition(condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// Check for equality
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			
			// Remove quotes if present
			right = strings.Trim(right, "\"'")
			
			return left == right
		}
	}
	
	// Check for inequality
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			
			// Remove quotes if present
			right = strings.Trim(right, "\"'")
			
			return left != right
		}
	}
	
	// By default, return false (don't skip)
	return false
}
