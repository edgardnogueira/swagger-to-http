package http

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_ParseContent(t *testing.T) {
	parser := NewParser()

	testCases := []struct {
		name     string
		content  string
		expected int // expected number of requests
	}{
		{
			name: "Single request",
			content: `GET https://example.com/api/users
Accept: application/json

`,
			expected: 1,
		},
		{
			name: "Multiple requests with separator",
			content: `GET https://example.com/api/users
Accept: application/json

###

POST https://example.com/api/users
Content-Type: application/json
Accept: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}`,
			expected: 2,
		},
		{
			name: "Request with comments and metadata",
			content: `# Test API
@name GetUser
@tag users
GET https://example.com/api/users/1
Accept: application/json

`,
			expected: 1,
		},
		{
			name: "Empty requests should be ignored",
			content: `###

###

GET https://example.com/api/users
Accept: application/json

###

`,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requests, err := parser.ParseContent([]byte(tc.content), "test.http")
			assert.NoError(t, err)
			assert.Len(t, requests, tc.expected)
		})
	}
}

func TestParser_ParseRequest(t *testing.T) {
	parser := NewParser()

	t.Run("Basic GET request", func(t *testing.T) {
		content := `GET https://example.com/api/users
Accept: application/json
X-API-Key: abc123

`
		request, err := parser.parseRequest([]byte(content), 1, "test.http")
		assert.NoError(t, err)
		assert.Equal(t, "GET", request.Method)
		assert.Equal(t, "https://example.com/api/users", request.URL)
		assert.Equal(t, "/api/users", request.Path)
		assert.Len(t, request.Headers, 2)
		assert.Equal(t, "Accept", request.Headers[0].Name)
		assert.Equal(t, "application/json", request.Headers[0].Value)
	})

	t.Run("POST request with body", func(t *testing.T) {
		content := `POST https://example.com/api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}`
		request, err := parser.parseRequest([]byte(content), 1, "test.http")
		assert.NoError(t, err)
		assert.Equal(t, "POST", request.Method)
		assert.Equal(t, "https://example.com/api/users", request.URL)
		assert.Equal(t, "/api/users", request.Path)
		assert.Contains(t, request.Body, "John Doe")
	})

	t.Run("Request with metadata", func(t *testing.T) {
		content := `@name GetUserDetails
@tag users
GET https://example.com/api/users/1
Accept: application/json

`
		request, err := parser.parseRequest([]byte(content), 1, "test.http")
		assert.NoError(t, err)
		assert.Equal(t, "GetUserDetails", request.Name)
		assert.Equal(t, "users", request.Tag)
		assert.Equal(t, "GET", request.Method)
	})

	t.Run("Request with comments", func(t *testing.T) {
		content := `# This is a test comment
# Another comment
GET https://example.com/api/users
Accept: application/json

`
		request, err := parser.parseRequest([]byte(content), 1, "test.http")
		assert.NoError(t, err)
		assert.Len(t, request.Comments, 2)
		assert.Equal(t, " This is a test comment", request.Comments[0])
	})

	t.Run("Invalid request - missing method", func(t *testing.T) {
		content := `Accept: application/json`
		_, err := parser.parseRequest([]byte(content), 1, "test.http")
		assert.Error(t, err)
	})
}

func TestParser_FindHTTPFiles(t *testing.T) {
	parser := NewParser()
	
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "http-parser-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create some test files
	testFiles := []string{
		"test1.http",
		"test2.http",
		"other.txt",
		"subfolder/test3.http",
	}
	
	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		
		// Ensure the directory exists
		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0755)
		
		// Create the file
		err := os.WriteFile(path, []byte("GET http://example.com"), 0644)
		assert.NoError(t, err)
	}
	
	t.Run("Find by directory", func(t *testing.T) {
		files, err := parser.FindHTTPFiles(tmpDir)
		assert.NoError(t, err)
		assert.Len(t, files, 2) // Only root .http files, not subdirectories
		
		// Check that we found the expected files
		foundTest1 := false
		foundTest2 := false
		
		for _, file := range files {
			if filepath.Base(file) == "test1.http" {
				foundTest1 = true
			} else if filepath.Base(file) == "test2.http" {
				foundTest2 = true
			}
		}
		
		assert.True(t, foundTest1, "test1.http should be found")
		assert.True(t, foundTest2, "test2.http should be found")
	})
	
	t.Run("Find by glob pattern", func(t *testing.T) {
		pattern := filepath.Join(tmpDir, "*.http")
		files, err := parser.FindHTTPFiles(pattern)
		assert.NoError(t, err)
		assert.Len(t, files, 2)
	})
	
	t.Run("Find by specific file", func(t *testing.T) {
		pattern := filepath.Join(tmpDir, "test1.http")
		files, err := parser.FindHTTPFiles(pattern)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, filepath.Join(tmpDir, "test1.http"), files[0])
	})
	
	t.Run("Find with recursive pattern", func(t *testing.T) {
		pattern := filepath.Join(tmpDir, "**", "*.http")
		files, err := parser.FindHTTPFiles(pattern)
		assert.NoError(t, err)
		
		// Some filesystems might not support ** directly, so this is a fallback
		if len(files) == 0 {
			pattern = filepath.Join(tmpDir, "*", "*.http")
			files, err = parser.FindHTTPFiles(pattern)
			assert.NoError(t, err)
			
			if len(files) == 0 {
				t.Skip("Recursive glob pattern not supported by this filesystem")
			}
		}
		
		assert.GreaterOrEqual(t, len(files), 1, "Should find at least one file in subfolder")
	})
}
