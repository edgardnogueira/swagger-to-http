package executor

import (
	"net/http"
	"sync"
)

// SessionStore defines the interface for managing HTTP session cookies
type SessionStore interface {
	// SetCookie adds or updates a cookie for a specific host
	SetCookie(host string, cookie *http.Cookie)

	// GetCookies gets all cookies for a specific host
	GetCookies(host string) []*http.Cookie

	// ClearCookies removes all cookies for a specific host
	ClearCookies(host string)

	// ClearAllCookies removes all cookies from all hosts
	ClearAllCookies()

	// GetHosts returns a slice of all hosts that have cookies
	GetHosts() []string
}

// memorySessionStore implements SessionStore with in-memory storage
type memorySessionStore struct {
	cookies map[string][]*http.Cookie
	mu      sync.RWMutex
}

// newMemorySessionStore creates a new in-memory session store
func newMemorySessionStore() *memorySessionStore {
	return &memorySessionStore{
		cookies: make(map[string][]*http.Cookie),
	}
}

// SetCookie adds or updates a cookie for a specific host
func (s *memorySessionStore) SetCookie(host string, cookie *http.Cookie) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find and replace existing cookie with same name, or append new cookie
	cookies := s.cookies[host]
	replaced := false
	for i, c := range cookies {
		if c.Name == cookie.Name {
			cookies[i] = cookie
			replaced = true
			break
		}
	}
	if !replaced {
		cookies = append(cookies, cookie)
	}
	s.cookies[host] = cookies
}

// GetCookies gets all cookies for a specific host
func (s *memorySessionStore) GetCookies(host string) []*http.Cookie {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of the cookies slice to avoid mutation issues
	cookies := s.cookies[host]
	result := make([]*http.Cookie, len(cookies))
	copy(result, cookies)
	return result
}

// ClearCookies removes all cookies for a specific host
func (s *memorySessionStore) ClearCookies(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.cookies, host)
}

// ClearAllCookies removes all cookies from all hosts
func (s *memorySessionStore) ClearAllCookies() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cookies = make(map[string][]*http.Cookie)
}

// GetHosts returns a slice of all hosts that have cookies
func (s *memorySessionStore) GetHosts() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hosts := make([]string, 0, len(s.cookies))
	for host := range s.cookies {
		hosts = append(hosts, host)
	}
	return hosts
}

// fileSessionStore could be implemented for persisting cookies across runs
// This would be useful for long-term session management
type fileSessionStore struct {
	// TODO: Implement file-based session storage if needed
}
