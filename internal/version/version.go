package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the application version, set during build
	Version = "dev"
	
	// Commit is the git commit hash, set during build
	Commit = "unknown"
	
	// BuildDate is the build date, set during build
	BuildDate = "unknown"
)

// Info returns version information as a formatted string
func Info() string {
	return fmt.Sprintf("swagger-to-http version %s (commit: %s, built: %s, %s/%s)",
		Version, Commit, BuildDate, runtime.GOOS, runtime.GOARCH)
}

// ShortInfo returns a short version string
func ShortInfo() string {
	return fmt.Sprintf("swagger-to-http v%s", Version)
}
