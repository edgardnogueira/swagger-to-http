# API Reference

This document provides a comprehensive reference for the core APIs and interfaces in the `swagger-to-http` tool. It is intended for developers who want to extend or integrate with the tool.

## Domain Models

### Swagger Models

The core models for representing Swagger/OpenAPI documents are defined in `internal/domain/models/swagger.go`.

#### SwaggerDoc

Represents a complete Swagger/OpenAPI document.

```go
type SwaggerDoc struct {
	Version     string                 `json:"swagger" yaml:"swagger"`
	OpenAPI     string                 `json:"openapi" yaml:"openapi"`
	Info        *Info                  `json:"info" yaml:"info"`
	Host        string                 `json:"host" yaml:"host"`
	BasePath    string                 `json:"basePath" yaml:"basePath"`
	Schemes     []string               `json:"schemes" yaml:"schemes"`
	Consumes    []string               `json:"consumes" yaml:"consumes"`
	Produces    []string               `json:"produces" yaml:"produces"`
	Paths       map[string]*PathItem   `json:"paths" yaml:"paths"`
	Definitions map[string]*Schema     `json:"definitions" yaml:"definitions"`
	Components  *Components            `json:"components" yaml:"components"`
	// Other fields...
}
```

#### PathItem

Represents an API path with its operations.

```go
type PathItem struct {
	Ref        string       `json:"$ref" yaml:"$ref"`
	Get        *Operation   `json:"get" yaml:"get"`
	Put        *Operation   `json:"put" yaml:"put"`
	Post       *Operation   `json:"post" yaml:"post"`
	Delete     *Operation   `json:"delete" yaml:"delete"`
	Options    *Operation   `json:"options" yaml:"options"`
	Head       *Operation   `json:"head" yaml:"head"`
	Patch      *Operation   `json:"patch" yaml:"patch"`
	Parameters []Parameter  `json:"parameters" yaml:"parameters"`
}
```

#### Operation

Represents an API operation (endpoint).

```go
type Operation struct {
	Tags        []string               `json:"tags" yaml:"tags"`
	Summary     string                 `json:"summary" yaml:"summary"`
	Description string                 `json:"description" yaml:"description"`
	OperationID string                 `json:"operationId" yaml:"operationId"`
	Consumes    []string               `json:"consumes" yaml:"consumes"`
	Produces    []string               `json:"produces" yaml:"produces"`
	Parameters  []Parameter            `json:"parameters" yaml:"parameters"`
	Responses   map[string]*Response   `json:"responses" yaml:"responses"`
	// Other fields...
}
```

### HTTP Models

The models for HTTP requests and responses are defined in `internal/domain/models/http.go`.

#### HTTPRequest

Represents an HTTP request.

```go
type HTTPRequest struct {
	Method      string              `json:"method"`
	Path        string              `json:"path"`
	Headers     map[string][]string `json:"headers"`
	Body        []byte              `json:"body"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
}
```

#### HTTPResponse

Represents an HTTP response.

```go
type HTTPResponse struct {
	StatusCode    int                 `json:"statusCode"`
	Status        string              `json:"status"`
	Headers       map[string][]string `json:"headers"`
	Body          []byte              `json:"body"`
	ContentType   string              `json:"contentType"`
	ContentLength int64               `json:"contentLength"`
	Duration      time.Duration       `json:"duration"`
	Request       *HTTPRequest        `json:"request"`
	RequestID     string              `json:"requestId"`
	Timestamp     time.Time           `json:"timestamp"`
}
```

#### HTTPFile

Represents a file containing HTTP requests.

```go
type HTTPFile struct {
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Requests []*HTTPRequest `json:"requests"`
}
```

#### HTTPCollection

Represents a collection of HTTP files.

```go
type HTTPCollection struct {
	Files map[string]*HTTPFile `json:"files"`
}
```

### Snapshot Models

The models for snapshot testing are defined in `internal/domain/models/snapshot.go`.

#### SnapshotData

Represents a stored snapshot.

```go
type SnapshotData struct {
	Metadata SnapshotMetadata `json:"metadata"`
	Content  string           `json:"content"`
}
```

#### SnapshotDiff

Represents the difference between a response and a snapshot.

```go
type SnapshotDiff struct {
	RequestPath   string       `json:"requestPath"`
	RequestMethod string       `json:"requestMethod"`
	StatusDiff    *StatusDiff  `json:"statusDiff"`
	HeaderDiff    *HeaderDiff  `json:"headerDiff"`
	BodyDiff      *BodyDiff    `json:"bodyDiff"`
	Equal         bool         `json:"equal"`
}
```

## Application Interfaces

### SwaggerParser

Interface for parsing Swagger/OpenAPI documents, defined in `internal/application/interfaces.go`.

```go
type SwaggerParser interface {
	Parse(ctx context.Context, data []byte) (*models.SwaggerDoc, error)
	ParseFile(ctx context.Context, filePath string) (*models.SwaggerDoc, error)
	ParseURL(ctx context.Context, url string) (*models.SwaggerDoc, error)
	Validate(ctx context.Context, doc *models.SwaggerDoc) error
}
```

### HTTPGenerator

Interface for generating HTTP requests from Swagger documents, defined in `internal/application/interfaces.go`.

```go
type HTTPGenerator interface {
	Generate(ctx context.Context, doc *models.SwaggerDoc) (*models.HTTPCollection, error)
	GenerateRequest(ctx context.Context, path string, pathItem *models.PathItem, method string, operation *models.Operation) (*models.HTTPRequest, error)
}
```

### HTTPExecutor

Interface for executing HTTP requests, defined in `internal/application/interfaces.go`.

```go
type HTTPExecutor interface {
	Execute(ctx context.Context, request *models.HTTPRequest, variables map[string]string) (*models.HTTPResponse, error)
	ExecuteFile(ctx context.Context, file *models.HTTPFile, variables map[string]string) ([]*models.HTTPResponse, error)
}
```

### SnapshotManager

Interface for managing response snapshots, defined in `internal/application/interfaces.go`.

```go
type SnapshotManager interface {
	SaveSnapshot(ctx context.Context, response *models.HTTPResponse, path string) error
	LoadSnapshot(ctx context.Context, path string) (*models.HTTPResponse, error)
	CompareWithSnapshot(ctx context.Context, response *models.HTTPResponse, snapshotPath string) (*models.SnapshotDiff, error)
}
```

## Command-Line Implementation

The command-line interface is implemented in the `internal/cli` package.

### Root Command

The root command sets up the application's CLI structure:

```go
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "swagger-to-http",
	Short: "A tool to convert Swagger/OpenAPI documents to HTTP files",
	Long: `swagger-to-http converts Swagger/OpenAPI documentation into organized 
HTTP request files with snapshot testing capabilities.`,
}
```

### Generate Command

The generate command handles the conversion of Swagger documents to HTTP files:

```go
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate HTTP files from a Swagger/OpenAPI document",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation...
	},
}
```

### Snapshot Commands

The snapshot commands enable snapshot testing functionality:

```go
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshot testing commands",
	Long:  "Commands for working with HTTP response snapshots",
}

