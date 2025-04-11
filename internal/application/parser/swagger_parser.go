package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"gopkg.in/yaml.v3"
)

// SwaggerParser implements the SwaggerParser interface
type SwaggerParser struct{}

// NewSwaggerParser creates a new SwaggerParser
func NewSwaggerParser() *SwaggerParser {
	return &SwaggerParser{}
}

// Parse parses a Swagger/OpenAPI document from a byte array
func (p *SwaggerParser) Parse(ctx context.Context, data []byte) (*models.SwaggerDoc, error) {
	var doc models.SwaggerDoc

	// Try JSON first
	err := json.Unmarshal(data, &doc)
	if err == nil {
		return &doc, p.Validate(ctx, &doc)
	}

	// Try YAML if JSON fails
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document as JSON or YAML: %w", err)
	}

	return &doc, p.Validate(ctx, &doc)
}

// ParseFile parses a Swagger/OpenAPI document from a file
func (p *SwaggerParser) ParseFile(ctx context.Context, filePath string) (*models.SwaggerDoc, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return p.Parse(ctx, data)
}

// ParseURL parses a Swagger/OpenAPI document from a URL
func (p *SwaggerParser) ParseURL(ctx context.Context, url string) (*models.SwaggerDoc, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response from URL %s: %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from URL %s: %w", url, err)
	}

	return p.Parse(ctx, data)
}

// Validate validates a Swagger/OpenAPI document
func (p *SwaggerParser) Validate(ctx context.Context, doc *models.SwaggerDoc) error {
	if doc == nil {
		return errors.New("document is nil")
	}

	// Check if it's OpenAPI 3.0.x or Swagger 2.0
	if doc.Version == "" && doc.SwaggerVersion == "" {
		return errors.New("invalid document: missing openapi or swagger version")
	}

	// Validate basic required fields
	if doc.Info.Title == "" {
		return errors.New("invalid document: missing info.title")
	}

	if doc.Info.Version == "" {
		return errors.New("invalid document: missing info.version")
	}

	if len(doc.Paths) == 0 {
		return errors.New("invalid document: no paths defined")
	}

	return nil
}

// DetectFormat detects the format of a Swagger/OpenAPI document
func DetectFormat(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	switch ext {
	case ".json":
		return "json", nil
	case ".yaml", ".yml":
		return "yaml", nil
	default:
		// Try to detect by reading the first few bytes
		f, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer f.Close()
		
		// Read the first 100 bytes to detect format
		buf := make([]byte, 100)
		_, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		
		// Check for JSON pattern
		if strings.HasPrefix(strings.TrimSpace(string(buf)), "{") {
			return "json", nil
		}
		
		// Assume YAML by default
		return "yaml", nil
	}
}
