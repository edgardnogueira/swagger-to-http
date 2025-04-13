package main

import (
	"fmt"
	"os"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application"
	"github.com/edgardnogueira/swagger-to-http/internal/cli"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/config"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/http"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/reporter"
	"github.com/edgardnogueira/swagger-to-http/internal/version"
)

func main() {
	// Initialize configuration
	configProvider := config.NewConfig()

	// Create HTTP parser
	httpParser := http.NewParser()

	// Create HTTP executor
	httpExecutor := http.NewExecutor(30*time.Second, nil)

	// Create test runner
	testRunner := application.NewTestRunnerService(httpExecutor, nil, nil)

	// Create test reporter
	testReporter := reporter.NewTestReporterService()

	// Initialize CLI
	if err := cli.Execute(configProvider, httpParser, httpExecutor, testRunner, testReporter); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
