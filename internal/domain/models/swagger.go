package models

// SwaggerDoc represents a Swagger/OpenAPI document
type SwaggerDoc struct {
	Version     string                 `json:"openapi,omitempty" yaml:"openapi,omitempty"`
	SwaggerVersion string              `json:"swagger,omitempty" yaml:"swagger,omitempty"`
	Info        Info                   `json:"info" yaml:"info"`
	Host        string                 `json:"host,omitempty" yaml:"host,omitempty"`
	BasePath    string                 `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	Schemes     []string               `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Consumes    []string               `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces    []string               `json:"produces,omitempty" yaml:"produces,omitempty"`
	Paths       map[string]PathItem    `json:"paths" yaml:"paths"`
	Components  *Components            `json:"components,omitempty" yaml:"components,omitempty"`
	Definitions map[string]interface{} `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	Parameters  map[string]Parameter   `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Tags        []Tag                  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Servers     []Server               `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// Info represents the metadata of a Swagger/OpenAPI document
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

// Contact represents contact information
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License represents license information
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// PathItem represents a path in the Swagger/OpenAPI paths section
type PathItem struct {
	Ref        string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Summary    string     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Get        *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Put        *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Post       *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Delete     *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options    *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Head       *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Patch      *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace      *Operation `json:"trace,omitempty" yaml:"trace,omitempty"`
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Operation represents an operation in a Swagger/OpenAPI path
type Operation struct {
	Tags        []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Consumes    []string               `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces    []string               `json:"produces,omitempty" yaml:"produces,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]Response    `json:"responses" yaml:"responses"`
	Security    []map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// Parameter represents a parameter in a Swagger/OpenAPI operation
type Parameter struct {
	Name            string      `json:"name" yaml:"name"`
	In              string      `json:"in" yaml:"in"`
	Description     string      `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Schema          *Schema     `json:"schema,omitempty" yaml:"schema,omitempty"`
	Type            string      `json:"type,omitempty" yaml:"type,omitempty"`
	Format          string      `json:"format,omitempty" yaml:"format,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Items           *Items      `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat string     `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Default         interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Maximum         *float64    `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum bool       `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum         *float64    `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum bool       `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength       *int64      `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength       *int64      `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern         string      `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems        *int64      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems        *int64      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems     bool        `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Enum            []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	MultipleOf      *float64    `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Content         map[string]MediaType `json:"content,omitempty" yaml:"content,omitempty"`
}

// RequestBody represents a request body in OpenAPI 3.0
type RequestBody struct {
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]MediaType   `json:"content" yaml:"content"`
	Required    bool                   `json:"required,omitempty" yaml:"required,omitempty"`
}

// Response represents a response in a Swagger/OpenAPI operation
type Response struct {
	Description string                 `json:"description" yaml:"description"`
	Schema      *Schema                `json:"schema,omitempty" yaml:"schema,omitempty"`
	Headers     map[string]Header      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`
	Content     map[string]MediaType   `json:"content,omitempty" yaml:"content,omitempty"`
}

// Schema represents a schema in Swagger/OpenAPI
type Schema struct {
	Ref                  string                 `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type                 string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Format               string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Title                string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Description          string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Default              interface{}            `json:"default,omitempty" yaml:"default,omitempty"`
	Items                *Items                 `json:"items,omitempty" yaml:"items,omitempty"`
	Maximum              *float64               `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              *float64               `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            *int64                 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int64                 `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             *int64                 `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             *int64                 `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          bool                   `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        *int64                 `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int64                 `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string               `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty" yaml:"enum,omitempty"`
	AdditionalProperties *Schema                `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Properties           map[string]*Schema     `json:"properties,omitempty" yaml:"properties,omitempty"`
	AllOf                []*Schema              `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []*Schema              `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []*Schema              `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  *Schema                `json:"not,omitempty" yaml:"not,omitempty"`
	AdditionalItems      *Schema                `json:"additionalItems,omitempty" yaml:"additionalItems,omitempty"`
	Example              interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
}

// Items represents items in a Schema
type Items struct {
	Ref                  string                 `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Type                 string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Format               string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Items                *Items                 `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat     string                 `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Default              interface{}            `json:"default,omitempty" yaml:"default,omitempty"`
	Maximum              *float64               `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                   `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              *float64               `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                   `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            *int64                 `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            *int64                 `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             *int64                 `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             *int64                 `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          bool                   `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Enum                 []interface{}          `json:"enum,omitempty" yaml:"enum,omitempty"`
}

// MediaType represents a media type in OpenAPI 3.0
type MediaType struct {
	Schema   *Schema                `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]Encoding    `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Encoding represents encoding information for a media type
type Encoding struct {
	ContentType   string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string            `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool              `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool              `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Header represents a header in a Swagger/OpenAPI response
type Header struct {
	Description     string      `json:"description,omitempty" yaml:"description,omitempty"`
	Type            string      `json:"type,omitempty" yaml:"type,omitempty"`
	Format          string      `json:"format,omitempty" yaml:"format,omitempty"`
	Items           *Items      `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat string     `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Default         interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Schema          *Schema     `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// Tag represents a tag in a Swagger/OpenAPI document
type Tag struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// Components represents components in an OpenAPI 3.0 document
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]interface{}    `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]Header         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]Link           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]Callback       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

// SecurityScheme represents a security scheme in OpenAPI 3.0
type SecurityScheme struct {
	Type             string           `json:"type" yaml:"type"`
	Description      string           `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string           `json:"name,omitempty" yaml:"name,omitempty"`
	In               string           `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string           `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string           `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows      `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string           `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows represents OAuth flows in OpenAPI 3.0
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlow represents an OAuth flow in OpenAPI 3.0
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}

// Link represents a link in OpenAPI 3.0
type Link struct {
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationID  string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *Server                `json:"server,omitempty" yaml:"server,omitempty"`
}

// Callback represents a callback in OpenAPI 3.0
type Callback map[string]PathItem

// Server represents a server in OpenAPI 3.0
type Server struct {
	URL         string                    `json:"url" yaml:"url"`
	Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable represents a server variable in OpenAPI 3.0
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}
