package models

import (
	"time"
)

// RequestResponse represents a request-response pair
type RequestResponse struct {
	Request  *HTTPRequest   `json:"request"`
	Response *HTTPResponse  `json:"response"`
	Duration time.Duration  `json:"duration"`
}
