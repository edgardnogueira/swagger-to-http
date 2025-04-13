package cli

import (
	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/http"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/test"
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
func Execute(
	configProvider application.ConfigProvider, 
	httpParser *http.Parser,
	httpExecutor application.HTTPExecutor,
	testRunner application.TestRunner,
	testReporter application.TestReporter,
	advancedTestRunner *test.AdvancedTestRunnerService,
	fileWriter application.FileWriter,
) error {
	// Add snapshot commands
	InitSnapshotCommands(rootCmd, configProvider, httpExecutor)
	
	// Add test commands
	AddTestCommands(rootCmd, configProvider, testRunner, testReporter)
	
	// Add advanced test commands
	AddAdvancedTestCommands(rootCmd, configProvider, advancedTestRunner, testReporter)

	return rootCmd.Execute()
}

func init() {
	// Add flags, subcommands, etc. here
	rootCmd.AddCommand(versionCmd)
}
