# httpreq - Go HTTP 请求库

[![Go Version](https://img.shields.io/badge/Go-1.23.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/GoEnthusiast/httpreq)](https://goreportcard.com/report/github.com/GoEnthusiast/httpreq)

一个功能简单且易于使用的 Go HTTP 请求库，支持单次提交请求、批量提交请求和流式提交请求，提供灵活的请求构建、代理支持和详细的响应信息。

## 📋 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [详细使用指南](#详细使用指南)
- [API 参考](#api-参考)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## ✨ 特性

- 🚀 **多种请求模式**: 单次提交请求、批量提交请求、流式提交请求
- 🔧 **灵活的请求构建**: 支持 JSON、表单、多部分表单、纯文本等多种内容类型
- 🌐 **代理支持**: 固定代理和动态代理
- ⏱️ **超时控制**: 可配置请求超时时间
- 📊 **详细响应信息**: 状态码、响应体、耗时等完整信息
- 🔒 **TLS 支持**: 自定义证书和密钥
- 🛠️ **易于使用**: 简洁的 API 设计，易于集成
- 📈 **性能优化**: 连接池管理和并发控制

## 📦 安装

```bash
go get github.com/GoEnthusiast/httpreq
go mod tidy
```

## 🚀 快速开始

### 基本使用

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqsingle"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestSingleGetMethod(t *testing.T) {
	// 创建单次提交请求器 (false 表示不启用 HTTP/2)
	requester := reqsingle.NewSingleRequester(false)

	// 发送 GET 请求
	req := &request.Request{
		Method: method.GET,
		URL:    "http://127.0.0.1:9000/testGetNoParams",
	}

	resp := requester.Do(req)
	if resp.Error != nil {
		t.Logf("请求错误: %v\n", resp.Error)
		return
	}

	t.Logf("状态码: %d\n", resp.ResponseStatusCode)
	t.Logf("响应内容: %s\n", string(resp.ResponseBody))
	t.Logf("请求开始时间: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
	t.Logf("请求结束时间: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
	t.Logf("请求耗时: %.2fs\n", resp.Duration)
}

```

### 批量提交请求示例

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqbatch"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestBatchGetMethod(t *testing.T) {
	// 创建批量提交请求器
	batchRequester := reqbatch.NewBatchRequester(false)

	// 准备多个请求
	requests := []*request.Request{}
	for i := 0; i < 10; i++ {
		requests = append(requests, &request.Request{
			Method: method.GET,
			URL:    "http://127.0.0.1:9000/testGetNoParams",
		})
	}

	// 执行批量提交请求
	responses := batchRequester.Do(requests)

	// 处理响应
	for _, resp := range responses {
		if resp.Error != nil {
			t.Logf("请求错误: %v\n", resp.Error)
			return
		}

		t.Logf("状态码: %d\n", resp.ResponseStatusCode)
		t.Logf("响应内容: %s\n", string(resp.ResponseBody))
		t.Logf("请求开始时间: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求结束时间: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求耗时: %.2fs\n", resp.Duration)
	}
}
```

### 流式提交请求示例

```go
package main

import (
	"github.com/GoEnthusiast/httpreq/method"
	"github.com/GoEnthusiast/httpreq/reqstream"
	"github.com/GoEnthusiast/httpreq/types/request"
	"testing"
)

func TestStreamGetMethod(t *testing.T) {
	// 创建流式提交请求器，设置并发数为 5
	streamRequester := reqstream.NewStreamRequester(false, 5)

	// 启动请求发送协程
	go func() {
		for i := 0; i < 20; i++ {
			req := &request.Request{
				Method: method.GET,
				URL:    "http://127.0.0.1:9000/testGetNoParams",
			}
			streamRequester.Do(req)
			t.Log("发送请求")
		}
	}()

	// 监听响应
	for i := 0; i < 20; i++ {
		resp := <-streamRequester.ResponseCh()
		if resp.Error != nil {
			t.Logf("请求错误: %v\n", resp.Error)
			return
		}

		t.Logf("状态码: %d\n", resp.ResponseStatusCode)
		t.Logf("响应内容: %s\n", string(resp.ResponseBody))
		t.Logf("请求开始时间: %s\n", resp.StartTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求结束时间: %s\n", resp.EndTime.Format("2006-01-02 15:04:05"))
		t.Logf("请求耗时: %.2fs\n", resp.Duration)
	}
}
```

## 📖 详细使用指南

### 1. GET 请求

```go
// 简单 GET 请求
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users",
}

// 带查询参数的 GET 请求
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users?page=1&limit=10",
}

// 带请求头的 GET 请求
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users",
    Header: http.Header{
        "Authorization": []string{"Bearer your-token"},
        "User-Agent":    []string{"MyApp/1.0"},
    },
}
```

### 2. POST 请求

**JSON 请求:**
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

**表单请求:**
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

**多部分表单请求 (文件上传):**
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

### 3. 代理设置

**在请求器中设置固定代理(适合持续使用长效代理)**
```go
// 创建单次提交请求器 (false 表示不启用 HTTP/2)
requester := reqsingle.NewSingleRequester(false)

// 若直接在请求器设置固定代理，请求参数中可以不设置代理
if err := requester.SetProxy("http://username:password@proxy.example.com:8080"); err != nil {
t.Logf("设置代理错误: %v\n", err)
}

// 发送 GET 请求
req := &request.Request{
Method: method.GET,
URL:    "https://httpbin.org/get",
}

resp := requester.Do(req)
```

**在请求器中设置动态代理(适合从自己的 IP 池中随机获取代理)**
```go
// 创建单次提交请求器 (false 表示不启用 HTTP/2)
requester := reqsingle.NewSingleRequester(false)

// 若直接在请求器设置固定代理，请求参数中可以不设置代理
if err := requester.SetProxy(func(r *http.Request) (*url.URL, error) {
    myProxy := getProxy()
    return url.Parse("http://" + myProxy)
}); err != nil {
    t.Logf("设置代理错误: %v\n", err)
}

// 发送 GET 请求
req := &request.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/get",
}

resp := requester.Do(req)
```

**在请求体中设置固定代理(适合持续使用长效代理):**
```go
req := &request.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy:  "http://username:password@proxy.example.com:8080",
}
```

**在请求体中设置动态代理(适合每个请求都从自己的 IP 池中随机获取代理):**
```go
req := &request.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy: func(req *http.Request) (*url.URL, error) {
        // 从代理池中动态获取代理
        proxyIP := getProxyFromPool()
        return url.Parse("http://" + proxyIP)
    },
}
```

### 4. 超时设置

```go
req := &request.Request{
    Method:  method.GET,
    URL:     "https://api.example.com/data",
    Timeout: 30 * time.Second, // 30秒超时
}
```

### 5. 设置 Cookiejar，实现会话保持。[非并发安全, 仅支持 single 请求器]

```go
// 使用默认 cookiejar
requester := reqsingle.NewSingleRequester(false)
requester.SetCookieJar(nil)

// 使用自定义 cookiejar
cookieJar, _ := cookiejar.New(nil)
...
requester.SetCookieJar(cookieJar)
```

## 📚 API 参考

### 请求结构体

```go
type Request struct {
    Method      method.HTTPMethod      // 请求方法 (GET, POST, PUT, DELETE)
    URL         string                 // 请求地址
    Header      http.Header            // 请求头
    Body        interface{}            // 请求体
    ContentType method.HTTPContentType // 请求内容类型
    Proxy       interface{}            // 代理设置
    Timeout     time.Duration          // 请求超时时间
    Meta        map[string]interface{} // 请求元数据
}
```

### 响应结构体

```go
type Response struct {
    Request            *Request  // 请求体
    ResponseStatusCode int       // 响应状态码
    ResponseBody       []byte    // 响应内容
    ResponseHeader     http.Header // 响应头
    Error              error     // 错误信息
    StartTime          time.Time // 开始时间
    EndTime            time.Time // 结束时间
    Duration           float64   // 耗时(秒)
}
```

### 支持的 HTTP 方法

```go
method.GET    // GET 请求
method.POST   // POST 请求
method.PUT    // PUT 请求
method.DELETE // DELETE 请求
```

### 支持的内容类型

```go
method.ContentTypeJSON  // application/json
method.ContentTypeForm  // application/x-www-form-urlencoded; charset=UTF-8
method.ContentTypeMulti // multipart/form-data
method.ContentTypeText  // text/plain
```

### 高级配置

#### TLS 配置

```go
requester := reqsingle.NewSingleRequester(false)
err := requester.SetTLS("cert.pem", "key.pem", "ca.pem")
if err != nil {
    log.Fatal(err)
}
```

#### 传输层配置

```go
requester := reqsingle.NewSingleRequester(false)

// 设置连接池参数
requester.SetMaxIdleConns(100)
requester.SetMaxIdleConnsPerHost(10)
requester.SetMaxConnsPerHost(100)

// 设置超时参数
requester.SetIdleConnTimeout(90 * time.Second)
requester.SetTLSHandshakeTimeout(10 * time.Second)
requester.SetExpectContinueTimeout(1 * time.Second)

// 设置 Keep-Alive 控制
requester.SetDisableKeepAlives(false) // 启用 Keep-Alive（默认）
// requester.SetDisableKeepAlives(true) // 禁用 Keep-Alive
```

## 🎯 最佳实践

### 1. 错误处理

```go
resp := requester.Do(req)
if resp.Error != nil {
    // 处理网络错误
    if strings.Contains(resp.Error.Error(), "timeout") {
        fmt.Println("请求超时")
    } else if strings.Contains(resp.Error.Error(), "connection refused") {
        fmt.Println("连接被拒绝")
    } else {
        fmt.Printf("请求失败: %v\n", resp.Error)
    }
    return
}

// 检查 HTTP 状态码
if resp.ResponseStatusCode >= 400 {
    fmt.Printf("HTTP 错误: %d\n", resp.ResponseStatusCode)
    return
}
```

### 2. 超时设置

```go
// 为不同类型的请求设置合适的超时时间
req := &reqsingle.Request{
    Method:  method.GET,
    URL:     "https://api.example.com/data",
    Timeout: 10 * time.Second, // 短请求
}

req := &reqsingle.Request{
    Method:  method.POST,
    URL:     "https://api.example.com/upload",
    Timeout: 60 * time.Second, // 长请求
}
```

### 3. 资源管理

```go
// 及时关闭文件句柄
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

### 4. 并发控制

```go
// 在批量提交请求中控制并发数量
batchRequester := reqbatch.NewBatchRequester(false)

// 分批处理大量请求
const batchSize = 100
for i := 0; i < len(allRequests); i += batchSize {
    end := i + batchSize
    if end > len(allRequests) {
        end = len(allRequests)
    }
    
    batch := allRequests[i:end]
    responses := batchRequester.Do(batch)
    
    // 处理响应...
}
```


## ❓ 常见问题

### Q: 如何处理 HTTPS 证书验证？

A: 默认情况下，库会验证 HTTPS 证书。如果需要跳过验证或使用自定义证书：

```go
// 跳过证书验证 (不推荐用于生产环境)
transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
requester.SetTransport(transport)

// 使用自定义证书
err := requester.SetTLS("cert.pem", "key.pem", "ca.pem")
```

### Q: 如何设置请求重试？

A: 目前库不内置重试机制，但可以轻松实现：

```go
func doWithRetry(requester reqsingle.SingleRequester, req *reqsingle.Request, maxRetries int) *response.Response {
    for i := 0; i < maxRetries; i++ {
        resp := requester.Do(req)
        if resp.Error == nil && resp.ResponseStatusCode < 500 {
            return resp
        }
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return requester.Do(req) // 最后一次尝试
}
```

### Q: 如何处理大文件上传？

A: 对于大文件，建议设置合适的超时时间：

```go
req := &reqsingle.Request{
    Method:  method.POST,
    URL:     "https://api.example.com/upload",
    Timeout: 300 * time.Second, // 5分钟超时
    Body: map[string]interface{}{
        "file": file,
    },
    ContentType: method.ContentTypeMulti,
}
```

### Q: 如何优化性能？

A: 可以通过以下方式优化性能：

1. **启用 HTTP/2**:
```go
requester := reqsingle.NewSingleRequester(true) // 启用 HTTP/2
```

2. **调整连接池参数**:
```go
requester.SetMaxIdleConns(100)
requester.SetMaxIdleConnsPerHost(10)
requester.SetDisableKeepAlives(false) // 启用 Keep-Alive 提高性能
```

3. **使用流式提交请求进行并发处理**:
```go
streamRequester := reqstream.NewStreamRequester(false, 10) // 10个并发
```

