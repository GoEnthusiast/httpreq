package reqsingle

import (
	"github.com/GoEnthusiast/httpreq/core"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

type SingleRequesterImpl struct {
	*core.RequestHandler
}

func (s *SingleRequesterImpl) Do(req *request.Request) *response.Response {
	return s.RequestHandler.ProcessRequest(req)
}

func NewSingleRequester(enableHttp2 bool) SingleRequester {
	return &SingleRequesterImpl{
		RequestHandler: core.NewRequestHandler(enableHttp2),
	}
}
