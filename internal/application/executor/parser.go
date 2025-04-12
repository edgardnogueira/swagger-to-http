package executor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// HTTPFileParser is responsible for parsing HTTP files
type HTTPFileParser struct {
	logger Logger
}

// NewHTTPFileParser creates a new HTTP file parser
func NewHTTPFileParser(logger Logger) *HTTPFileParser {
	if logger == nil {
		logger = newDefaultLogger()
	}
	return &HTTPFileParser{
		logger: logger,
	}
}

// ParseFile parses an HTTP file from the filesystem
func (p *HTTPFileParser) ParseFile(filePath string) (*models.HTTPFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create HTTP file with filename from path
	httpFile := &models.HTTPFile{
		Filename: filepath.Base(filePath),
	}

	// Parse the file content
	requests, err := p.parseRequests(file)
	if err != nil {
		return nil, err
	}
	httpFile.Requests = requests

	return httpFile, nil
}

// ParseString parses an HTTP file from a string
func (p *HTTPFileParser) ParseString(content, filename string) (*models.HTTPFile, error) {
	// Create HTTP file with given filename
	httpFile := &models.HTTPFile{
		Filename: filename,
	}

	// Parse the content
	requests, err := p.parseRequests(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	httpFile.Requests = requests

	return httpFile, nil
}

// parseRequests parses HTTP requests from a reader
func (p *HTTPFileParser) parseRequests(reader io.Reader) ([]models.HTTPRequest, error) {
	var requests []models.HTTPRequest
	var currentRequest *models.HTTPRequest
	var currentComments []string
	var bodyLines []string
	var inBody bool
	var boundaryLine string

	scanner := bufio.NewScanner(reader)
	lineNum := 0

	// Regex for HTTP method and URL
	methodURLRegex := regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+(.+)$`)
	// Regex for header
	headerRegex := regexp.MustCompile(`^([^:]+):\s*(.*)$`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check for request boundary
		if strings.HasPrefix(trimmedLine, "###") {
			// Finalize current request if we have one
			if currentRequest != nil {
				if inBody && len(bodyLines) > 0 {
					currentRequest.Body = strings.Join(bodyLines, "\n")
				}
				requests = append(requests, *currentRequest)
			}

			// Start a new request
			currentRequest = nil
			currentComments = nil
			bodyLines = nil
			inBody = false
			boundaryLine = trimmedLine
			continue
		}

		// If we're not currently processing a request
		if currentRequest == nil {
			// Check if line starts with a comment
			if strings.HasPrefix(trimmedLine, "//") || strings.HasPrefix(trimmedLine, "#") {
				commentText := strings.TrimSpace(trimmedLine[2:])
				if strings.HasPrefix(trimmedLine, "#") {
					commentText = strings.TrimSpace(trimmedLine[1:])
				}
				currentComments = append(currentComments, commentText)
				continue
			}

			// Check if line contains a method and URL
			matches := methodURLRegex.FindStringSubmatch(trimmedLine)
			if len(matches) == 3 {
				// Create a new request
				currentRequest = &models.HTTPRequest{
					Method:   matches[1],
					URL:      matches[2],
					Headers:  []models.HTTPHeader{},
					Comments: currentComments,
				}
				// Extract path from URL
				path, err := extractPath(currentRequest.URL)
				if err != nil {
					p.logger.Warnf("Failed to extract path from URL %s: %v", currentRequest.URL, err)
				} else {
					currentRequest.Path = path
				}
				currentComments = nil
				continue
			}

			// Ignore empty lines when not in a request
			if trimmedLine == "" {
				continue
			}

			// If we got here, there's content outside of a request
			p.logger.Warnf("Line %d: Content outside of request boundary: %s", lineNum, trimmedLine)
			continue
		}

		// We're processing a request
		if !inBody {
			// Check if this is a header
			headerMatches := headerRegex.FindStringSubmatch(trimmedLine)
			if len(headerMatches) == 3 {
				// Add header to current request
				currentRequest.Headers = append(currentRequest.Headers, models.HTTPHeader{
					Name:  headerMatches[1],
					Value: headerMatches[2],
				})
				continue
			}

			// Empty line marks the end of headers and start of body
			if trimmedLine == "" {
				inBody = true
				continue
			}

			// If we got here, there's unrecognized content in the headers section
			p.logger.Warnf("Line %d: Unrecognized content in headers: %s", lineNum, trimmedLine)
			continue
		}

		// We're in the body, add all content
		bodyLines = append(bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning HTTP file: %w", err)
	}

	// Finalize the last request if we have one
	if currentRequest != nil {
		if inBody && len(bodyLines) > 0 {
			currentRequest.Body = strings.Join(bodyLines, "\n")
		}
		requests = append(requests, *currentRequest)
	}

	return requests, nil
}

// extractPath extracts the path from a URL
// e.g., "https://api.example.com/users/1?q=test" -> "/users/1"
func extractPath(url string) (string, error) {
	// Remove query string if present
	pathPart := url
	if i := strings.Index(url, "?"); i >= 0 {
		pathPart = url[:i]
	}

	// Remove protocol and domain
	if strings.HasPrefix(pathPart, "http://") || strings.HasPrefix(pathPart, "https://") {
		if i := strings.Index(pathPart[8:], "/"); i >= 0 {
			// For http://, add 7 to the index; for https://, add 8
			slashIndex := i
			if strings.HasPrefix(pathPart, "https://") {
				slashIndex += 8
			} else {
				slashIndex += 7
			}
			pathPart = pathPart[slashIndex:]
		} else {
			pathPart = "/"
		}
	}

	return pathPart, nil
}
