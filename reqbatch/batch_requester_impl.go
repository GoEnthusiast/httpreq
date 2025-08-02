package reqbatch

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

type BatchRequesterImpl struct {
	*core.RequestHandler
}

func (s *BatchRequesterImpl) Do(reqs []*request.Request) []*response.Response {
	var (
		respCh = make(chan *response.Response, len(reqs)) // 带缓冲的通道，收集响应
	)

	for i := range reqs {
		req := reqs[i] // 避免 goroutine 闭包引用错误
		go func(r *request.Request) {
			resp := s.RequestHandler.ProcessRequest(r)
			respCh <- resp
		}(req)
	}

	// 收集结果
	responses := make([]*response.Response, 0, len(reqs))
	for i := 0; i < len(reqs); i++ {
		responses = append(responses, <-respCh)
	}
	return responses
}

func NewBatchRequester(enableHttp2 bool) BatchRequester {
	return &BatchRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
	}
}
