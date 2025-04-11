package snapshot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// Manager handles the saving, loading, and comparison of snapshots
type Manager struct {
	basePath   string
	formatters map[string]ResponseFormatter
}

// NewManager creates a new snapshot manager
func NewManager(basePath string) *Manager {
	manager := &Manager{
		basePath:   basePath,
		formatters: make(map[string]ResponseFormatter),
	}

	// Register default formatters
	manager.RegisterFormatter("application/json", &JSONFormatter{})
	manager.RegisterFormatter("application/xml", &XMLFormatter{})
	manager.RegisterFormatter("text/html", &HTMLFormatter{})
	manager.RegisterFormatter("text/plain", &TextFormatter{})
	manager.RegisterFormatter("*/*", &BinaryFormatter{})

	return manager
}

// RegisterFormatter registers a formatter for a specific content type
func (m *Manager) RegisterFormatter(contentType string, formatter ResponseFormatter) {
	m.formatters[contentType] = formatter
}

// SaveSnapshot saves a response as a snapshot
func (m *Manager) SaveSnapshot(ctx context.Context, response *models.HTTPResponse, path string) error {
	if response == nil {
		return fmt.Errorf("cannot save nil response")
	}

	// Ensure directory exists
	snapshotDir := filepath.Join(m.basePath, filepath.Dir(path))
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Generate filename for the snapshot
	snapshotPath := m.getSnapshotPath(path, response)

	// Format response for storage
	formatter := m.getFormatter(response.ContentType)
	formattedBody, err := formatter.Format(response.Body)
	if err != nil {
		return fmt.Errorf("failed to format response body: %w", err)
	}

	// Create snapshot data structure
	snapshot := models.SnapshotData{
		Metadata: models.SnapshotMetadata{
			RequestPath:   response.Request.Path,
			RequestMethod: response.Request.Method,
			ContentType:   response.ContentType,
			StatusCode:    response.StatusCode,
			Headers:       response.Headers,
			CreatedAt:     time.Now(),
		},
		Content: string(formattedBody),
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize snapshot: %w", err)
	}

	// Write to file
	if err := os.WriteFile(snapshotPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// LoadSnapshot loads a snapshot from the file system
func (m *Manager) LoadSnapshot(ctx context.Context, path string) (*models.HTTPResponse, error) {
	snapshotPath := filepath.Join(m.basePath, path)

	// Check if snapshot exists
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("snapshot does not exist: %s", path)
	}

	// Read snapshot file
	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	// Deserialize
	var snapshot models.SnapshotData
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot file: %w", err)
	}

	// Create HTTP Response from snapshot
	response := &models.HTTPResponse{
		StatusCode:    snapshot.Metadata.StatusCode,
		Status:        fmt.Sprintf("%d %s", snapshot.Metadata.StatusCode, statusText(snapshot.Metadata.StatusCode)),
		Headers:       snapshot.Metadata.Headers,
		ContentType:   snapshot.Metadata.ContentType,
		ContentLength: int64(len(snapshot.Content)),
		Body:          []byte(snapshot.Content),
		Request: &models.HTTPRequest{
			Method: snapshot.Metadata.RequestMethod,
			Path:   snapshot.Metadata.RequestPath,
		},
		Timestamp: snapshot.Metadata.CreatedAt,
	}

	return response, nil
}

// CompareWithSnapshot compares a response with a snapshot
func (m *Manager) CompareWithSnapshot(ctx context.Context, response *models.HTTPResponse, snapshotPath string) (*models.SnapshotDiff, error) {
	// Load snapshot
	snapshotResponse, err := m.LoadSnapshot(ctx, snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load snapshot for comparison: %w", err)
	}

	// Create diff result
	diff := &models.SnapshotDiff{
		RequestPath:   response.Request.Path,
		RequestMethod: response.Request.Method,
		Equal:         true,
	}

	// Compare status
	diff.StatusDiff = compareStatus(snapshotResponse.StatusCode, response.StatusCode)
	if !diff.StatusDiff.Equal {
		diff.Equal = false
	}

	// Compare headers
	diff.HeaderDiff = compareHeaders(snapshotResponse.Headers, response.Headers)
	if !diff.HeaderDiff.Equal {
		diff.Equal = false
	}

	// Compare body
	formatter := m.getFormatter(response.ContentType)
	diff.BodyDiff = formatter.Compare(snapshotResponse.Body, response.Body)
	if !diff.BodyDiff.Equal {
		diff.Equal = false
	}

	return diff, nil
}

