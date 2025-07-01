package salesforce

import "errors"

// configuration is now private to enforce functional configuration pattern
type configuration struct {
	compressionHeaders bool // compress request and response if true to save bandwidth
	apiVersion         string
	batchSizeMax       int
	bulkBatchSizeMax   int
}

// setDefaults sets the default configuration values
func (c *configuration) setDefaults() {
	c.compressionHeaders = false
	c.apiVersion = apiVersion
	c.batchSizeMax = batchSizeMax
	c.bulkBatchSizeMax = bulkBatchSizeMax
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
