package core

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GoEnthusiast/httpreq/builder"
	"github.com/GoEnthusiast/httpreq/transportsetting"
	"github.com/GoEnthusiast/httpreq/types/request"
	"github.com/GoEnthusiast/httpreq/types/response"
)

// RequestHandler 公共请求处理器
type RequestHandler struct {
	*transportsetting.TransportSetting
	client *http.Client
}

// NewRequestHandler 创建新的请求处理器
func NewRequestHandler(enableHttp2 bool) *RequestHandler {
	transportSetting := transportsetting.NewTransportSetting(enableHttp2)
	return &RequestHandler{
		TransportSetting: transportSetting,
		client: &http.Client{
			Transport: transportSetting.GetTransport(),
		},
	}
}

// ProcessRequest 处理单个请求，返回响应
func (h *RequestHandler) ProcessRequest(req *request.Request) *response.Response {
	startTime := time.Now()
	resp := &response.Response{
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
	if proxyE := h.TransportSetting.SetProxy(req.Proxy); proxyE != nil {
		resp.Error = fmt.Errorf("set proxy error: %s", proxyE.Error())
		return resp
	}

	h.client.Transport = h.TransportSetting.GetTransport()
	// 设置请求超时时间
	h.client.Timeout = req.Timeout

	// 发送请求
	httpResp, err := h.client.Do(httpReq)
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

// SetTLS 设置 TLS 配置
func (h *RequestHandler) SetTLS(certPath, keyPath, caPath string) error {
	return h.TransportSetting.SetTLS(certPath, keyPath, caPath)
}

// SetTransport 设置自定义 transport
func (h *RequestHandler) SetTransport(transport *http.Transport) {
	h.TransportSetting.SetTransport(transport)
}

// SetProxy 设置代理
func (h *RequestHandler) SetProxy(proxies interface{}) error {
	return h.TransportSetting.SetProxy(proxies)
}

// SetMaxIdleConns 设置最大空闲连接数
func (h *RequestHandler) SetMaxIdleConns(maxIdleConns int) {
	h.TransportSetting.SetMaxIdleConns(maxIdleConns)
}

// SetMaxIdleConnsPerHost 设置每个主机允许的最大空闲连接数
func (h *RequestHandler) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	h.TransportSetting.SetMaxIdleConnsPerHost(maxIdleConnsPerHost)
}

// SetMaxConnsPerHost 设置每个主机允许的最大连接数
func (h *RequestHandler) SetMaxConnsPerHost(maxConnsPerHost int) {
	h.TransportSetting.SetMaxConnsPerHost(maxConnsPerHost)
}

// SetIdleConnTimeout 设置空闲连接超时时间
func (h *RequestHandler) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	h.TransportSetting.SetIdleConnTimeout(idleConnTimeout)
}

// SetTLSHandshakeTimeout 设置TLS握手超时时间
func (h *RequestHandler) SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	h.TransportSetting.SetTLSHandshakeTimeout(tlsHandshakeTimeout)
}

// SetExpectContinueTimeout 设置Expect: 100-continue 机制的超时时间
func (h *RequestHandler) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	h.TransportSetting.SetExpectContinueTimeout(expectContinueTimeout)
}
