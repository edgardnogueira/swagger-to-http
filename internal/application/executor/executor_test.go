package executor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// Set up a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "GET" && r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle different paths
		switch r.URL.Path {
		case "/api/users":
			if r.Method == "GET" {
				// Return list of users
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`))
			} else if r.Method == "POST" {
				// Read request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read body", http.StatusBadRequest)
					return
				}
				
				// Return created user with ID
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"id":3,"name":"New User","data":` + string(body) + `}`))
			}
		case "/api/auth/token":
			// Simple OAuth token endpoint
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"bearer","expires_in":3600}`))
		case "/api/retry":
			// Test retry mechanism with a counter
			sessionCookie, err := r.Cookie("retry-counter")
			count := 0
			if err == nil {
				// Parse counter from cookie
				count = 1
				if sessionCookie.Value == "1" {
					count = 2
				}
			}
			
			if count < 2 {
				// Set cookie and return error for retry
				http.SetCookie(w, &http.Cookie{
					Name:  "retry-counter",
					Value: string(rune('0' + count)),
					Path:  "/",
				})
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			
			// Success after retry
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"retries":2}`))
		case "/api/timeout":
			// Simulate a timeout
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create executor service
	executor := NewService(
		WithClient(&http.Client{
			Timeout: 50 * time.Millisecond, // Short timeout for testing
		}),
	)

	t.Run("GET request", func(t *testing.T) {
		// Create request
		request := &models.HTTPRequest{
			Method: "GET",
			URL:    server.URL + "/api/users",
			Headers: []models.HTTPHeader{
				{Name: "Accept", Value: "application/json"},
			},
		}

		// Execute request
		resp, err := executor.Execute(context.Background(), request, nil)
		
		// Assert success
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.ContentType)
		assert.Contains(t, string(resp.Body), `"name":"John"`)
	})

	t.Run("POST request", func(t *testing.T) {
		// Create request with body
		request := &models.HTTPRequest{
			Method: "POST",
			URL:    server.URL + "/api/users",
			Headers: []models.HTTPHeader{
				{Name: "Content-Type", Value: "application/json"},
				{Name: "Accept", Value: "application/json"},
			},
			Body: `{"name":"New User","email":"new@example.com"}`,
		}

		// Execute request
		resp, err := executor.Execute(context.Background(), request, nil)
		
		// Assert success
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "application/json", resp.ContentType)
		assert.Contains(t, string(resp.Body), `"id":3`)
		assert.Contains(t, string(resp.Body), `"email":"new@example.com"`)
	})

	t.Run("With variables", func(t *testing.T) {
		// Create request with variables
		request := &models.HTTPRequest{
			Method: "GET",
			URL:    server.URL + "/api/users?filter={{filter}}",
			Headers: []models.HTTPHeader{
				{Name: "Accept", Value: "application/json"},
				{Name: "X-Custom-{{custom}}", Value: "{{value}}"},
			},
		}

		// Add variables
		variables := map[string]string{
			"filter": "active",
			"custom": "Header",
			"value":  "test-value",
		}

		// Create custom executor to inspect the request
		customServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check variable substitution in URL
			assert.Equal(t, "/api/users?filter=active", r.URL.Path+"?"+r.URL.RawQuery)
			
			// Check variable substitution in headers
			assert.Equal(t, "test-value", r.Header.Get("X-Custom-Header"))
			
			// Return success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true}`))
		}))
		defer customServer.Close()

		// Update URL to use custom server
		request.URL = customServer.URL + "/api/users?filter={{filter}}"

		// Execute request
		resp, err := executor.Execute(context.Background(), request, variables)
		
		// Assert success
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Timeout", func(t *testing.T) {
		// Create request to timeout endpoint
		request := &models.HTTPRequest{
			Method: "GET",
			URL:    server.URL + "/api/timeout",
		}

		// Execute request
		resp, err := executor.Execute(context.Background(), request, nil)
		
		// Assert timeout error
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "timeout")
	})
}

func TestHTTPFileParser(t *testing.T) {
	// Sample HTTP file content
	content := `### Get users
GET https://api.example.com/users
Accept: application/json

### Create user
POST https://api.example.com/users
Content-Type: application/json
Accept: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

###`

	// Create parser
	parser := NewHTTPFileParser(newNullLogger())

	// Parse content
	file, err := parser.ParseString(content, "test.http")

	// Assert success
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "test.http", file.Filename)
	assert.Len(t, file.Requests, 2)

	// Check first request
	assert.Equal(t, "GET", file.Requests[0].Method)
	assert.Equal(t, "https://api.example.com/users", file.Requests[0].URL)
	assert.Equal(t, "/users", file.Requests[0].Path)
	assert.Len(t, file.Requests[0].Headers, 1)
	assert.Equal(t, "Accept", file.Requests[0].Headers[0].Name)
	assert.Equal(t, "application/json", file.Requests[0].Headers[0].Value)
	assert.Empty(t, file.Requests[0].Body)

	// Check second request
	assert.Equal(t, "POST", file.Requests[1].Method)
	assert.Equal(t, "https://api.example.com/users", file.Requests[1].URL)
	assert.Equal(t, "/users", file.Requests[1].Path)
	assert.Len(t, file.Requests[1].Headers, 2)
	assert.Equal(t, "Content-Type", file.Requests[1].Headers[0].Name)
	assert.Equal(t, "application/json", file.Requests[1].Headers[0].Value)
	assert.Contains(t, file.Requests[1].Body, "John Doe")
}

