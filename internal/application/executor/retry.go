package executor

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

// RetryPolicy defines a policy for HTTP request retries
type RetryPolicy struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// InitialBackoff is the initial backoff duration
	InitialBackoff time.Duration
	// MaxBackoff is the maximum backoff duration
	MaxBackoff time.Duration
	// BackoffFactor is the factor to multiply backoff by after each retry
	BackoffFactor float64
	// Jitter is the factor by which to randomize backoff
	Jitter float64
	// RetryableStatusCodes are HTTP status codes that should trigger a retry
	RetryableStatusCodes []int
	// RetryableErrors are specific errors that should trigger a retry
	RetryableErrors []error
}

// DefaultRetryPolicy returns a reasonable default retry policy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
		Jitter:         0.2,
		RetryableStatusCodes: []int{
			http.StatusRequestTimeout,      // 408
			http.StatusTooManyRequests,     // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout,      // 504
		},
		RetryableErrors: []error{
			context.DeadlineExceeded,
		},
	}
}

// RetryableTransport is an http.RoundTripper that retries requests
type RetryableTransport struct {
	// Transport is the underlying transport
	Transport http.RoundTripper
	// Policy is the retry policy
	Policy *RetryPolicy
	// Logger is the logger to use
	Logger Logger
}

// NewRetryableTransport creates a new retryable transport
func NewRetryableTransport(transport http.RoundTripper, policy *RetryPolicy, logger Logger) *RetryableTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	if policy == nil {
		policy = DefaultRetryPolicy()
	}
	if logger == nil {
		logger = newDefaultLogger()
	}
	return &RetryableTransport{
		Transport: transport,
		Policy:    policy,
		Logger:    logger,
	}
}

// RoundTrip implements the http.RoundTripper interface
func (rt *RetryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp         *http.Response
		err          error
		retryCount   int
		shouldRetry  bool
		backoff      time.Duration
		originalBody []byte
	)

	// Save original body for retries if present
	if req.Body != nil {
		originalBodyReader, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		defer originalBodyReader.Close()
	}

	startTime := time.Now()
	for {
		// Create a new request to ensure it's fresh (especially body)
		if retryCount > 0 {
			if req.Body != nil {
				bodyReader, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				req.Body = bodyReader
			}
		}

		// Execute the request
		resp, err = rt.Transport.RoundTrip(req)

		// Check if we should retry
		shouldRetry, backoff = rt.shouldRetry(req.Context(), resp, err, retryCount)
		if !shouldRetry {
			break
		}

		// Close the response body if we're going to retry
		if resp != nil {
			resp.Body.Close()
		}

		// Log the retry
		rt.Logger.Debugf("Retrying request to %s after error: %v (retry %d/%d, backoff %s)",
			req.URL.String(), err, retryCount+1, rt.Policy.MaxRetries, backoff)

		// Increment retry counter
		retryCount++

		// Apply backoff
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoff):
			// Continue with retry
		}
	}

	if retryCount > 0 && err == nil {
		rt.Logger.Debugf("Request succeeded after %d retries (%s)", retryCount, time.Since(startTime))
	} else if retryCount > 0 {
		rt.Logger.Errorf("Request failed after %d retries (%s): %v", retryCount, time.Since(startTime), err)
	}

	return resp, err
}

// shouldRetry determines if a request should be retried
func (rt *RetryableTransport) shouldRetry(ctx context.Context, resp *http.Response, err error, retryCount int) (bool, time.Duration) {
	// Check if we've reached the maximum number of retries
	if retryCount >= rt.Policy.MaxRetries {
		return false, 0
	}

	// Check if the context is already canceled or deadline exceeded
	if ctx.Err() != nil {
		return false, 0
	}

	// Check for specific network errors that are retryable
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && (netErr.Timeout() || netErr.Temporary()) {
			return true, rt.calculateBackoff(retryCount)
		}

		var urlErr *url.Error
		if errors.As(err, &urlErr) && (urlErr.Timeout() || urlErr.Temporary()) {
			return true, rt.calculateBackoff(retryCount)
		}

		// Check for retryable errors in policy
		for _, retryableErr := range rt.Policy.RetryableErrors {
			if errors.Is(err, retryableErr) {
				return true, rt.calculateBackoff(retryCount)
			}
		}
	}

	// Check for retryable status codes
	if resp != nil {
		for _, code := range rt.Policy.RetryableStatusCodes {
			if resp.StatusCode == code {
				return true, rt.calculateBackoff(retryCount)
			}
		}
	}

	return false, 0
}

// calculateBackoff calculates the backoff duration for a retry
func (rt *RetryableTransport) calculateBackoff(retryCount int) time.Duration {
	// Calculate base backoff with exponential backoff
	backoff := rt.Policy.InitialBackoff * time.Duration(math.Pow(rt.Policy.BackoffFactor, float64(retryCount)))
	
	// Apply jitter to avoid thundering herd problem
	if rt.Policy.Jitter > 0 {
		backoff = time.Duration(float64(backoff) * (1.0 + rt.Policy.Jitter*(rand.Float64()*2-1)))
	}
	
	// Cap at max backoff
	if backoff > rt.Policy.MaxBackoff {
		backoff = rt.Policy.MaxBackoff
	}
	
	return backoff
}

// NewRetryableClient creates a new HTTP client with retry capabilities
func NewRetryableClient(policy *RetryPolicy, logger Logger) *http.Client {
	transport := NewRetryableTransport(http.DefaultTransport, policy, logger)
	return &http.Client{
		Transport: transport,
		Timeout:   time.Minute, // Default longer timeout to allow for retries
	}
}