var testCmd = &cobra.Command{
	Use:   "test [file-pattern]",
	Short: "Run snapshot tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation...
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [file-pattern]",
	Short: "Update snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implementation...
	},
}
```

## Extension Points

### Adding New Formatters

To add a new response formatter, implement the `ResponseFormatter` interface:

```go
type ResponseFormatter interface {
	Format(body []byte) ([]byte, error)
	Compare(expected, actual []byte) *models.BodyDiff
}
```

Then register it with the `SnapshotManager`:

```go
manager.RegisterFormatter("application/custom-type", &CustomFormatter{})
```

### Custom HTTP Execution

To implement custom HTTP request execution behavior, create a new `HTTPExecutor`:

```go
type CustomExecutor struct {
	// Custom fields...
}

func (e *CustomExecutor) Execute(ctx context.Context, request *models.HTTPRequest, variables map[string]string) (*models.HTTPResponse, error) {
	// Custom implementation...
}

func (e *CustomExecutor) ExecuteFile(ctx context.Context, file *models.HTTPFile, variables map[string]string) ([]*models.HTTPResponse, error) {
	// Custom implementation...
}
```

### Custom Swagger Parsing

To support additional Swagger/OpenAPI formats or versions, implement a new `SwaggerParser`:

```go
type CustomParser struct {
	// Custom fields...
}

func (p *CustomParser) Parse(ctx context.Context, data []byte) (*models.SwaggerDoc, error) {
	// Custom implementation...
}

// Implement other required methods...
```

## Configuration

The configuration system is implemented using Viper:

```go
// ConfigProvider defines the interface for retrieving configuration
type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringSlice(key string) []string
}
```

## Utility Functions

Several utility functions are available for working with the API:

### HTTP Utilities

```go
// MergeHeaders merges two sets of headers
func MergeHeaders(base, override map[string][]string) map[string][]string

// FormatHTTPRequest formats an HTTP request as a string
func FormatHTTPRequest(request *models.HTTPRequest) string

// ParseHTTPFile parses an HTTP file into individual requests
func ParseHTTPFile(content []byte) ([]*models.HTTPRequest, error)
```

### Swagger Utilities

```go
// GetOperationByMethod returns the Operation for a specific method
func GetOperationByMethod(pathItem *models.PathItem, method string) *models.Operation

// GetParameterByName finds a parameter by its name
func GetParameterByName(parameters []models.Parameter, name string) *models.Parameter

// GetSchemaByRef resolves a schema reference
func GetSchemaByRef(doc *models.SwaggerDoc, ref string) *models.Schema
```

## Error Handling

The application uses structured errors with context:

```go
// WrapError wraps an error with additional context
func WrapError(err error, message string) error

// NewError creates a new error with a message
func NewError(message string) error

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool
```

## Integration Examples

### Using the API in Your Code

```go
// Create a parser
parser := parser.NewSwaggerParser()

// Parse a Swagger document
doc, err := parser.ParseFile(context.Background(), "swagger.json")
if err != nil {
	return err
}

// Create a generator
generator := generator.NewHTTPGenerator(&config.Config{
	IndentJSON: true,
	DefaultTag: "default",
})

// Generate HTTP files
collection, err := generator.Generate(context.Background(), doc)
if err != nil {
	return err
}

// Create a file writer
writer := fs.NewFileWriter("http-requests")

// Write the HTTP files
err = writer.WriteCollection(context.Background(), collection)
if err != nil {
	return err
}
```

### Snapshot Testing Integration

```go
// Create an HTTP executor
executor := executor.NewHTTPExecutor()

// Create a snapshot manager
manager := snapshot.NewManager(".snapshots")

// Execute a request
response, err := executor.Execute(context.Background(), request, nil)
if err != nil {
	return err
}

// Compare with snapshot
diff, err := manager.CompareWithSnapshot(context.Background(), response, "api/test.snap.json")
if err != nil {
	return err
}

// Check if equal
if !diff.Equal {
	// Handle differences...
}
```

## Further Reading

- See the [Contributing Guide](contributing.md) for information on extending the tool
- Browse the source code on [GitHub](https://github.com/edgardnogueira/swagger-to-http)
- Check the [GoDoc documentation](https://godoc.org/github.com/edgardnogueira/swagger-to-http)
