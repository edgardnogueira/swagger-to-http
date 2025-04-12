package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Executor implements the HTTPExecutor interface for executing HTTP requests
type Executor struct {
	client      *http.Client
	environment map[string]string
}

// NewExecutor creates a new HTTP executor with the given options
func NewExecutor(timeout time.Duration, environment map[string]string) *Executor {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Executor{
		client:      client,
		environment: environment,
	}
}

// Execute executes an HTTP request and returns the response
func (e *Executor) Execute(ctx context.Context, request *models.HTTPRequest, variables map[string]string) (*models.HTTPResponse, error) {
	// Process variables - combine environment variables with request variables
	vars := e.combineVariables(variables)

	// Process request parts with variable substitution
	url := e.processVariables(request.URL, vars)
	body := e.processVariables(request.Body, vars)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, request.Method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	for _, header := range request.Headers {
		req.Header.Add(header.Name, e.processVariables(header.Value, vars))
	}

	// For form submissions, ensure the right content-type if not explicitly set
	if request.Method == "POST" && len(request.Body) > 0 {
		if !hasHeader(request.Headers, "Content-Type") {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Execute the request
	startTime := time.Now()
	resp, err := e.client.Do(req)
	duration := time.Since(startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response object
	response := &models.HTTPResponse{
		StatusCode:    resp.StatusCode,
		Status:        resp.Status,
		Headers:       make(map[string][]string),
		Body:          respBody,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		Duration:      duration,
		Request:       request,
		RequestID:     fmt.Sprintf("%s-%s", request.Method, request.Path),
		Timestamp:     time.Now(),
	}

	// Copy headers
	for name, values := range resp.Header {
		response.Headers[name] = values
	}

	return response, nil
}

// ExecuteFile executes all requests in an HTTP file
func (e *Executor) ExecuteFile(ctx context.Context, file *models.HTTPFile, variables map[string]string) ([]*models.HTTPResponse, error) {
	responses := make([]*models.HTTPResponse, 0, len(file.Requests))

	for _, request := range file.Requests {
		response, err := e.Execute(ctx, &request, variables)
		if err != nil {
			return responses, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// processVariables replaces variable references in the given text with their values
func (e *Executor) processVariables(text string, variables map[string]string) string {
	result := text

	// Replace {{variable}} references
	for name, value := range variables {
		result = strings.ReplaceAll(result, "{{"+name+"}}", value)
	}

	return result
}

// combineVariables merges environment variables with request-specific variables
func (e *Executor) combineVariables(requestVars map[string]string) map[string]string {
	// Create a new map to avoid modifying the environment
	vars := make(map[string]string)

	// Copy environment variables
	for k, v := range e.environment {
		vars[k] = v
	}

	// Add/override with request variables
	if requestVars != nil {
		for k, v := range requestVars {
			vars[k] = v
		}
	}

	return vars
}

// hasHeader checks if a specific header exists in the headers slice
func hasHeader(headers []models.HTTPHeader, name string) bool {
	lowerName := strings.ToLower(name)
	for _, h := range headers {
		if strings.ToLower(h.Name) == lowerName {
			return true
		}
	}
	return false
}
