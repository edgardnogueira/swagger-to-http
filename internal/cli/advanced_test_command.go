package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/test"
)

// AddAdvancedTestCommands adds advanced test-related commands to the root command
func AddAdvancedTestCommands(
	rootCmd *cobra.Command,
	configProvider application.ConfigProvider,
	advancedTestRunner *test.AdvancedTestRunnerService,
	testReporter application.TestReporter,
) {

	// Schema validation command
	validateCmd := &cobra.Command{
		Use:   "validate [file-patterns]",
		Short: "Validate responses against OpenAPI schema",
		Long:  `Execute HTTP requests and validate responses against OpenAPI/Swagger schema definitions`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create common test options using the same code from test_command.go
			options, err := createTestRunOptions(cmd)
			if err != nil {
				return err
			}

			// Get schema validation specific flags
			swaggerFile, _ := cmd.Flags().GetString("swagger-file")
			ignoreProps, _ := cmd.Flags().GetString("ignore-props")
			ignoreAddProps, _ := cmd.Flags().GetBool("ignore-add-props")
			ignoreFormats, _ := cmd.Flags().GetBool("ignore-formats")
			ignorePatterns, _ := cmd.Flags().GetBool("ignore-patterns")
			reqPropsOnly, _ := cmd.Flags().GetBool("req-props-only")
			ignoreNullable, _ := cmd.Flags().GetBool("ignore-nullable")

			// Parse ignore properties
			var ignoredProps []string
			if ignoreProps != "" {
				ignoredProps = strings.Split(ignoreProps, ",")
				for i := range ignoredProps {
					ignoredProps[i] = strings.TrimSpace(ignoredProps[i])
				}
			}

			// Create validation options
			validationOptions := models.ValidationOptions{
				IgnoreAdditionalProperties: ignoreAddProps,
				IgnoreFormats:             ignoreFormats,
				IgnorePatterns:             ignorePatterns,
				RequiredPropertiesOnly:     reqPropsOnly,
				IgnoreNullable:             ignoreNullable,
				IgnoredProperties:          ignoredProps,
			}

			// Add schema validation options
			options.ValidateSchema = true
			options.ValidationOptions = validationOptions

			// Load the Swagger file
			if swaggerFile == "" {
				return fmt.Errorf("swagger file is required for schema validation")
			}

			// Parse the Swagger file
			fmt.Printf("Loading Swagger file: %s\n", swaggerFile)
			swaggerDoc, err := loadSwaggerDoc(context.Background(), swaggerFile)
			if err != nil {
				return fmt.Errorf("failed to load Swagger file: %w", err)
			}

			// Run tests with schema validation
			report, err := advancedTestRunner.RunWithSchemaValidation(
				context.Background(),
				args,
				options,
				swaggerDoc,
			)
			if err != nil {
				return fmt.Errorf("failed to run tests with schema validation: %w", err)
			}

			// Print report to console
			consoleOptions := options.ReportOptions
			consoleOptions.Format = "console"
			err = testReporter.PrintReport(context.Background(), report, consoleOptions, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to print report: %w", err)
			}

			// Generate report file if output path specified
			if reportOutput := options.ReportOptions.OutputPath; reportOutput != "" {
				err = testReporter.SaveReport(context.Background(), report, options.ReportOptions)
				if err != nil {
					return fmt.Errorf("failed to save report: %w", err)
				}
				fmt.Printf("Report saved to %s\n", reportOutput)
			}

			// Return non-zero exit code if any tests failed
			if report.Summary.FailedTests > 0 || report.Summary.ErrorTests > 0 {
				return fmt.Errorf(
					"tests failed: %d failed, %d errors, %d schema validation failures",
					report.Summary.FailedTests,
					report.Summary.ErrorTests,
					report.Summary.SchemaFailed,
				)
			}

			return nil
		},
	}

	// Sequence command
	sequenceCmd := &cobra.Command{
		Use:   "sequence [file-patterns]",
		Short: "Run test sequences",
		Long:  `Execute test sequences with support for variable extraction and dependencies`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create common test options
			options, err := createTestRunOptions(cmd)
			if err != nil {
				return err
			}

			// Get sequence specific flags
			variablesPath, _ := cmd.Flags().GetString("variables-path")
			saveVars, _ := cmd.Flags().GetBool("save-vars")
			varFormat, _ := cmd.Flags().GetString("var-format")
			failFast, _ := cmd.Flags().GetBool("fail-fast")
			validateSchema, _ := cmd.Flags().GetBool("validate-schema")
			swaggerFile, _ := cmd.Flags().GetString("swagger-file")

			// Update options for sequences
			options.SequentialRun = true
			options.ExtractVariables = true
			options.SaveVariables = saveVars
			options.VariablesPath = variablesPath
			options.VariableFormat = varFormat
			options.FailFast = failFast
			options.ValidateSchema = validateSchema
			options.EnableAssertions = true

			// Load swagger doc if specified
			swaggerDoc := &models.SwaggerDoc{}
			if validateSchema && swaggerFile != "" {
				var loadErr error
				swaggerDoc, loadErr = loadSwaggerDoc(context.Background(), swaggerFile)
				if loadErr != nil {
					return fmt.Errorf("failed to load Swagger file: %w", loadErr)
				}
			}

			// Run sequences
			report, err := advancedTestRunner.RunSequences(context.Background(), args, options)
			if err != nil {
				return fmt.Errorf("failed to run sequences: %w", err)
			}

			// Print report to console
			consoleOptions := options.ReportOptions
			consoleOptions.Format = "console"
			consoleOptions.IncludeExtracted = true
			consoleOptions.IncludeAssertions = true
			err = testReporter.PrintReport(context.Background(), report, consoleOptions, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to print report: %w", err)
			}

			// Generate report file if output path specified
			if options.ReportOptions.OutputPath != "" {
				err = testReporter.SaveReport(context.Background(), report, options.ReportOptions)
				if err != nil {
					return fmt.Errorf("failed to save report: %w", err)
				}
				fmt.Printf("Report saved to %s\n", options.ReportOptions.OutputPath)
			}

			// Return non-zero exit code if any sequences failed
			if report.Summary.SequencesFailed > 0 {
				return fmt.Errorf("sequences failed: %d of %d", 
					report.Summary.SequencesFailed, 
					report.Summary.SequencesTotal)
			}

			return nil
		},
	}

	// Add flags to validate command
	validateCmd.Flags().String("swagger-file", "", "Path to Swagger/OpenAPI file")
	validateCmd.Flags().String("ignore-props", "", "Comma-separated properties to ignore in validation")
	validateCmd.Flags().Bool("ignore-add-props", false, "Ignore additional properties not in schema")
	validateCmd.Flags().Bool("ignore-formats", false, "Ignore format validation (e.g., date, email)")
	validateCmd.Flags().Bool("ignore-patterns", false, "Ignore pattern validation")
	validateCmd.Flags().Bool("req-props-only", false, "Validate only required properties")
	validateCmd.Flags().Bool("ignore-nullable", false, "Ignore nullable field validation")
	validateCmd.MarkFlagRequired("swagger-file")

	// Add flags to sequence command
	sequenceCmd.Flags().String("variables-path", "", "Path to load/save variables")
	sequenceCmd.Flags().Bool("save-vars", false, "Save extracted variables to file")
	sequenceCmd.Flags().String("var-format", "${%s}", "Variable format (default: ${varname})")
	sequenceCmd.Flags().Bool("fail-fast", false, "Stop sequence on first failure")
	sequenceCmd.Flags().Bool("validate-schema", false, "Validate responses against schema")
	sequenceCmd.Flags().String("swagger-file", "", "Path to Swagger/OpenAPI file")

	// Add commands to test command
	commands := rootCmd.Commands()
	for _, cmd := range commands {
		if cmd.Use == "test [file-patterns]" {
			cmd.AddCommand(validateCmd)
			cmd.AddCommand(sequenceCmd)
			break
		}
	}
}

