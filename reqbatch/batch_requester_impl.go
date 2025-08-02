package reqsingle

import (
	"fmt"
	"github.com/GoEnthusiast/httpreq/builder"
	"github.com/GoEnthusiast/httpreq/transportsetting"
	"io"
	"net/http"
	"time"
)

type BatchRequesterImpl struct {
	transportSetting *transportsetting.TransportSetting
	client           *http.Client
}

func (s *BatchRequesterImpl) GetTransportSetting() *transportsetting.TransportSetting {
	return s.transportSetting
}

func (s *BatchRequesterImpl) GetClient() *http.Client {
	return s.client
}

func (s *BatchRequesterImpl) Do(reqs []*Request) []*Response {
	var (
		respCh = make(chan *Response, len(reqs)) // 带缓冲的通道，收集响应
	)

	for i := range reqs {
		req := reqs[i] // 避免 goroutine 闭包引用错误
		go func(r *Request) {
			startTime := time.Now()
			resp := &Response{
				Request:   r,
				StartTime: startTime,
			}
			defer func() {
				resp.EndTime = time.Now()
				resp.Duration = resp.EndTime.Sub(startTime).Seconds()
				respCh <- resp
			}()

			// 处理请求参数
			body, contentType, bodyE := builder.BuildRequestBody(r.ContentType, r.Body)
			if bodyE != nil {
				resp.Error = fmt.Errorf("build request body error: %s", bodyE.Error())
				return
			}
			// 构造 http 请求
			httpReq, err := http.NewRequest(string(r.Method), r.URL, body)
			if err != nil {
				resp.Error = fmt.Errorf("new http request error: %s", err.Error())
				return
			}
			// 设置请求头
			if r.Header != nil {
				httpReq.Header = r.Header
			} else {
				httpReq.Header = make(http.Header)
			}
			// 设置 content-type
			if contentType != "" {
				httpReq.Header.Set("Content-Type", contentType)
			}
			// 设置代理 IP
			proxyE := s.transportSetting.SetProxy(r.Proxy)
			if proxyE != nil {
				resp.Error = fmt.Errorf("set proxy error: %s", proxyE.Error())
				return
			}

			s.client.Transport = s.transportSetting.GetTransport()
			// 设置请求超时时间
			s.client.Timeout = req.Timeout
			// 发送请求
			httpResp, err := s.client.Do(httpReq)
			if err != nil {
				resp.Error = fmt.Errorf("do http request error: %s", err.Error())
				return
			}
			defer httpResp.Body.Close()

			respBody, err := io.ReadAll(httpResp.Body)
			if err != nil {
				resp.Error = fmt.Errorf("read response body error: %s", err.Error())
				return
			}

			resp.ResponseStatusCode = httpResp.StatusCode
			resp.ResponseBody = respBody
		}(req)
	}

	// 收集结果
	responses := make([]*Response, 0, len(reqs))
	for i := 0; i < len(reqs); i++ {
		responses = append(responses, <-respCh)
	}
	return responses
}

func NewBatchRequester(enableHttp2 bool) BatchRequester {
	transportSetting := transportsetting.NewTransportSetting(enableHttp2)
	return &BatchRequesterImpl{
		transportSetting: transportSetting,
		client: &http.Client{
			Transport: transportSetting.GetTransport(),
		},
	}
}
