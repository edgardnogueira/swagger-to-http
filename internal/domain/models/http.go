package models

// HTTPRequest represents an HTTP request in .http file format
type HTTPRequest struct {
	Name     string
	Method   string
	URL      string
	Headers  []Header
	Body     string
	Comments []string
	Tag      string
	Path     string
}

// HTTPHeader represents an HTTP header
type HTTPHeader struct {
	Name  string
	Value string
}

// HTTPFile represents a collection of HTTP requests to be written to a .http file
type HTTPFile struct {
	Filename string
	Requests []HTTPRequest
}

// HTTPDirectory represents a directory containing HTTP files
type HTTPDirectory struct {
	Name  string
	Path  string
	Files []HTTPFile
}

// HTTPCollection represents a collection of directories and files
type HTTPCollection struct {
	RootDir      string
	Directories  []HTTPDirectory
	RootFiles    []HTTPFile
}
