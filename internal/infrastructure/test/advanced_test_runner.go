package test

import (
	"context"
	"fmt"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/asserter"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/extractor"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/sequencer"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/validator"
)

// AdvancedTestRunnerService extends the basic TestRunnerService with advanced testing features
type AdvancedTestRunnerService struct {
	*TestRunnerService
	schemaValidator   *validator.SchemaValidatorService
	variableExtractor *extractor.VariableExtractorService
	assertionEvaluator *asserter.AssertionEvaluatorService
	sequenceRunner    *sequencer.SequenceRunnerService
}

// NewAdvancedTestRunnerService creates a new AdvancedTestRunnerService
func NewAdvancedTestRunnerService(
	executor application.HTTPExecutor,
	snapshotManager application.SnapshotManager,
	fileWriter application.FileWriter,
) *AdvancedTestRunnerService {
	baseRunner := NewTestRunnerService(executor, snapshotManager, fileWriter)
	schemaValidator := validator.NewSchemaValidatorService()
	
	return &AdvancedTestRunnerService{
		TestRunnerService:  baseRunner,
		schemaValidator:    schemaValidator,
		variableExtractor:  extractor.NewVariableExtractorService(),
		assertionEvaluator: asserter.NewAssertionEvaluatorService(),
		sequenceRunner:     sequencer.NewSequenceRunnerService(executor, schemaValidator),
	}
}

// RunWithSchemaValidation runs tests with schema validation
func (s *AdvancedTestRunnerService) RunWithSchemaValidation(
	ctx context.Context,
	patterns []string,
	options models.TestRunOptions,
	swaggerDoc *models.SwaggerDoc,
) (*models.TestReport, error) {
	// Set the swagger doc in options for use in RunTest
	optionsCopy := options
	optionsCopy.SwaggerDoc = swaggerDoc
	
	// Use the base runner to run the tests
	report, err := s.TestRunnerService.RunTests(ctx, patterns, optionsCopy)
	if err != nil {
		return nil, err
	}
	
	// Update summary statistics for schema validation
	report.Summary.SchemaValidated = 0
	report.Summary.SchemaFailed = 0
	
	for _, result := range report.Results {
		if result.SchemaResult != nil {
			report.Summary.SchemaValidated++
			if !result.SchemaResult.Valid {
				report.Summary.SchemaFailed++
			}
		}
	}
	
	return report, nil
}

// RunSequences runs all test sequences matching patterns
func (s *AdvancedTestRunnerService) RunSequences(
	ctx context.Context,
	patterns []string,
	options models.TestRunOptions,
) (*models.TestReport, error) {
	// Find sequences matching the patterns
	sequences, err := s.sequenceRunner.FindSequences(ctx, patterns, options.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find sequences: %w", err)
	}
	
	// Create report
	report := &models.TestReport{
		Name:        "Sequence Tests",
		Summary:     models.TestSummary{},
		Results:     []models.TestResult{},
		Environment: options.EnvironmentVars,
		CreatedAt:   time.Now(),
		Sequences:   []models.TestSequenceResult{},
	}
	
	// Set start time
	report.Summary.StartTime = time.Now()
	
	// Run each sequence
	for _, sequence := range sequences {
		// Load variables from file if configured
		if options.VariablesPath != "" {
			vars, err := s.variableExtractor.LoadVariables(ctx, options.VariablesPath)
			if err == nil {
				// Merge with existing variables, giving priority to loaded ones
				for k, v := range vars {
					options.EnvironmentVars[k] = v
				}
			}
		}
		
		// Run the sequence
		sequenceResult, err := s.RunSequence(ctx, sequence, options)
		if err != nil {
			return nil, fmt.Errorf("failed to run sequence %s: %w", sequence.Name, err)
		}
		
		// Add sequence result to the report
		report.Sequences = append(report.Sequences, *sequenceResult)
		
		// Update summary statistics
		report.Summary.SequencesTotal++
		if sequenceResult.Success {
			report.Summary.SequencesPassed++
		} else {
			report.Summary.SequencesFailed++
		}
		
		// Convert sequence step results to test results for compatibility
		for _, stepResult := range sequenceResult.StepResults {
			testResult := models.TestResult{
				Name:     fmt.Sprintf("%s - %s", sequence.Name, stepResult.Name),
				FilePath: sequence.FilePath,
				Request:  stepResult.Request,
				Response: stepResult.Response,
				Status:   stepResult.Status,
				Error:    stepResult.Error,
				Tags:     sequence.Tags,
				MetaData: map[string]string{
					"sequence": sequence.Name,
					"step":     stepResult.Name,
				},
			}
			
			if len(stepResult.Variables) > 0 {
				testResult.ExtractedVars = stepResult.Variables
			}
			
			if stepResult.SchemaResult != nil {
				testResult.SchemaResult = stepResult.SchemaResult
			}
			
			if len(stepResult.AssertionResults) > 0 {
				testResult.AssertionResults = stepResult.AssertionResults
			}
			
			report.Results = append(report.Results, testResult)
		}
		
		// Stop if sequence failed and we need to stop on failure
		if !sequenceResult.Success && options.StopOnFailure {
			break
		}
	}
	
	// Set end time and calculate duration
	report.Summary.EndTime = time.Now()
	report.Summary.DurationMs = report.Summary.EndTime.Sub(report.Summary.StartTime).Milliseconds()
	
	// Update test summary statistics
	report.Summary.TotalTests = len(report.Results)
	for _, result := range report.Results {
		switch result.Status {
		case models.TestStatusPassed:
			report.Summary.PassedTests++
		case models.TestStatusFailed:
			report.Summary.FailedTests++
		case models.TestStatusSkipped:
			report.Summary.SkippedTests++
		case models.TestStatusError:
			report.Summary.ErrorTests++
		}
		
		if result.SchemaResult != nil {
			report.Summary.SchemaValidated++
			if !result.SchemaResult.Valid {
				report.Summary.SchemaFailed++
			}
		}
	}
	
	return report, nil
}

