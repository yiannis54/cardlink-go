package vposxml

import "net/http"

type Option func(*Client)

// WithHTTPClient sets the HTTP client to use for the client.
// If not set, a default HTTP client with a 30 second timeout will be used.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
