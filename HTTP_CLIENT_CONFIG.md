# HTTP Client Configuration

## Overview

The Salesforce Go client now supports custom HTTP client configuration, allowing you to:

1. **Provide a custom `http.Client`** - Full control over HTTP client settings including timeouts, TLS configuration, connection pooling, etc.
2. **Provide a custom `http.RoundTripper`** - Lower-level control over HTTP request/response cycle
3. **Use default HTTP client** - Sensible defaults are provided when no custom configuration is specified

## Usage Examples

### Custom HTTP Client

```go
package main

import (
    "crypto/tls"
    "net/http"
    "time"
    
    salesforce "github.com/mutovkin/go-salesforce/v2"
)

func main() {
    // Create a custom HTTP client with specific timeout and TLS settings
    customClient := &http.Client{
        Timeout: 60 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
                MinVersion:         tls.VersionTLS12,
            },
            MaxIdleConns:       20,
            IdleConnTimeout:    90 * time.Second,
            DisableCompression: false,
        },
    }

    creds := salesforce.Creds{
        Domain:         "your-domain.my.salesforce.com",
        Username:       "your-username", 
        Password:       "your-password",
        SecurityToken:  "your-security-token",
        ConsumerKey:    "your-consumer-key",
        ConsumerSecret: "your-consumer-secret",
    }

    // Initialize with custom HTTP client
    sf, err := salesforce.Init(creds, salesforce.WithHTTPClient(customClient))
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the client normally...
}
```

### Custom Round Tripper

```go
func main() {
    // Create a custom round tripper for fine-grained HTTP control
    customRoundTripper := &http.Transport{
        TLSClientConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        DisableCompression:  false,
        DisableKeepAlives:   false,
        MaxIdleConnsPerHost: 5,
    }

    creds := salesforce.Creds{
        // ... your credentials
    }

    // Initialize with custom round tripper
    sf, err := salesforce.Init(creds, 
        salesforce.WithRoundTripper(customRoundTripper),
        salesforce.WithAPIVersion("v64.0"),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Combining Multiple Configuration Options

```go
func main() {
    creds := salesforce.Creds{
        // ... your credentials  
    }

    customClient := &http.Client{
        Timeout: 45 * time.Second,
    }

    // Combine multiple configuration options
    sf, err := salesforce.Init(creds,
        salesforce.WithHTTPClient(customClient),
        salesforce.WithCompressionHeaders(true),
        salesforce.WithAPIVersion("v65.0"),
        salesforce.WithBatchSizeMax(150),
        salesforce.WithBulkBatchSizeMax(5000),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Options

### HTTP Client Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithHTTPClient(client *http.Client)` | Set a custom HTTP client | Default client with 30s timeout |
| `WithRoundTripper(rt http.RoundTripper)` | Set a custom round tripper | Default transport |

### Other Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithAPIVersion(version string)` | Set Salesforce API version | v63.0 |
| `WithCompressionHeaders(enabled bool)` | Enable/disable compression | false |
| `WithBatchSizeMax(size int)` | Set max batch size for collections | 200 |
| `WithBulkBatchSizeMax(size int)` | Set max batch size for bulk operations | 10000 |

## Default HTTP Client Configuration

When no custom HTTP client or round tripper is provided, the library uses:

```go
&http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: false,
    },
}
```

## Accessing Configuration

You can retrieve the current configuration using getter methods:

```go
sf, _ := salesforce.Init(creds, salesforce.WithHTTPClient(customClient))

// Get the configured HTTP client
client := sf.GetHTTPClient()

// Get other configuration values
apiVersion := sf.GetAPIVersion()
batchSizeMax := sf.GetBatchSizeMax()
bulkBatchSizeMax := sf.GetBulkBatchSizeMax()
compressionEnabled := sf.GetCompressionHeaders()
```

## Implementation Details

- The `doRequest` function now uses the configured HTTP client instead of `http.DefaultClient`
- The API version from configuration is used in all endpoint URLs
- When both `WithHTTPClient()` and `WithRoundTripper()` are used, the last one takes precedence
- Custom configurations are validated during initialization to prevent runtime errors

## Migration from Previous Versions

This change is fully backward compatible. Existing code will continue to work without any modifications, using the default HTTP client configuration.
