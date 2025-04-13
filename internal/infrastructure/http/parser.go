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

// Parser represents an HTTP file parser
type Parser struct {
	commentPattern *regexp.Regexp
	tagPattern     *regexp.Regexp
	namePattern    *regexp.Regexp
	headerPattern  *regexp.Regexp
	methodPattern  *regexp.Regexp
}

// NewParser creates a new HTTP file parser
func NewParser() *Parser {
	return &Parser{
		commentPattern: regexp.MustCompile(`^#\s*(.*)$`),
		tagPattern:     regexp.MustCompile(`^@tag\s+(.+)$`),
		namePattern:    regexp.MustCompile(`^@name\s+(.+)$`),
		headerPattern:  regexp.MustCompile(`^([^:]+):\s*(.+)$`),
		methodPattern:  regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(.+)$`),
	}
}

// ParseFile parses an HTTP file from the file system
func (p *Parser) ParseFile(filePath string) (*models.HTTPFile, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create an HTTP file with the filename
	httpFile := &models.HTTPFile{
		Filename: filePath,
		Requests: []models.HTTPRequest{},
	}

	scanner := bufio.NewScanner(file)
	var currentRequest *models.HTTPRequest
	var currentBody []string
	var readingBody bool
	var comments []string

	// Parse the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a request separator
		if strings.HasPrefix(line, "###") {
			// Save the current request if there is one
			if currentRequest != nil {
				currentRequest.Body = strings.Join(currentBody, "\n")
				httpFile.Requests = append(httpFile.Requests, *currentRequest)
			}

			// Reset for the next request
			currentRequest = nil
			currentBody = []string{}
			readingBody = false
			comments = []string{}
			continue
		}

		// Check if this is a comment
		if matches := p.commentPattern.FindStringSubmatch(line); len(matches) > 1 {
			comment := matches[1]
			comments = append(comments, comment)
			continue
		}

		// Check if this is a tag
		if matches := p.tagPattern.FindStringSubmatch(line); len(matches) > 1 {
			tag := matches[1]
			if currentRequest == nil {
				// Create a new request if needed
				currentRequest = &models.HTTPRequest{
					Comments: comments,
					Tag:      tag,
					Path:     filePath,
				}
				comments = []string{}
			} else {
				currentRequest.Tag = tag
			}
			continue
		}

		// Check if this is a name
		if matches := p.namePattern.FindStringSubmatch(line); len(matches) > 1 {
			name := matches[1]
			if currentRequest == nil {
				// Create a new request if needed
				currentRequest = &models.HTTPRequest{
					Comments: comments,
					Name:     name,
					Path:     filePath,
				}
				comments = []string{}
			} else {
				currentRequest.Name = name
			}
			continue
		}

		// Handle HTTP method line (GET, POST, etc.)
		if matches := p.methodPattern.FindStringSubmatch(line); len(matches) > 2 && !readingBody {
			// Save the current request if there is one
			if currentRequest != nil {
				currentRequest.Body = strings.Join(currentBody, "\n")
				httpFile.Requests = append(httpFile.Requests, *currentRequest)
			}

			// Create a new request
			method := matches[1]
			url := matches[2]

			currentRequest = &models.HTTPRequest{
				Method:   method,
				URL:      url,
				Headers:  []models.HTTPHeader{},
				Comments: comments,
				Path:     filePath,
			}

			// If no explicit name was set, use the path as the name
			if currentRequest.Name == "" {
				currentRequest.Name = fmt.Sprintf("%s %s", method, p.simplifyPath(url))
			}

			// Reset for the new request
			currentBody = []string{}
			readingBody = false
			comments = []string{}
			continue
		}

		// Handle headers if not already reading the body
		if !readingBody && currentRequest != nil {
			if matches := p.headerPattern.FindStringSubmatch(line); len(matches) > 2 {
				name := matches[1]
				value := matches[2]
				currentRequest.Headers = append(currentRequest.Headers, models.HTTPHeader{
					Name:  name,
					Value: value,
				})
				continue
			}
		}

		// Empty line after headers marks the start of the body or the end of the request
		if line == "" && currentRequest != nil && !readingBody {
			if len(currentRequest.Headers) > 0 {
				readingBody = true
			}
			continue
		}

		// If we're reading the body, add the line to the body
		if readingBody && currentRequest != nil {
			currentBody = append(currentBody, line)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Save the last request if there is one
	if currentRequest != nil {
		currentRequest.Body = strings.Join(currentBody, "\n")
		httpFile.Requests = append(httpFile.Requests, *currentRequest)
	}

	return httpFile, nil
}

// simplifyPath returns a simplified version of a URL path for use as a name
func (p *Parser) simplifyPath(url string) string {
	// Remove query parameters
	path := strings.SplitN(url, "?", 2)[0]

	// Extract just the path part if there's a full URL
	if strings.HasPrefix(path, "http") {
		parts := strings.SplitN(path, "/", 4)
		if len(parts) >= 4 {
			path = "/" + parts[3]
		}
	}

	// Remove trailing slash
	path = strings.TrimSuffix(path, "/")

	// If path is empty, use the base URL
	if path == "" {
		path = "root"
	}

	// Replace special characters
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")

	return path
}

// ParseDirectory parses all .http files in a directory
func (p *Parser) ParseDirectory(dirPath string) ([]*models.HTTPFile, error) {
	// Find all .http files in the directory
	matches, err := filepath.Glob(filepath.Join(dirPath, "*.http"))
	if err != nil {
		return nil, fmt.Errorf("failed to list HTTP files: %w", err)
	}

	// Parse each file
	var files []*models.HTTPFile
	for _, match := range matches {
		file, err := p.ParseFile(match)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", match, err)
		}
		files = append(files, file)
	}

	return files, nil
}

// ParseContent parses raw HTTP file content into requests
func (p *Parser) ParseContent(content []byte, filePath string) ([]models.HTTPRequest, error) {
	// Create a temporary file with the content
	tempDir, err := os.MkdirTemp("", "swagger-to-http")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "temp.http")
	if err := os.WriteFile(tempFile, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	// Parse the temporary file
	httpFile, err := p.ParseFile(tempFile)
	if err != nil {
		return nil, err
	}

	// Update the file paths to the original
	for i := range httpFile.Requests {
		httpFile.Requests[i].Path = filePath
	}

	return httpFile.Requests, nil
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
