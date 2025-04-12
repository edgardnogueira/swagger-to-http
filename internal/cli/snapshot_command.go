package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/application/snapshot"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/http"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddSnapshotCommands adds snapshot-related commands to the root command
func AddSnapshotCommands(rootCmd *cobra.Command, configProvider application.ConfigProvider) {
	// Base snapshot command
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Snapshot testing commands",
		Long:  "Commands for working with HTTP response snapshots",
	}
	
	// Snapshot test command
	testCmd := &cobra.Command{
		Use:   "test [file-pattern]",
		Short: "Run snapshot tests",
		Long:  "Execute HTTP requests and compare with stored snapshots",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			updateMode, _ := cmd.Flags().GetString("update")
			ignoreHeadersStr, _ := cmd.Flags().GetString("ignore-headers")
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			failOnMissing, _ := cmd.Flags().GetBool("fail-on-missing")
			cleanup, _ := cmd.Flags().GetBool("cleanup")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			
			// Parse ignore headers
			var ignoreHeaders []string
			if ignoreHeadersStr != "" {
				ignoreHeaders = strings.Split(ignoreHeadersStr, ",")
				for i, h := range ignoreHeaders {
					ignoreHeaders[i] = strings.TrimSpace(h)
				}
			}
			
			// Create snapshot options
			options := models.SnapshotOptions{
				UpdateMode:     updateMode,
				IgnoreHeaders:  ignoreHeaders,
				BasePath:       snapshotDir,
				UpdateExisting: updateMode == "all" || updateMode == "failed",
			}
			
			// Determine file pattern
			pattern := "**/*.http"
			if len(args) > 0 {
				pattern = args[0]
			}
			
			return runSnapshotTests(cmd, pattern, options, failOnMissing, cleanup, timeout)
		},
	}
	
	// Add flags to test command
	testCmd.Flags().String("update", "none", "Update mode: none, all, failed, missing")
	testCmd.Flags().String("ignore-headers", "Date,Set-Cookie", "Comma-separated headers to ignore in comparison")
	testCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	testCmd.Flags().Bool("fail-on-missing", false, "Fail when snapshot is missing")
	testCmd.Flags().Bool("cleanup", false, "Remove unused snapshots after testing")
	testCmd.Flags().Duration("timeout", 30*time.Second, "HTTP request timeout")
	
	// Snapshot update command
	updateCmd := &cobra.Command{
		Use:   "update [file-pattern]",
		Short: "Update snapshots",
		Long:  "Execute HTTP requests and update stored snapshots",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			timeout, _ := cmd.Flags().GetDuration("timeout")
			
			// Create snapshot options with update mode set to "all"
			options := models.SnapshotOptions{
				UpdateMode:     "all",
				BasePath:       snapshotDir,
				UpdateExisting: true,
			}
			
			// Determine file pattern
			pattern := "**/*.http"
			if len(args) > 0 {
				pattern = args[0]
			}
			
			return runSnapshotTests(cmd, pattern, options, false, false, timeout)
		},
	}
	
	// Add flags to update command
	updateCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	updateCmd.Flags().Duration("timeout", 30*time.Second, "HTTP request timeout")
	
	// Snapshot list command
	listCmd := &cobra.Command{
		Use:   "list [directory]",
		Short: "List snapshots",
		Long:  "List all snapshots in a directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			
			// Determine directory
			dir := ""
			if len(args) > 0 {
				dir = args[0]
			}
			
			return listSnapshots(cmd, snapshotDir, dir)
		},
	}
	
	// Add flags to list command
	listCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	
	// Snapshot cleanup command
	cleanupCmd := &cobra.Command{
		Use:   "cleanup [directory]",
		Short: "Cleanup snapshots",
		Long:  "Remove unused or orphaned snapshots",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			
			// Determine directory
			dir := ""
			if len(args) > 0 {
				dir = args[0]
			}
			
			return cleanupSnapshots(cmd, snapshotDir, dir)
		},
	}
	
	// Add flags to cleanup command
	cleanupCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	
	// Add commands to snapshot command
	snapshotCmd.AddCommand(testCmd)
	snapshotCmd.AddCommand(updateCmd)
	snapshotCmd.AddCommand(listCmd)
	snapshotCmd.AddCommand(cleanupCmd)
	
	// Add snapshot command to root
	rootCmd.AddCommand(snapshotCmd)
}

