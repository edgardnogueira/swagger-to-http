package http

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Parser handles parsing of .http files
type Parser struct{}

// NewParser creates a new HTTP file parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses an HTTP file into a collection of requests
func (p *Parser) ParseFile(filePath string) (*models.HTTPFile, error) {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse content
	requests, err := p.ParseContent(content, filePath)
	if err != nil {
		return nil, err
	}

	// Create HTTP file
	httpFile := &models.HTTPFile{
		Filename: filepath.Base(filePath),
		Requests: requests,
	}

	return httpFile, nil
}

// ParseContent parses raw HTTP file content into requests
func (p *Parser) ParseContent(content []byte, filePath string) ([]models.HTTPRequest, error) {
	var requests []models.HTTPRequest
	
	// Use request separator pattern: ###
	separator := []byte("###")
	parts := bytes.Split(content, separator)
	
	for i, part := range parts {
		// Skip empty parts
		if len(bytes.TrimSpace(part)) == 0 {
			continue
		}
		
		// Parse the request
		request, err := p.parseRequest(part, i+1, filePath)
		if err != nil {
			return nil, err
		}
		
		requests = append(requests, *request)
	}
	
	return requests, nil
}

// parseRequest parses a single HTTP request section
func (p *Parser) parseRequest(content []byte, index int, filePath string) (*models.HTTPRequest, error) {
	// Create scanner for reading lines
	scanner := bufio.NewScanner(bytes.NewReader(content))
	
	// Create request with default values
	request := &models.HTTPRequest{
		Name:     fmt.Sprintf("Request %d", index),
		Path:     filepath.Base(filePath),
		Comments: []string{},
		Headers:  []models.HTTPHeader{},
	}
	
	// State variables
	lineNum := 0
	inBody := false
	var bodyLines []string
	
	// Regular expressions for parsing
	requestLineRegex := regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(.+)$`)
	headerRegex := regexp.MustCompile(`^([^:]+):\s*(.*)$`)
	commentRegex := regexp.MustCompile(`^#(.*)$`)
	nameRegex := regexp.MustCompile(`^@name\s+(.+)$`)
	tagRegex := regexp.MustCompile(`^@tag\s+(.+)$`)
	
	// Process each line
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		
		// If we're in body mode, collect body lines
		if inBody {
			bodyLines = append(bodyLines, line)
			continue
		}
		
		// Empty line indicates start of body section
		if strings.TrimSpace(line) == "" {
			inBody = true
			continue
		}
		
		// Try to match different line types
		if match := requestLineRegex.FindStringSubmatch(line); match != nil {
			request.Method = match[1]
			request.URL = match[2]
			continue
		}
		
		if match := headerRegex.FindStringSubmatch(line); match != nil {
			header := models.HTTPHeader{
				Name:  match[1],
				Value: match[2],
			}
			request.Headers = append(request.Headers, header)
			continue
		}
		
		if match := commentRegex.FindStringSubmatch(line); match != nil {
			request.Comments = append(request.Comments, match[1])
			continue
		}
		
		if match := nameRegex.FindStringSubmatch(line); match != nil {
			request.Name = match[1]
			continue
		}
		
		if match := tagRegex.FindStringSubmatch(line); match != nil {
			request.Tag = match[1]
			continue
		}
		
		// If we reach here, we couldn't parse the line
		return nil, fmt.Errorf("invalid line %d: %s", lineNum, line)
	}
	
	// Combine body lines
	if len(bodyLines) > 0 {
		request.Body = strings.Join(bodyLines, "\n")
	}
	
	// Check that we have the minimum required fields
	if request.Method == "" || request.URL == "" {
		return nil, fmt.Errorf("request is missing method or URL")
	}
	
	// Extract path from URL (simplified)
	urlParts := strings.Split(request.URL, "://")
	if len(urlParts) > 1 {
		pathParts := strings.SplitN(urlParts[1], "/", 2)
		if len(pathParts) > 1 {
			request.Path = "/" + pathParts[1]
		}
	} else if strings.HasPrefix(request.URL, "/") {
		request.Path = request.URL
	}
	
	return request, nil
}

// FindHTTPFiles finds all .http files in a directory (or matching a glob pattern)
func (p *Parser) FindHTTPFiles(pattern string) ([]string, error) {
	// Handle glob pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	if len(matches) > 0 {
		return matches, nil
	}

	// If no matches and pattern doesn't have wildcard, try as directory
	if !strings.Contains(pattern, "*") {
		fileInfo, err := os.Stat(pattern)
		if err == nil && fileInfo.IsDir() {
			// It's a directory, find all .http files
			var files []string
			err = filepath.Walk(pattern, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".http") {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to walk directory: %w", err)
			}
			return files, nil
		}
	}

	return []string{}, nil
}
