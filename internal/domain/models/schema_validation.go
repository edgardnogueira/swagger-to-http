package models

// SchemaValidationResult represents the result of validating a response against a schema
type SchemaValidationResult struct {
	Valid          bool              `json:"valid"`
	Errors         []ValidationError `json:"errors,omitempty"`
	SchemaPath     string            `json:"schemaPath,omitempty"`
	ResponseStatus int               `json:"responseStatus"`
	ContentType    string            `json:"contentType"`
}

// ValidationError represents a schema validation error
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
	Schema  string `json:"schema,omitempty"`
}

// ValidationOptions contains options for schema validation
type ValidationOptions struct {
	IgnoreAdditionalProperties bool     `json:"ignoreAdditionalProperties"`
	IgnoreFormats              bool     `json:"ignoreFormats"`
	IgnorePatterns             bool     `json:"ignorePatterns"`
	RequiredPropertiesOnly     bool     `json:"requiredPropertiesOnly"`
	IgnoreNullable             bool     `json:"ignoreNullable"`
	IgnoredProperties          []string `json:"ignoredProperties,omitempty"`
}
