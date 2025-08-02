package reqstream

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

type StreamRequesterImpl struct {
	*core.RequestHandler
	reqCh  chan *request.Request
	respCh chan *response.Response
}

func (s *StreamRequesterImpl) worker() {
	for req := range s.reqCh {
		s.handleRequest(req)
	}
}

func (s *StreamRequesterImpl) handleRequest(req *request.Request) {
	resp := s.RequestHandler.ProcessRequest(req)
	s.respCh <- resp
}

func (s *StreamRequesterImpl) Do(req *request.Request) {
	s.reqCh <- req
}

func (s *StreamRequesterImpl) ResponseCh() <-chan *response.Response {
	return s.respCh
}

func NewStreamRequester(enableHttp2 bool, concurrency int) StreamRequester {
	s := &StreamRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
		reqCh:          make(chan *request.Request, concurrency), // 可调节的缓冲通道
		respCh:         make(chan *response.Response),
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	// 启动固定数量的 worker
	for i := 0; i < concurrency; i++ {
		go s.worker()
	}
	return s
}
