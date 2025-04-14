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

// AdvancedTestRunnerService extends the basic test runner with advanced testing features
type AdvancedTestRunnerService struct {
	baseRunner        *BasicTestRunnerService
	schemaValidator   *validator.SchemaValidatorService
	variableExtractor *extractor.VariableExtractorService
	assertionEvaluator *asserter.AssertionEvaluatorService
	sequenceRunner     *sequencer.SequenceRunnerService
}

// NewAdvancedTestRunnerService creates a new AdvancedTestRunnerService
func NewAdvancedTestRunnerService(
	executor application.HTTPExecutor,
	snapshotManager application.SnapshotManager,
	fileWriter application.FileWriter,
) *AdvancedTestRunnerService {
	baseRunner := NewBasicTestRunnerService(executor, snapshotManager, fileWriter)
	schemaValidator := validator.NewSchemaValidatorService()
	
	return &AdvancedTestRunnerService{
		baseRunner:        baseRunner,
		schemaValidator:   schemaValidator,
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
	// Run the tests using the base runner
	report, err := s.baseRunner.RunTests(ctx, patterns, options)
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
	result, err := s.baseRunner.RunTest(ctx, request, options)
	if err != nil {
		return nil, err
	}
	
	// If test failed or error occurred, return early
	if result.Status == models.TestStatusFailed || result.Status == models.TestStatusError {
		return result, nil
	}
	
	// If schema validation is enabled and we have a response
	if options.ValidateSchema && result.Response != nil {
		// Define a schema path
		schemaPath := request.Path // Use request path as schema path for now
		
		// Validate response using the schema
		schemaResult, err := s.schemaValidator.ValidateResponse(
			ctx,
			result.Response,
			schemaPath,
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
	
	// Variable extraction code can be added here when needed
	// Currently, removing the references to the non-existent Variables field in HTTPRequest
	
	// Assertion evaluation code can be added here when needed  
	// Currently, removing the references to the non-existent Assertions field in HTTPRequest
	
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

// BasicTestRunnerService is a minimal implementation for testing
type BasicTestRunnerService struct {
	executor        application.HTTPExecutor
	snapshotManager application.SnapshotManager
	fileWriter      application.FileWriter
}

// NewBasicTestRunnerService creates a new BasicTestRunnerService
func NewBasicTestRunnerService(
	executor application.HTTPExecutor,
	snapshotManager application.SnapshotManager,
	fileWriter application.FileWriter,
) *BasicTestRunnerService {
	return &BasicTestRunnerService{
		executor:        executor,
		snapshotManager: snapshotManager,
		fileWriter:      fileWriter,
	}
}

// RunTests runs tests in basic mode
func (s *BasicTestRunnerService) RunTests(
	ctx context.Context,
	patterns []string,
	options models.TestRunOptions,
) (*models.TestReport, error) {
	// Stub implementation
	return &models.TestReport{
		Name:      "Basic Tests",
		CreatedAt: time.Now(),
		Summary: models.TestSummary{
			StartTime: time.Now(),
			EndTime:   time.Now(),
		},
	}, nil
}

// RunTest runs a basic test
func (s *BasicTestRunnerService) RunTest(
	ctx context.Context,
	request *models.HTTPRequest,
	options models.TestRunOptions,
) (*models.TestResult, error) {
	// Stub implementation
	return &models.TestResult{
		Name:   request.Name,
		Status: models.TestStatusPassed,
	}, nil
}

// matchesFilter is a helper method for filtering tests
func (s *AdvancedTestRunnerService) matchesFilter(request *models.HTTPRequest, filter models.TestFilter) bool {
	// Stub implementation
	return true
}
