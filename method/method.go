// Package method provides HTTP method constants and validation utilities.
// 包 method 提供 HTTP 方法常量和验证工具
package method

// HTTPMethod represents an HTTP request method (GET, POST, PUT, DELETE, etc.)
// HTTPMethod 表示 HTTP 请求方法 (GET, POST, PUT, DELETE 等)
type HTTPMethod string

// HTTP method constants
// HTTP 方法常量
const (
	GET    HTTPMethod = "GET"    // HTTP GET method / HTTP GET 方法
	POST   HTTPMethod = "POST"   // HTTP POST method / HTTP POST 方法
	PUT    HTTPMethod = "PUT"    // HTTP PUT method / HTTP PUT 方法
	DELETE HTTPMethod = "DELETE" // HTTP DELETE method / HTTP DELETE 方法
)

// IsValid checks if the HTTP method is valid (supported by this library)
// IsValid 检查 HTTP 方法是否有效（本库支持的）
func (m HTTPMethod) IsValid() bool {
	switch m {
	case GET, POST, PUT, DELETE:
		return true
	}
	return false
}
