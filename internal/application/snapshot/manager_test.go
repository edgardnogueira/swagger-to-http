package snapshot

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_SaveAndLoadSnapshot(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager
	manager := NewManager(tempDir)

	// Create a test response
	response := &models.HTTPResponse{
		StatusCode:    200,
		Status:        "200 OK",
		Headers:       map[string][]string{"Content-Type": {"application/json"}},
		Body:          []byte(`{"name":"test","value":123}`),
		ContentType:   "application/json",
		ContentLength: 27,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Save the snapshot
	err = manager.SaveSnapshot(context.Background(), response, "api/test.http")
	require.NoError(t, err)

	// Verify the snapshot file exists
	snapshotPath := filepath.Join(tempDir, "api", "test_get.snap.json")
	_, err = os.Stat(snapshotPath)
	require.NoError(t, err)

	// Load the snapshot
	loadedResponse, err := manager.LoadSnapshot(context.Background(), "api/test_get.snap.json")
	require.NoError(t, err)

	// Verify the loaded response matches the original
	assert.Equal(t, response.StatusCode, loadedResponse.StatusCode)
	assert.Equal(t, response.ContentType, loadedResponse.ContentType)
	assert.Equal(t, response.Request.Method, loadedResponse.Request.Method)
	assert.Equal(t, response.Request.Path, loadedResponse.Request.Path)
	assert.JSONEq(t, string(response.Body), string(loadedResponse.Body))
}

func TestManager_CompareWithSnapshot_Equal(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager
	manager := NewManager(tempDir)

	// Create a test response
	response := &models.HTTPResponse{
		StatusCode:    200,
		Status:        "200 OK",
		Headers:       map[string][]string{"Content-Type": {"application/json"}},
		Body:          []byte(`{"name":"test","value":123}`),
		ContentType:   "application/json",
		ContentLength: 27,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Save the snapshot
	err = manager.SaveSnapshot(context.Background(), response, "api/test.http")
	require.NoError(t, err)

	// Compare with the same response
	diff, err := manager.CompareWithSnapshot(context.Background(), response, "api/test_get.snap.json")
	require.NoError(t, err)

	// Verify the comparison shows equality
	assert.True(t, diff.Equal)
	assert.True(t, diff.StatusDiff.Equal)
	assert.True(t, diff.HeaderDiff.Equal)
	assert.True(t, diff.BodyDiff.Equal)
}

func TestManager_CompareWithSnapshot_Different(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager
	manager := NewManager(tempDir)

	// Create a test response
	response1 := &models.HTTPResponse{
		StatusCode:    200,
		Status:        "200 OK",
		Headers:       map[string][]string{"Content-Type": {"application/json"}},
		Body:          []byte(`{"name":"test","value":123}`),
		ContentType:   "application/json",
		ContentLength: 27,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Save the snapshot
	err = manager.SaveSnapshot(context.Background(), response1, "api/test.http")
	require.NoError(t, err)

	// Create a different response
	response2 := &models.HTTPResponse{
		StatusCode:    201,
		Status:        "201 Created",
		Headers:       map[string][]string{"Content-Type": {"application/json"}, "Location": {"/api/test/1"}},
		Body:          []byte(`{"name":"test","value":456}`),
		ContentType:   "application/json",
		ContentLength: 27,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Compare different response with snapshot
	diff, err := manager.CompareWithSnapshot(context.Background(), response2, "api/test_get.snap.json")
	require.NoError(t, err)

	// Verify the comparison shows differences
	assert.False(t, diff.Equal)
	assert.False(t, diff.StatusDiff.Equal)
	assert.False(t, diff.HeaderDiff.Equal)
	assert.False(t, diff.BodyDiff.Equal)

	// Verify specific differences
	assert.Equal(t, 200, diff.StatusDiff.Expected)
	assert.Equal(t, 201, diff.StatusDiff.Actual)
	assert.Contains(t, diff.HeaderDiff.ExtraHeaders, "Location")
	assert.NotNil(t, diff.BodyDiff.JsonDiff)
	assert.Contains(t, diff.BodyDiff.DiffContent, "456") // Should show value difference
}

func TestManager_ListSnapshots(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager
	manager := NewManager(tempDir)

	// Create test responses
	response1 := &models.HTTPResponse{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        []byte(`{"id":1}`),
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test/1",
		},
	}

	response2 := &models.HTTPResponse{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        []byte(`{"id":2}`),
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test/2",
		},
	}

	// Save snapshots
	err = manager.SaveSnapshot(context.Background(), response1, "api/test1.http")
	require.NoError(t, err)
	err = manager.SaveSnapshot(context.Background(), response2, "api/test2.http")
	require.NoError(t, err)

	// List snapshots
	snapshots, err := manager.ListSnapshots(context.Background(), "api")
	require.NoError(t, err)

	// Verify we have two snapshots
	assert.Len(t, snapshots, 2)
	assert.Contains(t, snapshots, "api/test1_get.snap.json")
	assert.Contains(t, snapshots, "api/test2_get.snap.json")
}

func TestManager_CleanupSnapshots(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager
	manager := NewManager(tempDir)

	// Create test responses
	response1 := &models.HTTPResponse{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        []byte(`{"id":1}`),
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test/1",
		},
	}

	response2 := &models.HTTPResponse{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        []byte(`{"id":2}`),
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test/2",
		},
	}

	// Save snapshots
	err = manager.SaveSnapshot(context.Background(), response1, "api/test1.http")
	require.NoError(t, err)
	err = manager.SaveSnapshot(context.Background(), response2, "api/test2.http")
	require.NoError(t, err)

	// Create a used snapshots map with only one snapshot
	usedSnapshots := map[string]bool{
		"api/test1_get.snap.json": true,
	}

	// Cleanup unused snapshots
	err = manager.CleanupSnapshots(context.Background(), "api", usedSnapshots)
	require.NoError(t, err)

	// List snapshots after cleanup
	snapshots, err := manager.ListSnapshots(context.Background(), "api")
	require.NoError(t, err)

	// Verify only the used snapshot remains
	assert.Len(t, snapshots, 1)
	assert.Contains(t, snapshots, "api/test1_get.snap.json")
	assert.NotContains(t, snapshots, "api/test2_get.snap.json")
}

func TestManager_WithOptions(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager with options
	manager := NewManager(tempDir)
	manager = manager.WithOptions(models.SnapshotOptions{
		UpdateMode:    "all",
		IgnoreHeaders: []string{"Date", "Server"},
	})

	// Create a test response
	response := &models.HTTPResponse{
		StatusCode:    200,
		Status:        "200 OK",
		Headers:       map[string][]string{"Content-Type": {"application/json"}},
		Body:          []byte(`{"name":"test","value":123}`),
		ContentType:   "application/json",
		ContentLength: 27,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Save the snapshot
	err = manager.SaveSnapshot(context.Background(), response, "api/test.http")
	require.NoError(t, err)

	// Create a different response with ignored headers
	response2 := &models.HTTPResponse{
		StatusCode:  200,
		Status:      "200 OK",
		Headers:     map[string][]string{
			"Content-Type": {"application/json"}, 
			"Date":         {"Mon, 01 Jan 2023 12:00:00 GMT"},
			"Server":       {"test-server"},
		},
		Body:        []byte(`{"name":"test","value":123}`),
		ContentType: "application/json",
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
	}

	// Compare with snapshot - should ignore Date and Server headers
	diff, err := manager.CompareWithSnapshot(context.Background(), response2, "api/test_get.snap.json")
	require.NoError(t, err)

	// Should be equal because the difference is only in ignored headers
	assert.True(t, diff.Equal)
}