// RunSequence runs a single test sequence
func (s *AdvancedTestRunnerService) RunSequence(
	ctx context.Context,
	sequence *models.TestSequence,
	options models.TestRunOptions,
) (*models.TestSequenceResult, error) {
	return s.sequenceRunner.RunSequence(ctx, sequence, options)
}

// RunTest overrides the base RunTest to add schema validation and variable extraction
func (s *AdvancedTestRunnerService) RunTest(
	ctx context.Context,
	request *models.HTTPRequest,
	options models.TestRunOptions,
) (*models.TestResult, error) {
	// Run the test using the base implementation
	result, err := s.TestRunnerService.RunTest(ctx, request, options)
	if err != nil {
		return nil, err
	}
	
	// If test failed or error occurred, return early
	if result.Status == models.TestStatusFailed || result.Status == models.TestStatusError {
		return result, nil
	}
	
	// If schema validation is enabled and we have a swagger doc
	if options.ValidateSchema && options.SwaggerDoc != nil && result.Response != nil {
		// Validate response against swagger schema
		schemaResult, err := s.schemaValidator.ValidateResponseWithSwagger(
			ctx,
			result.Response,
			options.SwaggerDoc,
			request.Path,
			request.Method,
			options.ValidationOptions,
		)
		
		if err != nil {
			result.Error = fmt.Sprintf("Schema validation error: %v", err)
			result.Status = models.TestStatusError
		} else {
			result.SchemaResult = schemaResult
			if !schemaResult.Valid {
				result.Error = fmt.Sprintf(
					"Schema validation failed with %d errors",
					len(schemaResult.Errors),
				)
				result.Status = models.TestStatusFailed
			}
		}
	}
	
	// Extract variables if enabled
	if options.ExtractVariables && result.Response != nil && len(request.Variables) > 0 {
		extractedVars, err := s.variableExtractor.Extract(ctx, result.Response, request.Variables)
		if err != nil {
			// Only set error if required variables failed
			if err.Error() == "required" {
				result.Error = fmt.Sprintf("Failed to extract required variables: %v", err)
				result.Status = models.TestStatusFailed
			}
		} else {
			result.ExtractedVars = extractedVars
			
			// Save variables to file if configured
			if options.SaveVariables && options.VariablesPath != "" {
				s.variableExtractor.SaveVariables(ctx, extractedVars, options.VariablesPath)
			}
		}
	}
	
	// Evaluate assertions if enabled
	if options.EnableAssertions && result.Response != nil && len(request.Assertions) > 0 {
		assertionResults, err := s.assertionEvaluator.Evaluate(ctx, result.Response, request.Assertions)
		if err != nil {
			result.Error = fmt.Sprintf("Error evaluating assertions: %v", err)
			result.Status = models.TestStatusError
		} else {
			result.AssertionResults = assertionResults
			
			// Check if any assertions failed
			for _, assertionResult := range assertionResults {
				if !assertionResult.Succeeded {
					result.Status = models.TestStatusFailed
					result.Error = fmt.Sprintf(
						"Assertion failed: %s - %s",
						assertionResult.Type,
						assertionResult.Message,
					)
					break
				}
			}
		}
	}
	
	return result, nil
}

// Override RunTestFile to use the extended RunTest method
func (s *AdvancedTestRunnerService) RunTestFile(
	ctx context.Context,
	file *models.HTTPFile,
	options models.TestRunOptions,
) ([]*models.TestResult, error) {
	var results []*models.TestResult
	
	for _, request := range file.Requests {
		// Check if the test meets the filter criteria
		if !s.matchesFilter(&request, options.Filter) {
			continue
		}
		
		// Set the file path in the request
		request.Path = file.Filename
		
		// Run the test
		result, err := s.RunTest(ctx, &request, options)
		if err != nil {
			return nil, err
		}
		
		results = append(results, result)
		
		// Stop on failure if configured
		if options.StopOnFailure && (result.Status == models.TestStatusFailed || result.Status == models.TestStatusError) {
			break
		}
	}
	
	return results, nil
}
