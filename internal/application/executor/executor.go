package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Service implements the application.HTTPExecutor interface
type Service struct {
	client        *http.Client
	sessionStore  SessionStore
	variableStore VariableStore
	logger        Logger
}

// ServiceOption defines a function type for configuring the Service
type ServiceOption func(*Service)

// WithClient sets a custom HTTP client
func WithClient(client *http.Client) ServiceOption {
	return func(s *Service) {
		s.client = client
	}
}

// WithSessionStore sets a custom session store
func WithSessionStore(store SessionStore) ServiceOption {
	return func(s *Service) {
		s.sessionStore = store
	}
}

// WithVariableStore sets a custom variable store
func WithVariableStore(store VariableStore) ServiceOption {
	return func(s *Service) {
		s.variableStore = store
	}
}

// WithLogger sets a custom logger
func WithLogger(logger Logger) ServiceOption {
	return func(s *Service) {
		s.logger = logger
	}
}

// NewService creates a new HTTP executor service
func NewService(options ...ServiceOption) *Service {
	// Default configuration
	service := &Service{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow up to 10 redirects
				if len(via) >= 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				return nil
			},
		},
		sessionStore:  newMemorySessionStore(),
		variableStore: newMemoryVariableStore(),
		logger:        newDefaultLogger(),
	}

	// Apply options
	for _, option := range options {
		option(service)
	}

	return service
}

// Execute executes an HTTP request and returns the response
func (s *Service) Execute(ctx context.Context, request *models.HTTPRequest, variables map[string]string) (*models.HTTPResponse, error) {
	// Check if the context is already cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Combine global variables and request-specific variables
	allVariables := s.combineVariables(variables)

	// Apply variable substitution to URL, headers, and body
	url := s.applyVariableSubstitution(request.URL, allVariables)
	body := s.applyVariableSubstitution(request.Body, allVariables)

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, request.Method, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for _, header := range request.Headers {
		headerName := s.applyVariableSubstitution(header.Name, allVariables)
		headerValue := s.applyVariableSubstitution(header.Value, allVariables)
		req.Header.Add(headerName, headerValue)
	}

	// Apply session cookies if using a session
	if err := s.applySessionCookies(req); err != nil {
		return nil, fmt.Errorf("failed to apply session cookies: %w", err)
	}

	// Set request ID for tracking
	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
	s.logger.Debugf("Executing request %s: %s %s", requestID, request.Method, url)

	// Execute the request
	startTime := time.Now()
	resp, err := s.client.Do(req)
	duration := time.Since(startTime)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract content type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Store any cookies for session management
	s.storeSessionCookies(resp)

	// Extract variables from the response if configured
	s.extractVariables(resp, responseBody, requestID)

	// Create HTTP response
	httpResponse := &models.HTTPResponse{
		StatusCode:     resp.StatusCode,
		Status:         resp.Status,
		Headers:        convertHTTPHeaders(resp.Header),
		Body:           responseBody,
		ContentType:    contentType,
		ContentLength:  resp.ContentLength,
		Duration:       duration,
		Request:        request,
		RequestID:      requestID,
		Timestamp:      time.Now(),
	}

	s.logger.Debugf("Response for %s: %s (%d bytes) in %s", requestID, resp.Status, len(responseBody), duration)
	return httpResponse, nil
}

// ExecuteFile executes all requests in an HTTP file
func (s *Service) ExecuteFile(ctx context.Context, file *models.HTTPFile, variables map[string]string) ([]*models.HTTPResponse, error) {
	// Store responses for return
	responses := make([]*models.HTTPResponse, 0, len(file.Requests))

	// Execute each request in sequence
	for _, request := range file.Requests {
		select {
		case <-ctx.Done():
			return responses, ctx.Err()
		default:
			// Execute the request
			response, err := s.Execute(ctx, &request, variables)
			if err != nil {
				s.logger.Errorf("Failed to execute request: %s %s: %v", request.Method, request.Path, err)
				// Continue with the next request even if this one failed
				continue
			}
			responses = append(responses, response)
		}
	}

	return responses, nil
}

// combineVariables combines global variables and request-specific variables
func (s *Service) combineVariables(requestVars map[string]string) map[string]string {
	// Start with global variables
	allVars := s.variableStore.GetAll()

	// Add request-specific variables, overriding global ones if needed
	for k, v := range requestVars {
		allVars[k] = v
	}

	return allVars
}

// applyVariableSubstitution replaces variable references in a string
// Variable format: {{variableName}}
func (s *Service) applyVariableSubstitution(input string, variables map[string]string) string {
	result := input
	for name, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", name)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// applySessionCookies adds session cookies to the request if using a session
func (s *Service) applySessionCookies(req *http.Request) error {
	cookies := s.sessionStore.GetCookies(req.URL.Host)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	return nil
}

// storeSessionCookies stores cookies from a response for session management
func (s *Service) storeSessionCookies(resp *http.Response) {
	for _, cookie := range resp.Cookies() {
		s.sessionStore.SetCookie(resp.Request.URL.Host, cookie)
	}
}

// extractVariables extracts variables from the response based on configuration
func (s *Service) extractVariables(resp *http.Response, body []byte, requestID string) {
	// TODO: Implement variable extraction from response (headers, body, etc.)
	// This will be implemented as part of the advanced testing features
}

// convertHTTPHeaders converts http.Header to our model's map[string][]string format
func convertHTTPHeaders(headers http.Header) map[string][]string {
	result := make(map[string][]string, len(headers))
	for name, values := range headers {
		result[name] = values
	}
	return result
}
