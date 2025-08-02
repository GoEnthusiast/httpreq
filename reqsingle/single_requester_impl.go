package reqsingle

import (
	"fmt"
	"github.com/GoEnthusiast/httpreq/builder"
	"github.com/GoEnthusiast/httpreq/transportsetting"
	"io"
	"net/http"
	"time"
)

type SingleRequesterImpl struct {
	TransportSetting *transportsetting.TransportSetting
	Client           *http.Client
}

func (s *SingleRequesterImpl) Do(req *Request) *Response {
	startTime := time.Now()
	resp := &Response{
		Request:   req,
		StartTime: startTime,
	}
	defer func() {
		resp.EndTime = time.Now()
		resp.Duration = resp.EndTime.Sub(startTime).Seconds()
	}()

	// 处理请求参数
	body, contentType, bodyE := builder.BuildRequestBody(req.ContentType, req.Body)
	if bodyE != nil {
		resp.Error = fmt.Errorf("build request body error: %s", bodyE.Error())
		return resp
	}
	// 构造 http 请求
	httpReq, err := http.NewRequest(string(req.Method), req.URL, body)
	if err != nil {
		resp.Error = fmt.Errorf("new http request error: %s", err.Error())
		return resp
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
	proxyE := s.TransportSetting.SetProxy(req.Proxy)
	if proxyE != nil {
		resp.Error = fmt.Errorf("set proxy error: %s", proxyE.Error())
		return resp
	}

	s.Client.Transport = s.TransportSetting.GetTransport()
	// 设置请求超时时间
	s.Client.Timeout = req.Timeout
	// 发送请求
	httpResp, err := s.Client.Do(httpReq)
	if err != nil {
		resp.Error = fmt.Errorf("do http request error: %s", err.Error())
		return resp
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Errorf("read response body error: %s", err.Error())
		return resp
	}

	resp.ResponseStatusCode = httpResp.StatusCode
	resp.ResponseBody = respBody
	return resp
}

func NewSingleRequester(enableHttp2 bool) SingleRequester {
	transportSetting := transportsetting.NewTransportSetting(enableHttp2)
	return &SingleRequesterImpl{
		TransportSetting: transportSetting,
		Client: &http.Client{
			Transport: transportSetting.GetTransport(),
		},
	}
}
