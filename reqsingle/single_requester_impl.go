package reqsingle

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

// SingleRequesterImpl implements the SingleRequester interface
// SingleRequesterImpl 实现 SingleRequester 接口
type SingleRequesterImpl struct {
	*core.RequestHandler // Embedded request handler / 嵌入的请求处理器
}

// Do executes a single HTTP request and returns the response
// Do 执行单个 HTTP 请求并返回响应
func (s *SingleRequesterImpl) Do(req *request.Request) *response.Response {
	return s.RequestHandler.ProcessRequest(req)
}

// NewSingleRequester creates a new single request handler with optional HTTP/2 support
// NewSingleRequester 创建一个新的单次请求处理器，支持可选的 HTTP/2
func NewSingleRequester(enableHttp2 bool) SingleRequester {
	return &SingleRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
	}
}
