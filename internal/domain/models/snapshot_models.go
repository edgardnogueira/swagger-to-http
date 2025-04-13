package models

// SnapshotDiff represents the difference between a response and a snapshot
type SnapshotDiff struct {
	// Simple version for test reports
	HasDiff    bool              `json:"hasDiff"`
	DiffString string            `json:"diffString"`
	HeaderDiff map[string][]string `json:"headerDiff"`
	BodyDiff   string            `json:"bodyDiff"`
	StatusDiff bool              `json:"statusDiff"`
	
	// Extended version for detailed analysis
	RequestPath   string       `json:"requestPath,omitempty"`
	RequestMethod string       `json:"requestMethod,omitempty"`
	StatusDiffExt  *StatusDiff  `json:"statusDiffExt,omitempty"`
	HeaderDiffExt  *HeaderDiff  `json:"headerDiffExt,omitempty"`
	BodyDiffExt    *BodyDiff    `json:"bodyDiffExt,omitempty"`
	Equal         bool         `json:"equal,omitempty"`
}

// StatusDiff represents the difference between two status codes
type StatusDiff struct {
	Expected int  `json:"expected"`
	Actual   int  `json:"actual"`
	Equal    bool `json:"equal"`
}

// HeaderDiff represents the difference between two sets of headers
type HeaderDiff struct {
	MissingHeaders  map[string][]string `json:"missingHeaders"`
	ExtraHeaders    map[string][]string `json:"extraHeaders"`
	DifferentValues map[string]HeaderValueDiff `json:"differentValues"`
	Equal           bool `json:"equal"`
}

// HeaderValueDiff represents the difference between two header values
type HeaderValueDiff struct {
	Expected []string `json:"expected"`
	Actual   []string `json:"actual"`
}

// BodyDiff represents the difference between two response bodies
type BodyDiff struct {
	ContentType     string `json:"contentType"`
	ExpectedSize    int    `json:"expectedSize"`
	ActualSize      int    `json:"actualSize"`
	ExpectedContent string `json:"expectedContent"`
	ActualContent   string `json:"actualContent"`
	DiffContent     string `json:"diffContent"`
	JsonDiff        *JsonDiff `json:"jsonDiff,omitempty"`
	Equal           bool   `json:"equal"`
}

// JsonDiff represents the difference between two JSON objects
type JsonDiff struct {
	MissingFields   []string `json:"missingFields"`
	ExtraFields     []string `json:"extraFields"`
	DifferentTypes  map[string]TypeDiff `json:"differentTypes"`
	DifferentValues map[string]ValueDiff `json:"differentValues"`
	Equal           bool `json:"equal"`
}

// TypeDiff represents a difference in types
type TypeDiff struct {
	ExpectedType string `json:"expectedType"`
	ActualType   string `json:"actualType"`
}

// ValueDiff represents a difference in values
type ValueDiff struct {
	Expected interface{} `json:"expected"`
	Actual   interface{} `json:"actual"`
}

// SnapshotResult represents the result of a snapshot comparison
type SnapshotResult struct {
	SnapshotPath  string       `json:"snapshotPath"`
	Exists        bool         `json:"exists"`
	Matches       bool         `json:"matches"`
	Diff          *SnapshotDiff `json:"diff,omitempty"`
	WasUpdated    bool         `json:"wasUpdated"`
	WasCreated    bool         `json:"wasCreated"`
	UpdateMode    string       `json:"updateMode"`
	Error         string       `json:"error,omitempty"`
}
