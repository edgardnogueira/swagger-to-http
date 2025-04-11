package snapshot

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_RunTest(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-service-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager and service
	manager := NewManager(tempDir)
	service := NewService(manager, models.SnapshotOptions{
		BasePath:   tempDir,
		UpdateMode: "all",
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

	// Run a test, should create a new snapshot
	result, err := service.RunTest(context.Background(), response, "api/test.http")
	require.NoError(t, err)
	assert.True(t, result.Passed)
	assert.True(t, result.Updated)

	// Get stats
	stats := service.GetStats()
	assert.Equal(t, 1, stats.Total)
	assert.Equal(t, 1, stats.Passed)
	assert.Equal(t, 0, stats.Failed)
	assert.Equal(t, 0, stats.Updated)
	assert.Equal(t, 1, stats.Created)

	// Run the test again, should pass and not update
	result, err = service.RunTest(context.Background(), response, "api/test.http")
	require.NoError(t, err)
	assert.True(t, result.Passed)
	assert.False(t, result.Updated)

	// Get updated stats
	stats = service.GetStats()
	assert.Equal(t, 2, stats.Total)
	assert.Equal(t, 2, stats.Passed)
	assert.Equal(t, 0, stats.Failed)
	assert.Equal(t, 0, stats.Updated)
	assert.Equal(t, 1, stats.Created)

	// Create a different response
	differentResponse := &models.HTTPResponse{
		StatusCode:    400,
		Status:        "400 Bad Request",
		Headers:       map[string][]string{"Content-Type": {"application/json"}},
		Body:          []byte(`{"error":"invalid request"}`),
		ContentType:   "application/json",
		ContentLength: 29,
		Request: &models.HTTPRequest{
			Method: "GET",
			Path:   "/api/test",
		},
		Timestamp: time.Now(),
	}

	// Run a test with different response, should fail but update
	result, err = service.RunTest(context.Background(), differentResponse, "api/test.http")
	require.NoError(t, err)
	assert.True(t, result.Passed) // Passes because it updated
	assert.True(t, result.Updated)

	// Get updated stats
	stats = service.GetStats()
	assert.Equal(t, 3, stats.Total)
	assert.Equal(t, 3, stats.Passed)
	assert.Equal(t, 0, stats.Failed)
	assert.Equal(t, 1, stats.Updated)
	assert.Equal(t, 1, stats.Created)

	// Reset stats
	service.ResetStats()
	stats = service.GetStats()
	assert.Equal(t, 0, stats.Total)
	assert.Equal(t, 0, stats.Passed)
	assert.Equal(t, 0, stats.Failed)
	assert.Equal(t, 0, stats.Updated)
	assert.Equal(t, 0, stats.Created)
}

func TestService_CleanupUnusedSnapshots(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-service-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager and service
	manager := NewManager(tempDir)
	service := NewService(manager, models.SnapshotOptions{
		BasePath:   tempDir,
		UpdateMode: "all",
	})

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

	// Save snapshots directly using manager
	err = manager.SaveSnapshot(context.Background(), response1, "api/test1.http")
	require.NoError(t, err)
	err = manager.SaveSnapshot(context.Background(), response2, "api/test2.http")
	require.NoError(t, err)

	// Run a test for only the first response, marking it as used
	_, err = service.RunTest(context.Background(), response1, "api/test1.http")
	require.NoError(t, err)

	// Cleanup unused snapshots
	err = service.CleanupUnusedSnapshots(context.Background(), "api")
	require.NoError(t, err)

	// List snapshots after cleanup
	snapshots, err := manager.ListSnapshots(context.Background(), "api")
	require.NoError(t, err)

	// Verify only the used snapshot remains
	assert.Len(t, snapshots, 1)
	assert.Contains(t, snapshots, "api/test1_get.snap.json")
	assert.NotContains(t, snapshots, "api/test2_get.snap.json")
}

func TestService_MultipleTests(t *testing.T) {
	// Create a temporary directory for snapshots
	tempDir, err := os.MkdirTemp("", "snapshot-service-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a manager and service
	manager := NewManager(tempDir)
	service := NewService(manager, models.SnapshotOptions{
		BasePath:   tempDir,
		UpdateMode: "all",
	})

	// Create multiple test responses
	responses := []*models.HTTPResponse{
		{
			StatusCode:  200,
			ContentType: "application/json",
			Body:        []byte(`{"id":1,"name":"item1"}`),
			Request: &models.HTTPRequest{
				Method: "GET",
				Path:   "/api/items/1",
			},
		},
		{
			StatusCode:  201,
			ContentType: "application/json",
			Body:        []byte(`{"id":2,"name":"item2"}`),
			Request: &models.HTTPRequest{
				Method: "POST",
				Path:   "/api/items",
			},
		},
		{
			StatusCode:  204,
			ContentType: "",
			Body:        []byte{},
			Request: &models.HTTPRequest{
				Method: "DELETE",
				Path:   "/api/items/3",
			},
		},
	}

	// Run tests for all responses
	for i, response := range responses {
		result, err := service.RunTest(context.Background(), response, "api/tests.http")
		require.NoError(t, err)
		assert.True(t, result.Passed)
		assert.True(t, result.Updated)

		// Get stats after each test
		stats := service.GetStats()
		assert.Equal(t, i+1, stats.Total)
		assert.Equal(t, i+1, stats.Passed)
		assert.Equal(t, 0, stats.Failed)
		assert.Equal(t, 0, stats.Updated)
		assert.Equal(t, i+1, stats.Created)
	}

	// List snapshots
	snapshots, err := manager.ListSnapshots(context.Background(), "api")
	require.NoError(t, err)

	// Verify we have three snapshots
	assert.Len(t, snapshots, 3)
}
