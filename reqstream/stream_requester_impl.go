package reqstream

import (
	"fmt"
	"github.com/GoEnthusiast/httpreq/builder"
	"github.com/GoEnthusiast/httpreq/transportsetting"
	"io"
	"net/http"
	"time"
)

type StreamRequesterImpl struct {
	*transportsetting.TransportSetting
	client *http.Client
	reqCh  chan *Request
	respCh chan *Response
}

func (s *StreamRequesterImpl) worker() {
	for req := range s.reqCh {
		s.handleRequest(req)
	}
}

func (s *StreamRequesterImpl) handleRequest(req *Request) {
	startTime := time.Now()
	resp := &Response{
		Request:   req,
		StartTime: startTime,
	}
	defer func() {
		resp.EndTime = time.Now()
		resp.Duration = resp.EndTime.Sub(startTime).Seconds()
		s.respCh <- resp
	}()
	// 处理请求参数
	body, contentType, bodyE := builder.BuildRequestBody(req.ContentType, req.Body)
	if bodyE != nil {
		resp.Error = fmt.Errorf("build request body error: %s", bodyE)
		return
	}
	// 构造 http 请求
	httpReq, err := http.NewRequest(string(req.Method), req.URL, body)
	if err != nil {
		resp.Error = fmt.Errorf("new http request error: %s", err)
		return
	}
	// 设置请求头
	if req.Header != nil {
		httpReq.Header = req.Header
	} else {
		httpReq.Header = make(http.Header)
	}
	// 设置 content-type
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	// 设置代理 IP
	if proxyE := s.TransportSetting.SetProxy(req.Proxy); proxyE != nil {
		resp.Error = fmt.Errorf("set proxy error: %s", proxyE)
		return
	}

	s.client.Transport = s.TransportSetting.GetTransport()
	// 设置请求超时时间
	s.client.Timeout = req.Timeout
	// 发送请求
	httpResp, err := s.client.Do(httpReq)
	if err != nil {
		resp.Error = fmt.Errorf("do http request error: %s", err)
		return
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Errorf("read response body error: %s", err)
		return
	}

	resp.ResponseStatusCode = httpResp.StatusCode
	resp.ResponseBody = respBody
}

func (s *StreamRequesterImpl) Do(req *Request) {
	s.reqCh <- req
}

func (s *StreamRequesterImpl) ResponseCh() <-chan *Response {
	return s.respCh
}

func NewStreamRequester(enableHttp2 bool, concurrency int) StreamRequester {
	transportSetting := transportsetting.NewTransportSetting(enableHttp2)
	s := &StreamRequesterImpl{
		TransportSetting: transportSetting,
		client: &http.Client{
			Transport: transportSetting.GetTransport(),
		},
		reqCh:  make(chan *Request, concurrency), // 可调节的缓冲通道
		respCh: make(chan *Response),
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
