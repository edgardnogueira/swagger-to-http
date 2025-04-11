package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/application/generator"
	"github.com/edgardnogueira/swagger-to-http/internal/application/parser"
	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/config"
	"github.com/edgardnogueira/swagger-to-http/internal/infrastructure/fs"
	"github.com/spf13/cobra"
)

var (
	inputFile    string
	inputURL     string
	outputDir    string
	baseURL      string
	defaultTag   string
	indentJSON   bool
	includeAuth  bool
	authHeader   string
	authToken    string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate HTTP files from a Swagger/OpenAPI document",
	Long: `Generate HTTP request files from a Swagger/OpenAPI document.
This command parses the document and creates .http files organized by tags.`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Required flags
	generateCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Swagger/OpenAPI file to process (required if url not provided)")
	generateCmd.Flags().StringVarP(&inputURL, "url", "u", "", "URL to Swagger/OpenAPI document (required if file not provided)")
	
	// Optional flags with default values from config
	cp := config.NewConfigProvider()
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", cp.GetString("output.directory"), "Output directory for HTTP files")
	generateCmd.Flags().StringVarP(&baseURL, "base-url", "b", cp.GetString("generator.base_url"), "Base URL for requests (overrides the one in the Swagger doc)")
	generateCmd.Flags().StringVarP(&defaultTag, "default-tag", "t", cp.GetString("generator.default_tag"), "Default tag for operations without tags")
	generateCmd.Flags().BoolVarP(&indentJSON, "indent-json", "i", cp.GetBool("generator.indent_json"), "Indent JSON in request bodies")
	generateCmd.Flags().BoolVar(&includeAuth, "auth", cp.GetBool("generator.include_auth"), "Include authentication header in requests")
	generateCmd.Flags().StringVar(&authHeader, "auth-header", cp.GetString("generator.auth_header"), "Authentication header name")
	generateCmd.Flags().StringVar(&authToken, "auth-token", cp.GetString("generator.auth_token"), "Authentication token value")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Validate input parameters
	if inputFile == "" && inputURL == "" {
		return fmt.Errorf("either --file or --url must be provided")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create parser
	swaggerParser := parser.NewSwaggerParser()

	// Parse document
	var swaggerDoc, err = parseDocument(ctx, swaggerParser, inputFile, inputURL)
	if err != nil {
		return err
	}

	// Create generator with options
	httpGenerator := generator.NewHTTPGenerator(
		generator.WithBaseURL(baseURL),
		generator.WithDefaultTag(defaultTag),
		generator.WithIndentJSON(indentJSON),
		generator.WithAuth(includeAuth, authHeader, authToken),
	)

	// Generate HTTP requests
	log.Println("Generating HTTP requests...")
	collection, err := httpGenerator.Generate(ctx, swaggerDoc)
	if err != nil {
		return fmt.Errorf("failed to generate HTTP requests: %w", err)
	}

	// Set the output directory
	collection.RootDir = outputDir

	// Create file writer
	fileWriter := fs.NewFileWriter()

	// Write the collection to files
	log.Printf("Writing HTTP files to directory: %s\n", outputDir)
	if err := fileWriter.WriteCollection(ctx, collection); err != nil {
		return fmt.Errorf("failed to write HTTP files: %w", err)
	}

	log.Println("Successfully generated HTTP files!")
	return nil
}

// parseDocument parses a Swagger/OpenAPI document from a file or URL
func parseDocument(ctx context.Context, swaggerParser *parser.SwaggerParser, filePath, url string) (*models.SwaggerDoc, error) {
	if filePath != "" {
		log.Printf("Parsing Swagger file: %s\n", filePath)
		
		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", filePath)
		}
		
		return swaggerParser.ParseFile(ctx, filePath)
	}
	
	log.Printf("Parsing Swagger from URL: %s\n", url)
	return swaggerParser.ParseURL(ctx, url)
}