// runSnapshotTests runs snapshot tests for the given file pattern
func runSnapshotTests(cmd *cobra.Command, pattern string, options models.SnapshotOptions, failOnMissing, cleanup bool, timeout time.Duration) error {
	// Create snapshot manager and service
	manager := snapshot.NewManager(options.BasePath)
	service := snapshot.NewService(manager, options)
	
	// Create HTTP parser
	parser := http.NewParser()
	
	// Find HTTP files matching pattern
	files, err := parser.FindHTTPFiles(pattern)
	if err != nil {
		return fmt.Errorf("failed to find HTTP files: %w", err)
	}
	
	if len(files) == 0 {
		fmt.Println("No HTTP files found matching pattern:", pattern)
		return nil
	}
	
	fmt.Printf("Found %d HTTP files to test\n", len(files))
	
	// Create HTTP executor with default environment variables
	env := loadEnvironmentVariables()
	executor := http.NewExecutor(timeout, env)
	
	// Process each file
	totalResults := []*models.SnapshotResult{}
	
	for _, file := range files {
		fmt.Printf("\nTesting file: %s\n", file)
		
		// Read and parse HTTP file
		httpFile, err := parser.ParseFile(file)
		if err != nil {
			fmt.Printf("  Error parsing file: %s\n", err)
			continue
		}
		
		// Process each request
		for i, request := range httpFile.Requests {
			fmt.Printf("  Request %d: %s %s\n", i+1, request.Method, request.Path)
			
			// Execute request
			response, err := executor.Execute(context.Background(), &request, nil)
			if err != nil {
				fmt.Printf("    Error executing request: %s\n", err)
				continue
			}
			
			// Run snapshot test
			result, err := service.RunTest(context.Background(), response, file)
			if err != nil {
				if strings.Contains(err.Error(), "snapshot does not exist") && !failOnMissing {
					fmt.Printf("    %s Snapshot does not exist (created)\n", color.GreenString("✓"))
					totalResults = append(totalResults, &models.SnapshotResult{
						RequestPath:   request.Path,
						RequestMethod: request.Method,
						Passed:        true,
						Updated:       true,
					})
					continue
				}
				
				fmt.Printf("    %s %s\n", color.RedString("✗"), err)
				totalResults = append(totalResults, &models.SnapshotResult{
					RequestPath:   request.Path,
					RequestMethod: request.Method,
					Passed:        false,
					Error:         err,
				})
				continue
			}
			
			if result.Passed {
				if result.Updated {
					fmt.Printf("    %s Snapshot updated\n", color.YellowString("⟳"))
				} else {
					fmt.Printf("    %s Snapshot matched\n", color.GreenString("✓"))
				}
			} else {
				fmt.Printf("    %s Snapshot comparison failed\n", color.RedString("✗"))
				
				// Print diff details
				if result.Diff != nil && result.Diff.StatusDiff != nil && !result.Diff.StatusDiff.Equal {
					fmt.Printf("      Status code: expected %d, got %d\n", 
						result.Diff.StatusDiff.Expected, 
						result.Diff.StatusDiff.Actual)
				}
				
				if result.Diff != nil && result.Diff.HeaderDiff != nil && !result.Diff.HeaderDiff.Equal {
					fmt.Println("      Headers differ:")
					if len(result.Diff.HeaderDiff.MissingHeaders) > 0 {
						fmt.Println("        Missing headers:")
						for h := range result.Diff.HeaderDiff.MissingHeaders {
							fmt.Printf("          - %s\n", h)
						}
					}
					if len(result.Diff.HeaderDiff.ExtraHeaders) > 0 {
						fmt.Println("        Extra headers:")
						for h := range result.Diff.HeaderDiff.ExtraHeaders {
							fmt.Printf("          + %s\n", h)
						}
					}
				}
				
				if result.Diff != nil && result.Diff.BodyDiff != nil && !result.Diff.BodyDiff.Equal {
					fmt.Printf("      Body content differs (expected %d bytes, got %d bytes)\n", 
						result.Diff.BodyDiff.ExpectedSize, 
						result.Diff.BodyDiff.ActualSize)
					
					// Print diff preview if available
					if result.Diff.BodyDiff.DiffContent != "" {
						fmt.Println("      Diff preview:")
						lines := strings.Split(result.Diff.BodyDiff.DiffContent, "\n")
						maxLines := 10
						if len(lines) > maxLines {
							lines = lines[:maxLines]
							fmt.Printf("        %s\n        ...(truncated)...\n", 
								strings.Join(lines, "\n        "))
						} else {
							fmt.Printf("        %s\n", strings.Join(lines, "\n        "))
						}
					}
				}
			}
			
			totalResults = append(totalResults, result)
		}
	}
	
	// Print summary
	stats := service.GetStats()
	printTestSummary(stats)
	
	// Cleanup unused snapshots if requested
	if cleanup {
		fmt.Println("\nCleaning up unused snapshots...")
		if err := service.CleanupUnusedSnapshots(context.Background(), ""); err != nil {
			fmt.Printf("Error during cleanup: %s\n", err)
		} else {
			fmt.Println("Cleanup completed successfully")
		}
	}
	
	// Return error if any tests failed
	if stats.Failed > 0 {
		return fmt.Errorf("%d of %d tests failed", stats.Failed, stats.Total)
	}
	
	return nil
}

