package main

import (
	"fmt"
	"os"

	"github.com/edgardnogueira/swagger-to-http/internal/cli"
	"github.com/edgardnogueira/swagger-to-http/internal/version"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
