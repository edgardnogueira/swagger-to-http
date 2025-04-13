package main

import (
	"fmt"
	"os"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/cli"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/config"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/extractor"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/fs"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/http"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/reporter"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/snapshot"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/test"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/validator"
	"github.com/edgardnogueira/swagger-to-http/internal/version"
)

func main() {
	// Initialize configuration
	configProvider := config.NewConfig()

	// Create HTTP parser
	httpParser := http.NewParser()

	// Create HTTP executor
	httpExecutor := http.NewExecutor(30*time.Second, nil)

	// Create file system services
	fileWriter := fs.NewFileWriter()
	snapshotManager := snapshot.NewSnapshotManager(fileWriter)

	// Create basic test services
	testRunner := test.NewTestRunnerService(httpExecutor, snapshotManager, fileWriter)
	testReporter := reporter.NewTestReporterService()

	// Create advanced test services
	variableExtractor := extractor.NewVariableExtractorService()
	schemaValidator := validator.NewSchemaValidatorService()
	advancedTestRunner := test.NewAdvancedTestRunnerService(httpExecutor, snapshotManager, fileWriter)

	// Initialize CLI
	if err := cli.Execute(
		configProvider,
		httpParser,
		httpExecutor,
		testRunner,
		testReporter,
		advancedTestRunner,
		fileWriter,
	); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
