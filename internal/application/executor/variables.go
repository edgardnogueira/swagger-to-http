package executor

import (
	"os"
	"regexp"
	"sync"
)

// VariableStore defines the interface for managing HTTP request variables
type VariableStore interface {
	// Get retrieves a variable value by name
	Get(name string) (string, bool)

	// Set stores a variable with the given name and value
	Set(name, value string)

	// Delete removes a variable
	Delete(name string)

	// Clear removes all variables
	Clear()

	// GetAll returns a copy of all variables
	GetAll() map[string]string

	// LoadFromEnvironment loads variables from environment variables with optional prefix
	LoadFromEnvironment(prefix string)

	// LoadFromMap loads variables from a map
	LoadFromMap(vars map[string]string)

	// HasVariable checks if a string contains variable references
	HasVariable(input string) bool

	// ExtractVariables extracts variable names from input string
	ExtractVariables(input string) []string
}

// memoryVariableStore implements VariableStore with in-memory storage
type memoryVariableStore struct {
	variables map[string]string
	mu        sync.RWMutex
	varRegex  *regexp.Regexp
}

// newMemoryVariableStore creates a new in-memory variable store
func newMemoryVariableStore() *memoryVariableStore {
	return &memoryVariableStore{
		variables: make(map[string]string),
		varRegex:  regexp.MustCompile(`{{([^{}]+)}}`),
	}
}

// Get retrieves a variable value by name
func (s *memoryVariableStore) Get(name string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.variables[name]
	return value, exists
}

// Set stores a variable with the given name and value
func (s *memoryVariableStore) Set(name, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.variables[name] = value
}

// Delete removes a variable
func (s *memoryVariableStore) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.variables, name)
}

// Clear removes all variables
func (s *memoryVariableStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.variables = make(map[string]string)
}

// GetAll returns a copy of all variables
func (s *memoryVariableStore) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]string, len(s.variables))
	for k, v := range s.variables {
		result[k] = v
	}
	return result
}

// LoadFromEnvironment loads variables from environment variables with optional prefix
func (s *memoryVariableStore) LoadFromEnvironment(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, env := range os.Environ() {
		// Split the environment variable into key=value
		for i := 0; i < len(env); i++ {
			if env[i] == '=' {
				key := env[:i]
				value := env[i+1:]

				// If prefix is provided, only load variables with that prefix
				if prefix != "" {
					if len(key) > len(prefix) && key[:len(prefix)] == prefix {
						// Remove prefix from key
						s.variables[key[len(prefix):]] = value
					}
				} else {
					s.variables[key] = value
				}
				break
			}
		}
	}
}

// LoadFromMap loads variables from a map
func (s *memoryVariableStore) LoadFromMap(vars map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range vars {
		s.variables[k] = v
	}
}

// HasVariable checks if a string contains variable references
func (s *memoryVariableStore) HasVariable(input string) bool {
	return s.varRegex.MatchString(input)
}

// ExtractVariables extracts variable names from input string
func (s *memoryVariableStore) ExtractVariables(input string) []string {
	matches := s.varRegex.FindAllStringSubmatch(input, -1)
	results := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			results = append(results, match[1])
		}
	}
	return results
}
