package executor

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	// ApplyAuth applies authentication to the request
	ApplyAuth(req *http.Request) error
	// RefreshAuth refreshes the authentication if needed
	RefreshAuth(ctx context.Context) error
	// GetType returns the type of authentication
	GetType() string
}

// BasicAuthProvider implements basic authentication
type BasicAuthProvider struct {
	username string
	password string
}

// NewBasicAuthProvider creates a new basic authentication provider
func NewBasicAuthProvider(username, password string) *BasicAuthProvider {
	return &BasicAuthProvider{
		username: username,
		password: password,
	}
}

// ApplyAuth implements the AuthProvider interface
func (p *BasicAuthProvider) ApplyAuth(req *http.Request) error {
	auth := p.username + ":" + p.password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encoded)
	return nil
}

// RefreshAuth implements the AuthProvider interface
func (p *BasicAuthProvider) RefreshAuth(ctx context.Context) error {
	// Basic auth doesn't need refreshing
	return nil
}

// GetType implements the AuthProvider interface
func (p *BasicAuthProvider) GetType() string {
	return "basic"
}

// BearerTokenProvider implements bearer token authentication
type BearerTokenProvider struct {
	token string
}

// NewBearerTokenProvider creates a new bearer token authentication provider
func NewBearerTokenProvider(token string) *BearerTokenProvider {
	return &BearerTokenProvider{
		token: token,
	}
}

// ApplyAuth implements the AuthProvider interface
func (p *BearerTokenProvider) ApplyAuth(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+p.token)
	return nil
}

// RefreshAuth implements the AuthProvider interface
func (p *BearerTokenProvider) RefreshAuth(ctx context.Context) error {
	// Simple bearer token doesn't need refreshing
	return nil
}

// GetType implements the AuthProvider interface
func (p *BearerTokenProvider) GetType() string {
	return "bearer"
}

// OAuthConfig contains OAuth configuration
type OAuthConfig struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	Scopes       []string
	GrantType    string
}

// OAuthProvider implements OAuth authentication
type OAuthProvider struct {
	config      OAuthConfig
	token       string
	expiresAt   time.Time
	refreshToken string
	client      *http.Client
	logger      Logger
}

// NewOAuthProvider creates a new OAuth authentication provider
func NewOAuthProvider(config OAuthConfig, logger Logger) *OAuthProvider {
	if logger == nil {
		logger = newDefaultLogger()
	}

	// Create a custom client for token requests
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Never skip verification by default
			},
		},
	}

	return &OAuthProvider{
		config: config,
		client: client,
		logger: logger,
	}
}

// ApplyAuth implements the AuthProvider interface
func (p *OAuthProvider) ApplyAuth(req *http.Request) error {
	// Check if the token needs refreshing
	if p.token == "" || (p.expiresAt.Before(time.Now()) && p.expiresAt.Year() > 1) {
		if err := p.RefreshAuth(req.Context()); err != nil {
			return fmt.Errorf("failed to refresh OAuth token: %w", err)
		}
	}

	req.Header.Set("Authorization", "Bearer "+p.token)
	return nil
}

// RefreshAuth implements the AuthProvider interface
func (p *OAuthProvider) RefreshAuth(ctx context.Context) error {
	// Create form data for token request
	data := url.Values{}

	// Handle different grant types
	if p.refreshToken != "" {
		// Refresh token grant
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", p.refreshToken)
	} else if p.config.GrantType == "password" {
		// Resource owner password credentials grant
		data.Set("grant_type", "password")
		data.Set("username", p.config.Username)
		data.Set("password", p.config.Password)
	} else if p.config.GrantType == "client_credentials" {
		// Client credentials grant
		data.Set("grant_type", "client_credentials")
	} else {
		return fmt.Errorf("unsupported grant type: %s", p.config.GrantType)
	}

	// Add common parameters
	data.Set("client_id", p.config.ClientID)
	if p.config.ClientSecret != "" {
		data.Set("client_secret", p.config.ClientSecret)
	}
	if len(p.config.Scopes) > 0 {
		data.Set("scope", strings.Join(p.config.Scopes, " "))
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the token response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Update token data
	p.token = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		p.refreshToken = tokenResp.RefreshToken
	}

	// Calculate expiration time if provided
	if tokenResp.ExpiresIn > 0 {
		// Set expiration with a small buffer to avoid edge cases
		p.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)
	}

	p.logger.Debugf("OAuth token refreshed, expires in %d seconds", tokenResp.ExpiresIn)
	return nil
}

// GetType implements the AuthProvider interface
func (p *OAuthProvider) GetType() string {
	return "oauth2"
}
