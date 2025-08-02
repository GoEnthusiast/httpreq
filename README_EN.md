# httpreq - Go HTTP Request Library

[![Go Version](https://img.shields.io/badge/Go-1.23.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/GoEnthusiast/httpreq)](https://goreportcard.com/report/github.com/GoEnthusiast/httpreq)

A simple and easy-to-use Go HTTP request library that supports single request submission, batch request submission, and streaming request submission, providing flexible request building, proxy support, and detailed response information.

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Detailed Usage Guide](#detailed-usage-guide)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)
- [FAQ](#faq)

## ✨ Features

- 🚀 **Multiple Request Modes**: Single request submission, batch request submission, streaming request submission
- 🔧 **Flexible Request Building**: Supports JSON, form, multipart form, plain text, and other content types
- 🌐 **Proxy Support**: Fixed proxy and dynamic proxy
- ⏱️ **Timeout Control**: Configurable request timeout
- 📊 **Detailed Response Information**: Status code, response body, duration, and other complete information
- 🔒 **TLS Support**: Custom certificates and keys
- 🛠️ **Easy to Use**: Clean API design, easy to integrate
- 📈 **Performance Optimized**: Connection pool management and concurrency control

## 📦 Installation

```bash
go get github.com/GoEnthusiast/httpreq
go mod tidy
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqsingle"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestSingleGetMethod(t *testing.T) {
	// Create single request handler (false means not enabling HTTP/2)
	requester := reqsingle.NewSingleRequester(false)

	// Send GET request
	req := &request.Request{
		Method: method.GET,
		URL:    "http://127.0.0.1:9000/testGetNoParams",
	}

	resp := requester.Do(req)
	if resp.Error != nil {
		t.Logf("Request error: %v\n", resp.Error)
		return
	}

	t.Logf("Status code: %d\n", resp.ResponseStatusCode)
	t.Logf("Response content: %s\n", string(resp.ResponseBody))
	t.Logf("Request start time: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
	t.Logf("Request end time: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
	t.Logf("Request duration: %.2fs\n", resp.Duration)
}
```

### Batch Request Example

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqbatch"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestBatchGetMethod(t *testing.T) {
	// Create batch request handler
	batchRequester := reqbatch.NewBatchRequester(false)

	// Prepare multiple requests
	requests := []*request.Request{}
	for i := 0; i < 10; i++ {
		requests = append(requests, &request.Request{
			Method: method.GET,
			URL:    "http://127.0.0.1:9000/testGetNoParams",
		})
	}

	// Execute batch requests
	responses := batchRequester.Do(requests)

	// Process responses
	for _, resp := range responses {
		if resp.Error != nil {
			t.Logf("Request error: %v\n", resp.Error)
			return
		}

		t.Logf("Status code: %d\n", resp.ResponseStatusCode)
		t.Logf("Response content: %s\n", string(resp.ResponseBody))
		t.Logf("Request start time: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
		t.Logf("Request end time: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
		t.Logf("Request duration: %.2fs\n", resp.Duration)
	}
}
```

### Streaming Request Example

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqstream"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestStreamGetMethod(t *testing.T) {
	// Create streaming request handler with concurrency of 5
	streamRequester := reqstream.NewStreamRequester(false, 5)

	// Start request sending goroutine
	go func() {
		for i := 0; i < 20; i++ {
			req := &request.Request{
				Method: method.GET,
				URL:    "http://127.0.0.1:9000/testGetNoParams",
			}
			streamRequester.Do(req)
			t.Log("Send request")
		}
	}()

	// Listen for responses
	for i := 0; i < 20; i++ {
		resp := <-streamRequester.ResponseCh()
		if resp.Error != nil {
			t.Logf("Request error: %v\n", resp.Error)
			return
		}

		t.Logf("Status code: %d\n", resp.ResponseStatusCode)
		t.Logf("Response content: %s\n", string(resp.ResponseBody))
		t.Logf("Request start time: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
		t.Logf("Request end time: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
		t.Logf("Request duration: %.2fs\n", resp.Duration)
	}
}
```

## 📖 Detailed Usage Guide

### 1. GET Request

```go
// Simple GET request
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users",
}

// GET request with query parameters
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users?page=1&limit=10",
}

// GET request with headers
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users",
    Header: http.Header{
        "Authorization": []string{"Bearer your-token"},
        "User-Agent":    []string{"MyApp/1.0"},
    },
}
```

### 2. POST Request

**JSON Request:**
```go
req := &reqsingle.Request{
    Method: method.POST,
    URL:    "https://api.example.com/users",
    Body: map[string]interface{}{
        "name": "张三",
        "age":  25,
        "email": "zhangsan@example.com",
    },
    ContentType: method.ContentTypeJSON,
}
```

**Form Request:**
```go
req := &reqsingle.Request{
    Method: method.POST,
    URL:    "https://api.example.com/login",
    Body: map[string]string{
        "username": "user123",
        "password": "pass123",
    },
    ContentType: method.ContentTypeForm,
}
```

**Multipart Form Request (File Upload):**
```go
file, _ := os.Open("document.pdf")
defer file.Close()

req := &reqsingle.Request{
    Method: method.POST,
    URL:    "https://api.example.com/upload",
    Body: map[string]interface{}{
        "file": file,
        "description": "重要文档",
    },
    ContentType: method.ContentTypeMulti,
}
```

### 3. Proxy Settings

**Set Fixed Proxy in Request Handler (Suitable for Long-term Proxy Usage)**
```go
// Create single request handler (false means not enabling HTTP/2)
requester := reqsingle.NewSingleRequester(false)

// If setting fixed proxy directly in request handler, proxy can be omitted in request parameters
if err := requester.SetProxy("http://username:password@proxy.example.com:8080"); err != nil {
t.Logf("Set proxy error: %v\n", err)
}

// Send GET request
req := &request.Request{
Method: method.GET,
URL:    "https://httpbin.org/get",
}

resp := requester.Do(req)
```

**Set Dynamic Proxy in Request Handler (Suitable for Randomly Getting Proxy from Your IP Pool)**
```go
// Create single request handler (false means not enabling HTTP/2)
requester := reqsingle.NewSingleRequester(false)

// If setting fixed proxy directly in request handler, proxy can be omitted in request parameters
if err := requester.SetProxy(func(r *http.Request) (*url.URL, error) {
    myProxy := getProxy()
    return url.Parse("http://" + myProxy)
}); err != nil {
    t.Logf("Set proxy error: %v\n", err)
}

// Send GET request
req := &request.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/get",
}

resp := requester.Do(req)
```

**Set Fixed Proxy in Request Body (Suitable for Long-term Proxy Usage):**
```go
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy:  "http://username:password@proxy.example.com:8080",
}
```

**Set Dynamic Proxy in Request Body (Suitable for Getting Proxy from Your IP Pool for Each Request):**
```go
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy: func(req *http.Request) (*url.URL, error) {
        // Dynamically get proxy from proxy pool
        proxyIP := getProxyFromPool()
        return url.Parse("http://" + proxyIP)
    },
}
```

### 4. Timeout Settings

```go
req := &reqsingle.Request{
    Method:  method.GET,
    URL:     "https://api.example.com/data",
    Timeout: 30 * time.Second, // 30 seconds timeout
}
```

## 📚 API Reference

### Request Structure

```go
type Request struct {
    Method      method.HTTPMethod      // Request method (GET, POST, PUT, DELETE)
    URL         string                 // Request URL
    Header      http.Header            // Request headers
    Body        interface{}            // Request body
    ContentType method.HTTPContentType // Request content type
    Proxy       interface{}            // Proxy settings
    Timeout     time.Duration          // Request timeout
    Meta        map[string]interface{} // Request metadata
}
```

### Response Structure

```go
type Response struct {
    Request            *Request  // Request body
    ResponseStatusCode int       // Response status code
    ResponseBody       []byte    // Response content
    Error              error     // Error information
    StartTime          time.Time // Start time
    EndTime            time.Time // End time
    Duration           float64   // Duration (seconds)
}
```

### Supported HTTP Methods

```go
method.GET    // GET request
method.POST   // POST request
method.PUT    // PUT request
method.DELETE // DELETE request
```

### Supported Content Types

```go
method.ContentTypeJSON  // application/json
method.ContentTypeForm  // application/x-www-form-urlencoded
method.ContentTypeMulti // multipart/form-data
method.ContentTypeText  // text/plain
```

### Advanced Configuration

#### TLS Configuration

```go
requester := reqsingle.NewSingleRequester(false)
err := requester.SetTLS("cert.pem", "key.pem", "ca.pem")
if err != nil {
    log.Fatal(err)
}
```

#### Transport Layer Configuration

```go
requester := reqsingle.NewSingleRequester(false)

// Set connection pool parameters
requester.SetMaxIdleConns(100)
requester.SetMaxIdleConnsPerHost(10)
requester.SetMaxConnsPerHost(100)

// Set timeout parameters
requester.SetIdleConnTimeout(90 * time.Second)
requester.SetTLSHandshakeTimeout(10 * time.Second)
requester.SetExpectContinueTimeout(1 * time.Second)

// Set Keep-Alive control
requester.SetDisableKeepAlives(false) // Enable Keep-Alive (default)
// requester.SetDisableKeepAlives(true) // Disable Keep-Alive
```

## 🎯 Best Practices

### 1. Error Handling

```go
resp := requester.Do(req)
if resp.Error != nil {
    // Handle network errors
    if strings.Contains(resp.Error.Error(), "timeout") {
        fmt.Println("Request timeout")
    } else if strings.Contains(resp.Error.Error(), "connection refused") {
        fmt.Println("Connection refused")
    } else {
        fmt.Printf("Request failed: %v\n", resp.Error)
    }
    return
}

// Check HTTP status code
if resp.ResponseStatusCode >= 400 {
    fmt.Printf("HTTP error: %d\n", resp.ResponseStatusCode)
    return
}
```

### 2. Timeout Settings

```go
// Set appropriate timeout for different types of requests
req := &reqsingle.Request{
    Method:  method.GET,
    URL:     "https://api.example.com/data",
    Timeout: 10 * time.Second, // Short request
}

req := &reqsingle.Request{
    Method:  method.POST,
    URL:     "https://api.example.com/upload",
    Timeout: 60 * time.Second, // Long request
}
```

### 3. Resource Management

```go
// Close file handles promptly
file, err := os.Open("document.pdf")
if err != nil {
    return err
}
defer file.Close()

req := &reqsingle.Request{
    Method: method.POST,
    URL:    "https://api.example.com/upload",
    Body: map[string]interface{}{
        "file": file,
    },
    ContentType: method.ContentTypeMulti,
}
```

### 4. Concurrency Control

```go
// Control concurrency in batch requests
batchRequester := reqbatch.NewBatchRequester(false)

// Process large requests in batches
const batchSize = 100
for i := 0; i < len(allRequests); i += batchSize {
    end := i + batchSize
    if end > len(allRequests) {
        end = len(allRequests)
    }
    
    batch := allRequests[i:end]
    responses := batchRequester.Do(batch)
    
    // Process responses...
}
```

## ❓ FAQ

### Q: How to handle HTTPS certificate verification?

A: By default, the library verifies HTTPS certificates. If you need to skip verification or use custom certificates:

```go
// Skip certificate verification (not recommended for production)
transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
requester.SetTransport(transport)

// Use custom certificates
err := requester.SetTLS("cert.pem", "key.pem", "ca.pem")
```

### Q: How to set request retry?

A: The library doesn't have built-in retry mechanism, but it can be easily implemented:

```go
func doWithRetry(requester reqsingle.SingleRequester, req *reqsingle.Request, maxRetries int) *response.Response {
    for i := 0; i < maxRetries; i++ {
        resp := requester.Do(req)
        if resp.Error == nil && resp.ResponseStatusCode < 500 {
            return resp
        }
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return requester.Do(req) // Last attempt
}
```

### Q: How to handle large file uploads?

A: For large files, it's recommended to set appropriate timeout:

```go
req := &reqsingle.Request{
    Method:  method.POST,
    URL:     "https://api.example.com/upload",
    Timeout: 300 * time.Second, // 5 minutes timeout
    Body: map[string]interface{}{
        "file": file,
    },
    ContentType: method.ContentTypeMulti,
}
```

### Q: How to optimize performance?

A: You can optimize performance in the following ways:

1. **Enable HTTP/2**:
```go
requester := reqsingle.NewSingleRequester(true) // Enable HTTP/2
```

2. **Adjust connection pool parameters**:
```go
requester.SetMaxIdleConns(100)
requester.SetMaxIdleConnsPerHost(10)
requester.SetDisableKeepAlives(false) // Enable Keep-Alive for better performance
```

3. **Use streaming requests for concurrent processing**:
```go
streamRequester := reqstream.NewStreamRequester(false, 10) // 10 concurrent
```

## 📝 Changelog

### v1.0.0
- ✨ Initial version release
- 🚀 Support for single request submission, batch request submission, and streaming request submission
- 🔧 Support for multiple content types and proxy settings
- 📊 Provide detailed response information 