package reqsingle

import (
	"errors"
	"fmt"
	"github.com/GoEnthusiast/httpreq/client"
	"io"
	"net/http"
	"time"
)

type SingleRequesterImpl struct {
	Client *client.Client
}

func (s *SingleRequesterImpl) Do(req *Request) (*Response, error) {
	startTime := time.Now()
	// 设置请求超时
	s.Client.SetTimeout(req.Timeout)
	// 设置代理 IP
	proxyE := s.Client.SetProxy(req.Proxy)
	if proxyE != nil {
		return nil, proxyE
	}
	// 处理请求参数
	body, contentType, bodyE := s.Client.BuildRequestBody(req.ContentType, req.Body)
	if bodyE != nil {
		return nil, bodyE
	}
	// 构造 http 请求
	httpReq, err := http.NewRequest(string(req.Method), req.URL, body)
	if err != nil {
		return nil, err
	}
	// 设置 content-type
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	// 设置请求头
	httpReq.Header = req.Header

	// 发送请求
	httpResp, err := s.Client.GetClient().Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("http status code is %d\nresponse body is %s", httpResp.StatusCode, string(respBody)))
	}

	endTime := time.Now()
	return &Response{
		Request:            req,
		ResponseStatusCode: httpResp.StatusCode,
		ResponseBody:       respBody,
		StartTime:          startTime,
		EndTime:            endTime,
		Duration:           endTime.Sub(startTime).Seconds(),
	}, nil
}

func NewSingleRequester(enableHttp2 bool) SingleRequester {
	c := client.New(enableHttp2)
	return &SingleRequesterImpl{
		Client: c,
	}
}