func TestExecuteFile(t *testing.T) {
	// Set up a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"path":"` + r.URL.Path + `"}`))
	}))
	defer server.Close()

	// Create a test HTTP file
	file := &models.HTTPFile{
		Filename: "test.http",
		Requests: []models.HTTPRequest{
			{
				Method: "GET",
				URL:    server.URL + "/api/users",
				Headers: []models.HTTPHeader{
					{Name: "Accept", Value: "application/json"},
				},
			},
			{
				Method: "GET",
				URL:    server.URL + "/api/products",
				Headers: []models.HTTPHeader{
					{Name: "Accept", Value: "application/json"},
				},
			},
		},
	}

	// Create executor service
	executor := NewService()

	// Execute file
	responses, err := executor.ExecuteFile(context.Background(), file, nil)

	// Assert success
	assert.NoError(t, err)
	assert.NotNil(t, responses)
	assert.Len(t, responses, 2)

	// Check first response
	assert.Equal(t, http.StatusOK, responses[0].StatusCode)
	assert.Contains(t, string(responses[0].Body), `"path":"/api/users"`)

	// Check second response
	assert.Equal(t, http.StatusOK, responses[1].StatusCode)
	assert.Contains(t, string(responses[1].Body), `"path":"/api/products"`)
}

func TestVariableStore(t *testing.T) {
	// Create variable store
	store := newMemoryVariableStore()

	// Set and get variables
	store.Set("key1", "value1")
	store.Set("key2", "value2")

	// Test get
	val, exists := store.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)

	// Test get all
	vars := store.GetAll()
	assert.Len(t, vars, 2)
	assert.Equal(t, "value1", vars["key1"])
	assert.Equal(t, "value2", vars["key2"])

	// Test delete
	store.Delete("key1")
	_, exists = store.Get("key1")
	assert.False(t, exists)

	// Test clear
	store.Clear()
	vars = store.GetAll()
	assert.Len(t, vars, 0)

	// Test extract variables
	vars = store.ExtractVariables("Hello {{name}}, welcome to {{service}}")
	assert.Len(t, vars, 2)
	assert.Equal(t, "name", vars[0])
	assert.Equal(t, "service", vars[1])

	// Test has variable
	assert.True(t, store.HasVariable("Hello {{name}}"))
	assert.False(t, store.HasVariable("Hello name"))
}

func TestSessionStore(t *testing.T) {
	// Create session store
	store := newMemorySessionStore()

	// Create test cookies
	cookie1 := &http.Cookie{Name: "session", Value: "abc123", Domain: "example.com"}
	cookie2 := &http.Cookie{Name: "user", Value: "john", Domain: "example.com"}
	cookie3 := &http.Cookie{Name: "session", Value: "xyz789", Domain: "api.example.com"}

	// Set cookies
	store.SetCookie("example.com", cookie1)
	store.SetCookie("example.com", cookie2)
	store.SetCookie("api.example.com", cookie3)

	// Get cookies for host
	cookies := store.GetCookies("example.com")
	assert.Len(t, cookies, 2)

	// Check updating existing cookie
	cookie1Updated := &http.Cookie{Name: "session", Value: "updated", Domain: "example.com"}
	store.SetCookie("example.com", cookie1Updated)
	cookies = store.GetCookies("example.com")
	assert.Len(t, cookies, 2)
	assert.Equal(t, "updated", cookies[0].Value)

	// Get hosts
	hosts := store.GetHosts()
	assert.Len(t, hosts, 2)
	assert.Contains(t, hosts, "example.com")
	assert.Contains(t, hosts, "api.example.com")

	// Clear cookies for host
	store.ClearCookies("example.com")
	cookies = store.GetCookies("example.com")
	assert.Len(t, cookies, 0)
	cookies = store.GetCookies("api.example.com")
	assert.Len(t, cookies, 1)

	// Clear all cookies
	store.ClearAllCookies()
	hosts = store.GetHosts()
	assert.Len(t, hosts, 0)
}
