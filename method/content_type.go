package method

// HTTPContentType represents the Content-Type header value for HTTP requests
// HTTPContentType 表示 HTTP 请求的 Content-Type 头部值
type HTTPContentType string

// HTTP content type constants
// HTTP 内容类型常量
const (
	ContentTypeJSON  HTTPContentType = "application/json"                  // JSON content type / JSON 内容类型
	ContentTypeForm  HTTPContentType = "application/x-www-form-urlencoded" // Form data content type / 表单数据内容类型
	ContentTypeMulti HTTPContentType = "multipart/form-data"               // Multipart form data content type / 多部分表单数据内容类型
	ContentTypeText  HTTPContentType = "text/plain"                        // Plain text content type / 纯文本内容类型
)

// IsValid checks if the content type is valid (supported by this library)
// IsValid 检查内容类型是否有效（本库支持的）
func (c HTTPContentType) IsValid() bool {
	switch c {
	case ContentTypeJSON, ContentTypeForm, ContentTypeMulti, ContentTypeText:
		return true
	}
	return false
}
