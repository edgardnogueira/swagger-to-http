package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestExecutor_Execute(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		assert.Equal(t, "GET", r.Method)

		// Check headers
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "TestValue", r.Header.Get("X-Custom-Header"))

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	// Create executor
	executor := NewExecutor(10*time.Second, map[string]string{
		"BASE_URL": server.URL,
	})

	// Create request
	request := &models.HTTPRequest{
		Method: "GET",
		URL:    "{{BASE_URL}}/api/test",
		Headers: []models.HTTPHeader{
			{Name: "Accept", Value: "application/json"},
			{Name: "X-Custom-Header", Value: "{{HEADER_VALUE}}"},
		},
		Path: "/api/test",
	}

	// Execute request
	response, err := executor.Execute(context.Background(), request, map[string]string{
		"HEADER_VALUE": "TestValue",
	})

	// Assert no error
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Check response
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.ContentType)
	assert.Equal(t, `{"status":"success"}`, string(response.Body))
	assert.Equal(t, request, response.Request)
}

func TestExecutor_ExecuteWithBody(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		assert.Equal(t, "POST", r.Method)

		// Check content type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Parse body
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		assert.Equal(t, `{"name":"TestValue"}`, string(buf))

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":123,"name":"TestValue"}`))
	}))
	defer server.Close()

	// Create executor
	executor := NewExecutor(10*time.Second, map[string]string{
		"BASE_URL": server.URL,
	})

	// Create request
	request := &models.HTTPRequest{
		Method: "POST",
		URL:    "{{BASE_URL}}/api/create",
		Headers: []models.HTTPHeader{
			{Name: "Content-Type", Value: "application/json"},
		},
		Body: `{"name":"{{NAME_VALUE}}"}`,
		Path: "/api/create",
	}

	// Execute request
	response, err := executor.Execute(context.Background(), request, map[string]string{
		"NAME_VALUE": "TestValue",
	})

	// Assert no error
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Check response
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "application/json", response.ContentType)
	assert.Equal(t, `{"id":123,"name":"TestValue"}`, string(response.Body))
}

func TestExecutor_ExecuteFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response based on path
		w.Header().Set("Content-Type", "application/json")
		
		if r.URL.Path == "/api/test1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":1}`))
		} else if r.URL.Path == "/api/test2" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":2}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
		}
	}))
	defer server.Close()

	// Create executor
	executor := NewExecutor(10*time.Second, map[string]string{
		"BASE_URL": server.URL,
	})

	// Create HTTP file
	file := &models.HTTPFile{
		Filename: "test.http",
		Requests: []models.HTTPRequest{
			{
				Method: "GET",
				URL:    "{{BASE_URL}}/api/test1",
				Headers: []models.HTTPHeader{
					{Name: "Accept", Value: "application/json"},
				},
				Path: "/api/test1",
			},
			{
				Method: "GET",
				URL:    "{{BASE_URL}}/api/test2",
				Headers: []models.HTTPHeader{
					{Name: "Accept", Value: "application/json"},
				},
				Path: "/api/test2",
			},
		},
	}

	// Execute file
	responses, err := executor.ExecuteFile(context.Background(), file, nil)

	// Assert no error
	assert.NoError(t, err)
	assert.Len(t, responses, 2)

	// Check responses
	assert.Equal(t, http.StatusOK, responses[0].StatusCode)
	assert.Equal(t, `{"id":1}`, string(responses[0].Body))
	
	assert.Equal(t, http.StatusOK, responses[1].StatusCode)
	assert.Equal(t, `{"id":2}`, string(responses[1].Body))
}

func TestExecutor_ProcessVariables(t *testing.T) {
	executor := &Executor{}
	
	variables := map[string]string{
		"HOST":     "example.com",
		"API_PATH": "/api/v1",
		"TOKEN":    "abc123",
	}
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No variables",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "Single variable",
			input:    "Host: {{HOST}}",
			expected: "Host: example.com",
		},
		{
			name:     "Multiple variables",
			input:    "{{HOST}}{{API_PATH}}?token={{TOKEN}}",
			expected: "example.com/api/v1?token=abc123",
		},
		{
			name:     "Unknown variable",
			input:    "{{UNKNOWN}}",
			expected: "{{UNKNOWN}}",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.processVariables(tt.input, variables)
			assert.Equal(t, tt.expected, result)
		})
	}
}
