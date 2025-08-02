# httpreq

一个功能简单的 Go HTTP 请求库，支持单次提交请求、批量提交请求和流式提交请求，提供灵活的请求构建和代理支持。

## 特性

- 🚀 **多种请求模式**: 支持单次提交请求、批量提交请求和流式提交请求
- 🔧 **灵活的请求构建**: 支持 JSON、表单、多部分表单和纯文本等多种内容类型
- 🌐 **代理支持**: 支持固定代理和动态代理
- ⏱️ **超时控制**: 可配置请求超时时间
- 📊 **详细响应信息**: 包含状态码、响应体、耗时等完整信息
- 🛠️ **易于使用**: 简洁的 API 设计，易于集成

## 安装

```bash
go get github.com/GoEnthusiast/httpreq
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/GoEnthusiast/httpreq/method"
    "github.com/GoEnthusiast/httpreq/reqsingle"
)

func main() {
    // 创建单次请求器(Bool 参数用于控制是否开启 HTTP2)
    requester := reqsingle.NewSingleRequester(false)
    
    // 发送 GET 请求
    req := &reqsingle.Request{
        Method: method.GET,
        URL:    "https://httpbin.org/get",
    }
    
    resp := requester.Do(req)
    if resp.Error != nil {
        fmt.Printf("请求错误: %v\n", resp.Error)
        return
    }
    
    fmt.Printf("状态码: %d\n", resp.ResponseStatusCode)
    fmt.Printf("响应内容: %s\n", string(resp.ResponseBody))
    fmt.Printf("请求耗时: %.2fms\n", resp.Duration)
}
```

### 批量请求示例

```go
package main

import (
    "fmt"
    "github.com/GoEnthusiast/httpreq/method"
    "github.com/GoEnthusiast/httpreq/reqbatch"
)

func main() {
    // 创建批量请求器(Bool 参数用于控制是否开启 HTTP2)
    batchRequester := reqbatch.NewBatchRequester(false)
    
    // 准备多个请求
    requests := []*reqbatch.Request{
        {
            Method: method.GET,
            URL:    "https://httpbin.org/get",
        },
        {
            Method: method.POST,
            URL:    "https://httpbin.org/post",
            Body: map[string]interface{}{
                "name": "张三",
                "age":  25,
            },
            ContentType: method.ContentTypeJSON,
        },
        {
            Method: method.GET,
            URL:    "https://httpbin.org/headers",
            Header: map[string][]string{
                "User-Agent": {"MyApp/1.0"},
            },
        },
    }
    
    // 执行批量请求
    responses := batchRequester.Do(requests)
    
    // 处理响应
    for i, resp := range responses {
        if resp.Error != nil {
            fmt.Printf("请求 %d 失败: %v\n", i+1, resp.Error)
            continue
        }
        fmt.Printf("请求 %d 成功 - 状态码: %d\n", i+1, resp.ResponseStatusCode)
        fmt.Printf("响应内容: %s\n", string(resp.ResponseBody))
        fmt.Printf("请求耗时: %.2fms\n", resp.Duration*1000)
        fmt.Println("------------------------")
    }
}
```

### 流式请求示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/GoEnthusiast/httpreq/method"
    "github.com/GoEnthusiast/httpreq/reqstream"
)

func main() {
    // 创建流式请求器，设置并发数为 3
    streamRequester := reqstream.NewStreamRequester(false, 3)
    
    // 启动请求发送协程
    go func() {
        for i := 0; i < 10; i++ {
            req := &reqstream.Request{
                Method: method.GET,
                URL:    "https://httpbin.org/get",
            }
            streamRequester.Do(req)
        }
    }()
    
    // 监听响应
    for i := 0; i < 10; i++ {
        resp := <-streamRequester.ResponseCh()
        if resp.Error != nil {
            fmt.Printf("请求错误: %v\n", resp.Error)
            continue
        }
        fmt.Printf("响应 %d: %s\n", i+1, string(resp.ResponseBody))
    }
}
```

## API 文档

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
    Error              error     // 错误信息
    StartTime          time.Time // 开始时间
    EndTime            time.Time // 结束时间
    Duration           float64   // 耗时(毫秒)
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
method.ContentTypeForm  // application/x-www-form-urlencoded
method.ContentTypeMulti // multipart/form-data
method.ContentTypeText  // text/plain
```

