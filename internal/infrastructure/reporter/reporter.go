package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/edgardnogueira/swagger-to-http/internal/domain/models"
)

// TestReporterService implements the TestReporter interface
type TestReporterService struct {}

// NewTestReporterService creates a new TestReporterService
func NewTestReporterService() *TestReporterService {
	return &TestReporterService{}
}

// GenerateReport generates a report in the specified format
func (s *TestReporterService) GenerateReport(ctx context.Context, report *models.TestReport, options models.TestReportOptions) (io.Reader, error) {
	switch options.Format {
	case "json":
		return s.generateJSONReport(report, options)
	case "html":
		return s.generateHTMLReport(report, options)
	case "junit":
		return s.generateJUnitReport(report, options)
	case "console":
		return s.generateConsoleReport(report, options)
	default:
		return s.generateJSONReport(report, options)
	}
}

// SaveReport saves a report to the file system
func (s *TestReporterService) SaveReport(ctx context.Context, report *models.TestReport, options models.TestReportOptions) error {
	// Generate the report
	reader, err := s.GenerateReport(ctx, report, options)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Create the output directory if it doesn't exist
	dir := filepath.Dir(options.OutputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the output file
	file, err := os.Create(options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Copy the report to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write report to file: %w", err)
	}

	return nil
}

// PrintReport prints a report to the provided writer
func (s *TestReporterService) PrintReport(ctx context.Context, report *models.TestReport, options models.TestReportOptions, writer io.Writer) error {
	// Generate the report
	reader, err := s.GenerateReport(ctx, report, options)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Copy the report to the writer
	_, err = io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// Helper methods to generate different report formats

// generateJSONReport generates a JSON report
func (s *TestReporterService) generateJSONReport(report *models.TestReport, options models.TestReportOptions) (io.Reader, error) {
	// Create a copy of the report for modification if needed
	reportCopy := *report

	// Optionally strip request/response details to reduce size
	if !options.IncludeRequests || !options.IncludeResponses {
		for i := range reportCopy.Results {
			if !options.IncludeRequests {
				reportCopy.Results[i].Request = nil
			}
			if !options.IncludeResponses {
				reportCopy.Results[i].Response = nil
			}
		}
	}

	// Marshal the report to JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if options.Detailed {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(reportCopy); err != nil {
		return nil, fmt.Errorf("failed to encode JSON report: %w", err)
	}

	return &buf, nil
}

// generateHTMLReport generates an HTML report
func (s *TestReporterService) generateHTMLReport(report *models.TestReport, options models.TestReportOptions) (io.Reader, error) {
	// Define the HTML template
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Name}} - Test Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1, h2, h3 {
            color: #444;
        }
        .summary {
            background-color: #f5f5f5;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
        }
        .stat {
            background: white;
            padding: 10px 15px;
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            min-width: 100px;
        }
        .stat-label {
            font-size: 0.8em;
            color: #777;
        }
        .stat-value {
            font-size: 1.5em;
            font-weight: bold;
        }
        .passed .stat-value { color: #4CAF50; }
        .failed .stat-value { color: #F44336; }
        .skipped .stat-value { color: #FF9800; }
        .error .stat-value { color: #9C27B0; }
        .results {
            margin-top: 20px;
        }
        .result {
            background: white;
            margin-bottom: 10px;
            padding: 15px;
            border-radius: 5px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            border-left: 5px solid #ddd;
        }
        .result-passed { border-left-color: #4CAF50; }
        .result-failed { border-left-color: #F44336; }
        .result-skipped { border-left-color: #FF9800; }
        .result-error { border-left-color: #9C27B0; }
        .result-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .result-name {
            font-weight: bold;
            font-size: 1.1em;
        }
        .result-status {
            font-size: 0.9em;
            padding: 3px 8px;
            border-radius: 3px;
            text-transform: uppercase;
            font-weight: bold;
        }
        .status-passed { background: #E8F5E9; color: #2E7D32; }
        .status-failed { background: #FFEBEE; color: #C62828; }
        .status-skipped { background: #FFF8E1; color: #F57F17; }
        .status-error { background: #F3E5F5; color: #6A1B9A; }
        .result-details {
            margin-top: 10px;
            font-size: 0.9em;
        }
        .result-detail {
            margin-bottom: 5px;
        }
        .result-detail-label {
            font-weight: bold;
            display: inline-block;
            width: 120px;
        }
        .result-detail-value {
            display: inline-block;
        }
        .result-error {
            background: #FFEBEE;
            padding: 10px;
            border-radius: 5px;
            margin-top: 10px;
            font-family: monospace;
            white-space: pre-wrap;
        }
        .snapshot-diff {
            background: #F5F5F5;
            padding: 10px;
            border-radius: 5px;
            margin-top: 10px;
            font-family: monospace;
            white-space: pre-wrap;
            overflow-x: auto;
        }
        .request-response {
            margin-top: 15px;
            border-top: 1px solid #eee;
            padding-top: 15px;
        }
        .toggle-button {
            background: #f5f5f5;
            border: none;
            padding: 5px 10px;
            border-radius: 3px;
            cursor: pointer;
            margin-right: 5px;
        }
        .hidden {
            display: none;
        }
        .timestamp {
            color: #777;
            font-size: 0.9em;
        }
        .duration {
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Name}}</h1>
        <div class="timestamp">
            Created: {{formatTime .CreatedAt}} | 
            Duration: <span class="duration">{{formatDuration .Summary.DurationMs}}ms</span>
        </div>
        
        <h2>Summary</h2>
        <div class="summary">
            <div class="stat passed">
                <div class="stat-label">Passed</div>
                <div class="stat-value">{{.Summary.PassedTests}}</div>
            </div>
            <div class="stat failed">
                <div class="stat-label">Failed</div>
                <div class="stat-value">{{.Summary.FailedTests}}</div>
            </div>
            <div class="stat skipped">
                <div class="stat-label">Skipped</div>
                <div class="stat-value">{{.Summary.SkippedTests}}</div>
            </div>
            <div class="stat error">
                <div class="stat-label">Errors</div>
                <div class="stat-value">{{.Summary.ErrorTests}}</div>
            </div>
            <div class="stat">
                <div class="stat-label">Total</div>
                <div class="stat-value">{{.Summary.TotalTests}}</div>
            </div>
            <div class="stat">
                <div class="stat-label">Snapshots Created</div>
                <div class="stat-value">{{.Summary.SnapshotsCreated}}</div>
            </div>
            <div class="stat">
                <div class="stat-label">Snapshots Updated</div>
                <div class="stat-value">{{.Summary.SnapshotsUpdated}}</div>
            </div>
        </div>
        
        <h2>Results</h2>
        <div class="results">
            {{range .Results}}
            <div class="result result-{{.Status}}">
                <div class="result-header">
                    <div class="result-name">{{.Name}}</div>
                    <div class="result-status status-{{.Status}}">{{.Status}}</div>
                </div>
                <div class="result-details">
                    <div class="result-detail">
                        <span class="result-detail-label">File:</span>
                        <span class="result-detail-value">{{.FilePath}}</span>
                    </div>
                    <div class="result-detail">
                        <span class="result-detail-label">Method:</span>
                        <span class="result-detail-value">{{.Request.Method}}</span>
                    </div>
                    <div class="result-detail">
                        <span class="result-detail-label">URL:</span>
                        <span class="result-detail-value">{{.Request.URL}}</span>
                    </div>
                    <div class="result-detail">
                        <span class="result-detail-label">Duration:</span>
                        <span class="result-detail-value">{{formatDuration .Duration}}ms</span>
                    </div>
                    {{if .Tags}}
                    <div class="result-detail">
                        <span class="result-detail-label">Tags:</span>
                        <span class="result-detail-value">{{join .Tags ", "}}</span>
                    </div>
                    {{end}}
                    {{if .Error}}
                    <div class="result-error">{{.Error}}</div>
                    {{end}}
                    {{if .SnapshotResult}}
                    {{if .SnapshotResult.Diff.HasDiff}}
                    <div class="result-detail">
                        <span class="result-detail-label">Snapshot:</span>
                        <span class="result-detail-value">
                            <button class="toggle-button" onclick="toggleElement('diff-{{$index}}')">Toggle Diff</button>
                        </span>
                    </div>
                    <div id="diff-{{$index}}" class="snapshot-diff hidden">{{.SnapshotResult.Diff.DiffString}}</div>
                    {{end}}
                    {{end}}
                    
                    {{if $.IncludeRequests}}
                    <div class="request-response">
                        <button class="toggle-button" onclick="toggleElement('request-{{$index}}')">Toggle Request</button>
                        <div id="request-{{$index}}" class="hidden">
                            <h4>Request Headers</h4>
                            <pre>{{range .Request.Headers}}{{.Name}}: {{.Value}}
{{end}}</pre>
                            {{if .Request.Body}}
                            <h4>Request Body</h4>
                            <pre>{{.Request.Body}}</pre>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                    
                    {{if and $.IncludeResponses .Response}}
                    <div class="request-response">
                        <button class="toggle-button" onclick="toggleElement('response-{{$index}}')">Toggle Response</button>
                        <div id="response-{{$index}}" class="hidden">
                            <h4>Response Status</h4>
                            <pre>{{.Response.StatusCode}} {{.Response.Status}}</pre>
                            <h4>Response Headers</h4>
                            <pre>{{range $key, $values := .Response.Headers}}{{$key}}: {{join $values ", "}}
{{end}}</pre>
                            {{if .Response.Body}}
                            <h4>Response Body</h4>
                            <pre>{{formatBody .Response.Body .Response.ContentType}}</pre>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
    </div>
    <script>
        function toggleElement(id) {
            const element = document.getElementById(id);
            if (element.classList.contains('hidden')) {
                element.classList.remove('hidden');
            } else {
                element.classList.add('hidden');
            }
        }
    </script>
</body>
</html>`

	// Create template functions
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDuration": func(d interface{}) string {
			switch v := d.(type) {
			case time.Duration:
				return fmt.Sprintf("%.2f", float64(v.Milliseconds()))
			case int64:
				return fmt.Sprintf("%.2f", float64(v))
			default:
				return "0.00"
			}
		},
		"join": func(s []string, sep string) string {
			return strings.Join(s, sep)
		},
		"formatBody": func(body []byte, contentType string) string {
			if strings.Contains(contentType, "application/json") {
				var out bytes.Buffer
				err := json.Indent(&out, body, "", "  ")
				if err == nil {
					return out.String()
				}
			}
			return string(body)
		},
		"$index": func() int {
			return 0 // Will be replaced in the loop
		},
	}

	// Parse the template
	tmplObj, err := template.New("report").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Create template data
	data := struct {
		*models.TestReport
		IncludeRequests  bool
		IncludeResponses bool
	}{
		TestReport:       report,
		IncludeRequests:  options.IncludeRequests,
		IncludeResponses: options.IncludeResponses,
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmplObj.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return &buf, nil
}

// generateJUnitReport generates a JUnit XML report
func (s *TestReporterService) generateJUnitReport(report *models.TestReport, options models.TestReportOptions) (io.Reader, error) {
	// Define JUnit XML structures
	type JUnitProperty struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}

	type JUnitProperties struct {
		Properties []JUnitProperty `xml:"property"`
	}

	type JUnitFailure struct {
		Message string `xml:"message,attr"`
		Type    string `xml:"type,attr"`
		Content string `xml:",cdata"`
	}

	type JUnitTestCase struct {
		Name      string        `xml:"name,attr"`
		Classname string        `xml:"classname,attr"`
		Time      float64       `xml:"time,attr"`
		Failure   *JUnitFailure `xml:"failure,omitempty"`
		Skipped   *struct{}     `xml:"skipped,omitempty"`
		SystemOut string        `xml:"system-out,omitempty"`
	}

	type JUnitTestSuite struct {
		Name       string           `xml:"name,attr"`
		Tests      int              `xml:"tests,attr"`
		Failures   int              `xml:"failures,attr"`
		Errors     int              `xml:"errors,attr"`
		Skipped    int              `xml:"skipped,attr"`
		Time       float64          `xml:"time,attr"`
		Timestamp  string           `xml:"timestamp,attr"`
		Properties JUnitProperties  `xml:"properties"`
		TestCases  []JUnitTestCase  `xml:"testcase"`
	}

	type JUnitTestSuites struct {
		XMLName    xml.Name         `xml:"testsuites"`
		TestSuites []JUnitTestSuite `xml:"testsuite"`
	}

	// Create JUnit test cases
	testCases := make([]JUnitTestCase, 0, len(report.Results))
	for _, result := range report.Results {
		testCase := JUnitTestCase{
			Name:      result.Name,
			Classname: result.FilePath,
			Time:      float64(result.Duration.Milliseconds()) / 1000.0,
		}

		switch result.Status {
		case models.TestStatusFailed:
			testCase.Failure = &JUnitFailure{
				Message: "Test failed",
				Type:    "failure",
				Content: result.Error,
			}
		case models.TestStatusError:
			testCase.Failure = &JUnitFailure{
				Message: "Test error",
				Type:    "error",
				Content: result.Error,
			}
		case models.TestStatusSkipped:
			testCase.Skipped = &struct{}{}
		}

		testCases = append(testCases, testCase)
	}

	// Create properties
	properties := make([]JUnitProperty, 0)
	for k, v := range report.Environment {
		properties = append(properties, JUnitProperty{
			Name:  k,
			Value: v,
		})
	}

	// Create test suite
	testSuite := JUnitTestSuite{
		Name:      report.Name,
		Tests:     report.Summary.TotalTests,
		Failures:  report.Summary.FailedTests,
		Errors:    report.Summary.ErrorTests,
		Skipped:   report.Summary.SkippedTests,
		Time:      float64(report.Summary.DurationMs) / 1000.0,
		Timestamp: report.CreatedAt.Format(time.RFC3339),
		Properties: JUnitProperties{
			Properties: properties,
		},
		TestCases: testCases,
	}

	// Create test suites
	testSuites := JUnitTestSuites{
		TestSuites: []JUnitTestSuite{testSuite},
	}

	// Marshal to XML
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")
	if err := encoder.Encode(testSuites); err != nil {
		return nil, fmt.Errorf("failed to encode JUnit report: %w", err)
	}

	return &buf, nil
}

// generateConsoleReport generates a console (text) report
func (s *TestReporterService) generateConsoleReport(report *models.TestReport, options models.TestReportOptions) (io.Reader, error) {
	var buf bytes.Buffer

	// Write report header
	fmt.Fprintf(&buf, "===============================================\n")
	fmt.Fprintf(&buf, "   %s\n", report.Name)
	fmt.Fprintf(&buf, "===============================================\n")
	fmt.Fprintf(&buf, "Started: %s\n", report.Summary.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&buf, "Duration: %.2f ms\n", float64(report.Summary.DurationMs))
	fmt.Fprintf(&buf, "\n")

	// Write summary
	fmt.Fprintf(&buf, "SUMMARY:\n")
	fmt.Fprintf(&buf, "  Total:   %d\n", report.Summary.TotalTests)
	fmt.Fprintf(&buf, "  Passed:  %d\n", report.Summary.PassedTests)
	fmt.Fprintf(&buf, "  Failed:  %d\n", report.Summary.FailedTests)
	fmt.Fprintf(&buf, "  Skipped: %d\n", report.Summary.SkippedTests)
	fmt.Fprintf(&buf, "  Errors:  %d\n", report.Summary.ErrorTests)
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "  Snapshots:\n")
	fmt.Fprintf(&buf, "    Created: %d\n", report.Summary.SnapshotsCreated)
	fmt.Fprintf(&buf, "    Updated: %d\n", report.Summary.SnapshotsUpdated)
	fmt.Fprintf(&buf, "\n")

	// Write results
	fmt.Fprintf(&buf, "RESULTS:\n")
	for i, result := range report.Results {
		// Format status with color if enabled
		status := string(result.Status)
		if options.ColorOutput {
			switch result.Status {
			case models.TestStatusPassed:
				status = "\x1b[32mPASSED\x1b[0m" // Green
			case models.TestStatusFailed:
				status = "\x1b[31mFAILED\x1b[0m" // Red
			case models.TestStatusSkipped:
				status = "\x1b[33mSKIPPED\x1b[0m" // Yellow
			case models.TestStatusError:
				status = "\x1b[35mERROR\x1b[0m" // Magenta
			}
		}

		fmt.Fprintf(&buf, "  %d. %s [%s]\n", i+1, result.Name, status)
		fmt.Fprintf(&buf, "     File: %s\n", result.FilePath)
		if result.Request != nil {
			fmt.Fprintf(&buf, "     Method: %s %s\n", result.Request.Method, result.Request.URL)
		}
		if result.Duration > 0 {
			fmt.Fprintf(&buf, "     Duration: %.2f ms\n", float64(result.Duration.Milliseconds()))
		}
		if len(result.Tags) > 0 {
			fmt.Fprintf(&buf, "     Tags: %s\n", strings.Join(result.Tags, ", "))
		}
		if result.Error != "" {
			fmt.Fprintf(&buf, "     Error: %s\n", result.Error)
		}
		if result.SnapshotResult != nil && result.SnapshotResult.Diff != nil && result.SnapshotResult.Diff.HasDiff {
			fmt.Fprintf(&buf, "     Snapshot Diff: %s\n", summarizeDiff(result.SnapshotResult.Diff.DiffString))
		}
		fmt.Fprintf(&buf, "\n")
	}

	return &buf, nil
}

// summarizeDiff returns a shortened version of a diff for console output
func summarizeDiff(diff string) string {
	// Limit to 5 lines for console output
	lines := strings.Split(diff, "\n")
	if len(lines) > 5 {
		return strings.Join(lines[:5], "\n") + "\n... (truncated)"
	}
	return diff
}
