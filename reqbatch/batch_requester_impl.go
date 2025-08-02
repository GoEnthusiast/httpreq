package reqbatch

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

// BatchRequesterImpl implements the BatchRequester interface
// BatchRequesterImpl 实现 BatchRequester 接口
type BatchRequesterImpl struct {
	*core.RequestHandler // Embedded request handler / 嵌入的请求处理器
}

// Do executes multiple HTTP requests concurrently and returns all responses
// Do 并发执行多个 HTTP 请求并返回所有响应
func (s *BatchRequesterImpl) Do(reqs []*request.Request) []*response.Response {
	var (
		// Buffered channel to collect responses
		// 带缓冲的通道，用于收集响应
		respCh = make(chan *response.Response, len(reqs))
	)

	// Launch goroutines for concurrent request processing
	// 启动协程进行并发请求处理
	for i := range reqs {
		req := reqs[i] // Avoid goroutine closure reference error / 避免 goroutine 闭包引用错误
		go func(r *request.Request) {
			resp := s.RequestHandler.ProcessRequest(r)
			respCh <- resp
		}(req)
	}

	// Collect all responses
	// 收集所有响应
	responses := make([]*response.Response, 0, len(reqs))
	for i := 0; i < len(reqs); i++ {
		responses = append(responses, <-respCh)
	}
	return responses
}

// NewBatchRequester creates a new batch request handler with optional HTTP/2 support
// NewBatchRequester 创建一个新的批量请求处理器，支持可选的 HTTP/2
func NewBatchRequester(enableHttp2 bool) BatchRequester {
	return &BatchRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
	}
}