## 使用示例

### 1. 单次请求 (Single Request)

#### GET 请求

```go
// 无参数 GET 请求
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users",
}

// 带参数 GET 请求
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://api.example.com/users?page=1&limit=10",
}
```

#### POST JSON 请求

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

#### POST 表单请求

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

#### 多部分表单请求 (文件上传)

```go
file, _ := os.Open("document.pdf")
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

#### 使用代理

```go
// 固定代理
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy:  "http://username:password@proxy.example.com:8080",
}

// 动态代理
req := &reqsingle.Request{
    Method: method.GET,
    URL:    "https://httpbin.org/ip",
    Proxy: func(req *http.Request) (*url.URL, error) {
        // 动态获取代理地址
        proxyIP := getProxyFromPool()
        return url.Parse("http://" + proxyIP)
    },
}
```

#### 设置超时和请求头

```go
req := &reqsingle.Request{
    Method:  method.POST,
    URL:     "https://api.example.com/data",
    Timeout: 30 * time.Second,
    Header: http.Header{
        "Authorization": []string{"Bearer your-token"},
        "User-Agent":    []string{"MyApp/1.0"},
    },
    Body: map[string]interface{}{
        "data": "some data",
    },
    ContentType: method.ContentTypeJSON,
}
```

### 2. 批量请求 (Batch Request)

```go
import "github.com/GoEnthusiast/httpreq/reqbatch"

// 创建批量请求器
batchRequester := reqbatch.NewBatchRequester(false)

// 准备多个请求
requests := []*reqbatch.Request{
    {
        Method: method.GET,
        URL:    "https://api.example.com/users/1",
    },
    {
        Method: method.GET,
        URL:    "https://api.example.com/users/2",
    },
    {
        Method: method.POST,
        URL:    "https://api.example.com/users",
        Body: map[string]interface{}{
            "name": "新用户",
        },
        ContentType: method.ContentTypeJSON,
    },
}

// 执行批量请求
responses := batchRequester.Do(requests)

// 处理响应
for i, resp := range responses {
    if resp.Error != nil {
        fmt.Printf("请求 %d 失败: %v\n", i, resp.Error)
        continue
    }
    fmt.Printf("请求 %d 成功: %s\n", i, string(resp.ResponseBody))
}
```

### 3. 流式请求 (Stream Request)

流式请求适用于需要持续发送请求并异步接收响应的场景，如压力测试、数据采集等。

**工作原理:**
- 使用固定数量的 worker 协程处理请求
- 通过通道异步发送请求和接收响应
- 支持并发控制，避免资源耗尽

```go
import (
    "fmt"
    "time"
    "github.com/GoEnthusiast/httpreq/method"
    "github.com/GoEnthusiast/httpreq/reqstream"
)

// 创建流式请求器，设置并发数为 5
streamRequester := reqstream.NewStreamRequester(false, 5)

// 启动请求发送协程
go func() {
    for {
        req := &reqstream.Request{
            Method: method.GET,
            URL:    "https://api.example.com/test",
        }
        streamRequester.Do(req)
    }
}()

// 监听响应
for {
    resp := <-streamRequester.ResponseCh()
    if resp.Error != nil {
        fmt.Printf("请求错误: %v\n", resp.Error)
        continue
    }
    fmt.Printf("收到响应: %s\n", string(resp.ResponseBody))
    fmt.Printf("请求耗时: %.2fms\n", resp.Duration*1000)
    fmt.Println("------------------------")
}
```



## 最佳实践

1. **错误处理**: 始终检查响应中的错误信息
2. **超时设置**: 为所有请求设置合理的超时时间
3. **资源管理**: 及时关闭文件句柄和连接
4. **并发控制**: 在批量请求中控制并发数量
5. **日志记录**: 记录请求和响应的关键信息

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

- v1.0.0: 初始版本，支持单次请求、批量请求和流式请求