// listSnapshots lists the snapshots in a directory
func listSnapshots(cmd *cobra.Command, basePath, directory string) error {
	manager := snapshot.NewManager(basePath)
	
	snapshots, err := manager.ListSnapshots(context.Background(), directory)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}
	
	if len(snapshots) == 0 {
		fmt.Println("No snapshots found")
		return nil
	}
	
	fmt.Printf("Found %d snapshots:\n", len(snapshots))
	for _, s := range snapshots {
		fmt.Printf("  %s\n", s)
	}
	
	return nil
}

// cleanupSnapshots removes orphaned snapshots
func cleanupSnapshots(cmd *cobra.Command, basePath, directory string) error {
	manager := snapshot.NewManager(basePath)
	
	// Create HTTP parser to find valid HTTP files
	parser := http.NewParser()
	
	// List all snapshots
	snapshots, err := manager.ListSnapshots(context.Background(), directory)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}
	
	if len(snapshots) == 0 {
		fmt.Println("No snapshots found")
		return nil
	}
	
	fmt.Printf("Found %d snapshots\n", len(snapshots))
	
	// Find all HTTP files
	httpFiles, err := parser.FindHTTPFiles("**/*.http")
	if err != nil {
		return fmt.Errorf("failed to find HTTP files: %w", err)
	}
	
	// Create a map of potential snapshot paths
	validPaths := make(map[string]bool)
	
	// Parse HTTP files and add possible snapshot paths
	for _, file := range httpFiles {
		httpFile, err := parser.ParseFile(file)
		if err != nil {
			continue
		}
		
		dir := filepath.Dir(file)
		fileBase := filepath.Base(file)
		fileExt := filepath.Ext(fileBase)
		fileName := fileBase
		if fileExt != "" {
			fileName = fileBase[:len(fileBase)-len(fileExt)]
		}
		
		for _, request := range httpFile.Requests {
			snapshotName := fmt.Sprintf("%s_%s.snap.json", fileName, strings.ToLower(request.Method))
			snapshotPath := filepath.Join(basePath, dir, snapshotName)
			validPaths[snapshotPath] = true
		}
	}
	
	// Check each snapshot against the valid paths
	var orphaned []string
	for _, s := range snapshots {
		fullPath := filepath.Join(basePath, s)
		if _, valid := validPaths[fullPath]; !valid {
			orphaned = append(orphaned, s)
		}
	}
	
	if len(orphaned) == 0 {
		fmt.Println("No orphaned snapshots found")
		return nil
	}
	
	fmt.Printf("Found %d orphaned snapshots:\n", len(orphaned))
	for _, s := range orphaned {
		fmt.Printf("  %s\n", s)
	}
	
	// Confirm deletion
	fmt.Print("Do you want to delete these snapshots? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)
	
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Cleanup cancelled")
		return nil
	}
	
	// Delete orphaned snapshots
	deleted := 0
	for _, s := range orphaned {
		fullPath := filepath.Join(basePath, s)
		if err := os.Remove(fullPath); err != nil {
			fmt.Printf("Error deleting %s: %s\n", s, err)
		} else {
			deleted++
		}
	}
	
	fmt.Printf("Deleted %d orphaned snapshots\n", deleted)
	return nil
}

// printTestSummary prints a summary of test results
func printTestSummary(stats *models.SnapshotStats) {
	duration := stats.EndTime.Sub(stats.StartTime)
	
	fmt.Println("\n=======================================")
	fmt.Println("Snapshot Test Summary")
	fmt.Println("=======================================")
	fmt.Printf("Total tests:    %d\n", stats.Total)
	fmt.Printf("Passed:         %s\n", color.GreenString("%d", stats.Passed))
	
	if stats.Failed > 0 {
		fmt.Printf("Failed:         %s\n", color.RedString("%d", stats.Failed))
	} else {
		fmt.Printf("Failed:         %d\n", stats.Failed)
	}
	
	if stats.Created > 0 {
		fmt.Printf("Created:        %s\n", color.YellowString("%d", stats.Created))
	} else {
		fmt.Printf("Created:        %d\n", stats.Created)
	}
	
	if stats.Updated > 0 {
		fmt.Printf("Updated:        %s\n", color.YellowString("%d", stats.Updated))
	} else {
		fmt.Printf("Updated:        %d\n", stats.Updated)
	}
	
	fmt.Printf("Duration:       %.2f seconds\n", duration.Seconds())
	fmt.Println("=======================================")
}

// loadEnvironmentVariables loads environment variables for HTTP requests
func loadEnvironmentVariables() map[string]string {
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