// getFormatter returns the appropriate formatter for the given content type
func (m *Manager) getFormatter(contentType string) ResponseFormatter {
	// Normalize content type by removing parameters
	normalizedType := contentType
	if idx := strings.Index(contentType, ";"); idx != -1 {
		normalizedType = contentType[:idx]
	}

	// Try to find an exact match
	if formatter, ok := m.formatters[normalizedType]; ok {
		return formatter
	}

	// Try to find a match by type/subtype
	parts := strings.Split(normalizedType, "/")
	if len(parts) == 2 {
		wildcardType := parts[0] + "/*"
		if formatter, ok := m.formatters[wildcardType]; ok {
			return formatter
		}
	}

	// Default to binary formatter
	return m.formatters["*/*"]
}

// getSnapshotPath generates a path for the snapshot file
func (m *Manager) getSnapshotPath(path string, response *models.HTTPResponse) string {
	// Create a sanitized filename
	method := strings.ToLower(response.Request.Method)
	
	// Remove the extension if present
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := base
	if ext != "" {
		name = base[:len(base)-len(ext)]
	}
	
	// Replace non-alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitizedName := re.ReplaceAllString(name, "_")
	
	// Create the path
	filename := fmt.Sprintf("%s_%s.snap.json", sanitizedName, method)
	return filepath.Join(m.basePath, filepath.Dir(path), filename)
}

// Helper function to get status text for a status code
func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

// compareStatus compares two status codes
func compareStatus(expected, actual int) *models.StatusDiff {
	return &models.StatusDiff{
		Expected: expected,
		Actual:   actual,
		Equal:    expected == actual,
	}
}

// compareHeaders compares two sets of headers
func compareHeaders(expected, actual map[string][string) *models.HeaderDiff {
	diff := &models.HeaderDiff{
		MissingHeaders:   make(map[string][]string),
		ExtraHeaders:     make(map[string][]string),
		DifferentValues:  make(map[string]models.HeaderValueDiff),
		Equal:            true,
	}

	// Check for missing or different headers
	for key, expectedValues := range expected {
		// Normalize header key for comparison
		normalizedKey := strings.ToLower(key)
		
		// Look for the header in the actual headers
		actualValues, found := findHeader(actual, normalizedKey)
		
		if !found {
			// Header is missing in actual
			diff.MissingHeaders[key] = expectedValues
			diff.Equal = false
		} else if !compareStringSlices(expectedValues, actualValues) {
			// Header exists but values differ
			diff.DifferentValues[key] = models.HeaderValueDiff{
				Expected: expectedValues,
				Actual:   actualValues,
			}
			diff.Equal = false
		}
	}

	// Check for extra headers in actual
	for key, actualValues := range actual {
		normalizedKey := strings.ToLower(key)
		_, found := findHeader(expected, normalizedKey)
		
		if !found {
			diff.ExtraHeaders[key] = actualValues
			diff.Equal = false
		}
	}

	return diff
}

// findHeader finds a header by its normalized key
func findHeader(headers map[string][]string, normalizedKey string) ([]string, bool) {
	for key, values := range headers {
		if strings.ToLower(key) == normalizedKey {
			return values, true
		}
	}
	return nil, false
}

// compareStringSlices compares two string slices regardless of order
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create frequency maps
	mapA := make(map[string]int)
	mapB := make(map[string]int)

	for _, val := range a {
		mapA[val]++
	}

	for _, val := range b {
		mapB[val]++
	}

	// Compare frequency maps
	for val, countA := range mapA {
		countB, exists := mapB[val]
		if !exists || countA != countB {
			return false
		}
	}

	return true
}