// Helper function to load a Swagger document
func loadSwaggerDoc(ctx context.Context, filePath string) (*models.SwaggerDoc, error) {
	// We'll need to implement or use a Swagger parser here
	// For now, return a placeholder
	return &models.SwaggerDoc{
		Paths: make(map[string]*models.PathItem),
	}, nil
}

// Helper function to create test run options based on command flags
func createTestRunOptions(cmd *cobra.Command) (models.TestRunOptions, error) {
	// Get flags
	updateMode, _ := cmd.Flags().GetString("update")
	ignoreHeaders, _ := cmd.Flags().GetString("ignore-headers")
	snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
	failOnMissing, _ := cmd.Flags().GetBool("fail-on-missing")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	parallel, _ := cmd.Flags().GetBool("parallel")
	maxConcurrent, _ := cmd.Flags().GetInt("max-concurrent")
	stopOnFailure, _ := cmd.Flags().GetBool("stop-on-failure")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	methods, _ := cmd.Flags().GetStringSlice("methods")
	paths, _ := cmd.Flags().GetStringSlice("paths")
	names, _ := cmd.Flags().GetStringSlice("names")
	reportFormat, _ := cmd.Flags().GetString("report-format")
	reportOutput, _ := cmd.Flags().GetString("report-output")
	detailed, _ := cmd.Flags().GetBool("detailed")
	watch, _ := cmd.Flags().GetBool("watch")
	watchInterval, _ := cmd.Flags().GetInt("watch-interval")

	// Parse ignore headers
	ignoreHeadersList := []string{"Date", "Set-Cookie"}
	if ignoreHeaders != "" {
		ignoreHeadersList = strings.Split(ignoreHeaders, ",")
		for i := range ignoreHeadersList {
			ignoreHeadersList[i] = strings.TrimSpace(ignoreHeadersList[i])
		}
	}

	// Create test filter
	filter := models.TestFilter{
		Tags:    tags,
		Methods: methods,
		Paths:   paths,
		Names:   names,
	}

	// Add snapshot directory to filter paths if provided
	if snapshotDir != "" {
		if len(filter.Paths) == 0 {
			filter.Paths = []string{snapshotDir}
		}
	}

	// Create test run options
	options := models.TestRunOptions{
		UpdateSnapshots: updateMode,
		FailOnMissing:   failOnMissing,
		IgnoreHeaders:   ignoreHeadersList,
		Parallel:        parallel,
		MaxConcurrent:   maxConcurrent,
		StopOnFailure:   stopOnFailure,
		Filter:          filter,
		EnvironmentVars: extractEnvironmentVars(),
		ContinuousMode:  watch,
		WatchIntervalMs: watchInterval,
		Timeout:         timeout,
		ReportOptions: models.TestReportOptions{
			Format:           reportFormat,
			OutputPath:       reportOutput,
			IncludeRequests:   detailed,
			IncludeResponses:  detailed,
			ColorOutput:       true,
			Detailed:          detailed,
			IncludeExtracted:  true,
			IncludeAssertions: true,
		},
	}

	return options, nil
}

// Helper function to extract environment variables
func extractEnvironmentVars() map[string]string {
	env := make(map[string]string)
	
	// Add environment variables from .env file if exists
	// TODO: Implement .env file loading
	
	// Add system environment variables with HTTP_ prefix
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], "HTTP_") {
			key := strings.TrimPrefix(parts[0], "HTTP_")
			env[key] = parts[1]
		}
	}
	
	return env
}
