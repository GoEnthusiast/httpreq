// Package response provides the core response structure for HTTP responses
// 包 response 提供 HTTP 响应的核心响应结构
package response

import (
	"time"

	"github.com/GoEnthusiast/httpreq/types/request"
)

// Response represents an HTTP response with timing and error information
// Response 表示包含时间和错误信息的 HTTP 响应
type Response struct {
	Request            *request.Request // Original request object / 原始请求对象
	ResponseStatusCode int              // HTTP response status code / HTTP 响应状态码
	ResponseBody       []byte           // Response body content / 响应体内容
	Error              error            // Error occurred during request / 请求过程中发生的错误
	StartTime          time.Time        // Request start time / 请求开始时间
	EndTime            time.Time        // Request end time / 请求结束时间
	Duration           float64          // Request duration in milliseconds / 请求耗时（毫秒）
}
