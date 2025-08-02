// Package request provides the core request structure for HTTP requests
// 包 request 提供 HTTP 请求的核心请求结构
package request

import (
	"net/http"
	"time"

	"github.com/GoEnthusiast/httpreq/method"
)

// Request represents an HTTP request with all necessary parameters
// Request 表示包含所有必要参数的 HTTP 请求
type Request struct {
	Method      method.HTTPMethod      // HTTP request method (GET, POST, PUT, DELETE) / HTTP 请求方法 (GET, POST, PUT, DELETE)
	URL         string                 // Request URL / 请求地址
	Header      http.Header            // HTTP request headers / HTTP 请求头
	Body        interface{}            // Request body data / 请求体数据
	ContentType method.HTTPContentType // Content-Type header value / 请求内容类型
	Proxy       interface{}            // Proxy configuration (string or function) / 代理配置 (字符串或函数)
	Timeout     time.Duration          // Request timeout duration / 请求超时时间
	Meta        map[string]interface{} // Request metadata for custom use / 请求元数据，供自定义使用
}
