package cli

import (
	"github.com/edgardnogueira/swagger-to-http/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "swagger-to-http",
	Short:   "Convert Swagger/OpenAPI docs to HTTP request files",
	Long:    `A tool to convert Swagger/OpenAPI documentation into organized HTTP request files with snapshot testing capabilities.`,
	Version: version.Version,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add flags, subcommands, etc. here
	rootCmd.AddCommand(versionCmd)
}
