package salesforce

import (
	"errors"
	"net/http"
)

// configuration is now private to enforce functional configuration pattern
type configuration struct {
	compressionHeaders           bool // compress request and response if true to save bandwidth
	apiVersion                   string
	batchSizeMax                 int
	bulkBatchSizeMax             int
	httpClient                   *http.Client      // HTTP client to use for requests
	roundTripper                 http.RoundTripper // Custom round tripper
	shouldValidateAuthentication bool              // Validate session on client creation
}

// setDefaults sets the default configuration values
func (c *configuration) setDefaults() {
	c.compressionHeaders = false
	c.apiVersion = apiVersion
	c.batchSizeMax = batchSizeMax
	c.bulkBatchSizeMax = bulkBatchSizeMax
}

func (c *configuration) configureHttpClient() {
	// Set default HTTP client if none provided
	if c.httpClient == nil && c.roundTripper == nil {
		c.httpClient = &http.Client{
			Timeout: httpDefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:       httpDefaultMaxIdleConnections,
				IdleConnTimeout:    httpDefaultIdleConnTimeout,
				DisableCompression: false,
			},
		}
	} else if c.roundTripper != nil {
		// Use custom round tripper with default timeout
		c.httpClient = &http.Client{
			Transport: c.roundTripper,
			Timeout:   httpDefaultIdleConnTimeout,
		}
	}
}

// Option is a functional configuration option that can return an error
type Option func(*configuration) error

// WithCompressionHeaders sets whether to compress request and response headers
func WithCompressionHeaders(compression bool) Option {
	return func(c *configuration) error {
		c.compressionHeaders = compression
		return nil
	}
}

// WithAPIVersion sets the Salesforce API version to use
func WithAPIVersion(version string) Option {
	return func(c *configuration) error {
		if version == "" {
			return errors.New("API version cannot be empty")
		}
		c.apiVersion = version
		return nil
	}
}

// WithBatchSizeMax sets the maximum batch size for collections
func WithBatchSizeMax(size int) Option {
	return func(c *configuration) error {
		if size < 1 || size > 200 {
			return errors.New("batch size max must be between 1 and 200")
		}
		c.batchSizeMax = size
		return nil
	}
}

// WithBulkBatchSizeMax sets the maximum batch size for bulk operations
func WithBulkBatchSizeMax(size int) Option {
	return func(c *configuration) error {
		if size < 1 || size > 10000 {
			return errors.New("bulk batch size max must be between 1 and 10000")
		}
		c.bulkBatchSizeMax = size
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client for making requests
func WithHTTPClient(client *http.Client) Option {
	return func(c *configuration) error {
		if client == nil {
			return errors.New("HTTP client cannot be nil")
		}
		c.httpClient = client
		c.roundTripper = nil // Clear round tripper if client is set
		return nil
	}
}

// WithRoundTripper sets a custom round tripper for HTTP requests
func WithRoundTripper(rt http.RoundTripper) Option {
	return func(c *configuration) error {
		if rt == nil {
			return errors.New("round tripper cannot be nil")
		}
		c.roundTripper = rt
		c.httpClient = nil // Will be set in setDefaults with the round tripper
		return nil
	}
}

// WithValidateAuthentication sets whether to validate the authentication session on client creation
func WithValidateAuthentication(validate bool) Option {
	return func(c *configuration) error {
		c.shouldValidateAuthentication = validate
		return nil
	}
}
