// Package core provides the core request handling functionality for HTTP requests
// 包 core 提供 HTTP 请求的核心请求处理功能
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

// RequestHandler is the core request processor that handles HTTP requests
// RequestHandler 是处理 HTTP 请求的核心请求处理器
type RequestHandler struct {
	*transportsetting.TransportSetting              // Transport configuration / 传输层配置
	client                             *http.Client // HTTP client instance / HTTP 客户端实例
}

// NewRequestHandler creates a new request handler with optional HTTP/2 support
// NewRequestHandler 创建一个新的请求处理器，支持可选的 HTTP/2
func NewRequestHandler(enableHttp2 bool) *RequestHandler {
	transportSetting := transportsetting.NewTransportSetting(enableHttp2)
	return &RequestHandler{
		TransportSetting: transportSetting,
		client: &http.Client{
			Transport: transportSetting.GetTransport(),
		},
	}
}

// ProcessRequest processes a single HTTP request and returns the response
// ProcessRequest 处理单个 HTTP 请求并返回响应
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

	// Build request body based on content type
	// 根据内容类型构建请求体
	body, contentType, bodyE := builder.BuildRequestBody(req.ContentType, req.Body)
	if bodyE != nil {
		resp.Error = fmt.Errorf("build request body error: %s", bodyE.Error())
		return resp
	}

	// Create HTTP request
	// 创建 HTTP 请求
	httpReq, err := http.NewRequest(string(req.Method), req.URL, body)
	if err != nil {
		resp.Error = fmt.Errorf("new http request error: %s", err.Error())
		return resp
	}

	// Set request headers
	// 设置请求头
	if req.Header != nil {
		httpReq.Header = req.Header
	} else {
		httpReq.Header = make(http.Header)
	}

	// Set content-type header
	// 设置 content-type 头部
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}

	// Configure proxy settings
	// 配置代理设置
	if proxyE := h.TransportSetting.SetProxy(req.Proxy); proxyE != nil {
		resp.Error = fmt.Errorf("set proxy error: %s", proxyE.Error())
		return resp
	}

	h.client.Transport = h.TransportSetting.GetTransport()
	// Set request timeout
	// 设置请求超时时间
	h.client.Timeout = req.Timeout

	// Execute HTTP request
	// 执行 HTTP 请求
	httpResp, err := h.client.Do(httpReq)
	if err != nil {
		resp.Error = fmt.Errorf("do http request error: %s", err.Error())
		return resp
	}
	defer httpResp.Body.Close()

	// Read response body
	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Errorf("read response body error: %s", err.Error())
		return resp
	}

	resp.ResponseStatusCode = httpResp.StatusCode
	resp.ResponseBody = respBody
	return resp
}

// SetTLS configures TLS settings with certificate files
// SetTLS 使用证书文件配置 TLS 设置
func (h *RequestHandler) SetTLS(certPath, keyPath, caPath string) error {
	return h.TransportSetting.SetTLS(certPath, keyPath, caPath)
}

// SetTransport sets a custom HTTP transport
// SetTransport 设置自定义 HTTP 传输层
func (h *RequestHandler) SetTransport(transport *http.Transport) {
	h.TransportSetting.SetTransport(transport)
}

// SetProxy configures proxy settings
// SetProxy 配置代理设置
func (h *RequestHandler) SetProxy(proxies interface{}) error {
	return h.TransportSetting.SetProxy(proxies)
}

// SetMaxIdleConns sets the maximum number of idle connections
// SetMaxIdleConns 设置最大空闲连接数
func (h *RequestHandler) SetMaxIdleConns(maxIdleConns int) {
	h.TransportSetting.SetMaxIdleConns(maxIdleConns)
}

// SetMaxIdleConnsPerHost sets the maximum number of idle connections per host
// SetMaxIdleConnsPerHost 设置每个主机的最大空闲连接数
func (h *RequestHandler) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) {
	h.TransportSetting.SetMaxIdleConnsPerHost(maxIdleConnsPerHost)
}

// SetMaxConnsPerHost sets the maximum number of connections per host
// SetMaxConnsPerHost 设置每个主机的最大连接数
func (h *RequestHandler) SetMaxConnsPerHost(maxConnsPerHost int) {
	h.TransportSetting.SetMaxConnsPerHost(maxConnsPerHost)
}

// SetIdleConnTimeout sets the idle connection timeout
// SetIdleConnTimeout 设置空闲连接超时时间
func (h *RequestHandler) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	h.TransportSetting.SetIdleConnTimeout(idleConnTimeout)
}

// SetTLSHandshakeTimeout sets the TLS handshake timeout
// SetTLSHandshakeTimeout 设置 TLS 握手超时时间
func (h *RequestHandler) SetTLSHandshakeTimeout(tlsHandshakeTimeout time.Duration) {
	h.TransportSetting.SetTLSHandshakeTimeout(tlsHandshakeTimeout)
}

// SetExpectContinueTimeout sets the Expect: 100-continue timeout
// SetExpectContinueTimeout 设置 Expect: 100-continue 超时时间
func (h *RequestHandler) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	h.TransportSetting.SetExpectContinueTimeout(expectContinueTimeout)
}

// SetDisableKeepAlives sets whether to disable HTTP Keep-Alive
// SetDisableKeepAlives 设置是否禁用 HTTP Keep-Alive
func (h *RequestHandler) SetDisableKeepAlives(disableKeepAlives bool) {
	h.TransportSetting.SetDisableKeepAlives(disableKeepAlives)
}
