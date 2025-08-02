package reqstream

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

// StreamRequesterImpl implements the StreamRequester interface
// StreamRequesterImpl 实现 StreamRequester 接口
type StreamRequesterImpl struct {
	*core.RequestHandler                         // Embedded request handler / 嵌入的请求处理器
	reqCh                chan *request.Request   // Channel for incoming requests / 接收请求的通道
	respCh               chan *response.Response // Channel for outgoing responses / 发送响应的通道
}

// worker processes requests from the request channel
// worker 处理来自请求通道的请求
func (s *StreamRequesterImpl) worker() {
	for req := range s.reqCh {
		s.handleRequest(req)
	}
}

// handleRequest processes a single request and sends the response
// handleRequest 处理单个请求并发送响应
func (s *StreamRequesterImpl) handleRequest(req *request.Request) {
	resp := s.RequestHandler.ProcessRequest(req)
	s.respCh <- resp
}

// Do submits a request to the stream for processing
// Do 将请求提交到流中进行处理
func (s *StreamRequesterImpl) Do(req *request.Request) {
	s.reqCh <- req
}

// ResponseCh returns a channel for receiving responses
// ResponseCh 返回用于接收响应的通道
func (s *StreamRequesterImpl) ResponseCh() <-chan *response.Response {
	return s.respCh
}

// NewStreamRequester creates a new stream request handler with configurable concurrency
// NewStreamRequester 创建一个新的流式请求处理器，支持可配置的并发数
func NewStreamRequester(enableHttp2 bool, concurrency int) StreamRequester {
	s := &StreamRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
		reqCh:          make(chan *request.Request, concurrency), // Adjustable buffered channel / 可调节的缓冲通道
		respCh:         make(chan *response.Response),
	}

	// Ensure minimum concurrency of 1
	// 确保最小并发数为 1
	if concurrency <= 0 {
		concurrency = 1
	}

	// Start fixed number of worker goroutines
	// 启动固定数量的 worker 协程
	for i := 0; i < concurrency; i++ {
		go s.worker()
	}
	return s
}
