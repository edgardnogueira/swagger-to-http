package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/application/executor"
	"github.com/edgardnogueira/swagger-to-http/internal/application/snapshot"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
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
				IgnoreHeaders: ignoreHeaders,
				BasePath:      snapshotDir,
				UpdateExisting: updateMode == "all" || updateMode == "failed",
			}
			
			// Determine file pattern
			pattern := "**/*.http"
			if len(args) > 0 {
				pattern = args[0]
			}
			
			return runSnapshotTests(cmd, pattern, options, failOnMissing, cleanup)
		},
	}
	
	// Add flags to test command
	testCmd.Flags().String("update", "none", "Update mode: none, all, failed, missing")
	testCmd.Flags().String("ignore-headers", "Date,Set-Cookie", "Comma-separated headers to ignore in comparison")
	testCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	testCmd.Flags().Bool("fail-on-missing", false, "Fail when snapshot is missing")
	testCmd.Flags().Bool("cleanup", false, "Remove unused snapshots after testing")
	
	// Snapshot update command
	updateCmd := &cobra.Command{
		Use:   "update [file-pattern]",
		Short: "Update snapshots",
		Long:  "Execute HTTP requests and update stored snapshots",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotDir, _ := cmd.Flags().GetString("snapshot-dir")
			
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
			
			return runSnapshotTests(cmd, pattern, options, false, false)
		},
	}
	
	// Add flags to update command
	updateCmd.Flags().String("snapshot-dir", ".snapshots", "Directory for snapshot storage")
	
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