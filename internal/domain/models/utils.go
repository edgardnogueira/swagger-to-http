package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StringUtils provides utility functions for string operations
// that help with compatibility issues between string and []byte

// BytesToString safely converts a byte slice to a string
func BytesToString(data []byte) string {
	if data == nil {
		return ""
	}
	return string(data)
}

// StringToBytes safely converts a string to a byte slice
func StringToBytes(data string) []byte {
	return []byte(data)
}

// FormatJSON formats a JSON string for pretty printing
func FormatJSON(jsonStr string) (string, error) {
	var obj interface{}
	
	if jsonStr == "" {
		return "", nil
	}
	
	err := json.Unmarshal(StringToBytes(jsonStr), &obj)
	if err != nil {
		return jsonStr, err
	}
	
	formattedJSON, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return jsonStr, err
	}
	
	return BytesToString(formattedJSON), nil
}

// IsValidJSON checks if a string is valid JSON
func IsValidJSON(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal(StringToBytes(jsonStr), &js) == nil
}

// StripWhitespace removes all whitespace from a string
func StripWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, s)
}

// TruncateString truncates a string to a maximum length and adds ellipsis if needed
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	return s[:maxLen-3] + "..."
}

// EnsurePathSeparator ensures a path ends with a separator
func EnsurePathSeparator(path string) string {
	if path == "" {
		return ""
	}
	
	if path[len(path)-1] != '/' {
		return path + "/"
	}
	
	return path
}

// JoinPaths joins path segments with appropriate separators
func JoinPaths(segments ...string) string {
	if len(segments) == 0 {
		return ""
	}
	
	var result strings.Builder
	result.WriteString(segments[0])
	
	for _, segment := range segments[1:] {
		if result.Len() > 0 && result.String()[result.Len()-1] != '/' && segment != "" && segment[0] != '/' {
			result.WriteString("/")
		}
		
		// Avoid double slashes
		if result.Len() > 0 && result.String()[result.Len()-1] == '/' && segment != "" && segment[0] == '/' {
			result.WriteString(segment[1:])
		} else {
			result.WriteString(segment)
		}
	}
	
	return result.String()
}

// FormatError formats an error message or returns an empty string if err is nil
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// StringToError converts a string to an error, or returns nil if the string is empty
func StringToError(errStr string) error {
	if errStr == "" {
		return nil
	}
	return fmt.Errorf("%s", errStr)
}

// SafeStringValue returns a safe string value, converting nil to ""
func SafeStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SafeBoolValue returns a safe bool value, converting nil to false
func SafeBoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// ConvertMapToStringMap converts a map with interface{} values to a map with string values
func ConvertMapToStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
