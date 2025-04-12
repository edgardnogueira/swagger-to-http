package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/reporter"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/watcher"
	"github.com/spf13/cobra"
)

// AddTestCommands adds test-related commands to the root command
func AddTestCommands(rootCmd *cobra.Command, configProvider application.ConfigProvider,
	testRunner application.TestRunner, testReporter application.TestReporter) {

	// Test command
	testCmd := &cobra.Command{
		Use:   "test [file-patterns]",
		Short: "Run HTTP tests",
		Long:  `Execute HTTP requests and compare responses with expected values or snapshots`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			updateMode, _ := cmd.Flags().GetString("update")
			ignoreHeaders, _ := cmd.Flags().GetString("ignore-headers")
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			failOnMissing, _ := cmd.Flags().GetBool("fail-on-missing")
			cleanup, _ := cmd.Flags().GetBool("cleanup")
			timeoutStr, _ := cmd.Flags().GetString("timeout")
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

			// Parse timeout
			timeout := 30 * time.Second
			if timeoutStr != "" {
				parsedTimeout, err := time.ParseDuration(timeoutStr)
				if err != nil {
					return fmt.Errorf("invalid timeout format: %w", err)
				}
				timeout = parsedTimeout
			}

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

			// Create test run options
			options := models.TestRunOptions{
				UpdateSnapshots: updateMode,
				FailOnMissing:   failOnMissing,
				IgnoreHeaders:   ignoreHeadersList,
				Timeout:         timeout,
				Parallel:        parallel,
				MaxConcurrent:   maxConcurrent,
				StopOnFailure:   stopOnFailure,
				Filter:          filter,
				EnvironmentVars: extractEnvironmentVars(),
				ReportOptions: models.TestReportOptions{
					Format:           reportFormat,
					OutputPath:       reportOutput,
					IncludeRequests:  detailed,
					IncludeResponses: detailed,
					ColorOutput:      true,
					Detailed:         detailed,
				},
				ContinuousMode:  watch,
				WatchIntervalMs: watchInterval,
			}

			// Add snapshot directory to filter paths if provided
			if snapshotDir != "" {
				if len(options.Filter.Paths) == 0 {
					options.Filter.Paths = []string{snapshotDir}
				}
			}

			// Run in watch mode if specified
			if watch {
				return handleWatchMode(context.Background(), args, options, testRunner, testReporter)
			}

			// Run tests
			report, err := testRunner.RunTests(context.Background(), args, options)
			if err != nil {
				return fmt.Errorf("failed to run tests: %w", err)
			}

			// Print report to console
			consoleOptions := options.ReportOptions
			consoleOptions.Format = "console"
			err = testReporter.PrintReport(context.Background(), report, consoleOptions, os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to print report: %w", err)
			}

			// Generate report file if output path specified
			if reportOutput != "" {
				err = testReporter.SaveReport(context.Background(), report, options.ReportOptions)
				if err != nil {
					return fmt.Errorf("failed to save report: %w", err)
				}
				fmt.Printf("Report saved to %s\n", reportOutput)
			}

			// Return non-zero exit code if any tests failed
			if report.Summary.FailedTests > 0 || report.Summary.ErrorTests > 0 {
				return fmt.Errorf("tests failed: %d failed, %d errors", report.Summary.FailedTests, report.Summary.ErrorTests)
			}

			return nil
		},
	}

	// Add flags to test command
	testCmd.Flags().String("update", "none", "Update mode: none, all, failed, missing")
	testCmd.Flags().String("ignore-headers", "Date,Set-Cookie", "Comma-separated headers to ignore in comparison")
	testCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	testCmd.Flags().Bool("fail-on-missing", false, "Fail when snapshot is missing")
	testCmd.Flags().Bool("cleanup", false, "Remove unused snapshots after testing")
	testCmd.Flags().String("timeout", "30s", "HTTP request timeout")
	testCmd.Flags().Bool("parallel", false, "Run tests in parallel")
	testCmd.Flags().Int("max-concurrent", 5, "Maximum number of concurrent tests")
	testCmd.Flags().Bool("stop-on-failure", false, "Stop testing after first failure")
	testCmd.Flags().StringSlice("tags", []string{}, "Filter tests by tags")
	testCmd.Flags().StringSlice("methods", []string{}, "Filter tests by HTTP methods")
	testCmd.Flags().StringSlice("paths", []string{}, "Filter tests by request paths")
	testCmd.Flags().StringSlice("names", []string{}, "Filter tests by test names")
	testCmd.Flags().String("report-format", "console", "Report format: console, json, html, junit")
	testCmd.Flags().String("report-output", "", "Path to write report file")
	testCmd.Flags().Bool("detailed", false, "Include detailed information in report")
	testCmd.Flags().Bool("watch", false, "Run in continuous (watch) mode")
	testCmd.Flags().Int("watch-interval", 1000, "Interval between watch checks in milliseconds")

	// List command
	listCmd := &cobra.Command{
		Use:   "list [file-patterns]",
		Short: "List available HTTP tests",
		Long:  `Find and list HTTP tests in the specified files`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			tags, _ := cmd.Flags().GetStringSlice("tags")
			methods, _ := cmd.Flags().GetStringSlice("methods")
			paths, _ := cmd.Flags().GetStringSlice("paths")
			names, _ := cmd.Flags().GetStringSlice("names")

			// Create test filter
			filter := models.TestFilter{
				Tags:    tags,
				Methods: methods,
				Paths:   paths,
				Names:   names,
			}

			// Find tests matching the filter
			files, err := testRunner.FindTests(context.Background(), args, filter)
			if err != nil {
				return fmt.Errorf("failed to find tests: %w", err)
			}

			// Print test information
			fmt.Printf("Found %d files with tests:\n\n", len(files))
			for _, file := range files {
				fmt.Printf("File: %s\n", file.Filename)
				fmt.Printf("  Tests: %d\n", len(file.Requests))

				// Group tests by tag
				testsByTag := make(map[string][]models.HTTPRequest)
				for _, req := range file.Requests {
					tag := req.Tag
					if tag == "" {
						tag = "default"
					}
					testsByTag[tag] = append(testsByTag[tag], req)
				}

				// Print tests by tag
				for tag, tests := range testsByTag {
					fmt.Printf("  Tag: %s\n", tag)
					for _, test := range tests {
						fmt.Printf("    %s %s\n", test.Method, test.Name)
					}
				}
				fmt.Println()
			}

			return nil
		},
	}

	// Add flags to list command
	listCmd.Flags().StringSlice("tags", []string{}, "Filter tests by tags")
	listCmd.Flags().StringSlice("methods", []string{}, "Filter tests by HTTP methods")
	listCmd.Flags().StringSlice("paths", []string{}, "Filter tests by request paths")
	listCmd.Flags().StringSlice("names", []string{}, "Filter tests by test names")

	// Add test commands to root command
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(listCmd)
}

// handleWatchMode runs tests in watch mode
func handleWatchMode(ctx context.Context, patterns []string, options models.TestRunOptions, 
	testRunner application.TestRunner, testReporter application.TestReporter) error {
	
	// Create a watcher service
	watcherService := watcher.NewTestWatcherService(testRunner, testReporter)

	// Start watching
	if err := watcherService.Watch(ctx, patterns, options); err != nil {
		return fmt.Errorf("failed to start watcher: %w", err)
	}

	fmt.Println("Watching for changes. Press Ctrl+C to stop...")

	// Wait for interrupt signal
	<-ctx.Done()
	
	// Stop the watcher
	watcherService.Stop()
	
	return nil
}

// extractEnvironmentVars extracts environment variables with HTTP_ prefix
func extractEnvironmentVars() map[string]string {
	vars := make(map[string]string)
	
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			name := parts[0]
			value := parts[1]
			
			// Check for HTTP_ prefix
			if strings.HasPrefix(name, "HTTP_") {
				// Remove the prefix and add to map
				key := strings.TrimPrefix(name, "HTTP_")
				vars[key] = value
			}
		}
	}
	
	return vars
}
